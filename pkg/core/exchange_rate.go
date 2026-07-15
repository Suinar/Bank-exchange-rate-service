package core

import "time"

type Ranking struct {
	CurrencyIsoFrom int32 `json:"currency_iso_from"`
	CurrencyIsoTo   int32 `json:"currency_iso_to"`

	Date time.Time `json:"date"`

	RateSell  float32 `json:"rate_sell"`
	RateBuy   float32 `json:"rate_buy"`
	RateCross float32 `json:"rate_cross"`
}
