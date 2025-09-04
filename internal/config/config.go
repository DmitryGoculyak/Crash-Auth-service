package config

import (
	"Crash-Auth-service/internal/clients/billing"
	"Crash-Auth-service/internal/clients/currency"
	"Crash-Auth-service/internal/transport/routers"
	"Crash-Auth-service/pkg/db"
	"Crash-Auth-service/pkg/jwt"
	"Crash-Auth-service/pkg/logger"
	"Crash-Auth-service/pkg/metrics"

	"github.com/spf13/viper"

	"fmt"
	"sync"
)

var (
	err    error
	config *Config
	s      sync.Once
)

type Config struct {
	DBConfig             *db.DBConfig
	JwtConfig            *jwt.JWTConfig
	LoggerConfig         *logger.Config
	MetricsConfig        *metrics.Config
	BillingClientConfig  *billing.BillingConfig
	CurrencyClientConfig *currency.CurrencyConfig
	ServerConfig         *routers.ServerConfig
}

func LoadConfig() (*Config, error) {

	s.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		DBConfig := viper.Sub("database")
		JwtConfig := viper.Sub("jwt")
		LoggerConfig := viper.Sub("logger")
		MetricsConfig := viper.Sub("metrics")
		BillingClientConfig := viper.Sub("billing_client")
		CurrencyClientConfig := viper.Sub("currency_client")
		ServerConfig := viper.Sub("server")

		if err = parseSubConfig(DBConfig, &config.DBConfig); err != nil {
			return
		}
		if err = parseSubConfig(JwtConfig, &config.JwtConfig); err != nil {
			return
		}
		if err = parseSubConfig(LoggerConfig, &config.LoggerConfig); err != nil {
			return
		}
		if err = parseSubConfig(MetricsConfig, &config.MetricsConfig); err != nil {
			return
		}
		if err = parseSubConfig(BillingClientConfig, &config.BillingClientConfig); err != nil {
			return
		}
		if err = parseSubConfig(CurrencyClientConfig, &config.CurrencyClientConfig); err != nil {
			return
		}
		if err = parseSubConfig(ServerConfig, &config.ServerConfig); err != nil {
			return
		}
	})
	return config, err
}

func parseSubConfig[T any](subConfig *viper.Viper, parseTo *T) error {
	if subConfig == nil {
		return fmt.Errorf("can not read %T config: subconfig is nil", parseTo)
	}

	if err = subConfig.Unmarshal(parseTo); err != nil {
		return err
	}
	return nil
}
