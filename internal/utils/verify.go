package utils

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenStr, secret, tokenType string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	jwtType, ok := claims["type"].(string)
	if !ok || jwtType != tokenType {
		return 0, errors.New("invalid token type")
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("invalid token sub")
	}

	switch v := sub.(type) {
	case float64:
		return uint(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("invalid sub format")
		}
		return uint(id), nil
	default:
		return 0, errors.New("unsupported sub type")
	}
}
