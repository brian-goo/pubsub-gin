package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func LogErr(c *gin.Context, err interface{}) {
	switch e := err.(type) {
	// this is for standard error type eg. from db lib
	case error:
		c.Error(e)
	case string:
		c.Error(errors.New(e))
	default:
		c.Error(errors.New("error occurred: details unknown"))
	}
}
