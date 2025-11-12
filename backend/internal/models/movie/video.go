package movie

import (
	"slices"
	"strings"
)

type Video struct {
	Id          string `json:"id"`
	ISO639_1    string `json:"iso_639_1"`
	ISO3166_1   string `json:"iso_3166_1"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Site        string `json:"site"`
	Size        uint   `json:"size"` // defaults to 0
	Type        string `json:"type"`
	Official    bool   `json:"official"` // defaults to true
	PublishedAt string `json:"published_at"`
}

type Videos []Video

func (v *Videos) FilterTrailersAndTeasersOnYT() {
	var filtered Videos
	for _, video := range *v {
		if video.Site == "YouTube" && (video.Type == "Trailer" || video.Type == "Teaser") && video.Official {
			filtered = append(filtered, video)
		}
	}

	*v = filtered
}

func (v Videos) SortByRelevance() {
	weight := func(v Video) int {
		if v.Type != "Trailer" {
			return 0
		}
		if !strings.Contains(strings.ToLower(v.Name), "official trailer") {
			return 2
		}
		return 1
	}

	slices.SortStableFunc(v, func(a, b Video) int {
		return weight(b) - weight(a)
	})
}
