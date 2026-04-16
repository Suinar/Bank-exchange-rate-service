package grps

import (
	"log"
	"net"

	handler "github.com/Suinar/Bank-exhange-rate-service/internal/transport/grps/handler"
	ranking "github.com/Suinar/Bank-exhange-rate-service/ranking"
	"google.golang.org/grpc"
)

func RanGrpsServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	ranking.RegisterRankingServiceServer(grpcServer, &handler.RankingHandler{})

	log.Println("gRPC server running on :50051")

	grpcServer.Serve(lis)
}
