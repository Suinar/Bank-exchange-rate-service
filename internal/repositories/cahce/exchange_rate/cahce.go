package exchange_rate

import (
	context "context"
	strconv "strconv"
	time "time"

	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/cache/exchange_rate"
	ranking "github.com/kVinsom/Bank-exhange-rate-service/pkg/core"
	errors "github.com/kVinsom/Bank-repository-service/pkg"
	redis "github.com/redis/go-redis/v9"
)

const (
	exchangeRatesSetKey        = "exchange_rates"
	exchangeRateFromIndexesKey = "exchange_rate_from_indexes"
	exchangeRateKeyPrefix      = "exchange_rate:"
	exchangeRatePairPrefix     = exchangeRateKeyPrefix + "pair:"
	exchangeRateFromPrefix     = exchangeRateKeyPrefix + "from:"
)

var replaceExchangeRatesScript = redis.NewScript(`
local old_rate_keys = redis.call("SMEMBERS", KEYS[1])
local old_from_keys = redis.call("SMEMBERS", KEYS[2])

if #old_rate_keys > 0 then
	redis.call("UNLINK", unpack(old_rate_keys))
end
if #old_from_keys > 0 then
	redis.call("UNLINK", unpack(old_from_keys))
end

redis.call("UNLINK", KEYS[1], KEYS[2])

for i = 3, #ARGV, 6 do
	local from_iso = ARGV[i]
	local to_iso = ARGV[i + 1]
	local rate_key = ARGV[1] .. from_iso .. ":" .. to_iso
	local from_key = ARGV[2] .. from_iso

	redis.call(
		"HSET",
		rate_key,
		"currency_iso_from", from_iso,
		"currency_iso_to", to_iso,
		"date", ARGV[i + 2],
		"rate_sell", ARGV[i + 3],
		"rate_buy", ARGV[i + 4],
		"rate_cross", ARGV[i + 5]
	)
	redis.call("SADD", from_key, rate_key)
	redis.call("SADD", KEYS[1], rate_key)
	redis.call("SADD", KEYS[2], from_key)
end

return #ARGV - 2
`)

// ExchangeRateCache stores exchange rates and their source-currency indexes in Redis.
type ExchangeRateCache struct {
	rdb *redis.Client
}

var _ IExchangeRateCache = (*ExchangeRateCache)(nil)

// NewExchangeRateCache creates a ready-to-use exchange-rate cache.
func NewExchangeRateCache(rdb *redis.Client) *ExchangeRateCache {
	return &ExchangeRateCache{rdb: rdb}
}

// AddAllExchangeRate replaces the cached exchange-rate collection.
func (c *ExchangeRateCache) AddAllExchangeRate(ctx context.Context, rates []ranking.Ranking) error {
	const operation = "add_all_exchange_rates"
	defer log.OperationStarted(operation)()
	args := make([]any, 2, 2+(len(rates)*6))
	args[0] = exchangeRatePairPrefix
	args[1] = exchangeRateFromPrefix
	for i := range rates {
		rate := &rates[i]
		args = append(
			args,
			rate.CurrencyIsoFrom,
			rate.CurrencyIsoTo,
			rate.Date.UnixNano(),
			rate.RateSell,
			rate.RateBuy,
			rate.RateCross,
		)
	}

	if err := replaceExchangeRatesScript.Run(
		ctx,
		c.rdb,
		[]string{exchangeRatesSetKey, exchangeRateFromIndexesKey},
		args...,
	).Err(); err != nil {
		log.DependencyFailed(operation, err)
		return errors.CacheSetError
	}
	log.Result(operation, len(rates))

	return nil
}

// GetExchangeRate returns the cached rate for an exact currency pair.
func (c *ExchangeRateCache) GetExchangeRate(
	ctx context.Context,
	currencyIsoFrom int32,
	currencyIsoTo int32,
) (*ranking.Ranking, error) {
	const operation = "get_exchange_rate"
	defer log.OperationStarted(operation)()
	data, err := c.rdb.HGetAll(ctx, c.GetPrimaryKey(currencyIsoFrom, currencyIsoTo)).Result()
	if err != nil {
		log.DependencyFailed(operation, err)
		return nil, errors.CacheGetError
	}

	if len(data) == 0 {
		return nil, nil
	}

	rate, err := c.MapToRanking(data)
	if err != nil {
		log.MappingFailed(operation, err)
		return nil, errors.CacheMapError
	}
	log.Result(operation, 1)

	return &rate, nil
}

