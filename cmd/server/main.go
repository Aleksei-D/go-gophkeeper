package main

import (
	"context"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/datasource"
	"go-gophkeeper/internal/domain/postgres"
	"go-gophkeeper/internal/logger"
	"go-gophkeeper/internal/server/app"
	"go-gophkeeper/internal/server/service"
	"go-gophkeeper/internal/utils/building"
	crypto2 "go-gophkeeper/internal/utils/crypto"
	"go-gophkeeper/internal/utils/token"
	"net"
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

	db, err := datasource.NewDatabase(*cfg.DatabaseDsn)
	if err != nil {
		logger.Log.Fatal("cannot init repo", zap.Error(err))
	}
	jwtManager := token.NewJWTManager(*cfg.Key, cfg.TokenDuration.Duration)
	serviceApp := service.NewService(
		postgres.NewPostgresUserRepository(db),
		postgres.NewVaultRepository(db),
		postgres.NewEventRepository(db),
		jwtManager,
		&crypto2.BcryptHasher{},
	)

	listener, err := net.Listen("tcp", *cfg.ServerAddr)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	serve, err := app.NewApp(listener, serviceApp, cfg, jwtManager)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		if err := serve.Run(); err != nil {
			logger.Log.Fatal("can not start handlers server", zap.Error(err))
		}
	}()

	<-serverCtx.Done()
	serve.Stop()
}
