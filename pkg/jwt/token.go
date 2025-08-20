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

func (cfg *JWTConfig) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.SigningKey), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}
