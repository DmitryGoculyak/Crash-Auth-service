package service

import (
	"Crash-Auth-service/internal/clients/billing"
	"Crash-Auth-service/internal/clients/currency"
	"Crash-Auth-service/internal/dto"
	"Crash-Auth-service/internal/entities"
	"Crash-Auth-service/internal/repository"
	"Crash-Auth-service/pkg/jwt"
	"Crash-Auth-service/pkg/transaction"
	"Crash-Auth-service/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type AuthServiceServer interface {
	ProcessRegistration(ctx context.Context, fullName, email, password, currencyCode string) (*entities.User, error)
	ProcessAuthorization(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
}

type AuthService struct {
	repo           repository.AuthRepository
	txManager      transaction.TransactionManager
	validator      *validator.Validate
	log            *zap.Logger
	jwtToken       *jwt.JWTConfig
	billingClient  *billing.BillingClient
	currencyClient *currency.CurrencyClient
}

func NewAuthService(
	repo repository.AuthRepository,
	txManager transaction.TransactionManager,
	validator *validator.Validate,
	log *zap.Logger,
	jwtToken *jwt.JWTConfig,
	billingClient *billing.BillingClient,
	currencyClient *currency.CurrencyClient,
) AuthServiceServer {
	return &AuthService{
		repo:           repo,
		txManager:      txManager,
		validator:      validator,
		log:            log,
		jwtToken:       jwtToken,
		billingClient:  billingClient,
		currencyClient: currencyClient,
	}
}

func (s *AuthService) ProcessRegistration(ctx context.Context, fullName, email, password, currencyCode string) (*entities.User, error) {
	input := dto.RegistrationInput{
		FullName:     fullName,
		Email:        email,
		Password:     password,
		CurrencyCode: currencyCode,
	}

	if err := s.validator.Struct(input); err != nil {
		s.log.Warn("Validation input error",
			zap.String("fullName", fullName),
			zap.String("email", email),
			zap.String("password", password),
			zap.String("currencyCode", currencyCode),
			zap.Error(err),
		)
		return nil, fmt.Errorf("validate input error: %w", err)
	}

	hash, err := utils.CreateHash(input.Password)
	if err != nil {
		s.log.Warn("Create hash error",
			zap.String("password", input.Password),
			zap.Error(err),
		)
		return nil, err
	}

	var createUser *entities.User
	err = s.txManager.WithTransaction(ctx, func(tx *sqlx.Tx) error {

		newUser, nameErr := s.repo.SaveUserName(ctx, tx, input.FullName)
		if nameErr != nil {
			s.log.Warn("Failed to save user name",
				zap.String("fullName", input.FullName),
				zap.Error(nameErr),
			)
			return nameErr
		}
		createUser = newUser

		_, emailErr := s.repo.SaveUserEmail(ctx, tx, createUser.Id, input.Email)
		if emailErr != nil {
			s.log.Warn("Failed to save user email",
				zap.String("userID", createUser.Id),
				zap.String("email", input.Email),
				zap.Error(emailErr),
			)
			return emailErr
		}

		_, passErr := s.repo.SaveUserPassword(ctx, tx, createUser.Id, hash)
		if passErr != nil {
			s.log.Warn("Failed to save user password",
				zap.String("userID", createUser.Id),
				zap.String("password", hash),
				zap.Error(passErr),
			)
			return passErr
		}

		if err = s.currencyClient.CheckCurrencyExists(ctx, input.CurrencyCode); err != nil {
			s.log.Error("Currency check error",
				zap.String("currencyCode", input.CurrencyCode),
				zap.Error(err),
			)
			return err
		}

		if err = s.billingClient.CreateWallet(ctx, createUser.Id, input.CurrencyCode); err != nil {
			s.log.Error("Failed to create wallet",
				zap.String("userID", createUser.Id),
				zap.String("currencyCode", input.CurrencyCode),
				zap.Error(err),
			)
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("Transaction error", zap.Error(err))
		return nil, err
	}
	return createUser, nil
}

func (s *AuthService) ProcessAuthorization(ctx context.Context, email, password string) (string, error) {
	input := dto.AuthorizationInput{
		Email:    email,
		Password: password,
	}

	if err := s.validator.Struct(input); err != nil {
		s.log.Warn("Validation input error",
			zap.String("email", email),
			zap.String("password", password),
			zap.Error(err),
		)
		return "", fmt.Errorf("validate input error: %w", err)
	}

	userId, hash, err := s.repo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		s.log.Warn("find user by email error",
			zap.String("email", input.Email),
			zap.Error(err),
		)
		return "", err
	}

	if !utils.CheckHash(hash, input.Password) {
		s.log.Warn("Invalid password")
		return "", errors.New("invalid password")
	}

	token, err := s.jwtToken.GenerateToken(userId)
	if err != nil {
		s.log.Warn("Generate token error",
			zap.String("userID", userId),
			zap.Error(err),
		)
		return "", err
	}
	return token, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	input := dto.PasswordInput{
		Email:       email,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	if err := s.validator.Struct(input); err != nil {
		s.log.Warn("Validation input error",
			zap.String("Email", email),
			zap.Error(err),
		)
		return fmt.Errorf("validate input error: %w", err)
	}

	userId, currentHash, err := s.repo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		s.log.Error("find user by email error",
			zap.String("email", input.Email),
			zap.Error(err),
		)
		return err
	}

	if !utils.CheckHash(currentHash, input.OldPassword) {
		s.log.Warn("Invalid password")
		return errors.New("invalid password")
	}

	newHash, err := utils.CreateHash(input.NewPassword)
	if err != nil {
		s.log.Error("Create hash error",
			zap.Error(err),
		)
		return err
	}

	err = s.repo.UpdatePassword(ctx, userId, newHash)
	if err != nil {
		s.log.Warn("Update password error",
			zap.String("userID", userId),
			zap.Error(err),
		)
		return err
	}
	return nil
}
