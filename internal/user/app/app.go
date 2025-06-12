package app

import (
	"database/sql"
	"log"
	"microservices/user-management/cmd/user/config"
	"microservices/user-management/internal/pkg/auth"
	"time"

	"microservices/user-management/internal/user/app/router"
	"microservices/user-management/internal/user/infras/repo"
	"microservices/user-management/internal/user/usecases/users"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
		log.Fatal(err)
	}

	r := gin.Default()

	userRepo := repo.NewUserRepository(db)
	usecase := users.NewUserUsecase(userRepo, "5f3d4923ba202dad5036098efa1fe856f2bb9492063eb978571bcbb4fd934edd")
	userServer := router.NewUserServer(usecase)

	// Public routes
	r.POST("/api/v1/register", userServer.Register)
	r.POST("/api/v1/login", userServer.Login)

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(auth.JWTVerifyMiddleware())
	{
		protected.GET("/users", userServer.GetAllUsers)
		protected.GET("/users/:id", userServer.GetUserByID)
	}

	return &App{router: r, db: db}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}

func (a *App) Close() error {
	return a.db.Close()
}
