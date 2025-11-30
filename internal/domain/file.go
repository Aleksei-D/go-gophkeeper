package domain

import (
	"context"

	"go-gophkeeper/internal/models"
)

type FileRepository interface {
	Get(ctx context.Context, user *models.File) (*models.File, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	IsExist(ctx context.Context, login string) (bool, error)
}
