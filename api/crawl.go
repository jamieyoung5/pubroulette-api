package api

import (
	"encoding/json"
	"github.com/jamieyoung5/pubroulette-api/pkg/roulette"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, r)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	setCORSHeaders(w, r)
	w.Header().Set("Content-Type", "application/json")

	lat, lon, rad, leng, err := parseCrawlQueryParams(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, ErrInvalidInput, err.Error())
		return
	}

	latitude, longitude, err := roulette.ValidateLocation(lon, lat)
	if err != nil {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		return
	}
	radius, err := roulette.ValidateRadius(rad)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}
	length, err := roulette.ValidateLength(leng)
	if err != nil {
		http.Error(w, "Invalid length", http.StatusBadRequest)
		return
	}

	pubs, err := roulette.Crawl(latitude, longitude, radius, length, logger)
	if err != nil {
		logger.Error("Failed to play roulette", zap.Error(err))
		code := roulette.GetErrorCode(err)
		errorResponse(w, http.StatusInternalServerError, strconv.Itoa(code), "")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(pubs); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, ErrServerError, "")
	}
}

func parseCrawlQueryParams(r *http.Request) (float64, float64, int, int, error) {
	query := r.URL.Query()

	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	radius, err := strconv.Atoi(query.Get("radius"))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	length, err := strconv.Atoi(query.Get("length"))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return lat, lon, radius, length, nil
}
