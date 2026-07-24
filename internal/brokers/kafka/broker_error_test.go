package kafka

import (
	errors "errors"
	testing "testing"

	"github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAdmin_EnsureTopics_ClientError(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)
	expectedErr := errors.New("broker unavailable")

	client.EXPECT().
		CreateTopics(gomock.Any(), gomock.Any()).
		Return(nil, expectedErr).
		Times(1)

	err := sut.EnsureTopics(ctx, fixture.NewKafkaTopics()...)

	require.ErrorIs(t, err, expectedErr)
	require.ErrorContains(t, err, "create Kafka topics")
}

func TestAdmin_EnsureTopics_TopicError(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)
	expectedErr := errors.New("invalid topic configuration")

	client.EXPECT().
		CreateTopics(gomock.Any(), gomock.Any()).
		Return(fixture.NewKafkaTopicErrorResponse(expectedErr), nil).
		Times(1)

	err := sut.EnsureTopics(ctx, fixture.NewKafkaTopics()...)

	require.ErrorIs(t, err, expectedErr)
	require.ErrorContains(t, err, "create Kafka topic")
}

func TestAdmin_EnsureTopics_MissingTopics(t *testing.T) {
	t.Parallel()

	_, sut, ctx := NewSUT(t)

	err := sut.EnsureTopics(ctx)

	require.EqualError(t, err, "at least one Kafka topic is required")
}

func TestAdmin_EnsureTopics_EmptyTopicName(t *testing.T) {
	t.Parallel()

	_, sut, ctx := NewSUT(t)

	err := sut.EnsureTopics(ctx, "")

	require.EqualError(t, err, "Kafka topic name cannot be empty")
}
