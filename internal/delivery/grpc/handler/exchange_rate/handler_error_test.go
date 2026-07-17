package exchange_rate_test

import (
	testing "testing"

	fixture "github.com/Suinar/Bank-exhange-rate-service/internal/test/fixture"
	errors "github.com/Suinar/Bank-repository-service/pkg"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestExchangeRateHandler_GetExchangeRate_ServiceError(t *testing.T) {
	t.Parallel()
	service, sut, ctx := NewSUT(t)
	req := fixture.NewGetRelativeRankingRequestProto()
	service.EXPECT().GetExchangeRate(gomock.Eq(ctx), gomock.Eq(req.CurrencyIsoFrom), gomock.Eq(req.CurrencyIsoTo)).Return(nil, errors.TestError).Times(1)

	result, err := sut.GetExchangeRate(ctx, req)

	require.ErrorIs(t, err, errors.TestError)
	assert.Nil(t, result)
}

func TestExchangeRateHandler_GetExchangeRate_SameCurrencies(t *testing.T) {
	t.Parallel()
	_, sut, ctx := NewSUT(t)
	req := fixture.NewGetRelativeRankingRequestProto()
	req.CurrencyIsoTo = req.CurrencyIsoFrom
	result, err := sut.GetExchangeRate(ctx, req)
	require.ErrorIs(t, err, errors.BadRequest)
	assert.Nil(t, result)
}

func TestExchangeRateHandler_GetExchangeRate_NilRequest(t *testing.T) {
	t.Parallel()
	_, sut, ctx := NewSUT(t)
	result, err := sut.GetExchangeRate(ctx, nil)
	require.ErrorIs(t, err, errors.BadRequest)
	assert.Nil(t, result)
}

func TestExchangeRateHandler_GetAllExchangeRate_ServiceError(t *testing.T) {
	t.Parallel()
	service, sut, ctx := NewSUT(t)
	req := fixture.NewGetAllRankingRequestProto()
	service.EXPECT().GetAllExchangeRate(gomock.Eq(ctx), gomock.Eq(req.CurrencyIsoFrom)).Return(nil, errors.TestError).Times(1)
	result, err := sut.GetAllExchangeRate(ctx, req)
	require.ErrorIs(t, err, errors.TestError)
	assert.Nil(t, result)
}

func TestExchangeRateHandler_GetAllExchangeRate_NilRequest(t *testing.T) {
	t.Parallel()
	_, sut, ctx := NewSUT(t)
	result, err := sut.GetAllExchangeRate(ctx, nil)
	require.ErrorIs(t, err, errors.BadRequest)
	assert.Nil(t, result)
}
