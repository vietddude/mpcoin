package config

import (
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
	Port int `envconfig:"PORT" default:"8080"`
}

type DBConfig struct {
	ConnStr string `envconfig:"CONN_STR"`
}

type JWTConfig struct {
	SecretKey     string `envconfig:"JWT_SECRET_KEY"`
	TokenDuration string `envconfig:"JWT_TOKEN_DURATION" default:"15m"`
}

type RedisConfig struct {
	Address  string `envconfig:"REDIS_ADDR"`
	Password string `envconfig:"REDIS_PASSWORD"`
	Username string `envconfig:"REDIS_USERNAME"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

type EthereumConfig struct {
	URL       string `envconfig:"ETHEREUM_URL"`
	SecretKey string `envconfig:"ETHEREUM_SECRET_KEY"`
}

type KafkaConfig struct {
	Brokers []string `envconfig:"KAFKA_BROKERS" split_words:"true"`
	Topic   string   `envconfig:"KAFKA_TOPIC"`
}

type MailConfig struct {
	SMTPHost      string `envconfig:"SMTP_HOST"`
	SMTPPort      int    `envconfig:"SMTP_PORT" default:"587"`
	SMTPUsername  string `envconfig:"SMTP_USERNAME"`
	SMTPPassword  string `envconfig:"SMTP_PASSWORD"`
	FromEmail     string `envconfig:"FROM_EMAIL"`
	OTPExpiration string `envconfig:"OTP_EXPIRATION"`
}

func Load(log *logrus.Logger) (*Config, error) {
	// Initialize Config struct
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Debugf("failed to process env: %+v", err)
	}

	// Debug log the config
	log.Debugf("Loaded configuration: %+v", config)

	return &config, nil
}
