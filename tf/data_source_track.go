package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSearchTrack() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSearchTrackRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"artists": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSearchTrackRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	artists := d.Get("artists").([]interface{})
	names := make([]string, len(artists))
	for i, v := range artists {
		names[i] = v.(string)
	}

	result, err := client.Search(
		fmt.Sprintf("%s %s",
			d.Get("name").(string),
			strings.Join(names, " "),
		),
		spotify.SearchTypeTrack,
	)
	if err != nil {
		return fmt.Errorf("Search: %w", err)
	}
	if len(result.Tracks.Tracks) == 0 {
		return fmt.Errorf("Track not found")
	}

	track := result.Tracks.Tracks[0]
	d.Set("name", track.Name)
	d.Set("artists", flattenServiceArtists(track.Artists))
	d.SetId(string(track.ID))

	return nil
}

func flattenServiceArtists(in []spotify.SimpleArtist) []interface{} {
	var out = make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v.Name
	}
	return out
}
