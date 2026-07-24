package exchange_rate_test

import (
	context "context"
	testing "testing"

	cacheMocks "github.com/kVinsom/Bank-exhange-rate-service/internal/mocks/cache"
	exchangeRate "github.com/kVinsom/Bank-exhange-rate-service/internal/services/exchange_rate"
	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
	fixture "github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestExchangeRateService_GetExchangeRate_Success(t *testing.T) {
	t.Parallel()

	cache, sut, ctx := NewSUT(t)
	expected := fixture.NewRankingProto()

	cache.EXPECT().
		GetExchangeRate(gomock.Eq(ctx), testConstant.TestIsoFrom, testConstant.TestIsoTo).
		Return(fixture.NewRanking(), nil).
		Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestExchangeRateService_GetAllExchangeRate_Success(t *testing.T) {
	t.Parallel()

	cache, sut, ctx := NewSUT(t)
	expected := fixture.NewRankingListProto()

	cache.EXPECT().
		GetAllExchangeRate(gomock.Eq(ctx), testConstant.TestIsoFrom).
		Return(fixture.NewRankings(), nil).
		Times(1)

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestExchangeRateService_MapperToProto(t *testing.T) {
	t.Parallel()

	_, sut, _ := NewSUT(t)

	result := sut.MapperToProto(fixture.NewRanking())

	require.NotNil(t, result)
	assert.Equal(t, fixture.NewRankingProto(), result)
}

func NewSUT(
	t *testing.T,
) (*cacheMocks.MockIExchangeRateCache, *exchangeRate.ExchangeRateService, context.Context) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	cache := cacheMocks.NewMockIExchangeRateCache(ctrl)
	sut := exchangeRate.NewExchangeRateService(cache)

	return cache, sut, context.Background()
}
