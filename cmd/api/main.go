package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kevinhartarto/market-be/internal/database"
	"github.com/kevinhartarto/market-be/internal/server"
)

var ctx = context.Background()

func main() {

	db := database.StartDB()
	redis := server.StartRedis()
	app := server.NewHandler(db, redis)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3030"
	}
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Server Listening on http://localhost%s\n", addr)
	log.Println(redis.Ping(ctx))

	log.Fatal(app.Listen(addr))
	db.Close()
}
