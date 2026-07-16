package monobank

// MonoExchangeRate represents an exchange-rate item returned by the Monobank API.
type MonoExchangeRate struct {
	CurrencyCodeA int32 `json:"currencyCodeA"`
	CurrencyCodeB int32 `json:"currencyCodeB"`

	Date int64 `json:"date"`

	RateSell  float32 `json:"rateSell"`
	RateBuy   float32 `json:"rateBuy"`
	RateCross float32 `json:"rateCross"`
}
