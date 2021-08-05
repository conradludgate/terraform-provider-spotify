package main

import (
	"github.com/conradludgate/terraform-provider-spotify/spotify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: spotify.Provider,
	})
}
