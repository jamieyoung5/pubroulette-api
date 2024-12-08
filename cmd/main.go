package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jamieyoung5/pooblet/internal/handler"
	"github.com/jamieyoung5/pooblet/internal/osm"
	"github.com/jamieyoung5/pooblet/internal/roulette"
	"github.com/jamieyoung5/pooblet/internal/verification"
	"log"
	"net/http"
	"os"
	"sort"
)

// Function to execute the Overpass query
func getPubs2(lat, lon float64, rad int) {
	latitude, longitude, err := verification.VerifyLocation(lon, lat)
	if err != nil {
		panic(err)
	}

	radius, err := verification.VerifyRadius(rad)
	if err != nil {
		panic(err)
	}

	response, err := osm.GetAmenitiesInRadius(latitude, longitude, radius, "bar")
	if err != nil {
		panic(err)
	}

	tags := make(map[string]int)
	for _, element := range response {
		for tag, _ := range element.Tags {
			tags[tag] += 1
		}
	}

	keys := make([]string, 0, len(tags))

	for key := range tags {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return tags[keys[i]] > tags[keys[j]]
	})

	for _, k := range keys {
		msg := fmt.Sprintf("%s: %d", k, tags[k])
		fmt.Println(msg)
	}

}

func getPubs(lat, lon float64, rad int) {
	pubRoulette := roulette.NewGame()
	latitude, longitude, err := verification.VerifyLocation(lon, lat)
	if err != nil {
		panic(err)
	}
	radius, err := verification.VerifyRadius(rad)
	if err != nil {
		panic(err)
	}

	pub, err := pubRoulette.Play(latitude, longitude, radius)
	if err != nil {
		panic(err)
	}

	fmt.Println(pub)
}

/*
func main() {
	lat, lon := 55.953251, -3.188267
	radius := 500 // in meters

	getPubs(lat, lon, radius)
}*/

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r := mux.NewRouter()
	r.HandleFunc("/getPub", handler.GetPubHandler).Methods(http.MethodGet)
	r.Use(handler.CORSMiddleware)

	log.Println("Server is running on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
