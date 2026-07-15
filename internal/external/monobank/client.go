package monobank

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ranking "github.com/Suinar/Bank-exhange-rate-service/pkg/core"
	resty "github.com/go-resty/resty/v2"
)

type MonobankClient struct {
	client           *resty.Client
	baseURL          string
	currencyEndpoint string
}

func NewMonobankClient(baseURL string, currencyEndpoint string) *MonobankClient {
	client := resty.New()

	client.
		SetTimeout(5*time.Second).
		SetHeader("Accept", "application/json")

	return &MonobankClient{
		client:           client,
		baseURL:          baseURL,
		currencyEndpoint: currencyEndpoint,
	}
}

func (c *MonobankClient) GetAllExchangeRate(ctx context.Context) ([]ranking.Ranking, error) {
	var protoRates []MonoExchangeRate

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&protoRates).
		Get(c.baseURL + c.currencyEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(
			"monobank unexpected status code: %d",
			resp.StatusCode(),
		)
	}

	rates := make([]ranking.Ranking, len(protoRates))

	for i := range protoRates {
		rates[i] = *c.MonobankMapper(&protoRates[i])
	}

	return rates, nil
}

func (c *MonobankClient) MonobankMapper(monoModel *MonoExchangeRate) *ranking.Ranking {
	return &ranking.Ranking{
		CurrencyIsoFrom: monoModel.CurrencyCodeA,
		CurrencyIsoTo:   monoModel.CurrencyCodeB,

		Date: time.Unix(monoModel.Date, 0),

		RateSell:  monoModel.RateSell,
		RateBuy:   monoModel.RateBuy,
		RateCross: monoModel.RateCross,
	}
}
