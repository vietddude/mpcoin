package usecase

import (
	"context"
	"errors"
	"fmt"
	"mpc/internal/domain"
	"mpc/internal/infrastructure/auth"
	"mpc/internal/repository"
	"mpc/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthUseCase interface {
	Signup(ctx context.Context, params domain.CreateUserParams) (domain.CreateUserResponse, domain.CreateWalletResponse, string, string, error)
	Login(ctx context.Context, email, password string) (domain.LoginUserResponse, string, string, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string) (string, string, error)
}

type authUseCase struct {
	userRepo   repository.UserRepository
	walletUC   WalletUseCase
	jwtService auth.JWTService
}

func NewAuthUC(userRepo repository.UserRepository, walletUC WalletUseCase, jwtService auth.JWTService) AuthUseCase {
	return &authUseCase{userRepo: userRepo, walletUC: walletUC, jwtService: jwtService}
}

var _ AuthUseCase = (*authUseCase)(nil)

func (uc *authUseCase) Signup(ctx context.Context, params domain.CreateUserParams) (domain.CreateUserResponse, domain.CreateWalletResponse, string, string, error) {
	var user domain.User
	var wallet domain.Wallet

	err := uc.userRepo.WithTx(ctx, func(tx pgx.Tx) error {
		// Check if user already exists
		existingUser, err := uc.userRepo.GetUserByEmail(ctx, params.Email)
		if err == nil && existingUser.ID != uuid.Nil {
			return errors.New("user with this email already exists")
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(params.Password)
		if err != nil {
			return err
		}

		// Create user with hashed password
		createParams := domain.CreateHashedUserParams{
			Email:        params.Email,
			PasswordHash: hashedPassword,
		}

		user, err = uc.userRepo.CreateUser(ctx, createParams)
		if err != nil {
			return err
		}

		// Create wallet for the user
		wallet, err = uc.walletUC.CreateWallet(ctx, user.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return domain.CreateUserResponse{}, domain.CreateWalletResponse{}, "", "", err
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		return domain.CreateUserResponse{}, domain.CreateWalletResponse{}, "", "", err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return domain.CreateUserResponse{}, domain.CreateWalletResponse{}, "", "", err
	}

	createWalletResponse := domain.CreateWalletResponse{
		ID:      wallet.ID,
		UserID:  wallet.UserID,
		Address: wallet.Address,
	}

	createUserResponse := domain.CreateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return createUserResponse, createWalletResponse, token, refreshToken, nil
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (domain.LoginUserResponse, string, string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		fmt.Println("error getting user by email", err)
		return domain.LoginUserResponse{}, "", "", errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		fmt.Println("error checking password hash", err)
		return domain.LoginUserResponse{}, "", "", errors.New("invalid credentials")
	}

	// Generate access token
	accessToken, err := uc.jwtService.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		return domain.LoginUserResponse{}, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := uc.jwtService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return domain.LoginUserResponse{}, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return domain.LoginUserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, accessToken, refreshToken, nil
}

func (uc *authUseCase) Logout(ctx context.Context, token string) error {
	fmt.Println("logout token", token)
	claims, err := uc.jwtService.ValidateToken(ctx, token, auth.AccessToken)
	if err != nil {
		return err
	}

	return uc.jwtService.InvalidateToken(ctx, claims.UserID, auth.AccessToken)
}

func (uc *authUseCase) RefreshToken(ctx context.Context, token string) (string, string, error) {
	fmt.Println("refresh token", token)
	accessToken, refreshToken, err := uc.jwtService.RefreshTokens(ctx, token)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
