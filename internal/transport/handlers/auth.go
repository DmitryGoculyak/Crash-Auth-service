package handlers

import (
	"Crash-Auth-service/internal/dto"
	"Crash-Auth-service/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	service service.AuthServiceServer
	log     *zap.Logger
}

func NewAuthHandler(
	service service.AuthServiceServer,
	log *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		service: service,
		log:     log,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("request body error",
			zap.Any("request body", req),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body error"})
		return
	}

	user, err := h.service.ProcessRegistration(ctx, req.FullName, req.Email, req.Password, req.CurrencyCode)
	if err != nil {
		h.log.Error("process registration error",
			zap.Any("request", req),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "process registration error"})
		return
	}

	h.log.Info("create user successfully",
		zap.String("UserID", user.Id),
		zap.String("FullName", req.FullName),
		zap.String("Email", req.Email),
		zap.String("Password", req.Password),
		zap.String("CurrencyCode", req.CurrencyCode),
		zap.String("CreatedAt", user.CreatedAt.String()),
	)
	c.JSON(http.StatusCreated, gin.H{
		"Id": user.Id,
	})
}
