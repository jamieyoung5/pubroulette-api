package redis_client

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

const (
	endpointENV = "REDIS_ENDPOINT"
	portENV     = "REDIS_PORT"
	passwordENV = "REDIS_PASSWORD"
)

func NewRedisDatabase() *redis.Client {
	url := buildConnectionUrl()

	opt, _ := redis.ParseURL(url)
	return redis.NewClient(opt)
}

func buildConnectionUrl() string {
	endpoint := os.Getenv(endpointENV)
	if endpoint == "" {
		panic(endpointENV + " not set")
	}

	port := os.Getenv(portENV)
	if port == "" {
		panic(portENV + " not set")
	}

	password := os.Getenv(passwordENV)
	if password == "" {
		panic(passwordENV + " not set")
	}

	return fmt.Sprintf("redis:default//%s:%s@%s:%s", password, password, endpoint, port)
}
