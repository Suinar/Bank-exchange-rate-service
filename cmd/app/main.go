package main

import (
	"context"

	"github.com/Suinar/Bank-exhange-rate-service/internal/external/monobank"
)

func main() {
	client := monobank.NewMonobankClient("https://api.monobank.ua", "/bank/currency")

	client.GetAllExchangeRate(context.Background())
}
