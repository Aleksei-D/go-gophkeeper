package interceptors

import (
	"context"
	"go-gophkeeper/internal/pb/auth"

	"go-gophkeeper/internal/utils/token"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor интерцептор проверки токена
type AuthInterceptor struct {
	jwtManager *token.JWTManager
}

// NewAuthInterceptor возврат нового AuthInterceptor
func NewAuthInterceptor(jwtManager *token.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == auth.AuthService_Register_FullMethodName || info.FullMethod == auth.AuthService_Login_FullMethodName {
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

		ctx = context.WithValue(ctx, "userClaims", claims)
		return handler(ctx, req)
	}
}
