package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URL struct {
    ShortURL string `json:"short_url" bson:"short_url`
    LongURL  string `json:"long_url" bson:"long_url"`
}

func ShortenURL(c *gin.Context) {
    var requestBody URL
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

	ctx := context.Background()
	var existingEntry bson.M
    err := urlCollection.FindOne(ctx, bson.M{"long_url": requestBody.LongURL}).Decode(&existingEntry)
	if err == nil {
        existingShortURL := "http://localhost:8080/" + existingEntry["short_url"].(string)
        c.JSON(http.StatusOK, gin.H{"short_url": existingShortURL})
        return
    } else if err != mongo.ErrNoDocuments {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    shortCode := uuid.New().String()[:6]
    shortURL := "http://localhost:8080/" + shortCode


    _, err = urlCollection.InsertOne(context.Background(), bson.M{
        "short_url":  shortCode,
        "long_url":   requestBody.LongURL,
        "created_at": time.Now(),
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

	// Store in Redis
	ctx = context.Background()
	err = redisClient.Set(ctx,shortCode,requestBody.LongURL, 24*time.Hour).Err()
	if err != nil {
		log.Println("Error setting the redis key", err)
	}

    c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func RedirectURL(c *gin.Context) {
    shortCode := c.Param("shortURL")
	ctx := context.Background()

	// check redis
	longURL, err := redisClient.Get(ctx,shortCode).Result()

    if err == redis.Nil {
        var result URL
        err = urlCollection.FindOne(context.Background(), bson.M{"short_url": shortCode}).Decode(&result)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
            return
        }
        longURL = result.LongURL

        err = redisClient.Set(ctx, shortCode, longURL, 24*time.Hour).Err()
        if err != nil {
            log.Println("Error caching URL in Redis:", err)
        }
    } else if err != nil {
        log.Println("Redis GET error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    fmt.Println("Redirecting to:", longURL)
    c.Redirect(http.StatusFound, longURL)
}
