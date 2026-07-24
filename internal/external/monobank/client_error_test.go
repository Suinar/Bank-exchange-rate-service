package monobank_test

import (
	"net/http"
	"testing"

	fixture "github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	errors "github.com/kVinsom/Bank-repository-service/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestMonobankClient_GetAllExchangeRate_UnexpectedStatusCode(t *testing.T) {

	t.Parallel()
	mux, _, _, sut, ctx := NewSUT(t)
	mux.HandleFunc("/bank/currency", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := sut.GetAllExchangeRate(ctx)

	require.EqualError(t, err, "monobank unexpected status code: 500")
}

func TestMonobankClient_GetAllExchangeRate_RequestError(t *testing.T) {
	t.Parallel()
	_, server, _, sut, ctx := NewSUT(t)
	server.Close()

	err := sut.GetAllExchangeRate(ctx)

	require.Error(t, err)
}

func TestMonobankClient_MonobankMapper_NilModel(t *testing.T) {
	t.Parallel()
	_, _, _, sut, _ := NewSUT(t)

	result := sut.MonobankMapper(nil)

	assert.Nil(t, result)
}

func TestMonobankClient_GetAllExchangeRate_CacheError(t *testing.T) {
	t.Parallel()

	mux, _, cache, sut, ctx := NewSUT(t)
	mux.HandleFunc("/bank/currency", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(fixture.NewMonoExchangeRateJSON())
		require.NoError(t, err)
	})
	cache.EXPECT().
		AddAllExchangeRate(gomock.Eq(ctx), gomock.Eq(fixture.NewRankings())).
		Return(errors.TestError).
		Times(1)

	err := sut.GetAllExchangeRate(ctx)

	require.ErrorIs(t, err, errors.TestError)
}
