package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateKey(email string) (string, error) {
	key, err := GetEnvNoCon("ACCESS_KEY")
	if err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenSting, err := claims.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return tokenSting, nil
}
