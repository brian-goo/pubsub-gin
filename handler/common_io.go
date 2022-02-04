package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getResponse(res interface{}) (int, interface{}) {
	return http.StatusOK, res
}

// func getErrorResponse(res interface{}) (int, interface{}) {
// 	return http.StatusInternalServerError, res
// }

func setResponse(c *gin.Context, res interface{}) {
	c.JSON(getResponse(res))
}

// func setErrorResponse(c *gin.Context, res interface{}) {
// 	c.JSON(getErrorResponse(res))
// }

// func encode(i interface{}) ([]byte, error) {
// 	data, err := json.Marshal(i)
// 	if err != nil {
// 		return []byte(`{"status": "json encode failed"}`), err
// 	}
// 	return data, nil
// }

func decodeToChInit(i *[]byte) (*ChInit, error) {
	var chInit ChInit
	err := json.Unmarshal(*i, &chInit)
	if err != nil {
		return &chInit, err
	}
	return &chInit, nil
}

func decodeToMsg(i *[]byte) (*Msg, error) {
	var msg Msg
	err := json.Unmarshal(*i, &msg)
	if err != nil {
		return &msg, err
	}
	return &msg, nil
}
