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

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("request body error",
			zap.String("FullName", req.FullName),
			zap.String("Email", req.Email),
			zap.String("CurrencyCode", req.CurrencyCode),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body error"})
		return
	}

	user, err := h.service.ProcessRegistration(ctx, req.FullName, req.Email, req.Password, req.CurrencyCode)
	if err != nil {
		h.log.Error("process registration error",
			zap.String("FullName", req.FullName),
			zap.String("Email", req.Email),
			zap.String("CurrencyCode", req.CurrencyCode),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "process registration error"})
		return
	}

	h.log.Info("create user successfully",
		zap.String("UserID", user.Id),
		zap.String("FullName", req.FullName),
		zap.String("Email", req.Email),
		zap.String("CurrencyCode", req.CurrencyCode),
		zap.String("CreatedAt", user.CreatedAt.String()),
	)
	c.JSON(http.StatusCreated, gin.H{
		"Id": user.Id,
	})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("request body error",
			zap.Any("request body", req),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body error"})
		return
	}

	token, err := h.service.ProcessAuthorization(ctx, req.Email, req.Password)
	if err != nil {
		h.log.Error("process authorization error",
			zap.String("Email", req.Email),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "process authorization error"})
		return
	}

	h.log.Info("authorization user successfully",
		zap.String("Email", req.Email),
		zap.String("JWTToken", token),
	)
	c.JSON(http.StatusAccepted, gin.H{
		"JWT": token,
	})
}
