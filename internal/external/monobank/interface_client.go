package monobank

import (
	"context"

	ranking "github.com/Suinar/Bank-exhange-rate-service/pkg/core"
)

//go:generate mockgen -source=interface_client.go -destination=../../../internal/mocks/external/monobank.go -package=mocks

type IMonobankClient interface {
	GetAllExchangeRate(ctx context.Context) ([]ranking.Ranking, error)
}
