package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLibrary() *schema.Resource {
	return &schema.Resource{
		Create: resourceLibraryCreate,
		Read:   resourceLibraryRead,
		Update: resourceLibraryUpdate,
		Delete: resourceLibraryDelete,

		Schema: map[string]*schema.Schema{
			"tracks": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceLibraryCreate(d *schema.ResourceData, m interface{}) error {
	// client := m.(*spotify.Client)

	return resourceLibraryTracksCreate(d, m)
}

func resourceLibraryRead(d *schema.ResourceData, m interface{}) error {
	// client := m.(*spotify.Client)

	return resourceLibraryTracksRead(d, m)
}

func resourceLibraryUpdate(d *schema.ResourceData, m interface{}) error {
	// client := m.(*spotify.Client)

	return resourceLibraryTracksUpdate(d, m)
}

func resourceLibraryDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
