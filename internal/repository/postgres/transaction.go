package postgres

import (
	"context"
	"mpc/internal/domain"
	sqlc "mpc/internal/infrastructure/db/sqlc"
	"mpc/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepository struct {
	repository.BaseRepository
}

func NewTransactionRepo(dbPool *pgxpool.Pool) repository.TransactionRepository {
	return &transactionRepository{
		BaseRepository: repository.NewBaseRepo(dbPool),
	}
}

// Ensure TransactionRepository implements TransactionRepository
var _ repository.TransactionRepository = (*transactionRepository)(nil)

func (r *transactionRepository) CreateTransaction(ctx context.Context, params domain.CreateTransactionParams) (domain.Transaction, error) {

	var transaction domain.Transaction
	err := r.WithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(tx)
		createdTransaction, err := q.CreateTransaction(ctx, sqlc.CreateTransactionParams{
			ID:        pgtype.UUID{Bytes: params.ID, Valid: true},
			WalletID:  pgtype.UUID{Bytes: params.WalletID, Valid: true},
			ChainID:   pgtype.UUID{Bytes: params.ChainID, Valid: true},
			ToAddress: params.ToAddress,
			Amount:    params.Amount,
			TokenID:   pgtype.UUID{Bytes: params.TokenID, Valid: true},
			GasPrice:  pgtype.Text{String: params.GasPrice, Valid: true},
			GasLimit:  pgtype.Text{String: params.GasLimit, Valid: true},
			Nonce:     pgtype.Int8{Int64: params.Nonce, Valid: true},
			Status:    string(params.Status),
		})
		if err != nil {
			return err
		}
		transaction = domain.Transaction{
			ID:        createdTransaction.ID.Bytes,
			WalletID:  createdTransaction.WalletID.Bytes,
			ChainID:   createdTransaction.ChainID.Bytes,
			ToAddress: createdTransaction.ToAddress,
			Amount:    createdTransaction.Amount,
			TokenID:   createdTransaction.TokenID.Bytes,
			GasPrice:  createdTransaction.GasPrice.String,
			GasLimit:  createdTransaction.GasLimit.String,
			Nonce:     createdTransaction.Nonce.Int64,
			Status:    domain.Status(createdTransaction.Status),
		}
		return nil
	})
	return transaction, err
}

func (r *transactionRepository) GetTransaction(ctx context.Context, id uuid.UUID) (domain.Transaction, error) {
	q := sqlc.New(r.DB())
	transaction, err := q.GetTransaction(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return domain.Transaction{}, err
	}
	return domain.Transaction{
		ID:        transaction.ID.Bytes,
		WalletID:  transaction.WalletID.Bytes,
		ChainID:   transaction.ChainID.Bytes,
		ToAddress: transaction.ToAddress,
		Amount:    transaction.Amount,
		TokenID:   transaction.TokenID.Bytes,
		TxHash:    transaction.TxHash.String,
		GasPrice:  transaction.GasPrice.String,
		GasLimit:  transaction.GasLimit.String,
		Nonce:     transaction.Nonce.Int64,
		Status:    domain.Status(transaction.Status),
	}, nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction domain.Transaction) error {
	q := sqlc.New(r.DB())
	_, err := q.UpdateTransaction(ctx, sqlc.UpdateTransactionParams{
		ID:       pgtype.UUID{Bytes: transaction.ID, Valid: true},
		Status:   string(transaction.Status),
		TxHash:   pgtype.Text{String: transaction.TxHash, Valid: transaction.TxHash != ""},
		GasPrice: pgtype.Text{String: transaction.GasPrice, Valid: transaction.GasPrice != ""},
		GasLimit: pgtype.Text{String: transaction.GasLimit, Valid: transaction.GasLimit != ""},
		Nonce:    pgtype.Int8{Int64: transaction.Nonce, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) GetTransactions(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	// Implement the database operation here
	panic("not implemented")
}

func (r *transactionRepository) GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]domain.Transaction, error) {
	// Implement the database operation here
	panic("not implemented")
}
