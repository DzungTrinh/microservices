package app

import (
	"context"
	"fmt"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/internal/pkg/middlewares"
	"microservices/user-management/internal/rbac/infras/rabbitmq"
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

func InterceptorChain() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch info.FullMethod {
		case "/rbac.v1.RBACService/CreateRole",
			"/rbac.v1.RBACService/UpdateRole",
			"/rbac.v1.RBACService/DeleteRole",
			"/rbac.v1.RBACService/CreatePermission",
			"/rbac.v1.RBACService/DeletePermission",
			"/rbac.v1.RBACService/AssignRolesToUser",
			"/rbac.v1.RBACService/AssignPermissionsToRole",
			"/rbac.v1.RBACService/AssignPermissionsToUser":
			authCtx, err := middlewares.JWTVerifyInterceptor(ctx, req, func(c context.Context, _ interface{}) (interface{}, error) {
				return c, nil
			})
			if err != nil {
				return nil, err
			}
			return middlewares.AdminOnlyInterceptor(authCtx.(context.Context), req, handler)

		case "/rbac.v1.RBACService/GetRoleByID",
			"/rbac.v1.RBACService/ListRoles",
			"/rbac.v1.RBACService/ListPermissions",
			"/rbac.v1.RBACService/ListPermissionsForRole",
			"/rbac.v1.RBACService/ListPermissionsForUser",
			"/rbac.v1.RBACService/ListRolesForUser":
			return middlewares.JWTVerifyInterceptor(ctx, req, handler)

		default:
			return handler(ctx, req)
		}
	}
}

func NewApp(cfg config.Config) *App {
	l := logger.GetInstance()
	l.WithName("rbac-service")

	deps := InitializeDependencies(cfg)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(InterceptorChain()))
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

	// Initialize gRPC client for consumer
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0%s", cfg.GRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Fatalf("Failed to dial gRPC: %v", err)
	}

	rbacClient := rbacv1.NewRBACServiceClient(conn)

	// Initialize RabbitMQ consumer
	consumer, err := rabbitmq.NewConsumer(rbacClient)
	if err != nil {
		l.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}

	// Start consumer in background
	go func() {
		if err := consumer.ConsumeEvents(context.Background()); err != nil {
			l.Errorf("Failed to consume events: %v", err)
		}
	}()

	gwmux := runtime.NewServeMux()
	grpcEndpoint := fmt.Sprintf("0.0.0.0%s", cfg.GRPCPort)
	err = rbacv1.RegisterRBACServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndpoint, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		l.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	r := gin.Default()
	r.Any("/api/v1/*path", func(c *gin.Context) {
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
