package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func resourceLibraryTracksCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	trackIDs := spotifyIdsInterface(d.Get("tracks").([]interface{}))

	for _, rng := range batches(len(trackIDs), 100) {
		if err := client.AddTracksToLibrary(trackIDs[rng.Start:rng.End]...); err != nil {
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

		addTrackIDs := spotifyIdsInterface(add)
		subTrackIDs := spotifyIdsInterface(sub)

		for _, rng := range batches(len(add), 100) {
			if err := client.AddTracksToLibrary(addTrackIDs[rng.Start:rng.End]...); err != nil {
				return fmt.Errorf("AddTracksToLibrary: %w", err)
			}
		}
		for _, rng := range batches(len(sub), 100) {
			if err := client.RemoveTracksFromLibrary(subTrackIDs[rng.Start:rng.End]...); err != nil {
				return fmt.Errorf("AddTracksToLibrary: %w", err)
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
