package handler

import (
	"context"

	ranking "github.com/Suinar/Bank-exhange-rate-service/gen/github.com/Suinar/Bank-exhange-rate-service/gen/ranking"
)

type RankingHandler struct {
	ranking.UnimplementedRankingServiceServer
}

func (h *RankingHandler) GetAllRanking(ctx context.Context, req *ranking.GetAllRankingRequest) (
	*ranking.GetAllRankingResponse, error) {
	// todo: implement method.

	panic("implement me")
}

func (h *RankingHandler) GetRelativeRanking(ctx context.Context, req *ranking.GetRelativeRankingRequest) (
	*ranking.GetRelativeRankingResponse, error) {
	// todo: implement method.

	panic("implement me")
}
