package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTConfig struct {
	SigningKey    string        `mapstructure:"secret_key"`
	TokenLifetime time.Duration `mapstructure:"token_lifetime"`
}

func (cfg *JWTConfig) GenerateToken(userID string) (string, error) {
	if cfg.SigningKey == "" {
		return "", fmt.Errorf("JWT_SIGNING_KEY is empty")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(cfg.TokenLifetime).Unix(),
		"iat":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.SigningKey))
	if err != nil {
		return "", fmt.Errorf("failed to signedToken: %w", err)
	}

	return signedToken, nil
}
