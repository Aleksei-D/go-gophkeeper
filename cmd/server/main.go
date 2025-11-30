package main

import (
	"context"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/datasource"
	"go-gophkeeper/internal/domain/infrastructure"
	"go-gophkeeper/internal/logger"
	grpc_server "go-gophkeeper/internal/server/grpc"
	"go-gophkeeper/internal/service"
	"go-gophkeeper/internal/utils/building"
	"go-gophkeeper/internal/utils/token"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var buildVersion, buildDate, buildCommit string

func main() {
	building.PrintBuildVersion(buildVersion, buildDate, buildCommit)

	err := logger.Initialize("INFO")
	if err != nil {
		logger.Log.Fatal("cannot initialize zap", zap.Error(err))
	}

	cfg, err := config.NewServerConfig()
	if err != nil {
		logger.Log.Fatal("cannot initialize config", zap.Error(err))
	}

	db, err := datasource.NewDatabase(*cfg.DatabaseDsn, datasource.Server)
	if err != nil {
		logger.Log.Fatal("cannot init repo", zap.Error(err))
	}

	serviceApp := service.NewService(service.NewAuthService(infrastructure.NewPostgresUserRepository(db), token.NewJWTManager(*cfg.CryptoKey, cfg.TokenDuration.Duration)))

	app, err := grpc_server.NewApp(serviceApp, cfg)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		if err := app.Run(); err != nil {
			logger.Log.Fatal("can not start grpc server", zap.Error(err))
		}
	}()

	<-serverCtx.Done()
	app.Stop()
}
