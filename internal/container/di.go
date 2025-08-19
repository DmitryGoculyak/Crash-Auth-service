package container

import (
	"Crash-Auth-service/internal/clients/billing"
	"Crash-Auth-service/internal/clients/currency"
	"Crash-Auth-service/internal/config"
	repo "Crash-Auth-service/internal/repository/pgsql"
	"Crash-Auth-service/internal/service"
	"Crash-Auth-service/internal/transport/handlers"
	"Crash-Auth-service/internal/transport/routers"
	"Crash-Auth-service/pkg/db"
	"Crash-Auth-service/pkg/logger"
	"Crash-Auth-service/pkg/metrics"
	"Crash-Auth-service/pkg/transaction"
	"Crash-Auth-service/pkg/validation"
	"go.uber.org/fx"
)

func Build() *fx.App {
	return fx.New(
		config.Module,
		db.Module,
		validation.Module,
		transaction.Module,
		logger.Module,
		metrics.Module,
		billing.Module,
		currency.Module,
		repo.Module,
		handlers.Module,
		service.Module,
		routers.Module,
	)
}
