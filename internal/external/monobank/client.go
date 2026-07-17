package monobank

import (
	context "context"
	fmt "fmt"
	http "net/http"
	time "time"

	ranking "github.com/Suinar/Bank-exhange-rate-service/pkg/core"
	resty "github.com/go-resty/resty/v2"
)

// MonobankClient communicates with the Monobank currency API.
type MonobankClient struct {
	client           *resty.Client
	baseURL          string
	currencyEndpoint string
}

// NewMonobankClient creates a client configured for the Monobank currency API.
func NewMonobankClient(baseURL string, currencyEndpoint string) *MonobankClient {
	client := resty.New()

	// Keep the timeout at the transport boundary so a slow upstream cannot
	// consume request handlers indefinitely.
	client.
		SetTimeout(5*time.Second).
		SetHeader("Accept", "application/json")

	return &MonobankClient{
		client:           client,
		baseURL:          baseURL,
		currencyEndpoint: currencyEndpoint,
	}
}

// GetAllExchangeRate fetches the latest currency snapshot from Monobank.
func (c *MonobankClient) GetAllExchangeRate(ctx context.Context) ([]ranking.Ranking, error) {
	var protoRates []MonoExchangeRate

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&protoRates).
		Get(c.baseURL + c.currencyEndpoint)

	if err != nil {
		return nil, err
	}

	// Resty treats non-2xx responses as valid HTTP exchanges, so enforce the
	// upstream contract before consuming the response body.
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf(
			"monobank unexpected status code: %d",
			resp.StatusCode(),
		)
	}

	rates := make([]ranking.Ranking, len(protoRates))

	for i := range protoRates {
		rates[i] = mapMonobankRate(&protoRates[i])
	}

	return rates, nil
}

// MonobankMapper converts the Monobank transport model into a domain ranking.
func (c *MonobankClient) MonobankMapper(monoModel *MonoExchangeRate) *ranking.Ranking {
	if monoModel == nil {
		return nil
	}

	rate := mapMonobankRate(monoModel)

	return &rate
}

func mapMonobankRate(monoModel *MonoExchangeRate) ranking.Ranking {
	return ranking.Ranking{
		CurrencyIsoFrom: monoModel.CurrencyCodeA,
		CurrencyIsoTo:   monoModel.CurrencyCodeB,

		Date: time.Unix(monoModel.Date, 0),

		RateSell:  monoModel.RateSell,
		RateBuy:   monoModel.RateBuy,
		RateCross: monoModel.RateCross,
	}
}
