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

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		h.log.Error("userID missing in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.ChangePassword(ctx, userID, req.CurrentPassword, req.NewPassword); err != nil {
		h.log.Error("failed to change password", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	h.log.Info("change password successfully",
		zap.String("UserID", userID),
	)
	c.JSON(http.StatusOK, gin.H{"message": "change password successfully"})
}

func (h *AuthHandler) ChangeEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("request body error",
			zap.String("newEmail", req.NewEmail),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body error"})
		return
	}

	err := h.service.ChangeEmail(ctx, userID.(string), req.Password, req.NewEmail)
	if err != nil {
		h.log.Error("change email error",
			zap.String("userID", userID.(string)),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("change email successfully",
		zap.String("userID", userID.(string)),
		zap.String("NewEmail", req.NewEmail),
	)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Change email successfully",
	})
}

func (h *AuthHandler) ChangeFullName(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID := c.GetString("userID")
	if userID == "" {
		h.log.Error("userID missing in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.ChangeFullNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("request body error",
			zap.String("FullName", req.NewFullName),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body error"})
		return
	}

	err := h.service.ChangeFullName(ctx, userID, req.NewFullName)
	if err != nil {
		h.log.Error("change full name error",
			zap.String("FullName", req.NewFullName),
			zap.String("userId", userID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "change full name error"})
		return
	}

	h.log.Info("change full name successfully",
		zap.String("FullName", req.NewFullName),
		zap.String("Email", userID),
	)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Change full name successfully",
	})
}
