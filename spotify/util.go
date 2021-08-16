package spotify

import "github.com/zmb3/spotify/v2"

type Range struct {
	Start int
	End   int
}

func batches(length, batch int) []Range {
	if length <= 0 {
		return nil
	}

	var ranges []Range
	i := 0
	for ; i < (length - batch); i += batch {
		ranges = append(ranges, Range{i, i + batch})
	}
	return append(ranges, Range{i, length})
}

func spotifyIdsInterface(s []interface{}) []spotify.ID {
	output := make([]spotify.ID, len(s))
	for i, v := range s {
		output[i] = spotify.ID(v.(string))
	}
	return output
}
