package roulette

import (
	"slices"
)

var whitelist = map[string][]string{
	"Vegetarian Options":           []string{"diet:vegetarian"},
	"Vegan Options":                []string{"diet:vegan"},
	"Wheelchair Accessible":        []string{"wheelchair"},
	"Outdoor Seating":              []string{"outdoor_seating"},
	"Food":                         []string{"food"},
	"Wheelchair Accessible Toilet": []string{"toilets:wheelchair"},
	"Real Ale":                     []string{"real ale"},
	"Brewery":                      []string{"brewery"},
	"Micro Brewery":                []string{"microbrewery"},
	"LGBTQ+":                       []string{"lgbtq"},
}

func filterTags(tags []string) []string {
	seen := make(map[string]bool)
	var filteredTags []string

	for _, name := range tags {
		if alias := checkTag(name); alias != "" && !seen[alias] {
			filteredTags = append(filteredTags, alias)
			seen[alias] = true
		}
	}

	return filteredTags
}

func checkTag(tag string) string {
	for alias, tags := range whitelist {
		if slices.Contains(tags, tag) {
			return alias
		}
	}
	return ""
}
