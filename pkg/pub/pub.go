package pub

type Pub struct {
	Tags         []string
	Longitude    float64
	Latitude     float64
	Address      string
	Name         Names
	Rating       *float64
	TotalRatings *int
}

type Names struct {
	AltName string
	Name    string
	OldName string
}

const (
	usePlacesEnvVar = "USE_GOOGLE_PLACES"
)

func Merge(subject *Pub, merging Pub) {
	if merging.Tags != nil {
		subject.Tags = merging.Tags
	}

	if merging.Name.OldName != "" {
		subject.Name.OldName = merging.Name.OldName
	}

	if merging.Name.AltName != "" {
		subject.Name.AltName = merging.Name.AltName
	}
}
