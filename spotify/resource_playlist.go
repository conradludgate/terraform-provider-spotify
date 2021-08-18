package spotify

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func resourcePlaylist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaylistCreate,
		ReadContext:   resourcePlaylistRead,
		UpdateContext: resourcePlaylistUpdate,
		DeleteContext: resourcePlaylistDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Resource to manage a spotify playlist.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resulting playlist",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the resulting playlist",
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the playlist can be accessed publically",
			},
			"tracks": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of tracks for the playlist to contain",
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePlaylistCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		return diag.Errorf("GetCurrentUser: %s", err.Error())
	}

	userID := string(user.ID)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	public := d.Get("public").(bool)

	playlist, err := client.CreatePlaylistForUser(ctx, userID, name, description, public, false)
	if err != nil {
		return diag.Errorf("CreatePlaylist: %s", err.Error())
	}

	d.SetId(string(playlist.ID))

	trackIDs := spotifyIdsInterface(d.Get("tracks").([]interface{}))

	snapshotID := playlist.SnapshotID
	for _, rng := range batches(len(trackIDs), 100) {
		var err error
		snapshotID, err = client.AddTracksToPlaylist(ctx, playlist.ID, trackIDs[rng.Start:rng.End]...)
		if err != nil {
			return diag.Errorf("AddTracksToPlaylist: %s", err.Error())
		}
	}

	d.Set("snapshot_id", snapshotID)

	return nil
}

func resourcePlaylistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	playlistID := spotify.ID(d.Id())
	playlist, err := client.GetPlaylist(ctx, playlistID)

	if err != nil {
		return diag.Errorf("GetPlaylist: %s", err.Error())
	}

	d.Set("name", playlist.Name)
	d.Set("description", playlist.Description)
	d.Set("public", playlist.IsPublic)
	d.Set("snapshot_id", playlist.SnapshotID)

	tracks, err := client.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return diag.Errorf("GetPlaylistTracks: %s", err.Error())
	}

	trackIDs := make([]string, 0, tracks.Total)
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs = append(trackIDs, string(track.Track.ID))
		}
		err = client.NextPage(ctx, tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

func resourcePlaylistUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	id := spotify.ID(d.Id())
	if d.HasChanges("name", "description", "public") {
		err := client.ChangePlaylistNameAccessAndDescription(
			ctx,
			id,
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("public").(bool),
		)

		if err != nil {
			return diag.Errorf("ChangePlaylist: %s", err.Error())
		}
	}

	if d.HasChange("tracks") {
		new := spotifyIdsInterface(d.Get("tracks").([]interface{}))

		var err error
		var snapshotID string
		for i, rng := range batches(len(new), 100) {
			if i == 0 {
				err = client.ReplacePlaylistTracks(ctx, id, new[rng.Start:rng.End]...)
			} else {
				snapshotID, err = client.AddTracksToPlaylist(ctx, id, new[rng.Start:rng.End]...)
			}

			if err != nil {
				return diag.Errorf("update playlist tracks: %s", err.Error())
			}
		}

		if snapshotID != "" {
			d.Set("snapshot_id", snapshotID)
		}
	}

	return nil
}

func resourcePlaylistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	id := spotify.ID(d.Id())
	if err := client.UnfollowPlaylist(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
