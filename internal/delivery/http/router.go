package http

import (
	"mpc/internal/delivery/http/handler"
	"mpc/internal/delivery/http/middleware"
	"mpc/internal/infrastructure/auth"
	"mpc/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(
	userUC *usecase.UserUseCase,
	walletUC *usecase.WalletUseCase,
	txnUC *usecase.TxnUseCase,
	authUC *usecase.AuthUseCase,
	jwtService *auth.JWTService,
	log *logrus.Logger,
) *gin.Engine {

	router := gin.Default()
	// router.Use(middleware.LoggerMiddleware(log))
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(*authUC)
	userHandler := handler.NewUserHandler(userUC)
	walletHandler := handler.NewWalletHandler(walletUC)
	txnHandler := handler.NewTxnHandler(*txnUC)

	v1 := router.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("/", healthHandler.HealthCheck)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
		}

		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(*jwtService))
		{
			users.GET("/:id", userHandler.GetUser)
		}

		wallets := v1.Group("/wallets")
		wallets.Use(middleware.AuthMiddleware(*jwtService))
		{
			wallets.POST("/", walletHandler.CreateWallet)
			wallets.GET("/:id", walletHandler.GetWallet)
		}

		transactions := v1.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(*jwtService))
		{
			transactions.GET("/", txnHandler.GetTransactions)
			transactions.POST("/create", txnHandler.CreateTransaction)
			transactions.POST("/submit", txnHandler.SubmitTransaction)
		}
	}
	return router
}
