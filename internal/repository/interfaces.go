package repository

import (
	"context"
	"mpc/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params domain.CreateHashedUserParams) (domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, params domain.CreateWalletParams) (domain.Wallet, error)
	GetWallet(ctx context.Context, id int64) (domain.Wallet, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, params domain.CreateTransactionParams) (domain.Transaction, error)
	GetTransaction(ctx context.Context, id int64) (domain.Transaction, error)
}
