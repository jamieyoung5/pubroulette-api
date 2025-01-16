package roulette

import (
	"github.com/jamieyoung5/pooblet/internal/osm"
	"github.com/jamieyoung5/pooblet/internal/pub"
	"maps"
	"math/rand"
	"slices"
)

type ScraperFunc func(name string) (pub.Pub, error)

const parsingAttempts = 3

func Play(lat, long string, radius string, scrapers []ScraperFunc) (*pub.Pub, error) {
	places := make(osm.Places)

	// get list of all pubs in radius of lat/lon for each osm pub-related amenity
	for _, amenity := range osm.PubAmenities {
		result, err := osm.GetAmenitiesInRadius(lat, long, radius, amenity)
		if err != nil {
			return nil, err
		}
		maps.Copy(places, result)
	}

	if len(places) <= 0 {
		return nil, ErrNoPubsFound
	}

	// find a random pub from gathered places,
	// with a max of 3 attempts to allow for potential data anomalies
	for i := range parsingAttempts {
		randomPub, err := getRandomPub(places, scrapers)
		if err == nil {
			return randomPub, nil
		}

		delete(places, i)
	}

	return nil, ErrParsingFailure
}

func getRandomPub(places osm.Places, scrapers []ScraperFunc) (*pub.Pub, error) {
	placeId := randomPlace(places)

	randPub, err := pub.OsmElementToPub(places[placeId])
	if err != nil {
		return nil, err
	}

	// scrape for additional data about pub
	for _, scraper := range scrapers {
		result, scrapingErr := scraper(randPub.Name.Name)
		if scrapingErr != nil {
			//TODO: log here
			continue
		}

		pub.Merge(randPub, result)
	}

	return randPub, nil
}

func randomPlace(places osm.Places) int {
	keys := slices.Collect(maps.Keys(places))
	return keys[rand.Intn(len(keys))]
}
