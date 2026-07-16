package main

import (
	context "context"
	fmt "fmt"
	log "log"
	net "net"
	os "os"
	osSignal "os/signal"
	syscall "syscall"
	time "time"

	configs "github.com/Suinar/Bank-exhange-rate-service/internal/configs"
	deliveryGRPC "github.com/Suinar/Bank-exhange-rate-service/internal/delivery/grpc"
	monobank "github.com/Suinar/Bank-exhange-rate-service/internal/external/monobank"
	services "github.com/Suinar/Bank-exhange-rate-service/internal/services"
	exchangeRateService "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg := configs.Load()

	// The same binary is used by Docker and Kubernetes probes to avoid
	// shipping an additional diagnostic utility in the distroless image.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := checkHealth(cfg.GRPC.Port); err != nil {
			log.Print(err)
			os.Exit(1)
		}

		return
	}

	// Assemble dependencies explicitly at the composition root so transport,
	// application, and external integration layers remain decoupled.
	monoClient := monobank.NewMonobankClient(cfg.Monobank.BaseURL, cfg.Monobank.CurrencyEndpoint)
	exchangeService := exchangeRateService.NewExchangeRateService(monoClient)
	appServices := services.InitServices(exchangeService)

	// SIGTERM is sent by container orchestrators during rolling updates. Passing
	// its cancellation downstream allows the gRPC server to drain active calls.
	ctx, stop := osSignal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("exchange-rate service listening on :%s", cfg.GRPC.Port)
	if err := deliveryGRPC.RunGrpcServer(ctx, cfg, appServices); err != nil {
		log.Fatal(err)
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
