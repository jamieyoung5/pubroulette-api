package osm

type Places map[int]Element

type Element struct {
	Type  string            `json:"type"`
	ID    int               `json:"id"`
	Lat   float64           `json:"lat"`
	Lon   float64           `json:"lon"`
	Nodes []int             `json:"nodes"`
	Tags  map[string]string `json:"tags"`
}

type Response struct {
	Elements []Element `json:"elements"`
}
