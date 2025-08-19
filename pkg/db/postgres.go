package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func InitDB(cfg *DBConfig, log *zap.Logger) (*sqlx.DB, error) {

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database", zap.Error(err))
	}

	migrationsDir := "./migrations/pgsql"

	if err = goose.Up(db.DB, migrationsDir); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	log.Info("Successfully connected to database", zap.String("dir", migrationsDir))
	return db, nil
}
