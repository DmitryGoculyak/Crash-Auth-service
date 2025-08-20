package middleware

import (
	"Crash-Auth-service/internal/repository"
	"Crash-Auth-service/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtCfg *jwt.JWTConfig, repo repository.AuthRepository, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwtCfg.ParseToken(tokenString)
		if err != nil {
			log.Warn("invalid JWT token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		currentEmail, _, err := repo.FindUserByEmail(c.Request.Context(), userID)
		if err != nil {
			log.Warn("failed to get user email from DB", zap.String("userID", userID), zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("currentEmail", currentEmail)

		c.Next()
	}
}
