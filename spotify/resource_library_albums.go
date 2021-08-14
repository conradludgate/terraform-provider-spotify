package spotify

import (
	"context"
	"fmt"

	"github.com/conradludgate/spotify/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLibraryAlbums() *schema.Resource {
	return &schema.Resource{
		Create: resourceLibraryAlbumsCreate,
		Read:   resourceLibraryAlbumsRead,
		Update: resourceLibraryAlbumsUpdate,
		Delete: resourceLibraryAlbumsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"albums": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of album IDs to save to your 'liked albums'. *Note, if used incorrectly you may unlike all of your albums - use with caution*",
			},
		},
	}
}

func resourceLibraryAlbumsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)
	ctx := context.Background()

	trackIDs := spotifyIdsInterface(d.Get("tracks").(*schema.Set).List())

	for _, rng := range batches(len(trackIDs), 100) {
		if err := client.AddTracksToLibrary(ctx, trackIDs[rng.Start:rng.End]...); err != nil {
			return fmt.Errorf("AddTracksToLibrary: %w", err)
		}
	}

	d.SetId("library")

	return nil
}

func resourceLibraryAlbumsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)
	ctx := context.Background()

	trackIDs := schema.NewSet(schema.HashString, nil)

	tracks, err := client.CurrentUsersTracks(ctx)
	if err != nil {
		return fmt.Errorf("CurrentUsersTracks: %w", err)
	}
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs.Add(string(track.ID))
		}
		err = client.NextPage(ctx, tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

func resourceLibraryAlbumsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)
	ctx := context.Background()

	if d.HasChange("tracks") {
		old, new := d.GetChange("tracks")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := newSet.Difference(oldSet).List()
		sub := oldSet.Difference(newSet).List()

		addTrackIDs := spotifyIdsInterface(add)
		subTrackIDs := spotifyIdsInterface(sub)

		for _, rng := range batches(len(add), 100) {
			if err := client.AddTracksToLibrary(ctx, addTrackIDs[rng.Start:rng.End]...); err != nil {
				return fmt.Errorf("AddTracksToLibrary: %w", err)
			}
		}
		for _, rng := range batches(len(sub), 100) {
			if err := client.RemoveTracksFromLibrary(ctx, subTrackIDs[rng.Start:rng.End]...); err != nil {
				return fmt.Errorf("AddTracksToLibrary: %w", err)
			}
		}
	}

	return nil
}

func resourceLibraryAlbumsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
