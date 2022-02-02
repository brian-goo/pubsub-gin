package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var hmacKey = []byte("test")

func checkJWT(jwtS string) error {
	token, err := jwt.Parse(jwtS, func(token *jwt.Token) (interface{}, error) {
		return hmacKey, nil
	})

	if token.Valid {
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errors.New("invalid token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return errors.New("token has expired")
		} else {
			return err
		}
	} else {
		return err
	}
}
