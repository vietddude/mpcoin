package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mpc/internal/infrastructure/config"
	"mpc/internal/infrastructure/kafka"
	"mpc/internal/infrastructure/mail"
	"mpc/internal/infrastructure/otp"
	"mpc/internal/infrastructure/redis"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	OTPTopic = "otp_emails"
)

type OTPMessage struct {
	Email string `json:"email"`
}

type MailWorkerConfig struct {
	NumWorkers int `mapstructure:"NUM_WORKERS"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	consumer, err := kafka.NewKafkaConsumer(cfg, kafka.WithTopic(OTPTopic))
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}
	defer redisClient.Close()

	otpService := otp.NewOTPService(redisClient, time.Duration(cfg.MailConfig.OTPExpiration)*time.Second)

	mailClient, err := mail.NewClient(&cfg.MailConfig)
	if err != nil {
		log.Fatalf("Failed to create mail client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mailWorkerConfig MailWorkerConfig
	var wg sync.WaitGroup
	for i := 0; i < mailWorkerConfig.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleOTPMessages(ctx, consumer, otpService, mailClient)
		}()
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
	cancel()
	wg.Wait()
}

func handleOTPMessages(ctx context.Context, consumer *kafka.Reader, otpService otp.OTPService, mailClient *mail.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := consumer.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("Failed to read message: %v", err)
				continue
			}

			var otpMsg OTPMessage
			if err := json.Unmarshal(m.Value, &otpMsg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			if err := processOTPMessage(ctx, otpService, mailClient, otpMsg); err != nil {
				log.Printf("Failed to process OTP message: %v", err)
			}
		}
	}
}

func processOTPMessage(ctx context.Context, otpService otp.OTPService, mailClient *mail.Client, otpMsg OTPMessage) error {
	otp, err := otpService.GenerateOTP(ctx, otpMsg.Email)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	if err := mailClient.SendMail(otpMsg.Email, "Your OTP", fmt.Sprintf("Your OTP is: %s", otp)); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	log.Printf("OTP sent successfully to %s", otpMsg.Email)
	return nil
}
