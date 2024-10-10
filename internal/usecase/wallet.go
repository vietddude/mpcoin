package usecase

import (
	"context"
	"mpc/internal/domain"
	"mpc/internal/repository"
)

type WalletUseCase struct {
	walletRepo repository.WalletRepository
}

func NewWalletUseCase(walletRepo repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{walletRepo: walletRepo}
}

func (uc *WalletUseCase) CreateWallet(ctx context.Context, params domain.CreateWalletParams) (domain.Wallet, error) {
	return uc.walletRepo.CreateWallet(ctx, params)
}

func (uc *WalletUseCase) GetWallet(ctx context.Context, id int64) (domain.Wallet, error) {
	return uc.walletRepo.GetWallet(ctx, id)
}
