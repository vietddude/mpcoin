package postgres

import (
	"context"
	"mpc/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, params domain.CreateTransactionParams) (domain.Transaction, error) {
	// Implement the database operation here
	panic("not implemented")
}

func (r *TransactionRepository) GetTransaction(ctx context.Context, id int64) (domain.Transaction, error) {
	// Implement the database operation here
	panic("not implemented")
}
