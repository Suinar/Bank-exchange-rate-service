package cahce

import exchangeRate "github.com/kVinsom/Bank-exhange-rate-service/internal/repositories/cahce/exchange_rate"

// Caches groups the Redis-backed repositories used by the application.
type Caches struct {
	ExchangeRateCache exchangeRate.IExchangeRateCache
}

// InitCaches assembles the cache repositories used by application services.
func InitCaches(exchangeRateCache exchangeRate.IExchangeRateCache) *Caches {
	return &Caches{
		ExchangeRateCache: exchangeRateCache,
	}
}
