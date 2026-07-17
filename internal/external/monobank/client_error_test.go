package monobank_test

import (
	http "net/http"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestMonobankClient_GetAllExchangeRate_UnexpectedStatusCode(t *testing.T) {
	t.Parallel()
	mux, _, sut, ctx := NewSUT(t)
	mux.HandleFunc("/bank/currency", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	result, err := sut.GetAllExchangeRate(ctx)

	require.EqualError(t, err, "monobank unexpected status code: 500")
	assert.Nil(t, result)
}

func TestMonobankClient_GetAllExchangeRate_RequestError(t *testing.T) {
	t.Parallel()
	_, server, sut, ctx := NewSUT(t)
	server.Close()

	result, err := sut.GetAllExchangeRate(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestMonobankClient_MonobankMapper_NilModel(t *testing.T) {
	t.Parallel()
	_, _, sut, _ := NewSUT(t)
	assert.Nil(t, sut.MonobankMapper(nil))
}
