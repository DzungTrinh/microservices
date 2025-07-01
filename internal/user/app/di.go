package app

import (
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/user/app/router"
	"microservices/user-management/internal/user/delivery/v1/auth"
	"microservices/user-management/internal/user/delivery/v1/refresh_token"
	"microservices/user-management/internal/user/infras/rabbitmq"
	"microservices/user-management/internal/user/infras/repo"
	authUC "microservices/user-management/internal/user/usecases/auth"
	rtUC "microservices/user-management/internal/user/usecases/refresh_token"
	"microservices/user-management/internal/user/usecases/user"
	"microservices/user-management/pkg/logger"
	"microservices/user-management/pkg/mysql"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
)

type Dependencies struct {
	AuthUC           authUC.AuthUseCase
	RefreshTokenUC   rtUC.RefreshTokenUseCase
	UserUC           user.UserUseCase
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

	// Initialize RBAC gRPC client
	conn, err := grpc.NewClient(cfg.RBACGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Fatalf("Failed to connect to RBAC gRPC service: %v", err)
	}
	rbacClient := rbacv1.NewRBACServiceClient(conn)

	userRepo := repo.NewUserRepository(db)
	credRepo := repo.NewCredentialRepository(db)
	rtRepo := repo.NewRefreshTokenRepository(db)
	outboxRepo := repo.NewOutboxRepository(db)
	txManager := repo.NewTxManager(db)

	authUseCase := authUC.NewAuthUseCase(userRepo, credRepo, rtRepo, outboxRepo, txManager, rbacClient)
	refreshTokenUC := rtUC.NewRefreshTokenUseCase(rtRepo, userRepo, rbacClient)
	userUseCase := user.NewUserUseCase(userRepo, credRepo, outboxRepo, txManager, rbacClient)

	authCtrl := auth.NewAuthController(authUseCase)
	refreshTokenCtrl := refresh_token.NewRefreshTokenController(refreshTokenUC)

	userGrpcHandler := router.NewUserGrpcServer(authUseCase, refreshTokenUC, userUseCase)

	publisher, err := rabbitmq.NewPublisher(outboxRepo)
	if err != nil {
		l.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}

	return &Dependencies{
		AuthUC:           authUseCase,
		RefreshTokenUC:   refreshTokenUC,
		UserUC:           userUseCase,
		AuthCtrl:         authCtrl,
		RefreshTokenCtrl: refreshTokenCtrl,
		UserGrpcHandler:  userGrpcHandler,
		Publisher:        publisher,
		DB:               db,
	}
}
