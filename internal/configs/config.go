package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GRPC struct {
		GRPCPort string
		Network  string
	}

	monoBank struct {
		baseUrl          string
		CurrencyEndpoint string
	}

	kafka struct {
		Brokers []string
	}
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}

	cfg.GRPC.GRPCPort = viper.GetString("GRPS_PORT")
	cfg.GRPC.Network = viper.GetString("GRPC_NETWORK")

	cfg.monoBank.baseUrl = viper.GetString("MONOBANK_BASE_URL")
	cfg.monoBank.CurrencyEndpoint = viper.GetString("MONOBANK_CURRENCY_ENDPOINT")

	cfg.kafka.Brokers = viper.GetStringSlice("KAFKA_BROKERS")

	return cfg
}
