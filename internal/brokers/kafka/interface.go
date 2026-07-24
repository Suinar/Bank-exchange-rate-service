package kafka

import (
	context "context"

	kafkaGo "github.com/segmentio/kafka-go"
)

//go:generate mockgen -source=interface.go -destination=../../mocks/brokers/kafka.go -package=mocks

// ITopicCreator defines the Kafka operation required to create service topics.
type ITopicCreator interface {
	CreateTopics(context.Context, *kafkaGo.CreateTopicsRequest) (*kafkaGo.CreateTopicsResponse, error)
}
