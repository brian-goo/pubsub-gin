package main

import (
	"log"
	"net/http"
	"os"
	"ps/handler"
	mdw "ps/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var Rd = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDR"),
})

func main() {
	// handle redis conn err
	handler.Rd = Rd

	r := gin.Default()
	r.Use(mdw.CORS())
	// r.Use(mdw.AuthAPIKey())
	// r.Use(mdw.Auth0())

	router(r)

	server := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Println("server starting...", os.Getenv("PORT"))
	log.Fatal("failed to start server: ", server.ListenAndServe())
}
