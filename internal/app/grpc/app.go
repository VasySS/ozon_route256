package grpc

import (
	"context"
	"net"
	"net/http"
	"os"

	"workshop-1/config"
	"workshop-1/internal/app"
	"workshop-1/internal/app/grpc/handler"
	"workshop-1/internal/app/grpc/middleware"
	"workshop-1/internal/app/kafka"
	"workshop-1/internal/app/kafka/event"
	"workshop-1/internal/app/logger"
	"workshop-1/internal/domain"
	"workshop-1/internal/storage/inmem"
	"workshop-1/internal/storage/postgres"
	"workshop-1/pkg/pvz/v1"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(ctx context.Context) error {
	closer := app.NewCloser()

	// valkey := app.NewValkeyCache(ctx, closer)
	orderCache := inmem.New[int, domain.Order](config.CacheTTL, inmem.LFU, config.CacheCapacity)
	orderReturnCache := inmem.New[int, domain.OrderReturn](config.CacheTTL, inmem.LRU, config.CacheCapacity)
	storage := app.NewStorageFacade(ctx, closer, orderCache, orderReturnCache)

	eventFactory := event.NewDefaultFactory()
	kafkaProducer, err := kafka.New(
		[]string{config.KafkaBrokers},
		eventFactory,
		kafka.WithIdempotent(),
		kafka.WithMaxRetries(3),
		kafka.WithRequiredAcks(sarama.WaitForAll),
	)
	if err != nil {
		return err
	}
	closer.AddWithError(kafkaProducer.Close)

	go startSwagger(closer)
	go startHTTP(ctx, closer)
	go startGRPC(storage, kafkaProducer, closer)

	<-ctx.Done()
	logger.Info("Завершение работы приложения...")

	shutdownCtx, cancelTimeout := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancelTimeout()

	if err := closer.Close(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func startHTTP(ctx context.Context, closer *app.Closer) {
	logger.Info("Запуск HTTP сервиса по адресу", config.HTTPHost)

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(middleware.AuthHeaderMatcher))
	err := pvz.RegisterPVZServiceHandlerFromEndpoint(ctx, mux, config.GRPCHost,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		})
	if err != nil {
		logger.Fatal(err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/", middleware.CORS(mux))

	srv := &http.Server{
		Addr:    config.HTTPHost,
		Handler: httpMux,
	}
	closer.AddWithCtx(srv.Shutdown)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func startGRPC(storage postgres.Facade, kafka *kafka.Kafka, closer *app.Closer) {
	logger.Info("Запуск gRPC сервиса по адресу ", config.GRPCHost)

	lis, err := net.Listen("tcp", config.GRPCHost)
	if err != nil {
		logger.Fatal(err)
	}
	closer.AddWithError(lis.Close)

	pvzHandler := handler.NewHandler(storage, kafka)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Logger, middleware.Auth),
	)
	closer.Add(grpcServer.GracefulStop)

	reflection.Register(grpcServer)
	pvz.RegisterPVZServiceServer(grpcServer, pvzHandler)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(err)
	}
}

func startSwagger(closer *app.Closer) {
	logger.Info("Запуск Swagger по адресу ", config.SwaggerHost)

	mux := chi.NewMux()
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(config.SwaggerPath)
		if err != nil {
			logger.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	mux.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    config.SwaggerHost,
		Handler: mux,
	}
	closer.AddWithCtx(srv.Shutdown)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}
