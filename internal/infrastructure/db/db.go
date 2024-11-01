package db

import (
	"context"
	"log"
	"mpc/internal/infrastructure/config"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

// InitDB initializes the database connection pool
func InitDB(cfg *config.DBConfig) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		dbPool, err = pgxpool.New(context.Background(), cfg.ConnStr)
		if err != nil {
			log.Printf("Unable to create connection pool: %v\n", err)
		}
	})
	return dbPool, err
}

// GetDB returns the database connection pool
func GetDB() *pgxpool.Pool {
	return dbPool
}

// CloseDB closes the database connection pool
func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
	}
}
