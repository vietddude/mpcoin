package main

import (
	"log"
	"mpc/internal/delivery/http"
	"mpc/internal/infrastructure/auth"
	"mpc/internal/infrastructure/config"
	"mpc/internal/infrastructure/db"
	"mpc/internal/infrastructure/logger"
	"mpc/internal/repository/postgres"
	"mpc/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	log := logger.NewLogger()

	jwtConfig := auth.NewJWTConfig(cfg.JWT.SecretKey, cfg.JWT.TokenDuration)
	jwtService := auth.NewJWTService(jwtConfig)

	// repository
	userRepo := postgres.NewUserRepository(dbPool)
	walletRepo := postgres.NewWalletRepository(dbPool)
	transactionRepo := postgres.NewTransactionRepository(dbPool)

	// usecase
	userUseCase := usecase.NewUserUseCase(userRepo)
	walletUseCase := usecase.NewWalletUseCase(walletRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo)

	// router
	router := http.NewRouter(userUseCase, walletUseCase, transactionUseCase, jwtService, log)

	log.Fatal(router.Run(":8080"))
}
