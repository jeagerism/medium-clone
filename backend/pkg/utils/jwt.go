package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken รับค่า secret และ duration แบบ dynamic
func GenerateAccessToken(userID int, email, role string, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken รับค่า secret และ duration แบบ dynamic
func GenerateRefreshToken(userID int, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken ตรวจสอบ Refresh Token
func ValidateToken(tokenStr, secret string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("invalid token: %w", err)
	}

	// ตรวจสอบว่า Claims เป็นแบบ MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token claims")
	}

	return token, claims, nil
}
