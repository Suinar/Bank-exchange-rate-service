package exchange_rate

import (
	"context"

	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
)

type ExchangeRateService struct {
}

func NewExchangeRateService() *ExchangeRateService {
	return &ExchangeRateService{}
}

func (s *ExchangeRateService) GetExchangeRate(ctx context.Context, currencyIdFrom int64, currencyIdTo int64,
) (*ecxhangeRateProto.RankingList, error) {
	return nil, nil
}

func (s *ExchangeRateService) GetAllExchangeRate(ctx context.Context, currencyIdFrom int64,
) (*ecxhangeRateProto.RankingList, error) {
	return nil, nil
}
