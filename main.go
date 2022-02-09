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
	// p := "_log/dev.log"
	// os.MkdirAll(filepath.Dir(p), 0777)
	// f, err := os.Create(p)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// gin.DefaultWriter = io.MultiWriter(f)
	// log.SetOutput(f)
	// defer f.Close()

	// handle redis conn err
	handler.Rd = Rd

	r := gin.Default()
	r.Use(mdw.CORS())
	// r.Use(mdw.JwtAuth())

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
