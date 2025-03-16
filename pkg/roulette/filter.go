package roulette

import (
	"slices"
)

var whitelist = map[string][]string{
	"Vegetarian Options":           {"diet:vegetarian"},
	"Vegan Options":                {"diet:vegan"},
	"Wheelchair Accessible":        {"wheelchair"},
	"Outdoor Seating":              {"outdoor_seating"},
	"Food":                         {"food", "restaurant", "cafe"},
	"Wheelchair Accessible Toilet": {"toilets:wheelchair"},
	"Real Ale":                     {"real ale"},
	"Brewery":                      {"brewery"},
	"Micro Brewery":                {"microbrewery"},
	"LGBTQ+":                       {"lgbtq"},
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
