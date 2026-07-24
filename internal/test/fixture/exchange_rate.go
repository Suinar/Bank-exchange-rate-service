package fixture

import (
	time "time"

	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
	ranking "github.com/kVinsom/Bank-exhange-rate-service/pkg/core"
	ecxhangeRateProto "github.com/kVinsom/Bank-proto/exchange_rate"
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
		Date: timestamppb.New(time.Unix(
			testConstant.TestTimestamp,
			testConstant.TestTimestampNanoseconds,
		)),
		RateSell:  testConstant.TestRateSell,
		RateBuy:   testConstant.TestRateBuy,
		RateCross: testConstant.TestRateCross,
	}
}

// NewRankings returns a domain ranking collection with canonical test values.
func NewRankings() []ranking.Ranking {
	return []ranking.Ranking{*NewRanking()}
}

// NewCacheRankings returns rates for testing source-currency cache indexes.
func NewCacheRankings() []ranking.Ranking {
	return []ranking.Ranking{
		{
			CurrencyIsoFrom: testConstant.TestIsoFrom,
			CurrencyIsoTo:   testConstant.TestIsoTo,
			Date: time.Unix(
				testConstant.TestTimestamp,
				testConstant.TestTimestampNanoseconds,
			),
			RateSell:  testConstant.TestFirstCacheRateSell,
			RateBuy:   testConstant.TestFirstCacheRateBuy,
			RateCross: testConstant.TestFirstCacheRateCross,
		},
		{
			CurrencyIsoFrom: testConstant.TestIsoFrom,
			CurrencyIsoTo:   testConstant.TestIsoAlternative,
			Date: time.Unix(
				testConstant.TestSecondTimestamp,
				testConstant.TestTimestampNanoseconds,
			),
			RateSell:  testConstant.TestSecondCacheRateSell,
			RateBuy:   testConstant.TestSecondCacheRateBuy,
			RateCross: testConstant.TestSecondCacheRateCross,
		},
		{
			CurrencyIsoFrom: testConstant.TestIsoTo,
			CurrencyIsoTo:   testConstant.TestIsoFrom,
			Date: time.Unix(
				testConstant.TestThirdTimestamp,
				testConstant.TestTimestampNanoseconds,
			),
			RateSell:  testConstant.TestThirdCacheRateSell,
			RateBuy:   testConstant.TestThirdCacheRateBuy,
			RateCross: testConstant.TestThirdCacheRateCross,
		},
	}
}

// NewReplacementRanking returns a rate with a different target currency.
func NewReplacementRanking() *ranking.Ranking {
	rate := *NewRanking()
	rate.CurrencyIsoTo = testConstant.TestIsoAlternative

	return &rate
}

// NewReplacementRankings returns a cache collection containing the replacement rate.
func NewReplacementRankings() []ranking.Ranking {
	return []ranking.Ranking{*NewReplacementRanking()}
}

// NewInvalidRankingCacheData returns a malformed Redis hash representation.
func NewInvalidRankingCacheData() map[string]any {
	return map[string]any{
		testConstant.TestCurrencyIsoFromCacheField: testConstant.TestInvalidCacheValue,
	}
}

// NewRankingListProto returns a protobuf ranking list with canonical test values.
func NewRankingListProto() *ecxhangeRateProto.RankingList {
	return &ecxhangeRateProto.RankingList{
		Rankings: []*ecxhangeRateProto.Ranking{
			NewRankingProto(),
		},
	}
}

// NewMonoExchangeRateJSON returns a Monobank HTTP response fixture.
func NewMonoExchangeRateJSON() []byte {
	return []byte(testConstant.TestMonoExchangeRateJSON)
}

// NewRanking returns a domain ranking with canonical test values.
func NewRanking() *ranking.Ranking {
	return &ranking.Ranking{
		CurrencyIsoFrom: testConstant.TestIsoFrom,
		CurrencyIsoTo:   testConstant.TestIsoTo,
		Date: time.Unix(
			testConstant.TestTimestamp,
			testConstant.TestTimestampNanoseconds,
		),
		RateSell:  testConstant.TestRateSell,
		RateBuy:   testConstant.TestRateBuy,
		RateCross: testConstant.TestRateCross,
	}
}
