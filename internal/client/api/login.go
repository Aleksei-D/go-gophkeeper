package api

import (
	"context"
	"go-gophkeeper/internal/models"
	"google.golang.org/grpc"

	pb "go-gophkeeper/internal/pb/auth"
)

// AuthClient клиент для авторизации
type AuthClient struct {
	authServiceClient pb.AuthServiceClient
}

// NewClientAuth возвращает клиент для авторизации
func NewClientAuth(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		authServiceClient: pb.NewAuthServiceClient(conn),
	}
}

// UserRegister обращется к методу UserRegister на сервере
func (c *AuthClient) UserRegister(ctx context.Context, user *models.User) (*models.User, error) {
	req := &pb.UserMessageRequest{Login: user.Login, Password: user.Password}
	resp, err := c.authServiceClient.Register(ctx, req)
	if err != nil {
		return user, err
	}
	user.Token = resp.Token
	return user, nil
}

// UserLogin обращется к методу Login на сервере
func (c *AuthClient) UserLogin(ctx context.Context, user *models.User) (*models.User, error) {
	req := &pb.UserMessageRequest{Login: user.Login, Password: user.Password}
	resp, err := c.authServiceClient.Login(ctx, req)
	if err != nil {
		return user, err
	}
	user.Token = resp.Token
	return user, nil
}
