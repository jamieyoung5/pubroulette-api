package roulette

import (
	"os"
	"reflect"

	"github.com/jamieyoung5/pubroulette-api/pkg/googleplaces"
	"github.com/jamieyoung5/pubroulette-api/pkg/osmoverpass"
	"github.com/jamieyoung5/pubroulette-api/pkg/pub"
	"go.uber.org/zap"
)

type pubFinder interface {
	GetRandomPub(lat, lon, radius string) (*pub.Pub, error)
	GetAllAvailablePubs(lat, lon, radius string) ([]*pub.Pub, error)
}

func Play(lat, lon string, radius string, logger *zap.Logger) (*pub.Pub, error) {

	finder := getFinder(logger)

	logger.Info(
		"Starting new game",
		zap.String("latitude", lat),
		zap.String("longitude", lon),
		zap.String("radius", radius),
		zap.String("finder", reflect.TypeOf(finder).String()),
	)

	pub, err := finder.GetRandomPub(lat, lon, radius)
	if err != nil {
		logger.Error(
			"Error occurred while trying to find pub",
			zap.Error(err),
		)

		return nil, err
	}

	pub.Tags = filterTags(pub.Tags)

	logger.Info(
		"Found pub",
		zap.Any("pub", pub),
	)

	return pub, nil
}

func getFinder(logger *zap.Logger) pubFinder {
	usePlaces := os.Getenv("USE_GOOGLE_PLACES")
	if usePlaces == "true" {
		return googleplaces.NewClient(logger)
	} else {
		return osmoverpass.NewOverpassClient(logger)
	}
}
