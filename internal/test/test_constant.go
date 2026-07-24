package test

const (
	TestIsoFrom        int32 = 840
	TestIsoTo          int32 = 978
	TestIsoAlternative int32 = 980

	TestTimestamp            int64 = 1710000000
	TestSecondTimestamp      int64 = 1710000100
	TestThirdTimestamp       int64 = 1710000200
	TestTimestampNanoseconds int64 = 0

	TestRateSell  float32 = 3
	TestRateBuy   float32 = 2
	TestRateCross float32 = 2.5

	TestFirstCacheRateSell   float32 = 1.1
	TestFirstCacheRateBuy    float32 = 1.2
	TestFirstCacheRateCross  float32 = 1.3
	TestSecondCacheRateSell  float32 = 2.1
	TestSecondCacheRateBuy   float32 = 2.2
	TestSecondCacheRateCross float32 = 2.3
	TestThirdCacheRateSell   float32 = 3.1
	TestThirdCacheRateBuy    float32 = 3.2
	TestThirdCacheRateCross  float32 = 3.3

	TestCurrencyIsoFromCacheField = "currency_iso_from"
	TestInvalidCacheValue         = "invalid"
	TestRedisImage                = "redis:7-alpine"
	TestMonoExchangeRateJSON      = `[{"currencyCodeA":840,"currencyCodeB":978,"date":1710000000,"rateSell":3,"rateBuy":2,"rateCross":2.5}]`
)
