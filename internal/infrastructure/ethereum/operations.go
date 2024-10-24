// Package ethereum provides an implementation of the EthereumRepository interface
// for interacting with the Ethereum blockchain.
package ethereum

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"mpc/internal/repository"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumClient represents a client for interacting with the Ethereum blockchain.
type EthereumClient struct {
	client    *ethclient.Client
	secretKey string
}

func NewEthereumClient(url, secretKey string) (*EthereumClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	return &EthereumClient{client: client, secretKey: secretKey}, nil
}

// Ensure EthereumClient implements EthereumRepository
var _ repository.EthereumRepository = (*EthereumClient)(nil)

// CreateWallet generates a new Ethereum wallet.
// It returns the private key, the associated Ethereum address, and any error encountered.
func (c *EthereumClient) CreateWallet() (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, address, nil
}

// GetBalance retrieves the balance of the given Ethereum address.
// It returns the balance as a big.Int and any error encountered.
func (c *EthereumClient) GetBalance(address common.Address) (*big.Int, error) {
	balance, err := c.client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// CreateUnsignedTransaction creates an unsigned Ethereum transaction.
// It takes the sender's address, recipient's address, and the amount to send.
// Returns the unsigned transaction and any error encountered.
func (c *EthereumClient) CreateUnsignedTransaction(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nonce, err := c.client.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	gasLimit, err := c.client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: amount,
	})
	if err != nil {
		// If estimation fails, use a default value
		gasLimit = uint64(21000)
	}

	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	}

	return types.NewTx(txData), nil
}

// SignTransaction signs an Ethereum transaction with the given private key.
// It returns the signed transaction and any error encountered.
func (c *EthereumClient) SignTransaction(tx *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	chainID, err := c.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}

	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

// SubmitTransaction submits a signed transaction to the Ethereum network.
// It returns the transaction hash and any error encountered.
func (c *EthereumClient) SubmitTransaction(signedTx *types.Transaction) (common.Hash, error) {
	err := c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash(), nil
}

// WaitForTxn waits for a transaction to be mined and returns its receipt.
// It takes the transaction hash and returns the transaction receipt and any error encountered.
func (c *EthereumClient) WaitForTxn(hash common.Hash) (*types.Receipt, error) {
	receipt, err := c.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	return receipt, nil
}

// EncryptPrivateKey encrypts the given private key data using AES-GCM encryption.
// It returns the encrypted data and any error encountered.
func (c *EthereumClient) EncryptPrivateKey(data []byte) ([]byte, error) {
	// Define or retrieve the AES key (this should be securely stored and managed)
	aesKey := []byte(c.secretKey)

	// Create a new AES cipher block
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// Wrap the AES cipher in Galois Counter Mode (GCM)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce (GCM standard requires a nonce size of 12 bytes)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the private key (data)
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil) // prepend nonce to ciphertext

	return ciphertext, nil
}

// DecryptPrivateKey decrypts the given encrypted private key data using AES-GCM.
// It returns the decrypted data and any error encountered.
func (c *EthereumClient) DecryptPrivateKey(ciphertext []byte) ([]byte, error) {
	aesKey := []byte(c.secretKey)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Separate nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *EthereumClient) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(5 * time.Second) // Poll every 5 seconds
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute) // Set a 5-minute timeout

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for transaction receipt")
		case <-ticker.C:
			receipt, err := c.client.TransactionReceipt(ctx, hash)
			if err == nil {
				return receipt, nil
			}
			if err != ethereum.NotFound {
				return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
			}

			// If the transaction is successful, return the receipt
			if receipt.Status == 1 {
				return receipt, nil
			}
			// If err == ethereum.NotFound, continue polling
		}
	}
}

func (c *EthereumClient) ParseTransactionReceipt(receipt *types.Receipt) {
	fmt.Printf("Transaction Hash: %s\n", receipt.TxHash.Hex())
	fmt.Printf("Status: %d\n", receipt.Status) // 1 for success, 0 for failure
	fmt.Printf("Block Number: %d\n", receipt.BlockNumber)
	fmt.Printf("Block Hash: %s\n", receipt.BlockHash.Hex())
	fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
	fmt.Printf("Cumulative Gas Used: %d\n", receipt.CumulativeGasUsed)

	if receipt.ContractAddress != (common.Address{}) {
		fmt.Printf("Contract Address: %s\n", receipt.ContractAddress.Hex())
	}

	fmt.Printf("Transaction Index: %d\n", receipt.TransactionIndex)

	if receipt.EffectiveGasPrice != nil {
		fmt.Printf("Effective Gas Price: %s\n", receipt.EffectiveGasPrice.String())
	}

	fmt.Printf("Logs: %d\n", len(receipt.Logs))
	for i, log := range receipt.Logs {
		fmt.Printf("  Log %d:\n", i)
		fmt.Printf("    Address: %s\n", log.Address.Hex())
		fmt.Printf("    Topics: %v\n", log.Topics)
		fmt.Printf("    Data: %x\n", log.Data)
	}
}
