package exchange_rate

import (
	context "context"

	ranking "github.com/kVinsom/Bank-exhange-rate-service/pkg/core"
)

//go:generate mockgen -source=intefrace.go -destination=../../../mocks/cache/exchangeRateCache.go -package=mocks

// IExchangeRateCache defines operations for storing and retrieving
// exchange rates from the cache.
type IExchangeRateCache interface {
	// AddAllExchangeRate stores all provided exchange rates.
	AddAllExchangeRate(ctx context.Context, rates []ranking.Ranking) error

	// GetExchangeRate returns an exchange rate for the specified currency pair.
	GetExchangeRate(ctx context.Context, currencyIsoFrom int32, currencyIsoTo int32) (*ranking.Ranking, error)

	// GetAllExchangeRate returns all exchange rates for the specified source currency.
	GetAllExchangeRate(ctx context.Context, currencyIsoFrom int32) ([]ranking.Ranking, error)
}