// GetAllExchangeRate returns all cached rates for a source currency.
func (c *ExchangeRateCache) GetAllExchangeRate(
	ctx context.Context,
	currencyIsoFrom int32,
) ([]ranking.Ranking, error) {
	const operation = "get_all_exchange_rates"
	defer log.OperationStarted(operation)()
	keys, err := c.rdb.SMembers(ctx, c.GetCurrencyFromKey(currencyIsoFrom)).Result()
	if err != nil {
		log.DependencyFailed(operation, err)
		return nil, errors.CacheGetError
	}

	if len(keys) == 0 {
		return []ranking.Ranking{}, nil
	}

	pipe := c.rdb.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipe.HGetAll(ctx, key))
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.DependencyFailed(operation, err)
		return nil, errors.CacheGetError
	}

	rates := make([]ranking.Ranking, 0, len(cmds))
	staleKeys := make([]any, 0)
	for i, cmd := range cmds {
		data := cmd.Val()
		if len(data) == 0 {
			staleKeys = append(staleKeys, keys[i])
			continue
		}

		rate, err := c.MapToRanking(data)
		if err != nil {
			log.MappingFailed(operation, err)
			return nil, errors.CacheMapError
		}

		rates = append(rates, rate)
	}

	if len(staleKeys) > 0 {
		if err := c.rdb.SRem(ctx, c.GetCurrencyFromKey(currencyIsoFrom), staleKeys...).Err(); err != nil {
			log.DependencyFailed(operation, err)
			return nil, errors.CacheGetError
		}
	}
	log.Result(operation, len(rates))

	return rates, nil
}

// MapToRanking converts a Redis hash into an exchange-rate domain model.
func (c *ExchangeRateCache) MapToRanking(data map[string]string) (ranking.Ranking, error) {
	defer log.MappingStarted("redis_hash_to_ranking")()
	currencyIsoFrom, err := strconv.ParseInt(data["currency_iso_from"], 10, 32)
	if err != nil {
		return ranking.Ranking{}, errors.CacheMapError
	}

	currencyIsoTo, err := strconv.ParseInt(data["currency_iso_to"], 10, 32)
	if err != nil {
		return ranking.Ranking{}, errors.CacheMapError
	}

	date, err := strconv.ParseInt(data["date"], 10, 64)
	if err != nil {
		return ranking.Ranking{}, errors.CacheMapError
	}

	rateSell, err := strconv.ParseFloat(data["rate_sell"], 32)
	if err != nil {
		return ranking.Ranking{}, errors.CacheMapError
	}

	rateBuy, err := strconv.ParseFloat(data["rate_buy"], 32)
	if err != nil {
		return ranking.Ranking{}, errors.CacheMapError
	}

	rateCross, err := strconv.ParseFloat(data["rate_cross"], 32)
	if err != nil {
		return ranking.Ranking{}, errors.CacheMapError
	}

	return ranking.Ranking{
		CurrencyIsoFrom: int32(currencyIsoFrom),
		CurrencyIsoTo:   int32(currencyIsoTo),
		Date:            time.Unix(0, date),
		RateSell:        float32(rateSell),
		RateBuy:         float32(rateBuy),
		RateCross:       float32(rateCross),
	}, nil
}

// GetPrimaryKey builds the Redis hash key for a currency pair.
func (c *ExchangeRateCache) GetPrimaryKey(currencyIsoFrom int32, currencyIsoTo int32) string {
	return exchangeRatePairPrefix +
		strconv.FormatInt(int64(currencyIsoFrom), 10) + ":" +
		strconv.FormatInt(int64(currencyIsoTo), 10)
}

// GetCurrencyFromKey builds the Redis set key for a source-currency index.
func (c *ExchangeRateCache) GetCurrencyFromKey(currencyIsoFrom int32) string {
	return exchangeRateFromPrefix + strconv.FormatInt(int64(currencyIsoFrom), 10)
}
