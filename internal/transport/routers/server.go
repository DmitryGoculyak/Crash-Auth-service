package routers

import (
	"Crash-Auth-service/internal/middleware"
	"Crash-Auth-service/internal/transport/handlers"
	"Crash-Auth-service/pkg/jwt"
	"Crash-Auth-service/pkg/metrics"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunServer(
	cfg *ServerConfig,
	handler *handlers.AuthHandler,
	log *zap.Logger,
	metric *metrics.MetricsHelper,
	jwtToken *jwt.JWTConfig,
) {
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		middleware.ZapLogger(log),
		middleware.Metrics(metric),
	)

	api := r.Group("/api")
	{
		api.POST("registration", handler.RegisterUser)
		api.POST("login", handler.LoginUser)

		auth := api.Group("/auth", middleware.AuthMiddleware(jwtToken, log))
		{
			auth.PUT("change/password", handler.ChangePassword)
			auth.PUT("change/email", handler.ChangeEmail)
			auth.PUT("change/fullname", handler.ChangeFullName)
			auth.DELETE("delete/account", handler.DeleteUserAccount)
		}
	}
	r.GET(metric.Path(), metric.Handler())

	log.Info("Starting server on port" + cfg.Port)
	server := &http.Server{
		Addr:         cfg.Host + cfg.Port,
		Handler:      r,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Error starting server", zap.Error(err))
	}
}
