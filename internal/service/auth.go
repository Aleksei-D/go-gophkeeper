package service

import (
	"context"

	"go-gophkeeper/internal/domain"
	"go-gophkeeper/internal/logger"
	"go-gophkeeper/internal/models"
	crypto2 "go-gophkeeper/internal/utils/crypto"
	errors2 "go-gophkeeper/internal/utils/errors"
	"go-gophkeeper/internal/utils/token"

	"go.uber.org/zap"
)

type AuthService struct {
	userRepo   domain.AuthRepository
	managerJWT *token.JWTManager
}

func NewAuthService(authRepository domain.AuthRepository, managerJWT *token.JWTManager) *AuthService {
	return &AuthService{userRepo: authRepository, managerJWT: managerJWT}
}

func (a *AuthService) Login(ctx context.Context, user *models.User) (string, error) {
	var token string
	existUser, err := a.userRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return token, err
	}
	if existUser == nil || !crypto2.CheckPasswordHash(user.Password, existUser.Password) {
		return token, errors2.ErrInvalidCredentials
	}

	return a.managerJWT.CreateToken(user.Login)
}

func (a *AuthService) CreateUser(ctx context.Context, newUser *models.User) (string, error) {
	var token string

	ok, err := a.userRepo.IsExist(ctx, newUser.Login)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return token, err
	}

	if ok {
		return token, errors2.ErrUserAlreadyExists
	}

	hashPassword, err := crypto2.HashPassword(newUser.Password)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return token, err
	}

	newUser.Password = hashPassword
	user, err := a.userRepo.Create(ctx, newUser)
	if err != nil {
		logger.Log.Warn("User Create Error", zap.Error(err))
		return token, err
	}

	token, err = a.managerJWT.CreateToken(user.Login)
	if err != nil {
		return token, err
	}

	return token, nil
}
