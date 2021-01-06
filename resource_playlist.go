package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func resourcePlaylist() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlaylistCreate,
		Read:   resourcePlaylistRead,
		Update: resourcePlaylistUpdate,
		Delete: resourcePlaylistDelete,

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
				Description: "A list of tracks for the playlist to contain",
			},
			// "snapshot_id": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
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

	// d.Set("snapshot_id", playlist.SnapshotID)

	d.SetId(string(playlist.ID))

	return resourcePlaylistTracksCreate(d, m)
}

func resourcePlaylistRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := d.Id()
	playlist, err := client.GetPlaylist(spotify.ID(playlistID))

	if err != nil {
		return fmt.Errorf("GetPlaylist: %w", err)
	}

	d.Set("name", playlist.Name)
	d.Set("description", playlist.Description)
	d.Set("public", playlist.IsPublic)

	return resourcePlaylistTracksRead(d, m)
}

func resourcePlaylistUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	d.Partial(true)

	if d.HasChanges("name", "description", "public") {
		err := client.ChangePlaylistNameAccessAndDescription(
			spotify.ID(d.Id()),
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("public").(bool),
		)

		if err != nil {
			return fmt.Errorf("ChangePlaylist: %w", err)
		}

		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("public")
	}

	return resourcePlaylistTracksUpdate(d, m)
}

func resourcePlaylistDelete(d *schema.ResourceData, m interface{}) error {
	return resourcePlaylistTracksDelete(d, m)
}
