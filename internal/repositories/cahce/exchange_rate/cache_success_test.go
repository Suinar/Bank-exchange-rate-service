package exchange_rate

import (
	context "context"
	fmt "fmt"
	os "os"
	sync "sync"
	testing "testing"

	testConstant "github.com/kVinsom/Bank-exhange-rate-service/internal/test"
	fixture "github.com/kVinsom/Bank-exhange-rate-service/internal/test/fixture"
	core "github.com/kVinsom/Bank-exhange-rate-service/pkg/core"
	redis "github.com/redis/go-redis/v9"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

var integrationRedisOptions *redis.Options

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := redisContainer.Run(ctx, testConstant.TestRedisImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis integration tests failed: start container: %v\n", err)
		os.Exit(1)
	}

	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis integration tests failed: get connection string: %v\n", err)
		_ = testcontainers.TerminateContainer(container)
		os.Exit(1)
	}

	integrationRedisOptions, err = redis.ParseURL(connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis integration tests failed: parse connection string: %v\n", err)
		_ = testcontainers.TerminateContainer(container)
		os.Exit(1)
	}

	exitCode := m.Run()
	if err := testcontainers.TerminateContainer(container); err != nil {
		fmt.Fprintf(os.Stderr, "terminate Redis integration container: %v\n", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}

func TestExchangeRateCache_AddAllExchangeRate_Success(t *testing.T) {
	_, sut, ctx := NewSUT(t)
	expected := fixture.NewRanking()

	err := sut.AddAllExchangeRate(ctx, fixture.NewRankings())

	require.NoError(t, err)

	result, err := sut.GetExchangeRate(
		ctx,
		expected.CurrencyIsoFrom,
		expected.CurrencyIsoTo,
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestExchangeRateCache_GetAllExchangeRate_Success(t *testing.T) {
	_, sut, ctx := NewSUT(t)
	rates := fixture.NewCacheRankings()
	expected := rates[:2]

	require.NoError(t, sut.AddAllExchangeRate(ctx, rates))

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.NoError(t, err)
	require.Len(t, result, len(expected))
	assert.ElementsMatch(t, expected, result)
}

func TestExchangeRateCache_GetMissingRates(t *testing.T) {
	_, sut, ctx := NewSUT(t)

	rate, err := sut.GetExchangeRate(ctx, testConstant.TestIsoFrom, testConstant.TestIsoTo)
	require.NoError(t, err)
	assert.Nil(t, rate)

	rates, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)
	require.NoError(t, err)
	assert.Empty(t, rates)
	assert.NotNil(t, rates)
}

func TestExchangeRateCache_GetAllExchangeRate_RemovesStaleIndexMembers(t *testing.T) {
	client, sut, ctx := NewSUT(t)
	fromKey := sut.GetCurrencyFromKey(testConstant.TestIsoFrom)
	staleRateKey := sut.GetPrimaryKey(testConstant.TestIsoFrom, testConstant.TestIsoTo)
	require.NoError(t, client.SAdd(ctx, fromKey, staleRateKey).Err())

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)

	require.NoError(t, err)
	assert.Empty(t, result)

	indexedKeys, err := client.SMembers(ctx, fromKey).Result()
	require.NoError(t, err)
	assert.Empty(t, indexedKeys)
}

func TestExchangeRateCache_AddAllExchangeRate_ReplacesExistingRates(t *testing.T) {
	_, sut, ctx := NewSUT(t)
	oldRate := fixture.NewRanking()
	expected := fixture.NewReplacementRanking()

	require.NoError(t, sut.AddAllExchangeRate(ctx, fixture.NewRankings()))
	require.NoError(t, sut.AddAllExchangeRate(ctx, fixture.NewReplacementRankings()))

	removed, err := sut.GetExchangeRate(ctx, oldRate.CurrencyIsoFrom, oldRate.CurrencyIsoTo)
	require.NoError(t, err)
	assert.Nil(t, removed)

	result, err := sut.GetExchangeRate(ctx, expected.CurrencyIsoFrom, expected.CurrencyIsoTo)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expected, result)
}

func TestExchangeRateCache_AddAllExchangeRate_EmptyRatesClearsCache(t *testing.T) {
	_, sut, ctx := NewSUT(t)
	existing := fixture.NewRanking()

	require.NoError(t, sut.AddAllExchangeRate(ctx, fixture.NewRankings()))

	err := sut.AddAllExchangeRate(ctx, nil)

	require.NoError(t, err)

	result, err := sut.GetExchangeRate(
		ctx,
		existing.CurrencyIsoFrom,
		existing.CurrencyIsoTo,
	)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExchangeRateCache_AddAllExchangeRate_ConcurrentReplacementsRemainAtomic(t *testing.T) {
	_, sut, ctx := NewSUT(t)
	snapshots := [][]core.Ranking{
		fixture.NewRankings(),
		fixture.NewReplacementRankings(),
	}

	errs := make(chan error, len(snapshots))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(snapshots))
	for i := range snapshots {
		snapshot := snapshots[i]
		go func() {
			defer waitGroup.Done()
			errs <- sut.AddAllExchangeRate(ctx, snapshot)
		}()
	}

	waitGroup.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	result, err := sut.GetAllExchangeRate(ctx, testConstant.TestIsoFrom)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Contains(t, []int32{testConstant.TestIsoTo, testConstant.TestIsoAlternative}, result[0].CurrencyIsoTo)
}

func NewSUT(t *testing.T) (*redis.Client, *ExchangeRateCache, context.Context) {
	t.Helper()

	options := *integrationRedisOptions
	client := redis.NewClient(&options)
	ctx := context.Background()
	require.NoError(t, client.FlushDB(ctx).Err())

	t.Cleanup(func() {
		_ = client.FlushDB(context.Background()).Err()
		_ = client.Close()
	})

	return client, NewExchangeRateCache(client), ctx
}

func NewClosedSUT(t *testing.T) (*ExchangeRateCache, context.Context) {
	t.Helper()

	client, sut, ctx := NewSUT(t)
	require.NoError(t, client.Close())

	return sut, ctx
}
