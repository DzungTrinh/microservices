package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Lỗi server nội bộ",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
