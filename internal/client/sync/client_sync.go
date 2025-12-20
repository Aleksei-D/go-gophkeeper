package sync

import (
	"context"
	"go-gophkeeper/internal/client/api"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/domain"
	"go-gophkeeper/internal/logger"
	"go-gophkeeper/internal/models"
	"go.uber.org/zap"
	"time"
)

// ClientSync клиент синхронизации с хранилищем
type ClientSync struct {
	cgf         *config.Config
	user        *models.User
	vaultClient *api.VaultClient
	vaultRepo   domain.VaultRepository
	eventRepo   domain.EventRepository
}

// NewClientSync возвращает ClientSync
func NewClientSync(cfg *config.Config, vaultClient *api.VaultClient, vaultRepo domain.VaultRepository, eventRepo domain.EventRepository) *ClientSync {
	return &ClientSync{
		cgf:         cfg,
		vaultClient: vaultClient,
		vaultRepo:   vaultRepo,
		eventRepo:   eventRepo,
	}
}

// Run запускает ClientSync
func (c *ClientSync) Run(ctx context.Context) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	syncTicker := time.NewTicker(c.cgf.SyncInterval.Duration)
	defer syncTicker.Stop()

	errorCh := make(chan error)
	defer close(errorCh)

	localCardCH := c.localDataGenerator(ctx, doneCh, errorCh, syncTicker)
	vaultCardCH := c.vaultDataGenerator(ctx, doneCh, errorCh, localCardCH)
	go c.updateLocalData(ctx, doneCh, errorCh, vaultCardCH)

	for err := range errorCh {
		logger.Log.Warn(err.Error(), zap.Error(err))
	}
}

func (c *ClientSync) localDataGenerator(ctx context.Context, doneCh chan struct{}, errorCh chan<- error, syncTicker *time.Ticker) <-chan models.VaultObjects {
	dataCH := make(chan models.VaultObjects)
	go func() {
		defer close(dataCH)
	newLoop:
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Worker: Context done, exiting.")
				return
			case <-doneCh:
				return
			case <-syncTicker.C:
				vaultObjects, err := c.vaultRepo.GetDataToSync(ctx, c.user.Login)
				if err != nil {
					errorCh <- err
					continue newLoop
				}
				dataCH <- vaultObjects
			}
		}
	}()
	return dataCH
}

func (c *ClientSync) vaultDataGenerator(ctx context.Context, doneCh chan struct{}, errorCh chan<- error, localCardCH <-chan models.VaultObjects) <-chan models.VaultObjects {
	dataCH := make(chan models.VaultObjects)
	go func() {
		defer close(dataCH)
	newLoop:
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Worker: Context done, exiting.")
				return
			case <-doneCh:
				return
			case vaultObjects := <-localCardCH:
				newData, err := c.vaultClient.SyncVault(ctx, c.user.Token, vaultObjects)
				if err != nil {
					errorCh <- err
					continue newLoop
				}

				dataCH <- newData
			}
		}
	}()
	return dataCH
}

func (c *ClientSync) updateLocalData(ctx context.Context, doneCh chan struct{}, errorCh chan<- error, vaultObjectCH <-chan models.VaultObjects) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Worker: Context done, exiting.")
			return
		case <-doneCh:
			return
		case vaultObjects := <-vaultObjectCH:
			for _, vaultObject := range vaultObjects {
				err := c.vaultRepo.Add(ctx, vaultObject)
				if err != nil {
					errorCh <- err
				}
			}

			err := c.eventRepo.Add(ctx, c.user.Login)
			if err != nil {
				errorCh <- err
			}
		}
	}
}

// SetUser установка данных юзера полс авторизации
func (c *ClientSync) SetUser(user *models.User) {
	c.user = user
}
