package monobank

import (
	context "context"
	fmt "fmt"
	http "net/http"
	time "time"

	resty "github.com/go-resty/resty/v2"
	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/client/monobank"
	exchangeRateCache "github.com/kVinsom/Bank-exhange-rate-service/internal/repositories/cahce/exchange_rate"
	ranking "github.com/kVinsom/Bank-exhange-rate-service/pkg/core"
)

// MonobankClient communicates with the Monobank currency API.
type MonobankClient struct {
	client           *resty.Client
	baseURL          string
	currencyEndpoint string
	cache            exchangeRateCache.IExchangeRateCache
}

// NewMonobankClient creates a client configured for the Monobank currency API.
func NewMonobankClient(
	baseURL string,
	currencyEndpoint string,
	cache exchangeRateCache.IExchangeRateCache,
) *MonobankClient {
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
		cache:            cache,
	}
}

// GetAllExchangeRate fetches the latest currency snapshot from Monobank.
func (c *MonobankClient) GetAllExchangeRate(ctx context.Context) error {
	const operation = "get_all_exchange_rates"
	defer log.OperationStarted(operation)()
	var protoRates []MonoExchangeRate

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&protoRates).
		Get(c.baseURL + c.currencyEndpoint)

	if err != nil {
		log.DependencyFailed(operation, "monobank_api", err)
		return err
	}

	// Resty treats non-2xx responses as valid HTTP exchanges, so enforce the
	// upstream contract before consuming the response body.
	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf(
			"monobank unexpected status code: %d",
			resp.StatusCode(),
		)
		log.DependencyFailed(operation, "monobank_api", err)
		return err
	}

	rates := make([]ranking.Ranking, len(protoRates))

	for i := range protoRates {
		rates[i] = mapMonobankRate(&protoRates[i])
	}

	if err := c.cache.AddAllExchangeRate(ctx, rates); err != nil {
		log.DependencyFailed(operation, "exchange_rate_cache", err)
		return err
	}
	log.Result(operation, len(rates))

	return nil
}

// Close releases idle HTTP connections held by the client transport.
func (c *MonobankClient) Close() {
	defer log.OperationStarted("close")()
	c.client.GetClient().CloseIdleConnections()
}

// MonobankMapper converts the Monobank transport model into a domain ranking.
func (c *MonobankClient) MonobankMapper(monoModel *MonoExchangeRate) *ranking.Ranking {
	const operation = "monobank_rate_to_ranking"
	defer log.MappingStarted(operation)()
	if monoModel == nil {
		log.NilInput(operation)
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
