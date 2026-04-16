package handler

import (
	"context"

	proto "github.com/Suinar/Bank-exhange-rate-service/proto"
)

type RankingHandler struct {
	proto.UnimplementedRankingServiceServer
}

func (h *RankingHandler) GetAllRanking(ctx context.Context, req *proto.GetAllRankingRequest) (
	*proto.GetAllRankingResponse, error) {
	// todo: implement method.

	panic("implement me")
}

func (h *RankingHandler) GetRelativeRanking(ctx context.Context, req *proto.GetRelativeRankingRequest) (
	*proto.GetRelativeRankingResponse, error) {
	// todo: implement method.

	panic("implement me")
}
