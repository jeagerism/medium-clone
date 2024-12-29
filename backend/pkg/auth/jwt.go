package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jeagerism/medium-clone/backend/config"
)

type JWTService struct {
	config *config.JWT
}

type RefreshTokenClaims struct {
}

func NewJWTService(cfg *config.JWT) *JWTService {
	return &JWTService{config: cfg}
}

func (j *JWTService) GenerateRefreshToken(userID int, secretKey string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     duration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// func (j *JWTService) ValidateToken(token string, secretKey string) (*jwt.Token, error) {
// 	tokenWithoutBearer := token
// }
