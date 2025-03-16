package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var urlCollection *mongo.Collection

func ConnectMongo() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    mongoURI := os.Getenv("MONGO_URI")
    mongoDBName := os.Getenv("MONGO_DB")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(mongoURI)

    mongoClient, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Error connecting to MongoDB:", err)
    }
    client = mongoClient

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Could not ping MongoDB:", err)
    }

    fmt.Println("Connected to MongoDB!")

    urlCollection = client.Database(mongoDBName).Collection("urls")
}

func DisconnectMongo() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Disconnect(ctx); err != nil {
        log.Fatal("Error disconnecting MongoDB:", err)
    }
    fmt.Println("MongoDB connection closed.")
}
