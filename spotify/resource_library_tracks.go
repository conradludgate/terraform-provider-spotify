package spotify

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func resourceLibraryTracks() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLibraryTracksCreate,
		ReadContext:   resourceLibraryTracksRead,
		UpdateContext: resourceLibraryTracksUpdate,
		DeleteContext: resourceLibraryTracksDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"tracks": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of track IDs to save to your 'liked tracks'. *Note, if used incorrectly you may unlike all of your tracks - use with caution*",
			},
		},
	}
}

func resourceLibraryTracksCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	trackIDs := spotifyIdsInterface(d.Get("tracks").(*schema.Set).List())

	for _, rng := range batches(len(trackIDs), 100) {
		if err := client.AddTracksToLibrary(ctx, trackIDs[rng.Start:rng.End]...); err != nil {
			return diag.Errorf("AddTracksToLibrary: %s", err.Error())
		}
	}

	d.SetId("library")

	return nil
}

func resourceLibraryTracksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	trackIDs := schema.NewSet(schema.HashString, nil)

	tracks, err := client.CurrentUsersTracks(ctx)
	if err != nil {
		return diag.Errorf("CurrentUsersTracks: %s", err.Error())
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

func resourceLibraryTracksUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

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
				return diag.Errorf("AddTracksToLibrary: %s", err.Error())
			}
		}
		for _, rng := range batches(len(sub), 100) {
			if err := client.RemoveTracksFromLibrary(ctx, subTrackIDs[rng.Start:rng.End]...); err != nil {
				return diag.Errorf("AddTracksToLibrary: %s", err.Error())
			}
		}
	}

	return nil
}

func resourceLibraryTracksDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
