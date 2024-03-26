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

func dataSourceArtist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArtistRead,

		Schema: map[string]*schema.Schema{
			"spotify_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"spotify_id", "url"},
				Description:  "Spotify ID of the artist",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"spotify_id", "url"},
				Description:  "Spotify URL of the artist",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Name of the artist",
			},
		},
	}
}

func dataSourceArtistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	var id spotify.ID
	if u, ok := d.GetOk("url"); ok {
		u, err := url.Parse(u.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if !strings.HasPrefix(u.Path, "/artist/") {
			return diag.FromErr(errors.New("URL did not point to a spotify artist"))
		}
		id = spotify.ID(strings.TrimPrefix(u.Path, "/artist/"))
	} else {
		id = spotify.ID(d.Get("spotify_id").(string))
	}

	artist, err := client.GetArtist(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", artist.Name); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(string(artist.ID))

	return nil
}
