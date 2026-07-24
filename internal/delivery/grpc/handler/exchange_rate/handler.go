package exchange_rate

import (
	context "context"

	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/handler/exchange_rate"
	service "github.com/kVinsom/Bank-exhange-rate-service/internal/services/exchange_rate"

	errors "github.com/kVinsom/Bank-repository-service/pkg"

	ecxhangeRateProto "github.com/kVinsom/Bank-proto/exchange_rate"
)

// ExchangeRateHandler exposes exchange-rate operations through gRPC.
type ExchangeRateHandler struct {
	service service.IExchangeRateService

	ecxhangeRateProto.UnimplementedRankingRepositoryServer
}

// NewExchangeRateHandler creates a gRPC handler backed by the exchange-rate service.
func NewExchangeRateHandler(service service.IExchangeRateService) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		service: service,
	}
}

// GetExchangeRate validates a currency pair and returns its current exchange rate.
func (h *ExchangeRateHandler) GetExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetRelativeRankingRequest,
) (*ecxhangeRateProto.Ranking, error) {
	const operation = "get_exchange_rate"
	defer log.RequestStarted(operation)()

	// A relative rate requires two distinct currencies; reject the request at
	// the transport boundary before invoking application logic.
	if req == nil || req.CurrencyIsoFrom == req.CurrencyIsoTo {
		log.ValidationFailed(operation, "request is nil or currency pair contains identical currencies")
		return nil, errors.BadRequest
	}

	result, err := h.service.GetExchangeRate(ctx, req.CurrencyIsoFrom, req.CurrencyIsoTo)
	if err != nil {
		log.RequestError(operation, err)
	}
	return result, err
}

// GetAllExchangeRate returns all rates available for the requested source currency.
func (h *ExchangeRateHandler) GetAllExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetAllRankingRequest,
) (*ecxhangeRateProto.RankingList, error) {
	const operation = "get_all_exchange_rates"
	defer log.RequestStarted(operation)()
	if req == nil {
		log.ValidationFailed(operation, "request is nil")
		return nil, errors.BadRequest
	}

	result, err := h.service.GetAllExchangeRate(ctx, req.CurrencyIsoFrom)
	if err != nil {
		log.RequestError(operation, err)
	}
	return result, err
}
