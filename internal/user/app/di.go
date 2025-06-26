package app

import (
	"database/sql"
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/user/app/router"
	"microservices/user-management/internal/user/delivery/v1/auth"
	"microservices/user-management/internal/user/delivery/v1/refresh_token"
	"microservices/user-management/internal/user/infras/rabbitmq"
	"microservices/user-management/internal/user/infras/repo"
	authUC "microservices/user-management/internal/user/usecases/auth"
	rtUC "microservices/user-management/internal/user/usecases/refresh_token"
	"microservices/user-management/pkg/logger"
	"microservices/user-management/pkg/mysql"
)

type Dependencies struct {
	AuthUC           authUC.AuthUseCase
	RefreshTokenUC   rtUC.RefreshTokenUseCase
	AuthCtrl         *auth.AuthController
	RefreshTokenCtrl *refresh_token.RefreshTokenController
	UserGrpcHandler  *router.UserGrpcServer
	Publisher        *rabbitmq.Publisher
	DB               *sql.DB
}

func InitializeDependencies(cfg config.Config) *Dependencies {
	l := logger.GetInstance()
	l.WithName("user-service")

	db, err := mysql.NewDB(cfg.DatabaseDSN)
	if err != nil {
		l.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repo.NewUserRepository(db)
	credRepo := repo.NewCredentialRepository(db)
	rtRepo := repo.NewRefreshTokenRepository(db)
	outboxRepo := repo.NewOutboxRepository(db)
	txManager := repo.NewTxManager(db)

	authUseCase := authUC.NewAuthUseCase(userRepo, credRepo, rtRepo, outboxRepo, txManager)
	refreshTokenUC := rtUC.NewRefreshTokenUseCase(rtRepo, userRepo)

	authCtrl := auth.NewAuthController(authUseCase)
	refreshTokenCtrl := refresh_token.NewRefreshTokenController(refreshTokenUC)

	userGrpcHandler := router.NewUserGrpcServer(authUseCase)

	publisher, err := rabbitmq.NewPublisher(outboxRepo)
	if err != nil {
		l.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}

	return &Dependencies{
		AuthUC:           authUseCase,
		RefreshTokenUC:   refreshTokenUC,
		AuthCtrl:         authCtrl,
		RefreshTokenCtrl: refreshTokenCtrl,
		UserGrpcHandler:  userGrpcHandler,
		Publisher:        publisher,
		DB:               db,
	}
}
