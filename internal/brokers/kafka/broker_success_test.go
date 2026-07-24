package kafka

import (
	context "context"
	testing "testing"

	mocks "github.com/kVinsom/Bank-exhange-rate-service/internal/mocks/brokers"
	"github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	kafkaGo "github.com/segmentio/kafka-go"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAdmin_EnsureTopics_Success(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)
	expected := fixture.NewKafkaCreateTopicsRequest()

	client.EXPECT().
		CreateTopics(gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			_ context.Context,
			request *kafkaGo.CreateTopicsRequest,
		) (*kafkaGo.CreateTopicsResponse, error) {
			assert.Equal(t, expected, request)
			return fixture.NewKafkaCreateTopicsResponse(), nil
		}).
		Times(1)

	err := sut.EnsureTopics(ctx, fixture.NewKafkaTopics()...)

	require.NoError(t, err)
}

func TestAdmin_EnsureTopics_TopicAlreadyExists(t *testing.T) {
	t.Parallel()

	client, sut, ctx := NewSUT(t)

	client.EXPECT().
		CreateTopics(gomock.Any(), gomock.Any()).
		Return(fixture.NewKafkaTopicAlreadyExistsResponse(), nil).
		Times(1)

	err := sut.EnsureTopics(ctx, fixture.NewKafkaTopics()...)

	require.NoError(t, err)
}

func NewSUT(t *testing.T) (*mocks.MockITopicCreator, *Admin, context.Context) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	client := mocks.NewMockITopicCreator(ctrl)
	sut := &Admin{client: client}

	return client, sut, context.Background()
}
