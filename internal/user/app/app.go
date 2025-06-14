package app

import (
	"context"
	"database/sql"
	"log"
	"time"

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
	router *gin.Engine
	db     *sql.DB
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

	r := gin.Default()

	userRepo := repo.NewUserRepository(db)
	usecase := users.NewUserUsecase(userRepo)
	userServer := router.NewUserServer(usecase)

	r.POST("/api/v1/register", userServer.Register)
	r.POST("/api/v1/login", userServer.Login)
	r.POST("/api/v1/refresh", userServer.Refresh)

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

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := usecase.CleanExpiredTokens(context.Background()); err != nil {
				log.Printf("Failed to clean expired tokens: %v", err)
			}
		}
	}()

	return &App{router: r, db: db}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}

func (a *App) Close() error {
	return a.db.Close()
}
