package exchange_rate

import (
	context "context"
	testing "testing"

	mocks "github.com/Suinar/Bank-exhange-rate-service/internal/mocks/external"
	exchangeRate "github.com/Suinar/Bank-exhange-rate-service/internal/services/exchange_rate"
	testConstant "github.com/Suinar/Bank-exhange-rate-service/internal/test"
	fixture "github.com/Suinar/Bank-exhange-rate-service/internal/test/fixture"
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

func TestExchangeRateService_MapperToProto(t *testing.T) {
	t.Parallel()

	_, sut, _ := NewSUT(t)

	result := sut.MapperToProto(fixture.NewRanking())

	require.NotNil(t, result)
	assert.Equal(t, fixture.NewRankingProto(), result)
}

func NewSUT(t *testing.T) (*mocks.MockIMonobankClient, *exchangeRate.ExchangeRateService, context.Context) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	client := mocks.NewMockIMonobankClient(ctrl)
	sut := exchangeRate.NewExchangeRateService(client)

	return client, sut, context.Background()
}
