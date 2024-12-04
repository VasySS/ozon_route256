package main

import (
	"context"
	"os/signal"
	"syscall"
	"workshop-1/internal/app/grpc"
	"workshop-1/internal/app/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := grpc.Run(ctx); err != nil {
		logger.Error(err)
	}
}
