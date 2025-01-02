package roulette

import (
	"errors"
	"github.com/jamieyoung5/pooblet/internal/osm"
	"math/rand"
)

func Play(lat, long string, radius string) (*Pub, error) {
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
