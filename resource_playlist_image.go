package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func resourcePlaylistImage() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlaylistImageCreate,
		Read:   resourcePlaylistImageRead,
		Update: resourcePlaylistImageUpdate,
		Delete: resourcePlaylistImageDelete,

		Schema: map[string]*schema.Schema{
			"playlist_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "base64 encoded JPEG image",
			},
		},
	}
}

func resourcePlaylistImageCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := d.Get("playlist_id").(string)
	imageData := d.Get("image_data").(string)
	imageReader := strings.NewReader(imageData)

	err := client.SetPlaylistImage(spotify.ID(playlistID), imageReader)
	if err != nil {
		return fmt.Errorf("Upload Image: %w", err)
	}

	hash := md5.Sum([]byte(imageData))
	id := base64.RawStdEncoding.EncodeToString(hash[:])

	d.SetId(id)

	return nil
}

func resourcePlaylistImageRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePlaylistImageUpdate(d *schema.ResourceData, m interface{}) error {
	return resourcePlaylistImageCreate(d, m)
}

func resourcePlaylistImageDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
