package usecase

import (
	"context"
	"mpc/internal/domain"
	"mpc/internal/repository"
)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{transactionRepo: transactionRepo}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, params domain.CreateTransactionParams) (domain.Transaction, error) {
	return uc.transactionRepo.CreateTransaction(ctx, params)
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, id int64) (domain.Transaction, error) {
	return uc.transactionRepo.GetTransaction(ctx, id)
}
