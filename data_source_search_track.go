package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSearchTrack() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSearchTrackRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"artists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"album": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"year": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"limit": {
				Type:     schema.TypeInt,
				Default:  10,
				Optional: true,
			},
			"tracks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"artists": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"album": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
			"track": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"artists": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"album": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
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

func strSlice(slice []interface{}) []string {
	output := make([]string, len(slice))
	for i, v := range slice {
		output[i] = v.(string)
	}
	return output
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
		return fmt.Errorf("Could not perform search [%v]: %w", queries, err)
	}

	if results.Tracks.Total == 0 {
		return fmt.Errorf("Could not find track")
	}

	var tracks []interface{}
	var ids []string
	for _, track := range results.Tracks.Tracks {
		var artists []interface{}
		for _, artist := range track.Artists {
			artists = append(artists, artist.ID.String())
		}
		tracks = append(tracks, map[string]interface{}{
			"id":      track.ID.String(),
			"name":    track.Name,
			"artists": artists,
			"album":   track.Album.ID.String(),
		})

		ids = append(ids, track.ID.String())
	}

	if *limit == 1 {
		d.Set("track", tracks[0])
	}

	d.Set("tracks", tracks)
	// Sets an id in the state
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
