package config

import (
	"log"
	"microservices/user-management/models"
	"os"

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

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate tables
	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate users table: %v", err)
	}

	// Seed admin user
	SeedAdmin(database)

	DB = database
}

// SeedAdmin - Create default admin user if not exists
func SeedAdmin(db *gorm.DB) {
	var admin models.User
	// Check if admin exists
	if err := db.Where("role = ?", models.RoleAdmin).First(&admin).Error; err == gorm.ErrRecordNotFound {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			log.Fatal("ADMIN_PASSWORD not set in .env")
		}

		admin = models.User{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: adminPassword, // Password will be hashed by BeforeCreate hook
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
