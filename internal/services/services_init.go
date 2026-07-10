package services

import exchangeRate "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"

type Services struct {
	ExchangeRateService exchangeRate.ExchangeRateService
}

func InitServices(exchangeRateService exchangeRate.ExchangeRateService) *Services {
	return &Services{
		ExchangeRateService: exchangeRateService,
	}
}
