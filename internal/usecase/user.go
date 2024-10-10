package usecase

import (
	"context"
	"errors"
	"mpc/internal/domain"
	"mpc/internal/repository"
	"mpc/pkg/utils"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetUserByEmail(ctx, params.Email)
	if err == nil && existingUser.ID != 0 {
		return domain.User{}, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return domain.User{}, err
	}

	// Create user with hashed password
	createParams := domain.CreateHashedUserParams{
		Email:        params.Email,
		PasswordHash: hashedPassword,
	}

	return uc.userRepo.CreateUser(ctx, createParams)
}

func (uc *UserUseCase) GetUser(ctx context.Context, id int64) (domain.User, error) {
	return uc.userRepo.GetUser(ctx, id)
}

func (uc *UserUseCase) AuthenticateUser(ctx context.Context, email, password string) (domain.User, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return domain.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, params domain.UpdateUserParams) (domain.User, error) {
	// Check if user exists
	existingUser, err := uc.userRepo.GetUser(ctx, params.ID)
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}

	// Update fields if provided
	if params.Email != "" {
		existingUser.Email = params.Email
	}
	if params.Password != "" {
		hashedPassword, err := utils.HashPassword(params.Password)
		if err != nil {
			return domain.User{}, err
		}
		existingUser.PasswordHash = hashedPassword
	}

	return uc.userRepo.UpdateUser(ctx, existingUser)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
	return uc.userRepo.DeleteUser(ctx, id)
}

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return uc.userRepo.GetUserByEmail(ctx, email)
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	return uc.userRepo.GetUser(ctx, id)
}
