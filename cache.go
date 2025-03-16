package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// ConnectRedis establishes a connection to Redis.
func ConnectRedis() {
    // Load environment variables (if not loaded already)
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not loaded")
    }

    redisAddr := os.Getenv("REDIS_URL")
    if redisAddr == "" {
        redisAddr = "localhost:6379" // fallback value
    }

    // Initialize the Redis client.
    redisClient = redis.NewClient(&redis.Options{
        Addr: redisAddr,
        // Add Password: "yourpassword" if needed.
        DB: 0, // default DB
    })

    // Test connectivity.
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Error connecting to Redis:", err)
    }
    fmt.Println("Connected to Redis at", redisAddr)
}

// DisconnectRedis cleanly closes the Redis connection.
func DisconnectRedis() {
    err := redisClient.Close()
    if err != nil {
        log.Println("Error closing Redis:", err)
    }
    fmt.Println("Redis connection closed.")
}
