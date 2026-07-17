package exchange_rate_test

import (
	testing "testing"

	testConstant "github.com/Suinar/Bank-exhange-rate-service/internal/test"
	fixture "github.com/Suinar/Bank-exhange-rate-service/internal/test/fixture"
	errors "github.com/Suinar/Bank-repository-service/pkg"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestExchangeRateService_GetExchangeRate_ClientError(t *testing.T) {
	t.Parallel()
	client, sut, ctx := NewSUT(t)
	client.EXPECT().GetAllExchangeRate(gomock.Any()).Return(nil, errors.TestError).Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.ErrorIs(t, err, errors.TestError)
	assert.Nil(t, result)
}

func TestExchangeRateService_GetExchangeRate_NotFound(t *testing.T) {
	t.Parallel()
	client, sut, ctx := NewSUT(t)
	client.EXPECT().GetAllExchangeRate(gomock.Any()).Return(fixture.NewRankings(), nil).Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, 0)

	require.ErrorIs(t, err, errors.NotFound)
	assert.Nil(t, result)
}

func TestExchangeRateService_GetAllExchangeRate_ClientError(t *testing.T) {
	t.Parallel()
	client, sut, ctx := NewSUT(t)
	client.EXPECT().GetAllExchangeRate(gomock.Any()).Return(nil, errors.TestError).Times(1)

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.ErrorIs(t, err, errors.TestError)
	assert.Nil(t, result)
}

func TestExchangeRateService_MapperToProto_NilRanking(t *testing.T) {
	t.Parallel()
	_, sut, _ := NewSUT(t)
	assert.Nil(t, sut.MapperToProto(nil))
}
