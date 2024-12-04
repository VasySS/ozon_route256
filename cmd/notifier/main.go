package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"workshop-1/config"
	"workshop-1/internal/app/kafka/consumer"
	"workshop-1/internal/app/logger"

	"github.com/IBM/sarama"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx := signalHandler(context.Background(), wg)

	signalHandler(ctx, wg)

	handler := consumer.NewHandler()
	cg, err := consumer.NewConsumerGroup(
		[]string{config.KafkaBrokers},
		config.KafkaConsumerGroupID,
		[]string{config.KafkaPVZTopic},
		handler,
		consumer.WithOffsetsInitial(sarama.OffsetNewest),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer cg.Close()

	errorHandler(ctx, cg, wg)

	cg.Run(ctx, wg)

	wg.Wait()
}

func signalHandler(ctx context.Context, wg *sync.WaitGroup) context.Context {
	sigCtx, cancel := context.WithCancel(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer signal.Stop(sigChan)
		defer wg.Done()
		defer cancel()

		select {
		case sig := <-sigChan:
			logger.Info("получен сигнал:", sig)
			return
		case <-sigCtx.Done():
			logger.Info("контекст отменен")
			return
		}
	}()

	return sigCtx
}

func errorHandler(ctx context.Context, cg sarama.ConsumerGroup, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case err, ok := <-cg.Errors():
				if !ok {
					return
				}

				logger.Error(err)
			case <-ctx.Done():
				return
			}
		}
	}()
}
