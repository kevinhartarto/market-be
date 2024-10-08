package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func StartRedis() *redis.Client {
	redisUrl := os.Getenv("REDIS_URL")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%v:%v", redisUrl, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}
	fmt.Println("Connected to redis: %v", pong)

	return redisClient
}
