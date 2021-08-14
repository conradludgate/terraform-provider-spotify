package spotify

import (
	"context"

	"github.com/conradludgate/spotify/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLibraryAlbums() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLibraryAlbumsCreate,
		ReadContext:   resourceLibraryAlbumsRead,
		UpdateContext: resourceLibraryAlbumsUpdate,
		DeleteContext: resourceLibraryAlbumsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"albums": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of track IDs to save to your 'liked albums'. *Note, if used incorrectly you may unlike all of your albums - use with caution*",
			},
		},
	}
}

func resourceLibraryAlbumsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	trackIDs := spotifyIdsInterface(d.Get("albums").(*schema.Set).List())

	for _, rng := range batches(len(trackIDs), 100) {
		if err := client.AddAlbumsToLibrary(ctx, trackIDs[rng.Start:rng.End]...); err != nil {
			return diag.Errorf("AddAlbumsToLibrary: %s", err.Error())
		}
	}

	d.SetId("library")

	return nil
}

func resourceLibraryAlbumsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	trackIDs := schema.NewSet(schema.HashString, nil)

	Albums, err := client.CurrentUsersAlbums(ctx)
	if err != nil {
		return diag.Errorf("CurrentUsersAlbums: %s", err.Error())
	}
	for err == nil {
		for _, track := range Albums.Albums {
			trackIDs.Add(string(track.ID))
		}
		err = client.NextPage(ctx, Albums)
	}

	d.Set("albums", trackIDs)

	return nil
}

func resourceLibraryAlbumsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	if d.HasChange("albums") {
		old, new := d.GetChange("albums")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := newSet.Difference(oldSet).List()
		sub := oldSet.Difference(newSet).List()

		addTrackIDs := spotifyIdsInterface(add)
		subTrackIDs := spotifyIdsInterface(sub)

		for _, rng := range batches(len(add), 100) {
			if err := client.AddAlbumsToLibrary(ctx, addTrackIDs[rng.Start:rng.End]...); err != nil {
				return diag.Errorf("AddAlbumsToLibrary: %s", err.Error())
			}
		}
		for _, rng := range batches(len(sub), 100) {
			if err := client.RemoveAlbumsFromLibrary(ctx, subTrackIDs[rng.Start:rng.End]...); err != nil {
				return diag.Errorf("AddAlbumsToLibrary: %s", err.Error())
			}
		}
	}

	return nil
}

func resourceLibraryAlbumsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
