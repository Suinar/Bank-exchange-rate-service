package exchange_rate

import (
	context "context"

	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/service/exchange_rate"
	exchangeRateCache "github.com/kVinsom/Bank-exhange-rate-service/internal/repositories/cahce/exchange_rate"
	ranking "github.com/kVinsom/Bank-exhange-rate-service/pkg/core"
	ecxhangeRateProto "github.com/kVinsom/Bank-proto/exchange_rate"
	errors "github.com/kVinsom/Bank-repository-service/pkg"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type ExchangeRateService struct {
	cache exchangeRateCache.IExchangeRateCache
}

func NewExchangeRateService(cache exchangeRateCache.IExchangeRateCache) *ExchangeRateService {
	return &ExchangeRateService{
		cache: cache,
	}
}

func (s *ExchangeRateService) GetExchangeRate(ctx context.Context, currencyIsoFrom int32, currencyIsoTo int32,
) (*ecxhangeRateProto.Ranking, error) {
	const operation = "get_exchange_rate"
	defer log.OperationStarted(operation)()

	rate, err := s.cache.GetExchangeRate(ctx, currencyIsoFrom, currencyIsoTo)
	if err != nil {
		log.DependencyFailed(operation, "exchange_rate_cache", err)
		return nil, err
	}
	if rate != nil {
		return s.MapperToProto(rate), nil
	}

	return nil, errors.NotFound
}

func (s *ExchangeRateService) GetAllExchangeRate(ctx context.Context, currencyIsoFrom int32,
) (*ecxhangeRateProto.RankingList, error) {
	const operation = "get_all_exchange_rates"
	defer log.OperationStarted(operation)()

	rates, err := s.cache.GetAllExchangeRate(ctx, currencyIsoFrom)
	if err != nil {
		log.DependencyFailed(operation, "exchange_rate_cache", err)
		return nil, err
	}
	log.Result(operation, len(rates))

	return s.MapperListToProto(rates), nil
}

func (s *ExchangeRateService) MapperListToProto(rates []ranking.Ranking) *ecxhangeRateProto.RankingList {
	defer log.MappingStarted("ranking_list_to_proto")()

	protoRates := make([]*ecxhangeRateProto.Ranking, 0, len(rates))
	for i := range rates {
		protoRates = append(protoRates, s.MapperToProto(&rates[i]))
	}

	return &ecxhangeRateProto.RankingList{
		Rankings: protoRates,
	}
}

func (s *ExchangeRateService) MapperToProto(rate *ranking.Ranking) *ecxhangeRateProto.Ranking {
	const operation = "ranking_to_proto"
	defer log.MappingStarted(operation)()
	if rate == nil {
		log.NilInput(operation)
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
