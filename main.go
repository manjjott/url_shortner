package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
    ConnectMongo()
	ConnectRedis()

    r := gin.Default()

    r.POST("/shorten", ShortenURL)
    r.GET("/:shortURL", RedirectURL)

    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        <-sigChan
        DisconnectMongo()
		DisconnectRedis()
        os.Exit(0)
    }()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Println("Server running on port", port)
    r.Run(":" + port)
}
