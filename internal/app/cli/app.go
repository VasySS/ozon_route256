package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"workshop-1/config"
	"workshop-1/internal/app"
	"workshop-1/internal/app/cli/cobra"
	"workshop-1/internal/domain"
	"workshop-1/internal/storage/inmem"
)

func Run(ctx context.Context, stop context.CancelFunc) error {
	closer := app.NewCloser()

	// valkey := app.NewValkeyCache(ctx, closer)
	orderCache := inmem.New[int, domain.Order](config.CacheTTL, inmem.LFU, config.CacheCapacity)
	orderReturnCache := inmem.New[int, domain.OrderReturn](config.CacheTTL, inmem.LRU, config.CacheCapacity)
	storage := app.NewStorageFacade(ctx, closer, orderCache, orderReturnCache)

	cobraCLI := cobra.New(ctx, storage)
	closer.AddWithCtx(cobraCLI.GracefulShutdown)

	reader := bufio.NewReader(os.Stdin)
	go cobraCLI.RunInteractive(stop, reader)

	<-ctx.Done()
	fmt.Println("Завершение работы приложения...")

	shutdownCtx, stop := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer stop()

	if err := closer.Close(shutdownCtx); err != nil {
		return err
	}

	return nil
}
