package handler

import (
	"context"

	proto "github.com/Suinar/Bank-exhange-rate-service/proto"
)

type IRankingHandler interface {
	GetAllRanking(ctx context.Context, req proto.GetAllRankingRequest) (*proto.GetAllRankingResponse, error)
	GetRelativeRanking(ctx context.Context, req *proto.GetRelativeRankingRequest) (*proto.GetRelativeRankingResponse, error)
}
