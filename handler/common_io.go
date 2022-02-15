package handler

import (
	"encoding/json"
	"errors"
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

type PubSubIO interface {
	decode(i *[]byte) (PubSubIO, error)
}

func decodeJSON(i *[]byte, ps PubSubIO) (PubSubIO, error) {
	err := json.Unmarshal(*i, ps)
	if err != nil {
		return ps, err
	}
	return ps, nil
}

func (m ChInit) decode(i *[]byte) (PubSubIO, error) {
	return decodeJSON(i, &m)
}

func decode(i *[]byte, t interface{}) (*ChInit, *Msg, error) {
	switch tp := t.(type) {
	case ChInit:
		err := json.Unmarshal(*i, &tp)
		if err != nil {
			return &tp, &Msg{}, err
		}
		return &tp, &Msg{}, nil
	case Msg:
		err := json.Unmarshal(*i, &tp)
		if err != nil {
			return &ChInit{}, &tp, err
		}
		return &ChInit{}, &tp, nil
	default:
		return &ChInit{}, &Msg{}, errors.New("decode error for message")
	}
}

// func decodeToChInit(i *[]byte) (*ChInit, error) {
// 	var chInit ChInit
// 	err := json.Unmarshal(*i, &chInit)
// 	if err != nil {
// 		return &chInit, err
// 	}
// 	return &chInit, nil
// }

// func decodeToMsg(i *[]byte) (*Msg, error) {
// 	var msg Msg
// 	err := json.Unmarshal(*i, &msg)
// 	if err != nil {
// 		return &msg, err
// 	}
// 	return &msg, nil
// }
