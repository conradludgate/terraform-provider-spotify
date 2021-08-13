package spotify

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func resourcePlaylist() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlaylistCreate,
		Read:   resourcePlaylistRead,
		Update: resourcePlaylistUpdate,
		Delete: resourcePlaylistDelete,
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

func resourcePlaylistCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	user, err := client.CurrentUser()
	if err != nil {
		return fmt.Errorf("GetCurrentUser: %w", err)
	}

	userID := string(user.ID)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	public := d.Get("public").(bool)

	playlist, err := client.CreatePlaylistForUser(userID, name, description, public)
	if err != nil {
		return fmt.Errorf("CreatePlaylist: %w", err)
	}

	d.SetId(string(playlist.ID))

	trackIDs := spotifyIdsInterface(d.Get("tracks").([]interface{}))

	snapshotID := playlist.SnapshotID
	for _, rng := range batches(len(trackIDs), 100) {
		var err error
		snapshotID, err = client.AddTracksToPlaylist(playlist.ID, trackIDs[rng.Start:rng.End]...)
		if err != nil {
			return fmt.Errorf("AddTracksToPlaylist: %w", err)
		}
	}

	d.Set("snapshot_id", snapshotID)

	return nil
}

func resourcePlaylistRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := spotify.ID(d.Id())
	playlist, err := client.GetPlaylist(playlistID)

	if err != nil {
		return fmt.Errorf("GetPlaylist: %w", err)
	}

	d.Set("name", playlist.Name)
	d.Set("description", playlist.Description)
	d.Set("public", playlist.IsPublic)
	d.Set("snapshot_id", playlist.SnapshotID)

	trackIDs := schema.NewSet(schema.HashString, nil)

	tracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		return fmt.Errorf("GetPlaylistTracks: %w", err)
	}
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs.Add(string(track.Track.ID))
		}
		err = client.NextPage(tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

func resourcePlaylistUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	id := spotify.ID(d.Id())
	if d.HasChanges("name", "description", "public") {
		err := client.ChangePlaylistNameAccessAndDescription(
			id,
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("public").(bool),
		)

		if err != nil {
			return fmt.Errorf("ChangePlaylist: %w", err)
		}
	}

	if d.HasChange("tracks") {
		new := spotifyIdsInterface(d.Get("tracks").([]interface{}))

		var err error
		var snapshotID string
		for i, rng := range batches(len(new), 100) {
			if i == 0 {
				err = client.ReplacePlaylistTracks(id, new[rng.Start:rng.End]...)
			} else {
				snapshotID, err = client.AddTracksToPlaylist(id, new[rng.Start:rng.End]...)
			}

			if err != nil {
				return fmt.Errorf("update playlist tracks: %w", err)
			}
		}

		if snapshotID != "" {
			d.Set("snapshot_id", snapshotID)
		}
	}

	return nil
}

func resourcePlaylistDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
