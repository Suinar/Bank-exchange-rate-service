package exchange_rate_test

import (
	"testing"

	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
	errors "github.com/kVinsom/Bank-repository-service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExchangeRateService_GetExchangeRate_CacheError(t *testing.T) {
	t.Parallel()

	cache, sut, ctx := NewSUT(t)
	cache.EXPECT().
		GetExchangeRate(gomock.Eq(ctx), testConstant.TestIsoFrom, testConstant.TestIsoTo).
		Return(nil, errors.TestError).
		Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.ErrorIs(t, err, errors.TestError)
	assert.Nil(t, result)
}

func TestExchangeRateService_GetExchangeRate_NotFound(t *testing.T) {
	t.Parallel()
	cache, sut, ctx := NewSUT(t)
	cache.EXPECT().GetExchangeRate(gomock.Any(), testConstant.TestIsoFrom, int32(0)).Return(nil, nil).Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.NotFound)
}

func TestExchangeRateService_GetAllExchangeRate_CacheError(t *testing.T) {
	t.Parallel()

	cache, sut, ctx := NewSUT(t)
	cache.EXPECT().
		GetAllExchangeRate(gomock.Eq(ctx), testConstant.TestIsoFrom).
		Return(nil, errors.TestError).
		Times(1)

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.ErrorIs(t, err, errors.TestError)
	assert.Nil(t, result)
}

func TestExchangeRateService_MapperToProto_NilRanking(t *testing.T) {
	t.Parallel()
	_, sut, _ := NewSUT(t)

	result := sut.MapperToProto(nil)

	assert.Nil(t, result)
}
