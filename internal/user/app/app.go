package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"microservices/user-management/internal/pkg/middlewares"
	"microservices/user-management/internal/user/cron"
	"microservices/user-management/pkg/logger"
	userv1 "microservices/user-management/proto/gen/user/v1"
	"time"

	"net"

	"github.com/gin-gonic/gin"
	"microservices/user-management/cmd/user/config"

	"microservices/user-management/internal/user/infras/seed"
)

type App struct {
	router     *gin.Engine
	grpcServer *grpc.Server
	db         *sql.DB
}

func NewApp(cfg config.Config) *App {
	l := logger.GetInstance()
	l.WithName("user-service")

	deps := InitializeDependencies(cfg)

	// Seed admin account
	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		if err := seed.SeedAdmin(context.Background(), deps.UserUC, cfg.AdminEmail, cfg.AdminPassword); err != nil {
			l.Fatalf("Failed to seed admin: %v", err)
		}
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, deps.UserGrpcHandler)

	go func() {
		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			logger.GetInstance().Fatalf("Failed to listen for gRPC: %v", err)
		}
		logger.GetInstance().Infof("gRPC server running on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.GetInstance().Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start outbox worker
	go deps.Publisher.StartOutboxWorker(context.Background(), time.Second)

	// gRPC-Gateway setup
	gwmux := runtime.NewServeMux()
	grpcEndpoint := fmt.Sprintf("0.0.0.0%s", cfg.GRPCPort)
	err := userv1.RegisterUserServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndpoint, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		l.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	// Gin router
	r := gin.Default()
	r.Use(middlewares.CORS())

	// Serve Swagger JSON
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("./third_party/OpenAPI/v1/user.swagger.json")
	})

	// Mount gRPC-Gateway handlers
	r.Any("/api/v1/*path", func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})

	// Start token cleanup cron
	cron.StartTokenCleanup(context.Background(), deps.RefreshTokenUC)

	return &App{
		router:     r,
		grpcServer: grpcServer,
	}
}

func (a *App) Run(addr string) error {
	logger.GetInstance().Printf("HTTP server running on %s", addr)
	return a.router.Run(addr)
}

func (a *App) Close() error {
	a.grpcServer.GracefulStop()
	return a.db.Close()
}
