package repository

import (
	"Crash-Auth-service/internal/entities"
	"context"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	SaveUserName(ctx context.Context, tx *sqlx.Tx, fullName string) (*entities.User, error)
	SaveUserEmail(ctx context.Context, tx *sqlx.Tx, userId, email string) (*entities.UserEmail, error)
	SaveUserPassword(ctx context.Context, tx *sqlx.Tx, userId, password string) (*entities.UserPass, error)
}
