package handler

import (
	"context"

	ranking "github.com/Suinar/Bank-exhange-rate-service/gen/github.com/Suinar/Bank-exhange-rate-service/gen/ranking"
)

type IRankingHandler interface {
	GetAllRanking(ctx context.Context, req ranking.GetAllRankingRequest) (*ranking.GetAllRankingResponse, error)
	GetRelativeRanking(ctx context.Context, req *ranking.GetRelativeRankingRequest) (*ranking.GetRelativeRankingResponse, error)
}
