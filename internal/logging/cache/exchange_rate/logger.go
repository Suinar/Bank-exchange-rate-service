package exchange_rate

import log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/operation"

func OperationStarted(operation string) func() {
	return log.Started("exchange_rate", "cache", operation)
}
func DependencyFailed(operation string, err error) {
	log.Failed("exchange_rate", "cache", operation, "redis", err)
}
func MappingStarted(operation string) func() {
	return log.Started("exchange_rate", "cache mapper", operation)
}
func MappingFailed(operation string, err error) {
	log.Failed("exchange_rate", "cache mapper", operation, "cached_data", err)
}
func Result(operation string, count int) {
	log.Result("exchange_rate", "cache", operation, count)
}
