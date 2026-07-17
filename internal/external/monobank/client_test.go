package monobank

package monobank_test

import (
	context "context"
	http "net/http"
	httptest "net/http/httptest"
	testing "testing"

	monobank "github.com/Suinar/Bank-exhange-rate-service/internal/external/monobank"
	fixture "github.com/Suinar/Bank-exhange-rate-service/internal/test/fixture"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestMonobankClient_GetAllExchangeRate_Success(t *testing.T) {
	t.Parallel()

	mux, _, sut, ctx := NewSUT(t)
	mux.HandleFunc("/bank/currency", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/bank/currency", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`[{"currencyCodeA":840,"currencyCodeB":978,"date":1710000000,"rateSell":3,"rateBuy":2,"rateCross":2.5}]`))
		require.NoError(t, err)
	})

	result, err := sut.GetAllExchangeRate(ctx)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, fixture.NewRanking(), &result[0])
}

func TestMonobankClient_MonobankMapper(t *testing.T) {
	t.Parallel()

	_, _, sut, _ := NewSUT(t)
	model := fixture.NewMonoExchangeRate()

	result := sut.MonobankMapper(model)

	require.NotNil(t, result)
	assert.Equal(t, fixture.NewRanking(), result)
}

func NewSUT(t *testing.T) (*http.ServeMux, *httptest.Server, *monobank.MonobankClient, context.Context) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	sut := monobank.NewMonobankClient(server.URL, "/bank/currency")

	return mux, server, sut, context.Background()
}