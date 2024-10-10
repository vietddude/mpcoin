package postgres

import (
	"context"
	"mpc/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, params domain.CreateWalletParams) (domain.Wallet, error) {
	// Implement the database operation here
	panic("not implemented")
}

func (r *WalletRepository) GetWallet(ctx context.Context, id int64) (domain.Wallet, error) {
	// Implement the database operation here
	panic("not implemented")
}