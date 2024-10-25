package otp

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"mpc/internal/infrastructure/redis"
)

type OTPService interface {
	GenerateOTP(ctx context.Context, email string) (string, error)
	VerifyOTP(ctx context.Context, email, otp string) error
}

type otpService struct {
	redisClient *redis.RedisClient
	expiration  time.Duration
}

func NewOTPService(redisClient *redis.RedisClient, expiration time.Duration) OTPService {
	return &otpService{
		redisClient: redisClient,
		expiration:  expiration,
	}
}

func (s *otpService) GenerateOTP(ctx context.Context, email string) (string, error) {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	err := s.redisClient.Set(ctx, fmt.Sprintf("otp:%s", email), otp, s.expiration)
	if err != nil {
		return "", err
	}
	return otp, nil
}

func (s *otpService) VerifyOTP(ctx context.Context, email, otp string) error {
	storedOTP, err := s.redisClient.Get(ctx, fmt.Sprintf("otp:%s", email))
	if err != nil {
		return err
	}
	if storedOTP != otp {
		return fmt.Errorf("invalid OTP")
	}
	return s.redisClient.Delete(ctx, fmt.Sprintf("otp:%s", email))
}

