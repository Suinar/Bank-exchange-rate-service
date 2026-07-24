package kafka

import (
	stdlog "log"
	time "time"

	operationlog "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/operation"
)

func OperationStarted(operation string) func() {
	return operationlog.Started("kafka", "broker", operation)
}
func DependencyFailed(operation string, err error) {
	operationlog.Failed("kafka", "broker", operation, "kafka", err)
}
func TopicsPreparing(topics []string) {
	stdlog.Printf("kafka broker: preparing topics topics=%v", topics)
}
func TopicsReady(topics []string, attempt int) {
	stdlog.Printf("kafka broker: topics ready topics=%v attempt=%d", topics, attempt)
}
func Retry(attempt, maxAttempts int, delay time.Duration, err error) {
	stdlog.Printf(
		"kafka broker: connection retry attempt=%d max_attempts=%d retry_delay=%s error=%q",
		attempt, maxAttempts, delay, err,
	)
}
func TopicsFailed(attempts int, err error) {
	stdlog.Printf("kafka broker: topic preparation failed attempts=%d error=%q", attempts, err)
}
