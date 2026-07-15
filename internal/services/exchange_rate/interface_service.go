package exchange_rate

import (
	"context"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
)

//go:generate mockgen -source=interface_service.go -destination=../../../internal/mocks/services/exchangeRate.go -package=mocks

type IExchangeRateService interface {
	GetExchangeRate(ctx context.Context, currencyIsoFrom int32, currencyIsoTo int32) (*ecxhangeRateProto.Ranking, error)
	GetAllExchangeRate(ctx context.Context, currencyIsoFrom int32) (*ecxhangeRateProto.RankingList, error)
}
