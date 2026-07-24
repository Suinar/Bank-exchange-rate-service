package exchange_rate

import (
	testing "testing"

	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
	fixture "github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	cacheErrors "github.com/kVinsom/Bank-repository-service/pkg"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestExchangeRateCache_GetExchangeRate_RedisError(t *testing.T) {
	sut, ctx := NewClosedSUT(t)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, cacheErrors.CacheGetError)
}

func TestExchangeRateCache_GetAllExchangeRate_RedisError(t *testing.T) {
	sut, ctx := NewClosedSUT(t)

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, cacheErrors.CacheGetError)
}

func TestExchangeRateCache_AddAllExchangeRate_RedisError(t *testing.T) {
	sut, ctx := NewClosedSUT(t)

	err := sut.AddAllExchangeRate(ctx, fixture.NewRankings())

	require.Error(t, err)
	assert.ErrorIs(t, err, cacheErrors.CacheSetError)
}

func TestExchangeRateCache_GetExchangeRate_InvalidCachedData(t *testing.T) {
	client, sut, ctx := NewSUT(t)
	key := sut.GetPrimaryKey(testConstant.TestIsoFrom, testConstant.TestIsoTo)
	require.NoError(t, client.HSet(ctx, key, fixture.NewInvalidRankingCacheData()).Err())

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, cacheErrors.CacheMapError)
}

func TestExchangeRateCache_GetAllExchangeRate_InvalidCachedData(t *testing.T) {
	client, sut, ctx := NewSUT(t)
	key := sut.GetPrimaryKey(testConstant.TestIsoFrom, testConstant.TestIsoTo)
	require.NoError(t, client.HSet(ctx, key, fixture.NewInvalidRankingCacheData()).Err())
	require.NoError(t, client.SAdd(ctx, sut.GetCurrencyFromKey(testConstant.TestIsoFrom), key).Err())

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.ErrorIs(t, err, cacheErrors.CacheMapError)
	assert.Nil(t, result)
}
