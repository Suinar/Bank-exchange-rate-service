package cahce

import (
	context "context"
	fmt "fmt"

	configs "github.com/kVinsom/Bank-exhange-rate-service/internal/configs"
	redis "github.com/redis/go-redis/v9"
)

// ConnectToCahce creates a Redis client and verifies the connection.
func ConnectToCahce(ctx context.Context, cfg *configs.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect to Redis at %s: %w", cfg.Redis.Addr, err)
	}

	return client, nil
}
