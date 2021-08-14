package spotify

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/conradludgate/spotify/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTrack() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrackRead,

		Schema: map[string]*schema.Schema{
			"spotify_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"spotify_id", "url"},
				Description:  "Spotify ID of the track",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"spotify_id", "url"},
				Description:  "Spotify URL of the track",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Name of the track",
			},
			"artists": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The spotify IDs of the artists",
			},
			"album": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The spotify ID of the album",
			},
		},
	}
}

func dataSourceTrackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	var id spotify.ID
	if u, ok := d.GetOk("url"); ok {
		u, err := url.Parse(u.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if !strings.HasPrefix(u.Path, "/track/") {
			return diag.FromErr(errors.New("URL did not point to a spotify track"))
		}
		id = spotify.ID(strings.TrimPrefix(u.Path, "/track/"))
	} else {
		id = spotify.ID(d.Get("spotify_id").(string))
	}

	track, err := client.GetTrack(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", track.Name)
	d.Set("album", string(track.Album.ID))

	artists := make([]interface{}, 0, len(track.Artists))
	for _, artist := range track.Artists {
		artists = append(artists, string(artist.ID))
	}
	d.Set("artists", artists)
	d.SetId(string(track.ID))

	return nil
}
