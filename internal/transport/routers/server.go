package routers

import (
	"Crash-Auth-service/internal/middleware"
	"Crash-Auth-service/internal/repository"
	"Crash-Auth-service/internal/transport/handlers"
	"Crash-Auth-service/pkg/jwt"
	"Crash-Auth-service/pkg/metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunServer(
	handler *handlers.AuthHandler,
	log *zap.Logger,
	metric *metrics.MetricsHelper,
	repo repository.AuthRepository,
	jwtToken *jwt.JWTConfig,
) {
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		middleware.ZapLogger(log),
		middleware.Metrics(metric),
	)

	api := r.Group("/api", middleware.AuthMiddleware(jwtToken, repo, log))
	{
		api.POST("register/user", handler.RegisterUser)
		api.POST("login/user", handler.LoginUser)
		api.PUT("change/password", handler.ChangePassword)
		api.PUT("change/email", handler.ChangeEmail)
	}

	r.GET(metric.Path(), metric.Handler())

	log.Info("Starting http server on port 9000")
	if err := r.Run(":9000"); err != nil {
		log.Fatal("Failed to run server", zap.Error(err))
	}
}
