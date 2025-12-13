package cli

import (
	"context"
	"fmt"
	"go-gophkeeper/internal/client/api"
	"go-gophkeeper/internal/client/interceptors"
	"go-gophkeeper/internal/client/sync"
	"go-gophkeeper/internal/client/tui"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/datasource"
	"go-gophkeeper/internal/domain/postgres"
	"go-gophkeeper/internal/utils/crypto"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
)

// Cli клиент для работы с хранилищем
type Cli struct {
	conn         *grpc.ClientConn
	userTerminal *tui.Terminal
	syncClient   *sync.ClientSync
}

// NewCli создание нового клиента
func NewCli(cfg *config.Config) (*Cli, error) {
	var opts []grpc.DialOption

	if cfg.CryptoKey != nil {
		tlsCreds, err := crypto.GenerateTLSCreds(*cfg.CryptoKey)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(tlsCreds))
	}
	opts = append(opts, grpc.WithUnaryInterceptor(interceptors.TokenClientInterceptor))

	conn, err := grpc.NewClient(*cfg.ServerAddr, opts...)
	if err != nil {
		return nil, err
	}

	db, err := datasource.NewDatabase(*cfg.DatabaseDsn)
	if err != nil {
		return nil, err
	}

	vaultRepo := postgres.NewVaultRepository(db)
	eventRepo := postgres.NewEventRepository(db)

	userTerminal := tui.NewTerminal(api.NewClientAuth(conn), vaultRepo)
	syncClient := sync.NewClientSync(cfg, api.NewVaultClient(conn), vaultRepo, eventRepo)

	return &Cli{
		conn:         conn,
		userTerminal: userTerminal,
		syncClient:   syncClient,
	}, nil
}

// ConnClose закрытие соединений
func (c Cli) ConnClose() {
	c.conn.Close()
}

// Run запуск клиента
func (c Cli) Run() error {
	defer c.ConnClose()
	agentCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	err := c.userTerminal.Login(agentCtx)
	if err != nil {
		return err
	}

	c.syncClient.SetUser(c.userTerminal.User)
	go c.syncClient.Run(agentCtx)

	err = c.userTerminal.Run(agentCtx)
	if err != nil {
		return err
	}

	fmt.Println("ВСЕГО ХОРОШЕГО")
	return nil
}
