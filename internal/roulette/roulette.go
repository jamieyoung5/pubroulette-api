package roulette

import (
	"errors"
	"github.com/jamieyoung5/pooblet/internal/osm"
	"math/rand"
	"strconv"
)

type CrawlOptions struct {
	Length int
	Radius int
}

func Crawl(lat, long string, radius string, options CrawlOptions) ([]*Pub, error) {
	var pubs []*Pub

	for i := 0; i < options.Length; i++ {
		pub, err := Classic(lat, long, string(rune(options.Radius)))
		if err != nil {
			return nil, err
		}

		lat = strconv.FormatFloat(pub.Latitude, 'g', -1, 64)
		long = strconv.FormatFloat(pub.Longitude, 'g', -1, 64)
		pubs = append(pubs, pub)
	}

	return pubs, nil
}

func Classic(lat, long string, radius string) (*Pub, error) {
	var results []osm.Places
	for _, amenity := range amenities {
		result, err := osm.GetAmenitiesInRadius(lat, long, radius, amenity)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	places := combinePlaces(results)
	if len(places) <= 0 {
		return nil, errors.New("no places found")
	}
	randomPlace := places[getRandomPlace(places)]

	return parsePlaceToPub(randomPlace)
}

func getRandomPlace(places osm.Places) int {
	placeIndex := rand.Intn(len(places))
	for id, _ := range places {
		if placeIndex == 0 {
			return id
		}
		placeIndex--
	}

	return 0
}

func combinePlaces(placesByAmenity []osm.Places) (combinedPlaces osm.Places) {
	combinedPlaces = make(osm.Places)

	for _, places := range placesByAmenity {
		for id, place := range places {
			if _, ok := combinedPlaces[id]; !ok {
				combinedPlaces[id] = place
			}
		}
	}

	return combinedPlaces
}
