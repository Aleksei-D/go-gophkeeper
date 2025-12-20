package app

import (
	"context"
	"go-gophkeeper/internal/models"
	vaultPB "go-gophkeeper/internal/pb/vault"
	"go-gophkeeper/internal/server/service"
)

// VaultServer сервер хранилища
type VaultServer struct {
	vaultPB.UnimplementedVaultServiceServer
	Service *service.Service
}

// SyncSecrets ручка синхранизация данных
func (v *VaultServer) SyncSecrets(ctx context.Context, req *vaultPB.VaultListMessage) (*vaultPB.VaultListMessage, error) {
	vaultObjects, err := models.NewVaultObjectListFromProto(req)
	if err != nil {
		return nil, err
	}

	dataToClient, err := v.Service.SyncVault(ctx, vaultObjects)
	if err != nil {
		return nil, err
	}

	return dataToClient.ToProto(), nil
}
