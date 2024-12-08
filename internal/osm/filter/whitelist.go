package filter

var ValidTags = []Tag{
	{
		Name:   "diet:vegetarian",
		Alias:  "Vegetarian Options",
		Filter: filterBoolean,
	},
	{
		Name:   "diet:vegan",
		Alias:  "Vegan Options",
		Filter: filterBoolean,
	},
	{
		Name:   "wheelchair",
		Alias:  "Wheelchair Accessible",
		Filter: filterBoolean,
	},
	{
		Name:   "outdoor_seating",
		Alias:  "Outdoor Seating",
		Filter: filterBoolean,
	},
	{
		Name:   "food",
		Alias:  "Food",
		Filter: filterBoolean,
	},
	{
		Name:   "toilets:wheelchair",
		Alias:  "Wheelchair Accessible Toilet",
		Filter: filterBoolean,
	},
	{
		Name:   "real ale",
		Alias:  "Real Ale",
		Filter: filterPresence,
	},
	{
		Name:   "brewery",
		Alias:  "Brewery",
		Filter: filterBoolean,
	},
	{
		Name:   "microbrewery",
		Alias:  "Micro Brewery",
		Filter: filterBoolean,
	},
	{
		Name:   "lgbtq",
		Alias:  "LGBTQ+",
		Filter: filterPrimary,
	},
}

type Tag struct {
	Name   string
	Alias  string
	Filter func(value string) bool
}

type Names struct {
	AltName string
	Name    string
	OldName string
}
