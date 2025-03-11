package googleplaces

import "github.com/jamieyoung5/pubroulette-api/pkg/pub"

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Viewport struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

type Geometry struct {
	Location Location `json:"location"`
	Viewport Viewport `json:"viewport"`
}

type Photo struct {
	Height           int      `json:"height"`
	HtmlAttributions []string `json:"html_attributions"`
	PhotoReference   string   `json:"photo_reference"`
	Width            int      `json:"width"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

type OpeningHours struct {
	OpenNow bool `json:"open_now"`
}

type Result struct {
	BusinessStatus      string       `json:"business_status"`
	Geometry            Geometry     `json:"geometry"`
	Icon                string       `json:"icon"`
	IconBackgroundColor string       `json:"icon_background_color"`
	IconMaskBaseUri     string       `json:"icon_mask_base_uri"`
	Name                string       `json:"name"`
	OpeningHours        OpeningHours `json:"opening_hours"`
	Photos              []Photo      `json:"photos"`
	PlaceID             string       `json:"place_id"`
	PlusCode            PlusCode     `json:"plus_code"`
	PriceLevel          int          `json:"price_level"`
	Rating              float64      `json:"rating"`
	Reference           string       `json:"reference"`
	Scope               string       `json:"scope"`
	Types               []string     `json:"types"`
	UserRatingsTotal    int          `json:"user_ratings_total"`
	Vicinity            string       `json:"vicinity"`
}

type PlacesAPIResponse struct {
	HTMLAttributions []string `json:"html_attributions"`
	Results          []Result `json:"results"`
	Status           string   `json:"status"`
}

func (r *Result) toPub() *pub.Pub {

	tags := []string{}
	if r.Types != nil {
		tags = r.Types
	}

	return &pub.Pub{
		Name: pub.Names{
			Name: r.Name,
		},
		Tags:         tags,
		Longitude:    r.Geometry.Location.Lng,
		Latitude:     r.Geometry.Location.Lat,
		Address:      r.Vicinity,
		Rating:       &r.Rating,
		TotalRatings: &r.UserRatingsTotal,
	}
}
