package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"workshop-1/internal/app/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := cli.Run(ctx, stop); err != nil {
		fmt.Println(err)
	}
}
