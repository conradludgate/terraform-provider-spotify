package spotify

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func dataSourceAlbum() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlbumRead,

		Schema: map[string]*schema.Schema{
			"spotify_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"spotify_id", "url"},
				Description:  "Spotify ID of the album",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"spotify_id", "url"},
				Description:  "Spotify URL of the album",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Name of the album",
			},
			"artists": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The spotify IDs of the artists",
			},
		},
	}
}

func dataSourceAlbumRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	var id spotify.ID
	if u, ok := d.GetOk("url"); ok {
		u, err := url.Parse(u.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if !strings.HasPrefix(u.Path, "/album/") {
			return diag.FromErr(errors.New("URL did not point to a spotify album"))
		}
		id = spotify.ID(strings.TrimPrefix(u.Path, "/album/"))
	} else {
		id = spotify.ID(d.Get("spotify_id").(string))
	}

	album, err := client.GetAlbum(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", album.Name)

	artists := make([]interface{}, 0, len(album.Artists))
	for _, artist := range album.Artists {
		artists = append(artists, string(artist.ID))
	}
	d.Set("artists", artists)
	d.SetId(string(album.ID))

	return nil
}
