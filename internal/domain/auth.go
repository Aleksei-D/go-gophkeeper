package domain

import (
	"context"

	"go-gophkeeper/internal/models"
)

type AuthRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	IsExist(ctx context.Context, login string) (bool, error)
}
