package grpc

import (
	"log"
	"net"

	"github.com/Suinar/Bank-exhange-rate-service/internal/configs"
	service "github.com/Suinar/Bank-exhange-rate-service/internal/services"

	ecxhangeRateHandler "github.com/Suinar/Bank-exhange-rate-service/internal/delivery/grpc/handler/exchange_rate"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"

	"google.golang.org/grpc"
)

func RunGrpcServer(cfg *configs.Config, services *service.Services) {
	lis, err := net.Listen(cfg.GRPC.GRPCPort, cfg.GRPC.Network)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	RegisterServices(grpcServer, services)

	log.Println("repository-service running on :" + cfg.GRPC.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func RegisterServices(
	grpcServer *grpc.Server,
	services *service.Services) {
	ecxhangeRateProto.RegisterRankingRepositoryServer(
		grpcServer,
		ecxhangeRateHandler.NewExchangeRateHandlerHandler(services.ExchangeRateService),
	)

}
