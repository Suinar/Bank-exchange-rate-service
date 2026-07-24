package exchange_rate

import (
	context "context"

	ecxhangeRateProto "github.com/kVinsom/Bank-proto/exchange_rate"
)

// IExchangeRateHandler defines the exchange-rate operations exposed by the transport layer.
type IExchangeRateHandler interface {
	GetExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetRelativeRankingRequest) (*ecxhangeRateProto.Ranking, error)
	GetAllExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetAllRankingRequest) (*ecxhangeRateProto.RankingList, error)
}
