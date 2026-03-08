package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func InitJWT(key string) {
	jwtKey = []byte(key)
}

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Name     string `json:"name"`
	Platform string `json:"X-Platform"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64, name, platform string) (string, error) {
	exp := time.Now().Add(2 * time.Hour)
	if platform != "web" && platform != "mobile" {
		return "", errors.New("invalid platform for token.")
	}

	claims := &CustomClaims{
		UserID:   userID,
		Name:     name,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   fmt.Sprint(userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenStr string) (int64, string, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return jwtKey, nil
	})
	if err != nil {
		return 0, "", "", err
	}

	if token.Valid == false {
		return 0, "", "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if (claims.Platform != "web" && claims.Platform != "mobile") || claims.UserID == 0 || claims.Name == "" {
			return 0, "", "", errors.New("invalid claims data")
		}
		return claims.UserID, claims.Name, claims.Platform, nil
	} else {
		return 0, "", "", errors.New("invalid claims")
	}
}
