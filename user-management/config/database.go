package config

import (
	"database/sql"
	"log"
	"microservices/user-management/models"
	"os"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set")
	}

	// Connect to database
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Apply migrations
	if err := applyMigrations(dsn); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Remove AutoMigrate (replaced by golang-migrate)
	// Seed admin user
	SeedAdmin(database)

	DB = database
}

func applyMigrations(dsn string) error {
	// Initialize migration driver
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{})
	if err != nil {
		return err
	}

	// Initialize migrate
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations applied successfully")
	return nil
}

// SeedAdmin - Create default admin user if not exists
func SeedAdmin(db *gorm.DB) {
	var admin models.User
	if err := db.Where("role = ?", models.RoleAdmin).First(&admin).Error; err == gorm.ErrRecordNotFound {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			log.Fatal("ADMIN_PASSWORD not set in .env")
		}

		admin = models.User{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: adminPassword,
			Role:     models.RoleAdmin,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}

		log.Printf("Created admin user successfully: email=%s, role=%s", admin.Email, admin.Role)
	} else if err != nil {
		log.Fatalf("Error checking admin user: %v", err)
	} else {
		log.Printf("Admin user already exists: email=%s", admin.Email)
	}
}
