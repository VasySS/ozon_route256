package middleware

import (
	"context"
	"workshop-1/internal/app/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		logger.Infof("метод: %s, метаданные запроса: %v", info.FullMethod, md)
	}

	rewReq, _ := protojson.Marshal((req).(proto.Message))
	logger.Infof("метод: %s, запрос: %s", info.FullMethod, string(rewReq))

	res, err := handler(ctx, req)
	if err != nil {
		logger.Infof("метод: %s, ошибка: %s", info.FullMethod, err.Error())
		return nil, err
	}

	respReq, _ := protojson.Marshal((res).(proto.Message))
	logger.Infof("метод: %s, ответ: %s", info.FullMethod, string(respReq))

	return res, nil
}
