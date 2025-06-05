package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Xử lý request
		c.Next()

		// Ghi log sau khi request hoàn tất
		duration := time.Since(start)
		log.Printf("[%s] %s | Duration: %v | Status: %d", method, path, duration, c.Writer.Status())
	}
}
