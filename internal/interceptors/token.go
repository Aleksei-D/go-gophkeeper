package interceptors

import (
	"context"

	"go-gophkeeper/internal/utils/token"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager *token.JWTManager
}

func NewAuthInterceptor(jwtManager *token.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/auth.AuthService/Login" { // todo check methods
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		accessToken := values[0]
		claims, err := interceptor.jwtManager.VerifyToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Optionally, add claims to context for handler access
		ctx = context.WithValue(ctx, "userClaims", claims)

		return handler(ctx, req)
	}
}
