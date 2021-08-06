package spotify

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSearchTrack() *schema.Resource {
	trackResource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the track",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the track",
			},
			"artists": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IDs of the artists",
			},
			"album": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the album that the track appears on",
			},
		},
	}

	return &schema.Resource{
		Read: dataSourceSearchTrackRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the track",
			},
			"artist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the artist",
			},
			"album": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the album",
			},
			"year": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Year of release",
			},
			"limit": {
				Type:     schema.TypeInt,
				Default:  10,
				Optional: true,
			},
			"explicit": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Filter to allow explicit tracks",
			},
			"tracks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trackResource,
				Description: "List of tracks found",
			},
		},
	}
}

func addSearchTerm(queries []string, key, field string) []string {
	if field == "" {
		return queries
	}
	return append(queries, fmt.Sprintf("%s:%s", key, field))
}

func dataSourceSearchTrackRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	var queries []string
	queries = addSearchTerm(queries, "track", d.Get("name").(string))
	queries = addSearchTerm(queries, "artist", d.Get("artist").(string))
	queries = addSearchTerm(queries, "album", d.Get("album").(string))
	queries = addSearchTerm(queries, "year", d.Get("year").(string))

	limit := d.Get("limit").(int)

	results, err := client.SearchOpt(strings.Join(queries, " "), spotify.SearchTypeTrack, &spotify.Options{
		Limit: &limit,
	})

	if err != nil {
		return fmt.Errorf("could not perform search [%v]: %w", queries, err)
	}

	var tracks []interface{}
	for _, track := range results.Tracks.Tracks {
		var artists []interface{}
		for _, artist := range track.Artists {
			artists = append(artists, artist.ID.String())
		}

		trackData := map[string]interface{}{
			"id":      track.ID.String(),
			"name":    track.Name,
			"artists": artists,
			"album":   track.Album.ID.String(),
		}
		if track.Explicit && d.Get("explicit").(bool) {
			tracks = append(tracks, trackData)
		} else if !track.Explicit {
			tracks = append(tracks, trackData)
		}
	}

	d.Set("tracks", tracks)

	// Sets an id in the state
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
