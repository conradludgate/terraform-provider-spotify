package main

import (
	"github.com/conradludgate/terraform-provider-spotify/spotify"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//go:generate tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return spotify.Provider()
		},
	})
}
