package postgres

import (
	"context"
	"mpc/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, params domain.CreateHashedUserParams) (domain.User, error) {
	// Implement the database operation here
	panic("not implemented")
}

func (r *UserRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	// Implement the database operation here
	panic("not implemented")
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	// Implement the database operation here
	panic("not implemented")
}

func (r *UserRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	// Implement the update logic here
	// You may need to extract the user ID and update params from the user object
	// and use them to update the database
	panic("not implemented")
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	// Implement the database operation here
	panic("not implemented")
}
