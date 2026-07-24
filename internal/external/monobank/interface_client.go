package monobank

import (
	context "context"
)

//go:generate mockgen -source=interface_client.go -destination=../../../internal/mocks/external/monobank.go -package=mocks

// IMonobankClient defines the exchange-rate data required from Monobank.
type IMonobankClient interface {
	// GetAllExchangeRate refreshes the exchange-rate snapshot in the cache.
	GetAllExchangeRate(ctx context.Context) error
}
