package main

import (
	"github.com/gin-gonic/gin"
	"microservices/user-management/config"
	"microservices/user-management/controllers"
	"microservices/user-management/middleware"
)

func main() {
	// Kết nối database
	config.ConnectDatabase()

	// Khởi tạo Gin router
	r := gin.Default()

	// Áp dụng middleware logging và recovery cho tất cả route
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Route không cần xác thực
	r.POST("/login", controllers.Login)

	// Nhóm route yêu cầu JWT
	protected := r.Group("/users")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.POST("", controllers.CreateUser)
		protected.GET("", controllers.GetUsers)
		protected.GET("/:id", controllers.GetUserByID)
		protected.PUT("/:id", controllers.UpdateUser)
		protected.DELETE("/:id", controllers.DeleteUser)
	}

	// Chạy server
	r.Run(":8080")
}
