package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xelathan/golang_backend/config"
)

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Now().Add(time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": expiration,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
