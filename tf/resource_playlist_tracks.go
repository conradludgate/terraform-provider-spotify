package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zmb3/spotify"
)

func resourcePlaylistTracksCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := d.Id()

	tracks := d.Get("tracks").([]interface{})
	trackIDs := make([]spotify.ID, len(tracks))
	for i, track := range tracks {
		trackIDs[i] = spotify.ID(track.(string))
	}

	for i := 0; i < len(tracks)/100; i++ {
		_, err := client.AddTracksToPlaylist(spotify.ID(playlistID), trackIDs[100*i:100*i+100]...)
		if err != nil {
			return fmt.Errorf("AddTracksToPlaylist: %w", err)
		}
	}

	if len(tracks)%100 != 0 {
		_, err := client.AddTracksToPlaylist(spotify.ID(playlistID), trackIDs[100*(len(tracks)/100):]...)
		if err != nil {
			return fmt.Errorf("AddTracksToPlaylist: %w", err)
		}
	}

	return nil
}

func resourcePlaylistTracksRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := d.Id()

	trackIDs := schema.NewSet(schema.HashString, nil)

	tracks, err := client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		return fmt.Errorf("GetPlaylistTracks: %w", err)
	}
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs.Add(string(track.Track.ID))
		}
		err = client.NextPage(tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

func resourcePlaylistTracksUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	d.Partial(true)

	playlistID := d.Id()

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
			_, err := client.AddTracksToPlaylist(spotify.ID(playlistID), addTrackIDs[100*i:100*i+100]...)
			if err != nil {
				return fmt.Errorf("AddTracksToPlaylist: %w", err)
			}
		}

		if len(add)%100 != 0 {
			_, err := client.AddTracksToPlaylist(spotify.ID(playlistID), addTrackIDs[100*(len(add)/100):]...)
			if err != nil {
				return fmt.Errorf("AddTracksToPlaylist: %w", err)
			}
		}

		for i := 0; i < len(sub)/100; i++ {
			_, err := client.RemoveTracksFromPlaylist(spotify.ID(playlistID), subTrackIDs[100*i:100*i+100]...)
			if err != nil {
				return fmt.Errorf("RemoveTracksFromPlaylist: %w", err)
			}
		}

		if len(sub)%100 != 0 {
			_, err := client.RemoveTracksFromPlaylist(spotify.ID(playlistID), subTrackIDs[100*(len(sub)/100):]...)
			if err != nil {
				return fmt.Errorf("RemoveTracksFromPlaylist: %w", err)
			}
		}
	}

	d.SetPartial("tracks")

	d.Partial(false)
	return nil
}

func resourcePlaylistTracksDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
