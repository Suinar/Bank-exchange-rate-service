package external

import monobank "github.com/Suinar/Bank-exhange-rate-service/internal/external/monobank"

type Clients struct {
	MonoClient monobank.MonobankClient
}

func InitClients(monoClient monobank.MonobankClient) *Clients {
	return &Clients{
		MonoClient: monoClient,
	}
}
