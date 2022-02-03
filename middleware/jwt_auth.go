package middleware

import (
	"errors"
	"net/http"
	"os"
	"ps/handler"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var hmacKey = []byte(os.Getenv("JWT_KEY"))
var issuer = os.Getenv("JWT_ISSUER")

type jwToken string

func (jw *jwToken) decode() (string, error) {
	token, err := jwt.Parse(string(*jw), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("wrong algorithm")
		}
		return hmacKey, nil
	})

	if token == nil {
		return "", errors.New("token is null")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && claims["iss"] == issuer && token.Valid {
		if id, ok := claims["id"].(string); ok {
			return id, nil
		} else {
			return "", errors.New("no id claim found")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", errors.New("invalid token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "", errors.New("token has expired")
		} else {
			return "", errors.New("token decode error")
		}
	} else {
		return "", errors.New("token decode error")
	}
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authKey string
		switch {
		case c.Request.Header.Get("Authorization") != "":
			authKey = c.Request.Header.Get("Authorization")
		case c.Request.Header.Get("authorization") != "":
			authKey = c.Request.Header.Get("authorization")
		case c.Query("wskey") != "":
			authKey = c.Query("wskey")
		}

		token := jwToken(strings.Replace(authKey, "Bearer ", "", 1))
		if id, err := token.decode(); err != nil {
			handler.LogErr(c, "Auth error: "+err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})
			return
		} else {
			c.Set("user", id)
		}

		c.Next()
	}
}
