package client

import (
	"context"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/models"

	pb "go-gophkeeper/internal/pb/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientGRPC struct {
	conn       *grpc.ClientConn
	authCLient pb.AuthServiceClient
}

func NewClientAgent(cfg *config.Config) (*ClientGRPC, error) {
	opts := []grpc.DialOption{}

	if cfg.CryptoKey != nil {
		tlsCreds, err := generateTLSCreds(*cfg.CryptoKey)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(tlsCreds))
	}

	conn, err := grpc.NewClient(*cfg.ServerAddr, opts...)
	if err != nil {
		return nil, err
	}

	return &ClientGRPC{
		conn:       conn,
		authCLient: pb.NewAuthServiceClient(conn),
	}, nil
}

func (c *ClientGRPC) ConnClose() {
	c.conn.Close()
}

func (c *ClientGRPC) UserRegister(ctx context.Context, user *models.User) (*models.User, error) {
	req := &pb.UserMessageRequest{Login: user.Login, Password: user.Password}
	resp, err := c.authCLient.Register(ctx, req)
	if err != nil {
		return user, err
	}
	user.Token = resp.Token
	return user, nil
}

func (c *ClientGRPC) UserLogin(ctx context.Context, user *models.User) (*models.User, error) {
	req := &pb.UserMessageRequest{Login: user.Login, Password: user.Password}
	resp, err := c.authCLient.Login(ctx, req)
	if err != nil {
		return user, err
	}
	user.Token = resp.Token
	return user, nil
}

func generateTLSCreds(certFile string) (credentials.TransportCredentials, error) {
	return credentials.NewClientTLSFromFile(certFile, "")
}
