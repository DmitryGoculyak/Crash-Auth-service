package config

import (
	"Crash-Auth-service/internal/clients/billing"
	"Crash-Auth-service/internal/clients/currency"
	"Crash-Auth-service/internal/transport/routers"
	"Crash-Auth-service/pkg/db"
	"Crash-Auth-service/pkg/jwt"
	"Crash-Auth-service/pkg/logger"
	"Crash-Auth-service/pkg/metrics"

	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		LoadConfig,
		func(cfg *Config) *db.DBConfig { return cfg.DBConfig },
		func(cfg *Config) *jwt.JWTConfig { return cfg.JwtConfig },
		func(cfg *Config) *logger.Config { return cfg.LoggerConfig },
		func(cfg *Config) *metrics.Config { return cfg.MetricsConfig },
		func(cfg *Config) *billing.BillingConfig { return cfg.BillingClientConfig },
		func(cfg *Config) *currency.CurrencyConfig { return cfg.CurrencyClientConfig },
		func(cfg *Config) *routers.ServerConfig { return cfg.ServerConfig },
	),
)
