package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	DB       DBConfig       `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
	Ethereum EthereumConfig `mapstructure:",squash"`
	Kafka    KafkaConfig    `mapstructure:",squash"`
	Mail     MailConfig     `mapstructure:",squash"`
}

type AppConfig struct {
	Port int `mapstructure:"PORT"`
}

type DBConfig struct {
	ConnStr string `mapstructure:"CONN_STR"`
}

type JWTConfig struct {
	SecretKey     string        `mapstructure:"JWT_SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"JWT_TOKEN_DURATION"`
}

type RedisConfig struct {
	Address  string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	Username string `mapstructure:"REDIS_USERNAME"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type EthereumConfig struct {
	URL       string `mapstructure:"ETHEREUM_URL"`
	SecretKey string `mapstructure:"ETHEREUM_SECRET_KEY"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"KAFKA_BROKERS"`
	Topic   string   `mapstructure:"KAFKA_TOPIC"`
}

type MailConfig struct {
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      int    `mapstructure:"SMTP_PORT"`
	SMTPUsername  string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword  string `mapstructure:"SMTP_PASSWORD"`
	FromEmail     string `mapstructure:"FROM_EMAIL"`
	OTPExpiration int    `mapstructure:"OTP_EXPIRATION"`
}

func Load() (*Config, error) {
	// Enable reading from environment variables first
	viper.AutomaticEnv()

	// Try to load .env file (optional)
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		// Ignore file not found error and continue
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only log a warning for missing .env file
			log.Printf("Warning: .env file not found, using environment variables")
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Handle Kafka brokers special case
	brokersStr := viper.GetString("KAFKA_BROKERS")
	if brokersStr != "" {
		config.Kafka.Brokers = strings.Split(brokersStr, ",")
	}

	return &config, nil
}
