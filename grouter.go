package main

import (
	"ps/handler"

	"github.com/gin-gonic/gin"
)

func router(r *gin.Engine) {
	r.GET("/ping", handler.GetPing)
	// r.POST("/ping", handler.GetPing)

	r.GET("/ws", handler.GetWs)
}
