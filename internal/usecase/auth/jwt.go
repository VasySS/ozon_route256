package auth

import (
	"errors"
	"time"
	"workshop-1/config"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(userID int) (string, error) {
	expirationTime := time.Now().UTC().Add(config.AccessTokenTTL)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": expirationTime.Unix(),
	})

	token, err := claims.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateAccessToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	tokenParsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}

		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !tokenParsed.Valid {
		return nil, errors.New("неверный токен")
	}

	return claims, nil
}
