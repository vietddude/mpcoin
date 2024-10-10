package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	JWT   JWTConfig
	Redis RedisConfig
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
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	// Set default values
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using environment variables and defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Set default values if not provided
	if config.JWT.TokenDuration == 0 {
		config.JWT.TokenDuration = 24 * time.Hour // Default to 24 hours
	}

	return &config, nil
}
