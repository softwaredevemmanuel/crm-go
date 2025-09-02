package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(cfg.JWTSecret) // 🔥 change this to env variable in prod

// Generate JWT token
func GenerateJWT(userID string, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email": email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * time.Duration(cfg.JWTExpire)).Unix(), // token expires in 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
