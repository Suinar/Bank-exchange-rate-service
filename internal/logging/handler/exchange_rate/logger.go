package exchange_rate

import log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/operation"

func RequestStarted(operation string) func() {
	return log.Started("exchange_rate", "gRPC handler", operation)
}
func RequestError(operation string, err error) {
	log.Failed("exchange_rate", "gRPC handler", operation, "exchange_rate_service", err)
}
func ValidationFailed(operation, reason string) {
	log.ValidationFailed("exchange_rate", "gRPC handler", operation, reason)
}
