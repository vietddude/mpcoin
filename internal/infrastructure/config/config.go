package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Ethereum EthereumConfig
}

type DBConfig struct {
	ConnStr string `mapstructure:"CONN_STR"`
}

type AppConfig struct {
	Port int `mapstructure:"PORT"`
}

type JWTConfig struct {
	SecretKey     string        `mapstructure:"SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

type EthereumConfig struct {
	URL       string `mapstructure:"ETHEREUM_URL"`
	SecretKey string `mapstructure:"ETHEREUM_SECRET_KEY"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"ADDR"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

// Define default values
var defaults = map[string]string{
	"DB.CONN_STR":        "postgres://viet:123@localhost:5432/mpcoin?sslmode=disable",
	"DB.MAX_CONNECTIONS": "10",
	"APP.PORT":           "8080",
	"APP.ENV":            "development",
	"JWT.SECRET_KEY":     "chirp-chirp",
	"JWT.TOKEN_DURATION": "1h",
	"REDIS.ADDR":         "localhost:6379",
	"REDIS.PASSWORD":     "",
	"REDIS.DB":           "0",
	// "ETHEREUM.URL":        "https://sepolia.infura.io/v3/<INFURA_PROJECT_ID>",
	// "ETHEREUM.SECRET_KEY": "<INFURA_SECRET_KEY>",
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set environment variable names to match .env file
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		fmt.Println("No .env file found. Using environment variables.")
	}

	// Set default values if not provided in .env or environment
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Manually set values for fields that don't match the default structure
	config.DB.ConnStr = viper.GetString("CONN_STR")
	config.Ethereum.URL = viper.GetString("ETHEREUM_URL")
	config.Ethereum.SecretKey = viper.GetString("ETHEREUM_SECRET_KEY")

	// Set default values if not provided
	if config.JWT.TokenDuration == 0 {
		config.JWT.TokenDuration = 24 * time.Hour // Default to 24 hours
	}

	log.Printf("Config loaded")
	return &config, nil
}
