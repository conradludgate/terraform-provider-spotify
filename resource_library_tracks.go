package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func resourceLibraryTracksCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	tracks := d.Get("tracks").([]interface{})
	trackIDs := make([]spotify.ID, len(tracks))
	for i, track := range tracks {
		trackIDs[i] = spotify.ID(track.(string))
	}

	for i := 0; i < len(tracks)/100; i++ {
		if err := client.AddTracksToLibrary(trackIDs[100*i : 100*i+100]...); err != nil {
			return fmt.Errorf("AddTracksToLibrary: %w", err)
		}
	}

	if len(tracks)%100 != 0 {
		if err := client.AddTracksToLibrary(trackIDs[100*(len(tracks)/100):]...); err != nil {
			return fmt.Errorf("AddTracksToLibrary: %w", err)
		}
	}

	return nil
}

func resourceLibraryTracksRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	trackIDs := schema.NewSet(schema.HashString, nil)

	tracks, err := client.CurrentUsersTracks()
	if err != nil {
		return fmt.Errorf("CurrentUsersTracks: %w", err)
	}
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs.Add(string(track.ID))
		}
		err = client.NextPage(tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

func resourceLibraryTracksUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	d.Partial(true)

	if d.HasChange("tracks") {
		old, new := d.GetChange("tracks")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := newSet.Difference(oldSet).List()
		sub := oldSet.Difference(newSet).List()

		addTrackIDs := make([]spotify.ID, len(add))
		for i, track := range add {
			addTrackIDs[i] = spotify.ID(track.(string))
		}

		subTrackIDs := make([]spotify.ID, len(sub))
		for i, track := range sub {
			subTrackIDs[i] = spotify.ID(track.(string))
		}

		for i := 0; i < len(add)/100; i++ {
			if err := client.AddTracksToLibrary(addTrackIDs[100*i : 100*i+100]...); err != nil {
				return fmt.Errorf("AddTracksToLibrary: %w", err)
			}
		}

		if len(add)%100 != 0 {
			if err := client.AddTracksToLibrary(addTrackIDs[100*(len(add)/100):]...); err != nil {
				return fmt.Errorf("AddTracksToLibrary: %w", err)
			}
		}

		for i := 0; i < len(sub)/100; i++ {
			if err := client.RemoveTracksFromLibrary(subTrackIDs[100*i : 100*i+100]...); err != nil {
				return fmt.Errorf("RemoveTracksFromLibrary: %w", err)
			}
		}

		if len(sub)%100 != 0 {
			if err := client.RemoveTracksFromLibrary(subTrackIDs[100*(len(sub)/100):]...); err != nil {
				return fmt.Errorf("RemoveTracksFromLibrary: %w", err)
			}
		}
	}

	d.SetPartial("tracks")

	d.Partial(false)
	return nil
}

func resourceLibraryTracksDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
