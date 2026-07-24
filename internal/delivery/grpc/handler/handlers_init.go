package handler

import handler "github.com/kVinsom/Bank-exhange-rate-service/internal/delivery/grpc/handler/exchange_rate"

// Handlers groups the transport handlers exposed by the application.
type Handlers struct {
	ExchangeRateHandler handler.IExchangeRateHandler
}

// InitHandlers assembles the transport handlers used by the gRPC server.
func InitHandlers(exchangeRateHandler handler.IExchangeRateHandler) *Handlers {
	return &Handlers{
		ExchangeRateHandler: exchangeRateHandler,
	}
}
