package config

import (
	"os"
	"user-management/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Tải file .env
	if err := godotenv.Load(); err != nil {
		panic("Không thể tải file .env")
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		panic("DATABASE_DSN không được thiết lập")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Không thể kết nối đến database")
	}

	//// Tự động tạo bảng
	//database.AutoMigrate(&models.User{})

	// Seed tài khoản admin
	SeedAdmin(database)

	DB = database
}

// SeedAdmin - Tạo tài khoản admin mặc định nếu chưa tồn tại
func SeedAdmin(db *gorm.DB) {
	var admin models.User
	// Kiểm tra xem có người dùng admin nào trong database chưa
	if err := db.Where("role = ?", "admin").First(&admin).Error; err == gorm.ErrRecordNotFound {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			panic("ADMIN_PASSWORD không được thiết lập trong .env")
		}

		admin = models.User{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: adminPassword,
			Role:     "admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			panic("Không thể tạo tài khoản admin")
		}
	}
}
