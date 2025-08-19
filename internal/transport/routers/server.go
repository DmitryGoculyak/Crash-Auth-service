package routers

import (
	"Crash-Auth-service/internal/middleware"
	"Crash-Auth-service/internal/transport/handlers"
	"Crash-Auth-service/pkg/metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunServer(handler *handlers.AuthHandler, log *zap.Logger, metric *metrics.MetricsHelper) {
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		middleware.ZapLogger(log),
		middleware.Metrics(metric),
	)

	api := r.Group("/api")
	{
		api.POST("register/user", handler.RegisterUser)
		api.POST("login/user", handler.LoginUser)
	}

	r.GET(metric.Path(), metric.Handler())

	log.Info("Starting http server on port 9000")
	if err := r.Run(":9000"); err != nil {
		log.Fatal("Failed to run server", zap.Error(err))
	}
}
