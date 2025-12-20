package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-gophkeeper/internal/config"
	"go-gophkeeper/internal/domain/mock"
	"go-gophkeeper/internal/logger"
	"go-gophkeeper/internal/models"
	"go-gophkeeper/internal/pb/auth"
	"go-gophkeeper/internal/server/service"
	"go-gophkeeper/internal/utils/crypto/mocks"
	errors2 "go-gophkeeper/internal/utils/errors"
	"go-gophkeeper/internal/utils/token"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const (
	validUsername   = "validUsername"
	validPassword   = "validPassword"
	newUsername     = "newUsername"
	invalidUsername = "invalidUsername"
	invalidPassword = "invalidPassword"
)

func TestRegisterMethod(t *testing.T) {
	tests := []struct {
		name     string
		error    bool
		code     codes.Code
		login    string
		password string
	}{
		{
			name:     "positive test for register",
			error:    false,
			code:     codes.OK,
			login:    newUsername,
			password: validPassword,
		},
		{
			name:     "negative test for register",
			error:    true,
			code:     codes.AlreadyExists,
			login:    validUsername,
			password: validPassword,
		},
	}

	err := config.InitDefaultEnv()
	assert.NoError(t, err)

	cfg, err := config.InitConfig()
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mock.NewMockAuthRepository(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	jwtManager := token.NewJWTManager(*cfg.Key, cfg.TokenDuration.Duration)
	serviceApp := service.NewService(
		authRepo,
		mock.NewMockVaultRepository(ctrl),
		mock.NewMockEventRepository(ctrl),
		jwtManager,
		mockHasher,
	)

	listener := bufconn.Listen(1024 * 1024)
	server, err := NewApp(listener, serviceApp, cfg, jwtManager)

	assert.NoError(t, err)
	defer server.Stop()

	mockHasher.EXPECT().HashPassword(validPassword).Return(validPassword, nil).AnyTimes()

	authRepo.EXPECT().IsExist(gomock.Any(), newUsername).Return(false, nil).AnyTimes()
	authRepo.EXPECT().IsExist(gomock.Any(), validUsername).Return(true, nil).AnyTimes()
	authRepo.EXPECT().Create(gomock.Any(), &models.User{Login: newUsername, Password: validPassword}).Return(&models.User{Login: newUsername, Password: validPassword}, nil).AnyTimes()

	go func() {
		if err := server.Run(); err != nil {
			logger.Log.Fatal("can not start handlers server", zap.Error(err))
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	authClient := auth.NewAuthServiceClient(conn)
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			resp, err := authClient.Register(context.Background(), &auth.UserMessageRequest{Login: v.login, Password: v.password})
			if v.error {
				assert.ErrorContains(t, err, "user already exists")
			} else {
				assert.NoError(t, err)
				claims, err := token.ExtractUserClaimsFromToken(resp.Token)
				assert.NoError(t, err)
				assert.Equal(t, v.login, claims.Login)
			}
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, v.code, st.Code())
		})
	}
}

func TestLoginMethod(t *testing.T) {
	tests := []struct {
		name      string
		error     bool
		code      codes.Code
		login     string
		password  string
		textError string
	}{
		{
			name:      "positive test for login",
			error:     false,
			code:      codes.OK,
			login:     validUsername,
			password:  validPassword,
			textError: "",
		},
		{
			name:      "negative test for login",
			error:     true,
			code:      codes.Unauthenticated,
			login:     validUsername,
			password:  invalidPassword,
			textError: "invalid credentials",
		},
		{
			name:      "negative test for login with new user",
			error:     true,
			code:      codes.InvalidArgument,
			login:     newUsername,
			password:  invalidPassword,
			textError: "no content",
		},
	}

	err := config.InitDefaultEnv()
	assert.NoError(t, err)

	cfg, err := config.InitConfig()
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mock.NewMockAuthRepository(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	jwtManager := token.NewJWTManager(*cfg.Key, cfg.TokenDuration.Duration)
	serviceApp := service.NewService(
		authRepo,
		mock.NewMockVaultRepository(ctrl),
		mock.NewMockEventRepository(ctrl),
		jwtManager,
		mockHasher,
	)

	listener := bufconn.Listen(1024 * 1024)
	server, err := NewApp(listener, serviceApp, cfg, jwtManager)

	assert.NoError(t, err)
	defer server.Stop()

	mockHasher.EXPECT().CheckPasswordHash(validPassword, validPassword).Return(true).AnyTimes()
	mockHasher.EXPECT().CheckPasswordHash(invalidPassword, validPassword).Return(false).AnyTimes()

	authRepo.EXPECT().GetByLogin(gomock.Any(), validUsername).Return(&models.User{Login: validUsername, Password: validPassword}, nil).AnyTimes()
	authRepo.EXPECT().GetByLogin(gomock.Any(), newUsername).Return(nil, errors2.ErrNoContent).AnyTimes()

	go func() {
		if err := server.Run(); err != nil {
			logger.Log.Fatal("can not start handlers server", zap.Error(err))
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	authClient := auth.NewAuthServiceClient(conn)
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			resp, err := authClient.Login(context.Background(), &auth.UserMessageRequest{Login: v.login, Password: v.password})
			if v.error {
				assert.ErrorContains(t, err, v.textError)
			} else {
				assert.NoError(t, err)
				claims, err := token.ExtractUserClaimsFromToken(resp.Token)
				assert.NoError(t, err)
				assert.Equal(t, v.login, claims.Login)
			}
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, v.code, st.Code())
		})
	}
}
