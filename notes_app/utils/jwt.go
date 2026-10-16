package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultSecret = "superSecretKey"

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET is not set. Using insecure default secret.")
		secret = defaultSecret
	}

	return []byte(secret)
}

func GenerateToken(email string, userID int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userID,
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
	})

	return token.SignedString(getJWTSecret())
}

func VerifyToken(token string) (string, int32, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return getJWTSecret(), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", 0, errors.New("Invalid Token!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, errors.New("There was an error getting the claims from the token. Try Again!")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", 0, errors.New("Missing email claim in token")
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		return "", 0, errors.New("Missing userId claim in token")
	}

	return email, int32(userID), nil
}
