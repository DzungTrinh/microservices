package app

import (
	"database/sql"
	"log"
	"microservices/user-management/internal/user/app/router"
	"microservices/user-management/internal/user/infras/repo"
	"microservices/user-management/internal/user/usecases/users"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	DatabaseDSN string
	JWTSecret   string
}

type App struct {
	router *gin.Engine
	db     *sql.DB
}

func NewApp(cfg Config) *App {
	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://infras/mysql/migrations", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Initialize dependencies
	userRepo := repo.NewUserRepository(db)
	userUsecase := users.NewUserUsecase(userRepo, cfg.JWTSecret)
	userHandler := router.NewUserServer(userUsecase)

	// Set up router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.GET("/user", userHandler.GetUser) // Requires JWT middleware
	}

	return &App{router: router, db: db}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}

func (a *App) Close() error {
	return a.db.Close()
}
