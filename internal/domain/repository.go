package domain

import (
	"context"
	"time"

	"go-gophkeeper/internal/models"
)

// AuthRepository интерфейс работы юзерами
type AuthRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	IsExist(ctx context.Context, login string) (bool, error)
}

// VaultRepository интерфейс работы с хранилищем
type VaultRepository interface {
	Add(ctx context.Context, vaultObject models.VaultObject) error
	Get(ctx context.Context, name, login, dataType string) (models.VaultObject, error)
	IsExist(ctx context.Context, name, login string) (bool, error)
	GetDataToSync(ctx context.Context, login string) (models.VaultObjects, error)
	AddList(ctx context.Context, vaultObjects models.VaultObjects) error
}

// EventRepository интерфейс работы с событиями
type EventRepository interface {
	Add(ctx context.Context, login string) error
	Get(ctx context.Context, login string) (time.Time, error)
}
