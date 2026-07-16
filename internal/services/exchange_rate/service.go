package exchange_rate

import (
	context "context"

	monobank "github.com/Suinar/Bank-exhange-rate-service/internal/external/monobank"
	ranking "github.com/Suinar/Bank-exhange-rate-service/pkg/core"
	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
	errors "github.com/Suinar/Bank-repository-service/pkg"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type ExchangeRateService struct {
	monoClient monobank.IMonobankClient
}

func NewExchangeRateService(monoClient monobank.IMonobankClient) *ExchangeRateService {
	return &ExchangeRateService{monoClient: monoClient}
}

func (s *ExchangeRateService) GetExchangeRate(ctx context.Context, currencyIsoFrom int32, currencyIsoTo int32,
) (*ecxhangeRateProto.Ranking, error) {
	rates, err := s.monoClient.GetAllExchangeRate(ctx)
	if err != nil {
		return nil, err
	}

	for i := range rates {
		if rates[i].CurrencyIsoFrom == currencyIsoFrom && rates[i].CurrencyIsoTo == currencyIsoTo {
			return s.MapperToProto(&rates[i]), nil
		}
	}

	// The upstream request succeeded, but it did not contain the requested
	// currency pair. Keep this distinct from transport-level client errors.
	return nil, errors.InternalServerError
}

func (s *ExchangeRateService) GetAllExchangeRate(ctx context.Context, currencyIsoFrom int32,
) (*ecxhangeRateProto.RankingList, error) {
	rates, err := s.monoClient.GetAllExchangeRate(ctx)
	if err != nil {
		return nil, err
	}

	// Preallocate for the worst case while preserving an empty, non-nil list
	// when the source currency has no matching rates.
	protoRates := make([]*ecxhangeRateProto.Ranking, 0, len(rates))
	for i := range rates {
		if rates[i].CurrencyIsoFrom == currencyIsoFrom {
			protoRates = append(protoRates, s.MapperToProto(&rates[i]))
		}
	}

	return &ecxhangeRateProto.RankingList{
		Rankings: protoRates,
	}, nil
}

func (s *ExchangeRateService) MapperToProto(rate *ranking.Ranking) *ecxhangeRateProto.Ranking {
	if rate == nil {
		return nil
	}

	return &ecxhangeRateProto.Ranking{
		CurrencyIsoFrom: rate.CurrencyIsoFrom,
		CurrencyIsoTo:   rate.CurrencyIsoTo,
		Date:            timestamppb.New(rate.Date),
		RateSell:        rate.RateSell,
		RateBuy:         rate.RateBuy,
		RateCross:       rate.RateCross,
	}
}
