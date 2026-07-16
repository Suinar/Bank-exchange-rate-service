package exchange_rate

import (
	context "context"

	service "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"

	errors "github.com/Suinar/Bank-repository-service/pkg"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
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

	// A relative rate requires two distinct currencies; reject the request at
	// the transport boundary before invoking application logic.
	if req == nil || req.CurrencyIsoFrom == req.CurrencyIsoTo {
		return nil, errors.BadRequest
	}

	return h.service.GetExchangeRate(ctx, req.CurrencyIsoFrom, req.CurrencyIsoTo)
}

// GetAllExchangeRate returns all rates available for the requested source currency.
func (h *ExchangeRateHandler) GetAllExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetAllRankingRequest,
) (*ecxhangeRateProto.RankingList, error) {
	if req == nil {
		return nil, errors.BadRequest
	}

	return h.service.GetAllExchangeRate(ctx, req.CurrencyIsoFrom)
}
