package main

import (
	"context"
	"encoding/json"
	"log"
	"mpc/internal/domain"
	"mpc/internal/infrastructure/config"
	"mpc/internal/infrastructure/db"
	"mpc/internal/infrastructure/ethereum"
	"mpc/internal/infrastructure/kafka"
	"mpc/internal/repository"
	"mpc/internal/repository/postgres"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

const (
	TopicTxReceipt = "tx_receipt_topic"
	TopicBalance   = "balance_topic"
)

type TxReceiptTask struct {
	TxHash string `json:"tx_hash"`
}

type BalanceTask struct {
	Address string `json:"address"`
}

// This worker is responsible for processing transaction receipts and updating the transaction status.
// It also listens to the balance topic to update the balance in the database.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	ethClient, err := ethereum.NewEthereumClient(cfg.Ethereum.URL, cfg.Ethereum.SecretKey)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %v", err)
	}

	// kafka
	producer, err := kafka.NewKafkaProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewKafkaConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// repositories
	txnRepo := postgres.NewTransactionRepo(dbPool)

	ctx := context.Background()

	go processTxReceiptTopic(ctx, consumer, txnRepo, ethClient)
	// go processBalanceTopic(ctx, cfg.Kafka.Brokers, ethRepo)

	select {}
}

func processTxReceiptTopic(ctx context.Context, consumer *kafka.Reader, txnRepo repository.TransactionRepository, ethClient *ethereum.EthereumClient) {
	for {
		m, err := kafka.ReadNewMessage(ctx, consumer)
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}
		log.Printf("Received new message")

		txnID := uuid.MustParse(string(m.Key))
		txnData := m.Value
		// TODO: Decode the transaction data
		var txn domain.TxnMessage
		err = json.Unmarshal(txnData, &txn)
		if err != nil {
			log.Fatalf("Failed to decode transaction data: %v", err)
		}

		// Fetch the transaction from the database
		txnFound, err := txnRepo.GetTransaction(ctx, txnID)
		if err != nil {
			log.Fatalf("Transaction not found: %v", err)
		}

		if txnFound.Status == domain.StatusSubmitted {
			// TODO: Process the transaction receipt
			receipt, err := ethClient.GetTransactionReceipt(ctx, common.HexToHash(txn.TxHash))
			if err != nil {
				log.Fatalf("Failed to fetch transaction receipt: %v", err)
			}
			ethClient.ParseTransactionReceipt(receipt)
			// Update the transaction status in the database
			txnFound.Status = domain.StatusSuccess
			txnFound.GasPrice = receipt.EffectiveGasPrice.String()
			txnFound.GasLimit = strconv.FormatUint(receipt.GasUsed, 10)
			txnFound.Nonce = int64(receipt.TransactionIndex)

			err = txnRepo.UpdateTransaction(ctx, txnFound)
			if err != nil {
				log.Fatalf("Failed to update transaction status: %v", err)
			}
			log.Printf("Transaction status updated: (%v, %v)", txnFound.ID, txnFound.Status)
		}
	}
}
