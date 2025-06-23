package app

import (
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "microservices/user-management/proto/gen"
	"net"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/pkg/auth"
	"microservices/user-management/internal/user/app/router"
	"microservices/user-management/internal/user/infras/mysql"
	"microservices/user-management/internal/user/infras/repo"
	"microservices/user-management/internal/user/infras/seed"
	"microservices/user-management/internal/user/usecases/users"
)

type App struct {
	router     *gin.Engine
	db         *sql.DB
	grpcServer *grpc.Server
}

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Database connection pool initialized")
	return db, nil
}

func NewApp(cfg config.Config) *App {
	db, err := NewDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed admin account
	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		if err := seed.SeedAdmin(context.Background(), mysql.New(db), cfg.AdminEmail, cfg.AdminPassword); err != nil {
			log.Fatalf("Failed to seed admin: %v", err)
		}
	}

	// Initialize use case and repository
	userRepo := repo.NewUserRepository(db)
	usecase := users.NewUserUsecase(userRepo)

	// gRPC server setup
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, router.NewUserGrpcServer(usecase))

	// gRPC-Gateway setup
	gwmux := runtime.NewServeMux()
	err = pb.RegisterAuthServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:8082", []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	// Gin router setup
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Mount gRPC-Gateway handlers
	r.Any("/api/v1/register", func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/api/v1/login", func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/api/v1/refresh", func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})

	// Existing REST routes (optional, can be removed if gRPC-Gateway is sufficient)
	userServer := router.NewUserRestServer(usecase)
	protected := r.Group("/api/v1")
	protected.Use(auth.JWTVerifyMiddleware())
	{
		adminProtected := protected.Group("")
		adminProtected.Use(auth.AdminOnlyMiddleware())
		{
			adminProtected.GET("/users", userServer.GetAllUsers)
			adminProtected.PUT("/users/:id/roles", userServer.UpdateUserRoles)
			adminProtected.GET("/users/:id", userServer.GetUserByID)
		}

		protected.GET("/users/me", userServer.GetCurrentUser)
	}

	// Token cleanup goroutine
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := usecase.CleanExpiredTokens(context.Background()); err != nil {
				log.Printf("Failed to clean expired tokens: %v", err)
			}
		}
	}()

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			log.Fatalf("Failed to listen for gRPC: %v", err)
		}
		log.Printf("gRPC server running on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	return &App{
		router:     r,
		grpcServer: grpcServer,
		db:         db,
	}
}

func (a *App) Run(addr string) error {
	log.Printf("HTTP server running on %s", addr)
	return a.router.Run(addr)
}

func (a *App) Close() error {
	a.grpcServer.GracefulStop()
	return a.db.Close()
}
