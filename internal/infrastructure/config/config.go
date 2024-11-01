package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
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
	var result map[string]interface{}
	var config Config

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// First unmarshal into a map
	if err := viper.Unmarshal(&result); err != nil {
		return nil, fmt.Errorf("unable to decode into map: %w", err)
	}

	// Then decode map into the struct using mapstructure
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &config,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating decoder: %w", err)
	}

	if err := decoder.Decode(result); err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	// Handle Kafka brokers separately if needed
	brokersStr := viper.GetString("KAFKA_BROKERS")
	if brokersStr != "" && len(config.Kafka.Brokers) == 0 {
		config.Kafka.Brokers = strings.Split(brokersStr, ",")
	}

	// Debug: Print final config
	log.Printf("Loaded config: %+v", config)

	return &config, nil
}
