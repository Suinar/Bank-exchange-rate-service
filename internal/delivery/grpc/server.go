package grpc

import (
	context "context"
	errors "errors"
	fmt "fmt"
	net "net"
	time "time"

	configs "github.com/Suinar/Bank-exhange-rate-service/internal/configs"
	exchangeRateHandler "github.com/Suinar/Bank-exhange-rate-service/internal/delivery/grpc/handler/exchange_rate"
	services "github.com/Suinar/Bank-exhange-rate-service/internal/services"
	exchangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
	grpc "google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// RunGrpcServer serves gRPC requests until the context is cancelled or serving fails.
func RunGrpcServer(ctx context.Context, cfg *configs.Config, appServices *services.Services) error {
	address := net.JoinHostPort("", cfg.GRPC.Port)
	listener, err := net.Listen(cfg.GRPC.Network, address)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", address, err)
	}

	server := grpc.NewServer()
	RegisterServices(server, appServices)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

		stopped := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(stopped)
		}()

		// Do not let a stuck client stream block pod termination indefinitely.
		select {
		case <-stopped:
		case <-time.After(25 * time.Second):
			server.Stop()
		}

		return nil
	case err = <-serveErr:
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}

		return fmt.Errorf("serve gRPC: %w", err)
	}
}

// RegisterServices binds application handlers to the gRPC server.
func RegisterServices(server *grpc.Server, appServices *services.Services) {
	exchangeRateProto.RegisterRankingRepositoryServer(
		server,
		exchangeRateHandler.NewExchangeRateHandler(appServices.ExchangeRateService),
	)
}
