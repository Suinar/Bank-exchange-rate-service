package services

import exchangeRate "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"

// Services groups the application services used by transport layers.
type Services struct {
	ExchangeRateService *exchangeRate.ExchangeRateService
}

// InitServices assembles the application service dependencies.
func InitServices(exchangeRateService *exchangeRate.ExchangeRateService) *Services {
	return &Services{
		ExchangeRateService: exchangeRateService,
	}
}
