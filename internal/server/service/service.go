package service

import (
	"context"
	"errors"
	"go-gophkeeper/internal/domain"
	"go-gophkeeper/internal/logger"
	"go-gophkeeper/internal/models"
	crypto2 "go-gophkeeper/internal/utils/crypto"
	errors2 "go-gophkeeper/internal/utils/errors"
	"go-gophkeeper/internal/utils/token"
	"go.uber.org/zap"
)

// Service сервисный слой
type Service struct {
	userRepo   domain.AuthRepository
	vaultRepo  domain.VaultRepository
	eventRepo  domain.EventRepository
	managerJWT *token.JWTManager
	hasher     crypto2.Hasher
}

// NewService возврат нового Service
func NewService(userRepo domain.AuthRepository, vaultRepo domain.VaultRepository, eventRepo domain.EventRepository, managerJWT *token.JWTManager, hasher crypto2.Hasher) *Service {
	return &Service{
		userRepo:   userRepo,
		vaultRepo:  vaultRepo,
		eventRepo:  eventRepo,
		managerJWT: managerJWT,
		hasher:     hasher,
	}
}

// Login слой логирования
func (s *Service) Login(ctx context.Context, user *models.User) (string, error) {
	var tokenJWT string
	existUser, err := s.userRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return tokenJWT, err
	}

	if existUser == nil || !s.hasher.CheckPasswordHash(user.Password, existUser.Password) {
		return tokenJWT, errors2.ErrInvalidCredentials
	}

	return s.managerJWT.CreateToken(user.Login)
}

// CreateUser слой регистрации
func (s *Service) CreateUser(ctx context.Context, newUser *models.User) (string, error) {
	var tokenJWT string

	ok, err := s.userRepo.IsExist(ctx, newUser.Login)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return tokenJWT, err
	}

	if ok {
		return tokenJWT, errors2.ErrUserAlreadyExists
	}

	hashPassword, err := s.hasher.HashPassword(newUser.Password)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return tokenJWT, err
	}

	newUser.Password = hashPassword
	user, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		logger.Log.Warn("User Create Error", zap.Error(err))
		return tokenJWT, err
	}

	tokenJWT, err = s.managerJWT.CreateToken(user.Login)
	if err != nil {
		return tokenJWT, err
	}

	return tokenJWT, nil
}

// SyncVault синхронизирует данные
func (s *Service) SyncVault(ctx context.Context, dataFromClient models.VaultObjects) (models.VaultObjects, error) {
	var newData models.VaultObjects
	claims := ctx.Value("userClaims")
	userClaims, ok := claims.(token.UserClaims)
	if !ok {
		return newData, errors.New("invalid user claims")
	}

	newSecrets, err := s.vaultRepo.GetDataToSync(ctx, userClaims.Login)
	if err != nil {
		return newSecrets, err
	}

	err = s.vaultRepo.AddList(ctx, dataFromClient)
	if err != nil {
		return newSecrets, err
	}

	err = s.eventRepo.Add(ctx, userClaims.Login)
	if err != nil {
		return newSecrets, err
	}

	return newSecrets, nil
}
