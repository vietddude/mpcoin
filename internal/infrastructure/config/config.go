package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Ethereum EthereumConfig
	Kafka    KafkaConfig
	Mail     MailConfig
}

type AppConfig struct {
	Port int `env:"PORT" default:"8080"`
}

type DBConfig struct {
	ConnStr string `env:"CONN_STR"`
}

type JWTConfig struct {
	SecretKey     string `env:"JWT_SECRET_KEY"`
	TokenDuration string `env:"JWT_TOKEN_DURATION" default:"15m"`
}

type RedisConfig struct {
	Address  string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	Username string `env:"REDIS_USERNAME"`
	DB       int    `env:"REDIS_DB" default:"0"`
}

type EthereumConfig struct {
	URL       string `env:"ETHEREUM_URL"`
	SecretKey string `env:"ETHEREUM_SECRET_KEY"`
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" envconfig:"KAFKA_BROKERS" split_words:"true"`
	Topic   string   `env:"KAFKA_TOPIC"`
}

type MailConfig struct {
	SMTPHost      string `env:"SMTP_HOST"`
	SMTPPort      int    `env:"SMTP_PORT" default:"587"`
	SMTPUsername  string `env:"SMTP_USERNAME"`
	SMTPPassword  string `env:"SMTP_PASSWORD"`
	FromEmail     string `env:"FROM_EMAIL"`
	OTPExpiration string `env:"OTP_EXPIRATION"`
}

func Load(log *logrus.Logger) (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		// Only log as info since .env file is optional when using environment variables
		log.Infof("Note: .env file not found: %v", err)
	}

	// Initialize Config struct
	var config Config
	err = envconfig.Process("", &config)
	if err != nil {
		return nil, fmt.Errorf("failed to process env: %w", err)
	}

	// Debug log the config
	log.Debugf("Loaded configuration: %+v", config)

	return &config, nil
}
