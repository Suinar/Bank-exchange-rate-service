package monobank

import (
	monobank "github.com/kVinsom/Bank-exhange-rate-service/internal/external/monobank"
	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
)

// NewMonoExchangeRate returns a Monobank model with canonical test values.
func NewMonoExchangeRate() *monobank.MonoExchangeRate {
	return &monobank.MonoExchangeRate{
		CurrencyCodeA: testConstant.TestIsoFrom,
		CurrencyCodeB: testConstant.TestIsoTo,
		Date:          testConstant.TestTimestamp,
		RateSell:      testConstant.TestRateSell,
		RateBuy:       testConstant.TestRateBuy,
		RateCross:     testConstant.TestRateCross,
	}
}
