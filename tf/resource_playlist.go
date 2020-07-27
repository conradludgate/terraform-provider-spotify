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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// "collaborative": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// },
			"tracks": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
	return nil
}
