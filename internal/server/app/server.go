package app

import (
	"crypto/tls"
	"fmt"
	"go-gophkeeper/internal/config"
	pbAuth "go-gophkeeper/internal/pb/auth"
	pbVault "go-gophkeeper/internal/pb/vault"
	"go-gophkeeper/internal/server/interceptors"
	"go-gophkeeper/internal/server/service"
	"go-gophkeeper/internal/utils/token"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// App приложение сервера
type App struct {
	Listener   net.Listener
	gRPCServer *grpc.Server
	addrGRPC   string
}

// Run запуск сервера
func (a *App) Run() error {
	if err := a.gRPCServer.Serve(a.Listener); err != nil {
		return err
	}

	return nil
}

// Stop остановка сервера
func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}

// NewApp возврат нового приложения
func NewApp(listener net.Listener, service *service.Service, cfg *config.Config, jwtManager *token.JWTManager) (*App, error) {
	var serverOpts []grpc.ServerOption
	var interceptorsOpts []grpc.UnaryServerInterceptor

	authInterceptor := interceptors.NewAuthInterceptor(jwtManager)
	interceptorsOpts = append(interceptorsOpts, interceptors.LoggingInterceptor, authInterceptor.Unary())

	if cfg.Cert != nil {
		tlsCredentials, err := loadTLSCredentials(*cfg.Cert, *cfg.CryptoKey)
		if err != nil {
			return nil, err
		}

		serverOpts = append(serverOpts, grpc.Creds(tlsCredentials))
	}

	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(interceptorsOpts...))
	gRPCServer := grpc.NewServer(serverOpts...)
	registerServer(gRPCServer, service)

	return &App{
		Listener:   listener,
		gRPCServer: gRPCServer,
		addrGRPC:   *cfg.ServerAddr,
	}, nil
}

func loadTLSCredentials(certFilePath, privateKeyFilePath string) (credentials.TransportCredentials, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	certPem, err := os.ReadFile(certFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading certificate key file: %v", err)
	}

	tlsCert, err := tls.X509KeyPair(certPem, privateKeyPEM)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(cfg), nil
}

func registerServer(gRPCServer *grpc.Server, service *service.Service) {
	pbAuth.RegisterAuthServiceServer(gRPCServer, &AuthServer{Service: service})
	pbVault.RegisterVaultServiceServer(gRPCServer, &VaultServer{Service: service})
}
