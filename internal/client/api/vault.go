package api

import (
	"context"
	"go-gophkeeper/internal/client/interceptors"
	"go-gophkeeper/internal/models"
	pb "go-gophkeeper/internal/pb/vault"

	"google.golang.org/grpc"
)

// VaultClient клиент для работы с хранилищем
type VaultClient struct {
	vaultServiceClient pb.VaultServiceClient
}

// NewVaultClient возвращает клиент для работы с хранилищем
func NewVaultClient(conn *grpc.ClientConn) *VaultClient {
	return &VaultClient{
		vaultServiceClient: pb.NewVaultServiceClient(conn),
	}
}

// SyncVault обращается к метооду SyncVault на сервере
func (c *VaultClient) SyncVault(ctx context.Context, token string, vaultObjects models.VaultObjects) (models.VaultObjects, error) {
	callOption := []grpc.CallOption{interceptors.TokenOption{Token: token}}
	resp, err := c.vaultServiceClient.SyncVault(ctx, vaultObjects.ToProto(), callOption...)
	if err != nil {
		return nil, err
	}
	newVaultObjects, err := models.NewVaultObjectListFromProto(resp)
	if err != nil {
		return nil, err
	}

	return newVaultObjects, nil
}
