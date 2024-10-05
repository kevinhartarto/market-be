package server

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func StartRedis() *redis.Client {
	redisUrl := os.Getenv("REDIS_URL")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%v:%v", redisUrl, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0})

	return redisClient
}
