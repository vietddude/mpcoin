package repository

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"mpc/internal/domain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params domain.CreateHashedUserParams) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DBTransaction
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, params domain.CreateWalletParams) (domain.Wallet, error)
	GetWallet(ctx context.Context, id uuid.UUID) (domain.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (domain.Wallet, error)
	GetWalletByAddress(ctx context.Context, address string) (domain.Wallet, error)
	DBTransaction
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, params domain.CreateTransactionParams) (domain.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (domain.Transaction, error)
	GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]domain.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction domain.Transaction) error
	DBTransaction
}

type EthereumRepository interface {
	CreateWallet() (*ecdsa.PrivateKey, common.Address, error)
	GetBalance(address common.Address) (*big.Int, error)
	CreateUnsignedTransaction(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error)
	SignTransaction(tx *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error)
	SubmitTransaction(signedTx *types.Transaction) (common.Hash, error)
	WaitForTxn(hash common.Hash) (*types.Receipt, error)
	EncryptPrivateKey(data []byte) ([]byte, error)
	DecryptPrivateKey(ciphertext []byte) ([]byte, error)
}
