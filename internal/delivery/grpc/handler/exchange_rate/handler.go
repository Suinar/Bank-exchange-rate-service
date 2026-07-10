package exchange_rate

import (
	"context"

	service "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"

	errors "github.com/Suinar/Bank-repository-service/pkg"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
)

type ExchangeReteHandler struct {
	service service.ExchangeRateService

	ecxhangeRateProto.UnimplementedRankingRepositoryServer
}

func NewExchangeRateHandlerHandler(service service.ExchangeRateService,
) *ExchangeReteHandler {
	return &ExchangeReteHandler{
		service: service,
	}
}

func (h *ExchangeReteHandler) GetExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetRelativeRankingRequest,
) (*ecxhangeRateProto.RankingList, error) {

	if req.CurrencyIdFrom == req.CurrencyIdTo {
		return nil, errors.BadRequest
	}

	return h.service.GetExchangeRate(ctx, req.CurrencyIdFrom, req.CurrencyIdTo)
}

func (h *ExchangeReteHandler) GetAllExchangeRate(ctx context.Context, req *ecxhangeRateProto.GetAllRankingRequest,
) (*ecxhangeRateProto.RankingList, error) {

	return h.service.GetAllExchangeRate(ctx, req.CurrencyIdFrom)
}
