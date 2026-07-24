package fixture

import (
	kafka "github.com/segmentio/kafka-go"
)

const (
	RelativeRankingRequestTopic = "exchange-rate.get-relative-ranking.request"
	AllRankingRequestTopic      = "exchange-rate.get-all-ranking.request"
)

// NewKafkaTopics returns the request topics owned by the exchange-rate service.
func NewKafkaTopics() []string {
	return []string{
		RelativeRankingRequestTopic,
		AllRankingRequestTopic,
	}
}

// NewKafkaCreateTopicsRequest returns the expected Kafka topic creation request.
func NewKafkaCreateTopicsRequest() *kafka.CreateTopicsRequest {
	return &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{
			{
				Topic:             RelativeRankingRequestTopic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
			{
				Topic:             AllRankingRequestTopic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		},
	}
}

// NewKafkaCreateTopicsResponse returns a successful topic creation response.
func NewKafkaCreateTopicsResponse() *kafka.CreateTopicsResponse {
	return &kafka.CreateTopicsResponse{
		Errors: map[string]error{},
	}
}

// NewKafkaTopicAlreadyExistsResponse returns an idempotent creation response.
func NewKafkaTopicAlreadyExistsResponse() *kafka.CreateTopicsResponse {
	return &kafka.CreateTopicsResponse{
		Errors: map[string]error{
			RelativeRankingRequestTopic: kafka.TopicAlreadyExists,
		},
	}
}

// NewKafkaTopicErrorResponse returns a response containing a topic-level error.
func NewKafkaTopicErrorResponse(err error) *kafka.CreateTopicsResponse {
	return &kafka.CreateTopicsResponse{
		Errors: map[string]error{
			RelativeRankingRequestTopic: err,
		},
	}
}
