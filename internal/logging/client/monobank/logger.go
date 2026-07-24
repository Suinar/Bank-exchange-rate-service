package monobank

import log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/operation"

func OperationStarted(operation string) func() {
	return log.Started("monobank", "client", operation)
}
func DependencyFailed(operation, dependency string, err error) {
	log.Failed("monobank", "client", operation, dependency, err)
}
func MappingStarted(operation string) func() {
	return log.Started("monobank", "mapper", operation)
}
func NilInput(operation string) {
	log.NilInput("monobank", "mapper", operation)
}
func Result(operation string, count int) {
	log.Result("monobank", "client", operation, count)
}
