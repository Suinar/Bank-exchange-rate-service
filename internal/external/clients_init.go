package external

import monobank "github.com/kVinsom/Bank-exhange-rate-service/internal/external/monobank"

// Clients groups integrations with external systems.
type Clients struct {
	MonoClient *monobank.MonobankClient
}

// InitClients assembles the external clients used by the application.
func InitClients(monoClient *monobank.MonobankClient) *Clients {
	return &Clients{
		MonoClient: monoClient,
	}
}
