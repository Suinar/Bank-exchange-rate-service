package handler

import handler "github.com/Suinar/Bank-exhange-rate-service/internal/delivery/grpc/handler/exchange_rate"

type Handlers struct {
	ExchangeReteHandler handler.ExchangeReteHandler
}

func InitHandlers(exchangeReteHandler handler.ExchangeReteHandler) *Handlers {
	return &Handlers{
		ExchangeReteHandler: exchangeReteHandler,
	}
}
