package handler

import (
	"github.com/gin-gonic/gin"
)

func GetPing(c *gin.Context) {

	LogErr(c, "tt")

	setResponse(c, gin.H{"status": "ok"})
	// c.JSON(getResponse(gin.H{"status": sub}))
}
