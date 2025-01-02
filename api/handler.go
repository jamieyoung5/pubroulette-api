package api

import (
	"encoding/json"
	"github.com/jamieyoung5/pooblet/internal/roulette"
	"github.com/jamieyoung5/pooblet/internal/verification"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	rad, err := strconv.Atoi(query.Get("radius"))
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	latitude, longitude, err := verification.VerifyLocation(lon, lat)
	if err != nil {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		return
	}
	radius, err := verification.VerifyRadius(rad)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	pub, err := roulette.Play(latitude, longitude, radius)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pub)
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
