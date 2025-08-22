package pgsql

import (
	"Crash-Auth-service/internal/entities"
	"Crash-Auth-service/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) repository.AuthRepository {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) SaveUserName(ctx context.Context, tx *sqlx.Tx, fullName string) (*entities.User, error) {
	var userName entities.User
	err := tx.GetContext(ctx, &userName, "INSERT INTO users (full_name) VALUES ($1) RETURNING id",
		fullName)
	if err != nil {
		return nil, err
	}
	return &userName, nil
}

func (r *AuthRepo) SaveUserEmail(ctx context.Context, tx *sqlx.Tx, userId, email string) (*entities.UserEmail, error) {
	var userEmail entities.UserEmail
	err := tx.GetContext(ctx, &userEmail, "INSERT INTO emails (user_id, email) VALUES ($1, $2) RETURNING user_id",
		userId, email)
	if err != nil {
		return nil, err
	}
	return &userEmail, nil
}

func (r *AuthRepo) SaveUserPassword(ctx context.Context, tx *sqlx.Tx, userId, password string) (*entities.UserPass, error) {
	var userPass entities.UserPass
	err := tx.GetContext(ctx, &userPass, "INSERT INTO passwords(user_id, hash) VALUES ($1, $2) RETURNING user_id",
		userId, password)
	if err != nil {
		return nil, err
	}
	return &userPass, nil
}

func (r *AuthRepo) FindUserByEmail(ctx context.Context, email string) (string, string, error) {
	var userId, hash string
	err := r.db.QueryRowxContext(ctx, `
		SELECT u.id, p.hash FROM users u
		JOIN emails e ON u.id = e.user_id
		JOIN passwords p ON u.id = p.user_id
		WHERE e.email = $1`,
		email).Scan(&userId, &hash)
	if err != nil {
		return "", "", err
	}
	return userId, hash, nil
}

func (r *AuthRepo) FindPasswordByUserID(ctx context.Context, userId string) (string, error) {
	var passwordHash string

	err := r.db.GetContext(ctx, &passwordHash, "SELECT hash FROM passwords WHERE user_id = $1 LIMIT 1",
		userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("password not found for userID=%s", userId)
		}
		return "", err
	}

	return passwordHash, nil
}

func (r *AuthRepo) UpdatePassword(ctx context.Context, userId, newPassword string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE passwords SET hash = $1 WHERE user_id = $2",
		newPassword, userId)
	if err != nil {
		return err
	}
	return err
}

func (r *AuthRepo) UpdateEmail(ctx context.Context, userId, newEmail string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE emails SET email = $1 WHERE user_id = $2",
		newEmail, userId)
	if err != nil {
		return err
	}
	return err
}

func (r *AuthRepo) UpdateFullName(ctx context.Context, userId, fullName string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET full_name = $1 WHERE id = $2",
		fullName, userId)
	if err != nil {
		return err
	}
	return err
}

func (r *AuthRepo) DeleteUserById(ctx context.Context, userId string) error {
	row, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
