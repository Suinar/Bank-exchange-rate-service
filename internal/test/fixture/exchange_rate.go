package fixture

import (
	time "time"

	monobank "github.com/Suinar/Bank-exhange-rate-service/internal/external/monobank"
	testConstant "github.com/Suinar/Bank-exhange-rate-service/internal/test"
	ranking "github.com/Suinar/Bank-exhange-rate-service/pkg/core"
	ecxhangeRateProto "github.com/Suinar/Bank-proto/exchange_rate"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// NewGetRelativeRankingRequestProto returns a valid request for a currency pair.
func NewGetRelativeRankingRequestProto() *ecxhangeRateProto.GetRelativeRankingRequest {
	return &ecxhangeRateProto.GetRelativeRankingRequest{
		CurrencyIsoFrom: testConstant.TestIsoFrom,
		CurrencyIsoTo:   testConstant.TestIsoTo,
	}
}

// NewGetAllRankingRequestProto returns a valid request for all source-currency rates.
func NewGetAllRankingRequestProto() *ecxhangeRateProto.GetAllRankingRequest {
	return &ecxhangeRateProto.GetAllRankingRequest{
		CurrencyIsoFrom: testConstant.TestIsoFrom,
	}
}

// NewRankingProto returns the canonical protobuf ranking used by tests.
func NewRankingProto() *ecxhangeRateProto.Ranking {
	return &ecxhangeRateProto.Ranking{
		CurrencyIsoFrom: testConstant.TestIsoFrom,
		CurrencyIsoTo:   testConstant.TestIsoTo,
		Date:            timestamppb.New(time.Unix(1710000000, 0)),
		RateSell:        testConstant.TestRateSell,
		RateBuy:         testConstant.TestRateBuy,
		RateCross:       testConstant.TestRateCross,
	}
}

// NewRankings returns a domain ranking collection with canonical test values.
func NewRankings() []ranking.Ranking {
	return []ranking.Ranking{*NewRanking()}
}

// NewRankingListProto returns a protobuf ranking list with canonical test values.
func NewRankingListProto() *ecxhangeRateProto.RankingList {
	return &ecxhangeRateProto.RankingList{
		Rankings: []*ecxhangeRateProto.Ranking{
			NewRankingProto(),
		},
	}
}

// NewMonoExchangeRate returns a Monobank model with canonical test values.
func NewMonoExchangeRate() *monobank.MonoExchangeRate {
	return &monobank.MonoExchangeRate{
		CurrencyCodeA: testConstant.TestIsoFrom,
		CurrencyCodeB: testConstant.TestIsoTo,
		Date:          1710000000,
		RateSell:      testConstant.TestRateSell,
		RateBuy:       testConstant.TestRateBuy,
		RateCross:     testConstant.TestRateCross,
	}
}

// NewRanking returns a domain ranking with canonical test values.
func NewRanking() *ranking.Ranking {
	return &ranking.Ranking{
		CurrencyIsoFrom: testConstant.TestIsoFrom,
		CurrencyIsoTo:   testConstant.TestIsoTo,
		Date:            time.Unix(1710000000, 0),
		RateSell:        testConstant.TestRateSell,
		RateBuy:         testConstant.TestRateBuy,
		RateCross:       testConstant.TestRateCross,
	}
}
