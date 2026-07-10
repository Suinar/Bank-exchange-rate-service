package exchange_rate

import (
	"context"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
)

type IExchangeRateHandler interface {
	GetExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetRelativeRankingRequest) (*ecxhangeRateProto.RankingList, error)
	GetAllExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetAllRankingRequest) (*ecxhangeRateProto.RankingList, error)
}

