package configs

import (
	errors "errors"
	fmt "fmt"
	os "os"
	time "time"

	log "github.com/kVinsom/Bank-exhange-rate-service/internal/logging/config"
	viper "github.com/spf13/viper"
)

// Config contains the runtime settings required by application components.
type Config struct {
	// GRPC contains settings for running the gRPC server.
	GRPC struct {
		Port    string
		Network string
	}

	// Monobank contains settings for requests to the Monobank currency API.
	Monobank struct {
		BaseURL          string
		CurrencyEndpoint string
		RefreshInterval  time.Duration
	}

	// Kafka contains broker addresses and topics used for exchange-rate requests.
	Kafka struct {
		Brokers                        []string
		GetRelativeRankingRequestTopic string
		GetAllRankingRequestTopic      string
	}

	// Redis contains the connection and pool settings used by the cache.
	Redis struct {
		Addr         string
		Password     string
		DB           int
		PoolSize     int
		MinIdleConns int
	}
}

// Load reads application settings from the environment and the local .env file.
func Load() (*Config, error) {
	const configFile = ".env"
	log.Loading(configFile)
	viper.SetDefault("GRPC_PORT", "55053")
	viper.SetDefault("GRPC_NETWORK", "tcp")
	viper.SetDefault("MONOBANK_BASE_URL", "https://api.monobank.ua")
	viper.SetDefault("MONOBANK_CURRENCY_ENDPOINT", "/bank/currency")
	viper.SetDefault("MONOBANK_REFRESH_INTERVAL", "5m")
	viper.SetDefault("GET_RELATIVE_RANKING_REQUEST_TOPIC", "exchange-rate.get-relative-ranking.request")
	viper.SetDefault("GET_ALL_RANKING_REQUEST_TOPIC", "exchange-rate.get-all-ranking.request")
	viper.SetDefault("REDIS_ADDR", "localhost:16380")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_POOL_SIZE", 10)
	viper.SetDefault("REDIS_MIN_IDLE_CONNS", 2)

	viper.SetConfigFile(configFile)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.LoadingFailed(configFile, err)
			return nil, fmt.Errorf("read configuration: %w", err)
		}
		log.FileNotFound(configFile)
	} else {
		log.FileLoaded(configFile)
	}

	cfg := &Config{}

	cfg.GRPC.Port = viper.GetString("GRPC_PORT")
	cfg.GRPC.Network = viper.GetString("GRPC_NETWORK")

	cfg.Monobank.BaseURL = viper.GetString("MONOBANK_BASE_URL")
	cfg.Monobank.CurrencyEndpoint = viper.GetString("MONOBANK_CURRENCY_ENDPOINT")
	cfg.Monobank.RefreshInterval = viper.GetDuration("MONOBANK_REFRESH_INTERVAL")
	if cfg.Monobank.RefreshInterval <= 0 {
		return nil, errors.New("MONOBANK_REFRESH_INTERVAL must be greater than zero")
	}

	cfg.Kafka.Brokers = viper.GetStringSlice("KAFKA_BROKERS")
	cfg.Kafka.GetRelativeRankingRequestTopic = viper.GetString("GET_RELATIVE_RANKING_REQUEST_TOPIC")
	cfg.Kafka.GetAllRankingRequestTopic = viper.GetString("GET_ALL_RANKING_REQUEST_TOPIC")

	cfg.Redis.Addr = viper.GetString("REDIS_ADDR")
	cfg.Redis.Password = viper.GetString("REDIS_PASSWORD")
	cfg.Redis.DB = viper.GetInt("REDIS_DB")
	cfg.Redis.PoolSize = viper.GetInt("REDIS_POOL_SIZE")
	cfg.Redis.MinIdleConns = viper.GetInt("REDIS_MIN_IDLE_CONNS")

	log.Loaded(cfg.GRPC.Port, cfg.Redis.Addr, cfg.Monobank.BaseURL)
	return cfg, nil
}
