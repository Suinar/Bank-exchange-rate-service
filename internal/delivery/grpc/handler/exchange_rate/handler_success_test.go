package exchange_rate_test

import (
	context "context"
	testing "testing"

	handler "github.com/kVinsom/Bank-exhange-rate-service/internal/delivery/grpc/handler/exchange_rate"
	mocks "github.com/kVinsom/Bank-exhange-rate-service/internal/mocks/services"
	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
	fixture "github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestExchangeRateHandler_GetExchangeRate_Success(t *testing.T) {
	t.Parallel()

	service, sut, ctx := NewSUT(t)
	req := fixture.NewGetRelativeRankingRequestProto()
	expected := fixture.NewRankingProto()

	service.EXPECT().
		GetExchangeRate(gomock.Eq(ctx), req.CurrencyIsoFrom, req.CurrencyIsoTo).
		Return(expected, nil).
		Times(1)

	result, err := sut.GetExchangeRate(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExchangeRateHandler_GetAllExchangeRate_Success(t *testing.T) {
	t.Parallel()

	service, sut, ctx := NewSUT(t)
	req := fixture.NewGetAllRankingRequestProto()
	expected := fixture.NewRankingListProto()

	service.EXPECT().
		GetAllExchangeRate(gomock.Eq(ctx), req.CurrencyIsoFrom).
		Return(expected, nil).
		Times(1)

	result, err := sut.GetAllExchangeRate(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExchangeRateHandler_GetAllExchangeRate_PassesCurrencyToService(t *testing.T) {
	t.Parallel()

	service, sut, ctx := NewSUT(t)
	req := fixture.NewGetAllRankingRequestProto()
	req.CurrencyIsoFrom = testConstant.TestIsoTo
	expected := fixture.NewRankingListProto()

	service.EXPECT().
		GetAllExchangeRate(gomock.Eq(ctx), testConstant.TestIsoTo).
		Return(expected, nil).
		Times(1)

	result, err := sut.GetAllExchangeRate(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func NewSUT(t *testing.T) (*mocks.MockIExchangeRateService, handler.IExchangeRateHandler, context.Context) {
	t.Helper()

	ctrl := gomock.NewController(t)
	service := mocks.NewMockIExchangeRateService(ctrl)
	sut := handler.NewExchangeRateHandler(service)

	return service, sut, context.Background()
}
