package kafka

import (
	context "context"
	fmt "fmt"
	time "time"

	configs "github.com/kVinsom/Bank-exhange-rate-service/internal/configs"
	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/broker/kafka"
)

// EnsureKafkaTopics connects to Kafka and prepares the topics owned by this service.
func EnsureKafkaTopics(ctx context.Context, cfg *configs.Config) error {
	const (
		attempts   = 10
		retryDelay = 2 * time.Second
	)

	admin := NewAdmin(cfg.Kafka.Brokers)
	topics := []string{
		cfg.Kafka.GetRelativeRankingRequestTopic,
		cfg.Kafka.GetAllRankingRequestTopic,
	}

	var err error
	log.TopicsPreparing(topics)
	for attempt := 1; attempt <= attempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = admin.EnsureTopics(attemptCtx, topics...)
		cancel()
		if err == nil {
			log.TopicsReady(topics, attempt)
			return nil
		}

		if attempt < attempts {
			log.Retry(attempt, attempts, retryDelay, err)
			timer := time.NewTimer(retryDelay)
			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return ctx.Err()
			case <-timer.C:
			}
		}
	}

	log.TopicsFailed(attempts, err)
	return fmt.Errorf("prepare Kafka topics after %d attempts: %w", attempts, err)
}
