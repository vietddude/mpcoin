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

type walletRepository struct {
	repository.BaseRepository
}

func NewWalletRepo(dbPool *pgxpool.Pool) repository.WalletRepository {
	return &walletRepository{
		BaseRepository: repository.NewBaseRepo(dbPool),
	}
}

// Ensure WalletRepository implements WalletRepository
var _ repository.WalletRepository = (*walletRepository)(nil)

func (r *walletRepository) CreateWallet(ctx context.Context, params domain.CreateWalletParams) (domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.WithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(tx)
		createdWallet, err := q.CreateWallet(ctx, sqlc.CreateWalletParams{
			UserID:              pgtype.UUID{Bytes: params.UserID, Valid: true},
			Address:             params.Address,
			EncryptedPrivateKey: params.EncryptedPrivateKey,
		})
		if err != nil {
			return err
		}

		wallet = domain.Wallet{
			ID:                  createdWallet.ID.Bytes,
			UserID:              createdWallet.UserID.Bytes,
			Address:             createdWallet.Address,
			EncryptedPrivateKey: createdWallet.EncryptedPrivateKey,
		}
		return nil
	})
	return wallet, err
}

func (r *walletRepository) GetWallet(ctx context.Context, id uuid.UUID) (domain.Wallet, error) {
	q := sqlc.New(r.DB())
	wallet, err := q.GetWallet(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return domain.Wallet{}, err
	}

	return domain.Wallet{
		ID:                  wallet.ID.Bytes,
		UserID:              wallet.UserID.Bytes,
		Address:             wallet.Address,
		EncryptedPrivateKey: wallet.EncryptedPrivateKey,
	}, nil
}

func (r *walletRepository) GetWalletByAddress(ctx context.Context, address string) (domain.Wallet, error) {
	q := sqlc.New(r.DB())
	wallet, err := q.GetWalletByAddress(ctx, address)
	if err != nil {
		return domain.Wallet{}, err
	}

	return domain.Wallet{
		ID:                  wallet.ID.Bytes,
		UserID:              wallet.UserID.Bytes,
		Address:             wallet.Address,
		EncryptedPrivateKey: wallet.EncryptedPrivateKey,
	}, nil
}

func (r *walletRepository) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (domain.Wallet, error) {
	q := sqlc.New(r.DB())
	wallet, err := q.GetWalletByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return domain.Wallet{}, err
	}
	return domain.Wallet{
		ID:                  wallet.ID.Bytes,
		UserID:              wallet.UserID.Bytes,
		Address:             wallet.Address,
		EncryptedPrivateKey: wallet.EncryptedPrivateKey,
	}, nil
}
