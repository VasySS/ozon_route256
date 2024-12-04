package middleware

import (
	"context"
	"workshop-1/internal/usecase/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]struct{}{
	"/pvz.PVZService/RegisterUser": {},
}

func Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if _, ok := publicMethods[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "не предоставлены метаданные")
	}

	token := md.Get("X-Api-Token")
	if len(token) == 0 {
		return nil, status.Error(codes.Unauthenticated, "не предоставлен токен")
	}

	if _, err := auth.ValidateAccessToken(token[0]); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return handler(ctx, req)
}

func AuthHeaderMatcher(key string) (string, bool) {
	switch key {
	case "X-Api-Token":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
