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

	var trackIDs []string

	tracks, err := client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		return fmt.Errorf("GetPlaylistTracks: %w", err)
	}
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs = append(trackIDs, string(track.Track.ID))
		}
		err = client.NextPage(tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

// Currently just replaces the entire playlists tracks with the new ones
// Would be better to use snapshot_id, create a diff
// work out where to insert and delete tracks.
// Might be more efficient depending on the diff, and would be safer as to not
// delete a users entire playlist if an api call fails half way
func resourcePlaylistTracksUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	d.Partial(true)

	playlistID := d.Id()

	if d.HasChange("tracks") {
		tracks := d.Get("tracks").([]interface{})

		trackIDs := make([]spotify.ID, len(tracks))
		for i, track := range tracks {
			trackIDs[i] = spotify.ID(track.(string))
		}

		if len(tracks) <= 100 {
			if err := client.ReplacePlaylistTracks(spotify.ID(playlistID), trackIDs...); err != nil {
				return fmt.Errorf("ReplacePlaylistTracks: %w", err)
			}
		} else {
			if err := client.ReplacePlaylistTracks(spotify.ID(playlistID), trackIDs[:100]...); err != nil {
				return fmt.Errorf("ReplacePlaylistTracks: %w", err)
			}

			for i := 1; i < len(tracks)/100; i++ {
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
		}
	}

	d.SetPartial("tracks")

	d.Partial(false)
	return nil
}

func resourcePlaylistTracksDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := d.Id()
	if err := client.ReplacePlaylistTracks(spotify.ID(playlistID)); err != nil {
		return fmt.Errorf("AddTracksToPlaylist: %w", err)
	}

	return nil
}
