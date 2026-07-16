package exchange_rate_test

import (
	context "context"
	testing "testing"

	mocks "github.com/Suinar/Bank-exhange-rate-service/internal/mocks/external"
	exchangeRate "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"
	testConstant "github.com/Suinar/Bank-exhange-rate-service/internal/test"
	fixture "github.com/Suinar/Bank-exhange-rate-service/internal/test/fixture"
	errors "github.com/Suinar/Bank-repository-service/pkg"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestExchangeRateService_GetExchangeRate_Success(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)
	expected := fixture.NewRankingProto()

	client.EXPECT().
		GetAllExchangeRate(gomock.Any()).
		Return(fixture.NewRankings(), nil).
		Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestExchangeRateService_GetExchangeRate_ClientError(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)

	client.EXPECT().
		GetAllExchangeRate(gomock.Any()).
		Return(nil, errors.TestError).
		Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.TestError)
}

func TestExchangeRateService_GetExchangeRate_NotFound(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)

	client.EXPECT().
		GetAllExchangeRate(gomock.Any()).
		Return(fixture.NewRankings(), nil).
		Times(1)

	result, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.InternalServerError)
}

func TestExchangeRateService_GetAllExchangeRate_Success(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)
	expected := fixture.NewRankingListProto()

	client.EXPECT().
		GetAllExchangeRate(gomock.Any()).
		Return(fixture.NewRankings(), nil).
		Times(1)

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestExchangeRateService_GetAllExchangeRate_ClientError(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)

	client.EXPECT().
		GetAllExchangeRate(gomock.Any()).
		Return(nil, errors.TestError).
		Times(1)

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.TestError)
}

func TestExchangeRateService_MapperToProto(t *testing.T) {
	t.Parallel()

	_, sut, _ := NewSUT(t)

	result := sut.MapperToProto(fixture.NewRanking())

	require.NotNil(t, result)
	assert.Equal(t, fixture.NewRankingProto(), result)
}

func TestExchangeRateService_MapperToProto_NilRanking(t *testing.T) {
	t.Parallel()

	_, sut, _ := NewSUT(t)

	result := sut.MapperToProto(nil)

	assert.Nil(t, result)
}

func NewSUT(t *testing.T) (*mocks.MockIMonobankClient, *exchangeRate.ExchangeRateService, context.Context) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	client := mocks.NewMockIMonobankClient(ctrl)
	sut := exchangeRate.NewExchangeRateService(client)

	return client, sut, context.Background()
}
