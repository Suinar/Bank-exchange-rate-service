package configs

import (
	errors "errors"
	log "log"
	os "os"

	viper "github.com/spf13/viper"
)

// Config contains the runtime settings required by application components.
type Config struct {
	GRPC struct {
		Port    string
		Network string
	}

	Monobank struct {
		BaseURL          string
		CurrencyEndpoint string
	}

	Kafka struct {
		Brokers []string
	}
}

// Load reads application settings from the environment and the local .env file.
func Load() *Config {
	viper.SetDefault("GRPC_PORT", "50053")
	viper.SetDefault("GRPC_NETWORK", "tcp")
	viper.SetDefault("MONOBANK_BASE_URL", "https://api.monobank.ua")
	viper.SetDefault("MONOBANK_CURRENCY_ENDPOINT", "/bank/currency")

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	cfg := &Config{}

	cfg.GRPC.Port = viper.GetString("GRPC_PORT")
	cfg.GRPC.Network = viper.GetString("GRPC_NETWORK")

	cfg.Monobank.BaseURL = viper.GetString("MONOBANK_BASE_URL")
	cfg.Monobank.CurrencyEndpoint = viper.GetString("MONOBANK_CURRENCY_ENDPOINT")

	cfg.Kafka.Brokers = viper.GetStringSlice("KAFKA_BROKERS")

	return cfg
}
