package main

import (
	"github.com/gin-gonic/gin"
	"microservices/user-management/cmd/api-gateway/config"
	"microservices/user-management/internal/pkg/middlewares"
	"microservices/user-management/pkg/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func reverseProxy(targetHost string) gin.HandlerFunc {
	target, _ := url.Parse(targetHost)
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		logger.GetInstance().Printf("Forwarding to: %s%s", req.URL.Host, req.URL.Path)
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	cfg := config.GetInstance()
	port := cfg.Port

	userServiceURL := cfg.UserHttpPort
	rbacServiceURL := cfg.RbacHttpPort

	router := gin.Default()
	router.Use(gin.Recovery(), middlewares.CORS())

	// User service endpoints
	userGroup := router.Group("/api/v1/user")
	{
		userGroup.POST("/register", reverseProxy(userServiceURL))
		userGroup.POST("/login", reverseProxy(userServiceURL))
		userGroup.POST("/refresh", reverseProxy(userServiceURL))
	}

	rbacGroup := router.Group("/api/v1/rbac")
	{
		// Apply auth to all rbac routes
		rbacGroup.Use(middlewares.JWTMiddleware())

		// Admin-only endpoints
		rbacGroup.POST("/roles", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/roles/:id", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/roles", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.PUT("/roles/:id", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/roles/:id", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/permissions", middlewares.RequirePermission("manage_permissions"), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/permissions/:id", middlewares.RequirePermission("manage_permissions"), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/user-roles", middlewares.RequirePermission("manage_users"), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/role-permissions", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/user-permissions", middlewares.RequirePermission("manage_users"), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/roles/:id/permissions", middlewares.RequirePermission("manage_roles"), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/permissions", middlewares.RequirePermission("manage_permissions"), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/users/:id/permissions", middlewares.RequirePermission("manage_users"), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/users/:id/roles", middlewares.RequirePermission("manage_users"), reverseProxy(rbacServiceURL))
	}

	if err := router.Run(":" + cfg.Port); err != nil {
		logger.GetInstance().Fatalf("Failed to run server: %v", err)
	}
	logger.GetInstance().Printf("API Gateway running on port %s", port)
}
