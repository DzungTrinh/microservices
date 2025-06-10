package router

import (
	"github.com/gin-gonic/gin"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/usecases/users"

	"net/http"
)

type UserServer struct {
	usecase *users.UserUsecase
}

func NewUserServer(usecase *users.UserUsecase) *UserServer {
	return &UserServer{usecase: usecase}
}

func (h *UserServer) Register(c *gin.Context) {
	var req domain.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.usecase.Register(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserServer) Login(c *gin.Context) {
	var req domain.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.usecase.Login(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserServer) GetUser(c *gin.Context) {
	id, exists := c.Get("user_id") // Assumes JWT middleware sets user_id
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	resp, err := h.usecase.GetUserByID(c, id.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
