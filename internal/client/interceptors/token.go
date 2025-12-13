package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TokenClientInterceptor интерцептор клиента для прокидывания токена авторизации
func TokenClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var token string
	for _, opt := range opts {
		if t, ok := opt.(TokenOption); ok {
			token = t.Token
			break
		}
	}

	if token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// TokenOption callOption для передачи токена в TokenClientInterceptor
type TokenOption struct {
	grpc.EmptyCallOption
	Token string
}
