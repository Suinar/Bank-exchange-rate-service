package main

import (
	context "context"
	fmt "fmt"
	net "net"
	os "os"
	osSignal "os/signal"
	syscall "syscall"
	time "time"

	kafkaBroker "github.com/kVinsom/Bank-exhange-rate-service/internal/brokers/kafka"
	configs "github.com/kVinsom/Bank-exhange-rate-service/internal/configs"
	deliveryGRPC "github.com/kVinsom/Bank-exhange-rate-service/internal/delivery/grpc"
	monobank "github.com/kVinsom/Bank-exhange-rate-service/internal/external/monobank"
	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/app"
	cahce "github.com/kVinsom/Bank-exhange-rate-service/internal/repositories/cahce"
	exchangeRateCache "github.com/kVinsom/Bank-exhange-rate-service/internal/repositories/cahce/exchange_rate"
	services "github.com/kVinsom/Bank-exhange-rate-service/internal/services"
	exchangeRateService "github.com/kVinsom/Bank-exhange-rate-service/internal/services/exchange_rate"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	os.Exit(run())
}

func run() int {
	log.Configure()
	log.Started()
	defer log.Stopped()

	log.ConfigurationLoading()
	cfg, err := configs.Load()
	if err != nil {
		log.ConfigurationFailed(err)
		return 1
	}
	log.ConfigurationLoaded(net.JoinHostPort("", cfg.GRPC.Port), cfg.Redis.Addr, cfg.Monobank.RefreshInterval)

	// The same binary is used by Docker and Kubernetes probes to avoid
	// shipping an additional diagnostic utility in the distroless image.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := checkHealth(cfg.GRPC.Port); err != nil {
			log.GRPCServerFailed(fmt.Errorf("health check: %w", err))
			return 1
		}

		return 0
	}

	// SIGTERM is sent by container orchestrators during rolling updates. Passing
	// its cancellation downstream allows the gRPC server to drain active calls.
	ctx, stop := osSignal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.RedisConnecting(cfg.Redis.Addr)
	redisCtx, redisCancel := context.WithTimeout(ctx, 5*time.Second)
	cacheClient, err := cahce.ConnectToCahce(redisCtx, cfg)
	redisCancel()
	if err != nil {
		log.RedisConnectionFailed(cfg.Redis.Addr, err)
		return 1
	}
	defer cacheClient.Close()
	log.RedisConnected(cfg.Redis.Addr)

	log.ComponentsInitializing()
	rateCache := exchangeRateCache.NewExchangeRateCache(cacheClient)
	appCaches := cahce.InitCaches(rateCache)

	// Assemble dependencies explicitly at the composition root so transport,
	// application, cache, and external integration layers remain decoupled.
	monoClient := monobank.NewMonobankClient(
		cfg.Monobank.BaseURL,
		cfg.Monobank.CurrencyEndpoint,
		appCaches.ExchangeRateCache,
	)
	defer monoClient.Close()
	log.InitialExchangeRatesLoading()
	if err := monoClient.GetAllExchangeRate(ctx); err != nil {
		log.InitialExchangeRatesLoadingFailed(err)
	} else {
		log.InitialExchangeRatesLoaded()
	}
	go refreshExchangeRates(ctx, monoClient, cfg.Monobank.RefreshInterval)

	exchangeService := exchangeRateService.NewExchangeRateService(appCaches.ExchangeRateCache)
	appServices := services.InitServices(exchangeService)
	log.ComponentsInitialized()

	if err := kafkaBroker.EnsureKafkaTopics(ctx, cfg); err != nil {
		log.GRPCServerFailed(err)
		return 1
	}

	log.GRPCServerStarting(net.JoinHostPort("", cfg.GRPC.Port))
	if err := deliveryGRPC.RunGrpcServer(ctx, cfg, appServices); err != nil {
		log.GRPCServerFailed(err)
		return 1
	}
	log.ShutdownSignalReceived()
	log.GRPCServerStopped()
	return 0
}

func refreshExchangeRates(ctx context.Context, client *monobank.MonobankClient, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.RefreshSchedulerStarted(interval)

	for {
		select {
		case <-ctx.Done():
			log.RefreshSchedulerStopped()
			return
		case <-ticker.C:
			startedAt := time.Now()
			if err := client.GetAllExchangeRate(ctx); err != nil {
				log.RefreshFailed(time.Since(startedAt), err)
				continue
			}
			log.RefreshCompleted(time.Since(startedAt))
		}
	}
}

// checkHealth verifies that the local gRPC server reports a serving state.
func checkHealth(port string) error {
	// Bound probe duration independently from the server lifecycle so a failed
	// connection cannot leave the container healthcheck hanging.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	address := net.JoinHostPort("127.0.0.1", port)
	connection, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer connection.Close()

	response, err := healthpb.NewHealthClient(connection).Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if response.Status != healthpb.HealthCheckResponse_SERVING {
		return fmt.Errorf("service is not healthy: %s", response.Status)
	}

	return nil
}
