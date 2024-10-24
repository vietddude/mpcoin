package usecase

import (
	"context"
	"mpc/internal/infrastructure/config"
	"mpc/internal/infrastructure/kafka"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

var mockConfig = config.Config{
	Kafka: config.KafkaConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "mpc",
	},
}

func TestPublishMessage(t *testing.T) {

	// Create a real Kafka writer
	writer, err := kafka.NewKafkaProducer(&mockConfig)
	if err != nil {
		t.Fatalf("Failed to create Kafka writer: %v", err)
	}
	defer writer.Close()

	uc := &txnUseCase{kafkaProducer: writer}

	ctx := context.Background()
	txnID := uuid.MustParse("48f8fc39-e2c5-467a-8429-14d68a1a4e37")
	chainID := uuid.MustParse("1ec0a60a-08fe-4fb2-a6b2-ca2f506f8275")
	txHash := common.HexToHash("0x1dd41fd648b9feb3b4ee6a8e7ed0e964d01d20188b9873650fcdded8b8e6a566")

	// Call the method you want to test
	uc.publishMessage(ctx, txnID, chainID, txHash)
}
