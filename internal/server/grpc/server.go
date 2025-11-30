package server

import (
	"crypto/tls"
	"fmt"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/interceptors"
	pb "go-gophkeeper/internal/pb/auth"
	"go-gophkeeper/internal/service"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	gRPCServer *grpc.Server
	addrGRPC   string
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", a.addrGRPC)
	if err != nil {
		return err
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}

func NewApp(service *service.Service, cfg *config.Config) (*App, error) {
	var serverOpts []grpc.ServerOption
	var interceptorsOpts []grpc.UnaryServerInterceptor

	interceptorsOpts = append(interceptorsOpts, interceptors.LoggingInterceptor)

	if cfg.Cert != nil {
		tlsCredentials, err := loadTLSCredentials(*cfg.Cert, *cfg.CryptoKey)
		if err != nil {
			return nil, err
		}

		serverOpts = append(serverOpts, grpc.Creds(tlsCredentials))
	}

	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(interceptorsOpts...))
	gRPCServer := grpc.NewServer(serverOpts...)
	registerAuthServer(gRPCServer, service.AuthService)

	return &App{
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

	config := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func registerAuthServer(gRPCServer *grpc.Server, service *service.AuthService) {
	pb.RegisterAuthServiceServer(gRPCServer, &AuthServer{service: service})
}

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service *service.AuthService
}
