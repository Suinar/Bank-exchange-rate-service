package exchange_rate

import log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/operation"

func OperationStarted(operation string) func() {
	return log.Started("exchange_rate", "service", operation)
}
func MappingStarted(operation string) func() {
	return log.Started("exchange_rate", "mapper", operation)
}
func DependencyFailed(operation, dependency string, err error) {
	log.Failed("exchange_rate", "service", operation, dependency, err)
}
func NilInput(operation string) {
	log.NilInput("exchange_rate", "mapper", operation)
}
func Result(operation string, count int) {
	log.Result("exchange_rate", "service", operation, count)
}
