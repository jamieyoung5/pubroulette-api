package osm

type tagFilter struct {
	alias  string
	filter func(value string) bool
}

var whitelist = map[string]tagFilter{
	"diet:vegetarian":    {alias: "Vegetarian Options", filter: boolean},
	"diet:vegan":         {alias: "Vegan Options", filter: boolean},
	"wheelchair":         {alias: "Wheelchair Accessible", filter: boolean},
	"outdoor_seating":    {alias: "Outdoor Seating", filter: boolean},
	"food":               {alias: "Food", filter: boolean},
	"toilets:wheelchair": {alias: "Wheelchair Accessible Toilet", filter: boolean},
	"real ale":           {alias: "Real Ale", filter: presence},
	"brewery":            {alias: "Brewery", filter: boolean},
	"microbrewery":       {alias: "Micro Brewery", filter: boolean},
	"lgbtq":              {alias: "LGBTQ+", filter: primary},
}

func FilterTags(tags map[string]string) {

	for name, alias := range tags {

		if _, ok := whitelist[name]; !ok {
			delete(tags, name)
			continue
		}

		if !whitelist[name].filter(alias) {
			delete(tags, name)
			continue
		}

		tags[name] = whitelist[name].alias
	}
}

func boolean(value string) bool {
	return value != "no"
}

func presence(value string) bool {
	return true
}

func primary(value string) bool {
	return value == "primary"
}
