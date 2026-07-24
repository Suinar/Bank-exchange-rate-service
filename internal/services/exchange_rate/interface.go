package exchange_rate

import (
	context "context"

	ecxhangeRateProto "github.com/kVinsom/Bank-proto/exchange_rate"
)

//go:generate mockgen -source=interface_service.go -destination=../../../internal/mocks/services/exchangeRate.go -package=mocks

// IExchangeRateService defines exchange-rate use cases available to transport layers.
type IExchangeRateService interface {
	// GetExchangeRate returns the rate for an exact source and target currency pair.
	GetExchangeRate(ctx context.Context, currencyIsoFrom int32, currencyIsoTo int32) (*ecxhangeRateProto.Ranking, error)

	// GetAllExchangeRate returns all available rates for a source currency.
	GetAllExchangeRate(ctx context.Context, currencyIsoFrom int32) (*ecxhangeRateProto.RankingList, error)
}
