package exchange_rate

import (
	"context"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
)

//go:generate mockgen -source=interface_service.go -destination=../../../internal/mocks/services/exchangeRate.go -package=mocks

type IExchangeRateService interface {
	GetExchangeRate(ctx context.Context, currencyIdFrom int64, currencyIdTo int64) (*ecxhangeRateProto.RankingList, error)
	GetAllExchangeRate(ctx context.Context, currencyIdFrom int64) (*ecxhangeRateProto.RankingList, error)
}
