package main

import (
	"log"
	_ "mpc/docs"
	"mpc/internal/delivery/http"
	_ "mpc/internal/domain"
	"mpc/internal/infrastructure/auth"
	"mpc/internal/infrastructure/config"
	"mpc/internal/infrastructure/db"
	"mpc/internal/infrastructure/ethereum"
	"mpc/internal/infrastructure/logger"
	"mpc/internal/infrastructure/redis"
	"mpc/internal/repository/postgres"
	"mpc/internal/usecase"
)

// @title MPC API
// @version 1.0
// @description This is the API documentation for the MPC project.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// db
	dbPool, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// redis
	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	// logger
	log := logger.NewLogger()

	// jwt
	jwtConfig := auth.NewJWTConfig(cfg.JWT.SecretKey, cfg.JWT.TokenDuration, cfg.JWT.TokenDuration*30)
	jwtService := auth.NewJWTService(jwtConfig, *redisClient)

	// ethereum
	ethClient, err := ethereum.NewEthereumClient(cfg.Ethereum.URL, cfg.Ethereum.SecretKey)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %v", err)
	}

	// repository
	userRepo := postgres.NewUserRepo(dbPool)
	walletRepo := postgres.NewWalletRepo(dbPool)
	transactionRepo := postgres.NewTransactionRepo(dbPool)

	// usecase
	walletUC := usecase.NewWalletUC(walletRepo, ethClient)
	authUC := usecase.NewAuthUC(userRepo, walletUC, *jwtService)
	userUC := usecase.NewUserUC(userRepo)
	txnUC := usecase.NewTxnUC(transactionRepo, ethClient, walletUC, *redisClient)

	// router
	router := http.NewRouter(&userUC, &walletUC, &txnUC, &authUC, jwtService, log)

	log.Fatal(router.Run(":8080"))
}
