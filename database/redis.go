package database

import (
	"log"

	"github.com/go-redis/redis"
)

var RDB *redis.Client

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping().Result()
	if err != nil {
		log.Panic(err)
	}
}
