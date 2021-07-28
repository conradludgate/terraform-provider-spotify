package spotify

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"artists": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names of the artists",
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
			"track": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        trackResource,
				Description: "Convenience option for tracks[0]. Only set if limit = 1",
			},
		},
	}
}

func addSearchTerm(queries []string, key, field string) []string {
	if field == "" {
		return queries
	}
	if strings.Contains(field, " ") {
		return append(queries, fmt.Sprintf("%s:\"%s\"", key, field))
	}
	return append(queries, fmt.Sprintf("%s:%s", key, field))
}

func dataSourceSearchTrackRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	var queries []string

	queries = addSearchTerm(queries, "track", d.Get("name").(string))

	artists := d.Get("artists").([]interface{})
	for _, artist := range artists {
		queries = addSearchTerm(queries, "artist", artist.(string))
	}

	queries = addSearchTerm(queries, "album", d.Get("album").(string))

	queries = addSearchTerm(queries, "year", d.Get("year").(string))

	var limit *int
	if lim, ok := d.GetOk("limit"); ok {
		lim := lim.(int)
		limit = &lim
	}

	results, err := client.SearchOpt(strings.Join(queries, " "), spotify.SearchTypeTrack, &spotify.Options{
		Limit: limit,
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

	if len(tracks) == 0 {
		return fmt.Errorf("could not find track")
	}

	if *limit == 1 {
		d.Set("track", tracks[0])
	}

	d.Set("tracks", tracks)
	// Sets an id in the state
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
