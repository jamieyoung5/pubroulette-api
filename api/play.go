package api

import (
	"encoding/json"
	"errors"
	"github.com/jamieyoung5/pooblet/pkg/osm"
	"github.com/jamieyoung5/pooblet/pkg/redis-client"
	"github.com/jamieyoung5/pooblet/pkg/roulette"
	"github.com/jamieyoung5/pooblet/pkg/whatpub"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// response codes
const (
	ErrGeneralRouletteError = "1"
	ErrNoPubsFound          = "2"
	ErrServerError          = "3"
	ErrInvalidInput         = "4"
)

var (
	allowedOrigins = map[string]bool{
		"https://www.pubroulette-web.vercel.app": true,
		"https://www.pubroulette.com":            true,
		"https://www.pubroulette.xyz":            true,
	}
	logger  *zap.Logger
	redisDb *redis.Client
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialise logger")
	}

	redisDb = redis_client.NewRedisDatabase()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}

	setCORSHeaders(w, r)
	w.Header().Set("Content-Type", "application/json")

	lat, lon, rad, err := parseQueryParams(r)
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

	scrapers := []roulette.Scraper{
		{Source: "whatpub.com", Scrape: whatpub.Scrape},
	}

	overpassApi := osm.NewOverpassApi(logger)

	game := roulette.NewGame(logger, scrapers, overpassApi, redisDb)

	pub, err := game.Play(latitude, longitude, radius)
	if err != nil {
		logger.Error("Failed to play roulette", zap.Error(err))
		code := roulette.GetErrorCode(err)
		errorResponse(w, http.StatusInternalServerError, strconv.Itoa(code), "")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(pub); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		errorResponse(w, http.StatusInternalServerError, ErrServerError, "")
	}
}

func parseQueryParams(r *http.Request) (float64, float64, int, error) {
	query := r.URL.Query()

	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		return 0, 0, 0, errors.New("invalid latitude")
	}

	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		return 0, 0, 0, errors.New("invalid longitude")
	}

	radius, err := strconv.Atoi(query.Get("radius"))
	if err != nil {
		return 0, 0, 0, errors.New("invalid radius")
	}

	return lat, lon, radius, nil
}

func errorResponse(w http.ResponseWriter, status int, code string, msg string) {
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{"error": code, "message": msg})
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}
