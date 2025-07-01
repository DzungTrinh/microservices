package main

import (
	"github.com/gin-gonic/gin"
	"microservices/user-management/cmd/api-gateway/config"
	"microservices/user-management/internal/pkg/middlewares"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func reverseProxy(targetHost string) gin.HandlerFunc {
	target, _ := url.Parse(targetHost)
	proxy := httputil.NewSingleHostReverseProxy(target)

	//remove duplicate CORS headers
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		return nil
	}

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		logger.GetInstance().Printf("Forwarding to: %s%s", req.URL.Host, req.URL.Path)
	}

	return func(c *gin.Context) {
		cb := getBreaker(targetHost)
		_, err := cb.Execute(func() (interface{}, error) {
			proxy.ServeHTTP(c.Writer, c.Request)
			return nil, nil
		})

		if err != nil {
			logger.GetInstance().Errorf("Circuit breaker triggered for %s: %v", targetHost, err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			c.Abort()
		}
	}
}

func main() {
	cfg := config.GetInstance()
	port := cfg.Port

	userServiceURL := cfg.UserHttpPort
	rbacServiceURL := cfg.RbacHttpPort

	router := gin.Default()
	router.Use(gin.Recovery(), middlewares.CORS())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test endpoint")
	})

	// Serve Swagger JSON files
	router.GET("/swagger/user.json", func(c *gin.Context) {
		c.File("third_party/OpenAPI/v1/user.swagger.json")
	})
	router.GET("/swagger/rbac.json", func(c *gin.Context) {
		c.File("third_party/OpenAPI/v1/rbac.swagger.json")
	})

	// User service endpoints
	userGroup := router.Group("/api/v1/user")
	{
		userGroup.POST("/register", reverseProxy(userServiceURL))
		userGroup.POST("/login", reverseProxy(userServiceURL))
		userGroup.POST("/refresh", reverseProxy(userServiceURL))

		userGroup.Use(middlewares.JWTMiddleware())
		userGroup.GET("/users", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(userServiceURL))
	}

	rbacGroup := router.Group("/api/v1/rbac")
	{
		rbacGroup.Use(middlewares.JWTMiddleware())
		rbacGroup.POST("/roles", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/roles/:id", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/roles", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.PUT("/roles/:id", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/roles/:id", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/permissions", middlewares.RequirePermission(constants.PermissionManagePermissions), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/permissions/:id", middlewares.RequirePermission(constants.PermissionManagePermissions), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/user-roles", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/role-permissions", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.POST("/user-permissions", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/roles/:id/permissions", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/permissions", middlewares.RequirePermission(constants.PermissionManagePermissions), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/users/:id/permissions", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(rbacServiceURL))
		rbacGroup.GET("/users/:id/roles", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/role-permissions/:id/:id2", middlewares.RequirePermission(constants.PermissionManageRoles), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/user-roles/:id/:id2", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(rbacServiceURL))
		rbacGroup.DELETE("/user-permissions/:id/:id2", middlewares.RequirePermission(constants.PermissionManageUsers), reverseProxy(rbacServiceURL))
	}

	if err := router.Run(":" + cfg.Port); err != nil {
		logger.GetInstance().Fatalf("Failed to run server: %v", err)
	}
	logger.GetInstance().Printf("API Gateway running on port %s", port)
}
