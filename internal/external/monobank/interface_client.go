package monobank

import (
	context "context"

	ranking "github.com/Suinar/Bank-exhange-rate-service/pkg/core"
)

//go:generate mockgen -source=interface_client.go -destination=../../../internal/mocks/external/monobank.go -package=mocks

// IMonobankClient defines the exchange-rate data required from Monobank.
type IMonobankClient interface {
	// GetAllExchangeRate returns the upstream snapshot as domain values so the
	// service layer does not depend on Monobank's transport model.
	GetAllExchangeRate(ctx context.Context) ([]ranking.Ranking, error)
}
