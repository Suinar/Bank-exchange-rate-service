package core

import "time"

type Ranking struct {
	CurrencyIDFrom int64 `json:"currency_id_from"`
	CurrencyIDTo   int64 `json:"currency_id_to"`

	Date time.Time `json:"date"`

	RateSell  int32 `json:"rate_sell"`
	RateBuy   int32 `json:"rate_buy"`
	RateCross int32 `json:"rate_cross"`
}

