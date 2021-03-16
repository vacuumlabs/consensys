package redis

import (
	"github.com/go-redis/redis/v8"
	"os"
)

var client *redis.Client

func init() {
	var REDIS_HOST string
	var REDIS_PORT string

	if REDIS_HOST = os.Getenv("REDIS_HOST"); REDIS_HOST == "" {
		REDIS_HOST = "localhost"
	}
	if REDIS_PORT = os.Getenv("REDIS_PORT"); REDIS_PORT == "" {
		REDIS_PORT = "6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: "",
		DB:       0,
	})
}
