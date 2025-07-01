package app

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/internal/rbac/app/seed"
	"microservices/user-management/pkg/logger"
	rbacv1 "microservices/user-management/proto/gen/rbac/v1"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	router     *gin.Engine
	grpcServer *grpc.Server
	deps       *Dependencies
	logger     *logger.LoggerService
}

func NewApp(cfg config.Config) *App {
	l := logger.GetInstance()
	l.WithName("rbac-service")

	deps := InitializeDependencies(cfg)

	// Seed roles
	if err := seed.SeedRoles(context.Background(), deps.RoleUC, deps.RolePermUC, deps.PermUC); err != nil {
		l.Errorf("Failed to seed roles: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(
					grpc_recovery.WithRecoveryHandler(func(p any) error {
						log.Printf("panic occurred: %v", p)
						return status.Errorf(codes.Internal, "internal server error")
					}),
				),
			),
		),
	)

	rbacv1.RegisterRBACServiceServer(grpcServer, deps.RBACGrpcHandler)

	go func() {
		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			l.Fatalf("Failed to listen for gRPC: %v", err)
		}
		l.Infof("gRPC server running on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			l.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	gwmux := runtime.NewServeMux()
	grpcEndpoint := fmt.Sprintf("0.0.0.0%s", cfg.GRPCPort)
	err := rbacv1.RegisterRBACServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndpoint, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		l.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	r := gin.Default()
	r.Any("/*path", func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})

	return &App{
		router:     r,
		grpcServer: grpcServer,
		deps:       deps,
		logger:     l,
	}
}

func (a *App) Run(addr string) error {
	a.logger.Infof("HTTP server running on address %s", addr)
	return a.router.Run(addr)
}

func (a *App) Close() error {
	a.logger.Info("Shutting down application")
	err := a.logger.Sync()
	if err != nil {
		return err
	}
	a.grpcServer.GracefulStop()
	return a.deps.DB.Close()
}
