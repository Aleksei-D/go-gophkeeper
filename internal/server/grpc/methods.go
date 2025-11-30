package server

import (
	"context"
	"errors"
	"go-gophkeeper/internal/models"
	pb "go-gophkeeper/internal/pb/auth"
	errors2 "go-gophkeeper/internal/utils/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthServer) Login(ctx context.Context, req *pb.UserMessageRequest) (*pb.UserMessageResponse, error) {
	var loginResponse pb.UserMessageResponse

	token, err := a.service.Login(ctx, &models.User{Login: req.Login, Password: req.Password})
	if err != nil {
		if errors.Is(err, errors2.ErrInvalidCredentials) {
			return &loginResponse, status.Error(codes.Unauthenticated, err.Error())
		}
		return &loginResponse, status.Error(codes.InvalidArgument, err.Error())
	}

	loginResponse.Token = token
	return &loginResponse, nil
}

func (a *AuthServer) Register(ctx context.Context, req *pb.UserMessageRequest) (*pb.UserMessageResponse, error) {
	var registerResponse pb.UserMessageResponse

	token, err := a.service.CreateUser(ctx, &models.User{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, errors2.ErrUserAlreadyExists) {
			return &registerResponse, status.Error(codes.AlreadyExists, err.Error())
		}
		return &registerResponse, status.Error(codes.InvalidArgument, err.Error())
	}

	registerResponse.Token = token
	return &registerResponse, nil
}
