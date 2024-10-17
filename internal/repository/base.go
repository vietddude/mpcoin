package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTransaction interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type BaseRepository struct {
	db *pgxpool.Pool
}

func NewBaseRepo(db *pgxpool.Pool) BaseRepository {
	return BaseRepository{db: db}
}

func (r *BaseRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *BaseRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
