package main

import (
	"os"
	"ps/handler"

	"github.com/gin-gonic/gin"
)

func router(r *gin.Engine) {
	if os.Getenv("APP_ENV") == "LOCAL" {
		r.Static("/js", "./_js")
	}

	r.GET("/ping", handler.GetPing)
	// r.POST("/ping", handler.GetPing)

	r.GET("/ws", handler.GetWs)
}
