package usecases

//
//import (
//	"fmt"
//	"github.com/microcosm-cc/bluemonday"
//	"github.com/ugorji/go/codec"
//	"gorm.io/gorm"
//	"log"
//	"microservices/user-management/user-management/models"
//
//	"net/http"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v5"
//	"golang.org/x/crypto/bcrypt"
//)
//
//// UserController defines the controller struct
//type UserController struct {
//	DB *gorm.DB
//}
//
//// NewUserController creates a new controller with dependencies
//func NewUserController(db *gorm.DB) *UserController {
//	return &UserController{DB: db}
//}
//
//// Login handles user authentication
//func (ctrl *UserController) Login(c *gin.Context) {
//	var req struct {
//		Email    string `json:"email" binding:"required,email"`
//		Password string `json:"password" binding:"required,min=6"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		log.Printf("Invalid login input: %v", err)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input format")
//		return
//	}
//
//	// Sanitize email
//	p := bluemonday.StrictPolicy()
//	req.Email = strings.ToLower(p.Sanitize(req.Email))
//
//	var user models.User
//	if err := ctrl.DB.Where("email = ?", req.Email).First(&user).Error; err == gorm.ErrRecordNotFound {
//		log.Printf("Email not found: %s", req.Email)
//		errors.HandleError(c, http.StatusUnauthorized, errors.ErrEmailNotFound, "Email not found")
//		return
//	} else if err != nil {
//		log.Printf("Database error during login: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrFetchUsersFailed, "Internal server error")
//		return
//	}
//
//	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
//		log.Printf("Password mismatch for email %s", req.Email)
//		errors.HandleError(c, http.StatusUnauthorized, errors.ErrIncorrectPassword, "Incorrect password")
//		return
//	}
//
//	token, err := generateJWT(user)
//	if err != nil {
//		log.Printf("Failed to generate JWT: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrTokenGeneration, "Failed to generate token")
//		return
//	}
//
//	var h codec.JsonHandle
//	var b []byte
//	enc := codec.NewEncoderBytes(&b, &h)
//
//	if err := enc.Encode(gin.H{"token": token}); err != nil {
//		log.Printf("Failed to encode token response: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInvalidInput, "Failed to process response")
//		return
//	}
//
//	c.Data(http.StatusOK, "application/json; charset=utf-8", b)
//}
//
//// CreateUser handles user creation
//func (ctrl *UserController) CreateUser(c *gin.Context) {
//	var req struct {
//		Name     string `json:"name" binding:"required"`
//		Email    string `json:"email" binding:"required,email"`
//		Password string `json:"password" binding:"required,min=6"`
//		Role     string `json:"role" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		log.Printf("Invalid create user input: %v", err)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input format")
//		return
//	}
//
//	// Sanitize inputs
//	p := bluemonday.StrictPolicy()
//	req.Name = p.Sanitize(req.Name)
//	req.Email = strings.ToLower(p.Sanitize(req.Email))
//	req.Role = p.Sanitize(req.Role)
//
//	if !models.IsValidRole(req.Role) {
//		log.Printf("Invalid role: %s", req.Role)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidRole, "Role must be 'user' or 'admin'")
//		return
//	}
//
//	user := models.User{
//		Name:     req.Name,
//		Email:    req.Email,
//		Password: req.Password,
//		Role:     models.Role(req.Role),
//	}
//
//	if err := ctrl.DB.Create(&user).Error; err != nil {
//		log.Printf("Failed to create user: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrCreateUserFailed, "Failed to create user")
//		return
//	}
//
//	var h codec.JsonHandle
//	var b []byte
//	enc := codec.NewEncoderBytes(&b, &h)
//
//	if err := enc.Encode(gin.H{
//		"id":    user.ID,
//		"name":  user.Name,
//		"email": user.Email,
//		"role":  user.Role,
//	}); err != nil {
//		log.Printf("Failed to encode create user response: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInvalidInput, "Failed to process response")
//		return
//	}
//
//	c.Data(http.StatusCreated, "application/json; charset=utf-8", b)
//}
//
//// GetUsers retrieves all users
//func (ctrl *UserController) GetUsers(c *gin.Context) {
//	var users []models.User
//	if err := ctrl.DB.Find(&users).Error; err != nil {
//		log.Printf("Failed to fetch users: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrFetchUsersFailed, "Failed to fetch users")
//		return
//	}
//
//	var h codec.JsonHandle
//	var b []byte
//	enc := codec.NewEncoderBytes(&b, &h)
//
//	if err := enc.Encode(users); err != nil {
//		log.Printf("Failed to encode users response: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInvalidInput, "Failed to process response")
//		return
//	}
//
//	c.Data(http.StatusOK, "application/json; charset=utf-8", b)
//}
//
//// GetUserByID retrieves a user by ID
//func (ctrl *UserController) GetUserByID(c *gin.Context) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		log.Printf("Invalid user ID: %v", err)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid user ID")
//		return
//	}
//
//	var user models.User
//	if err := ctrl.DB.First(&user, id).Error; err == gorm.ErrRecordNotFound {
//		log.Printf("User not found: ID=%d", id)
//		errors.HandleError(c, http.StatusNotFound, errors.ErrUserNotFound, "User not found")
//		return
//	} else if err != nil {
//		log.Printf("Failed to fetch user: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrFetchUsersFailed, "Failed to fetch user")
//		return
//	}
//
//	var h codec.JsonHandle
//	var b []byte
//	enc := codec.NewEncoderBytes(&b, &h)
//
//	if err := enc.Encode(gin.H{
//		"id":    user.ID,
//		"name":  user.Name,
//		"email": user.Email,
//		"role":  user.Role,
//	}); err != nil {
//		log.Printf("Failed to encode user response: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInvalidInput, "Failed to process response")
//		return
//	}
//
//	c.Data(http.StatusOK, "application/json; charset=utf-8", b)
//}
//
//// UpdateUser updates a user
//func (ctrl *UserController) UpdateUser(c *gin.Context) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		log.Printf("Invalid user ID: %v", err)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid user ID")
//		return
//	}
//
//	var req struct {
//		Name     string `json:"name"`
//		Email    string `json:"email" binding:"email"`
//		Password string `json:"password" binding:"min=6"`
//		Role     string `json:"role"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		log.Printf("Invalid update user input: %v", err)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid input format")
//		return
//	}
//
//	// Sanitize inputs
//	p := bluemonday.StrictPolicy()
//	req.Name = p.Sanitize(req.Name)
//	req.Email = strings.ToLower(p.Sanitize(req.Email))
//	req.Role = p.Sanitize(req.Role)
//
//	var user models.User
//	if err := ctrl.DB.First(&user, id).Error; err == gorm.ErrRecordNotFound {
//		log.Printf("User not found: ID=%d", id)
//		errors.HandleError(c, http.StatusNotFound, errors.ErrUserNotFound, "User not found")
//		return
//	} else if err != nil {
//		log.Printf("Failed to fetch user: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrFetchUsersFailed, "Failed to fetch user")
//		return
//	}
//
//	if req.Name != "" {
//		user.Name = req.Name
//	}
//	if req.Email != "" {
//		user.Email = req.Email
//	}
//	if req.Password != "" {
//		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
//		if err != nil {
//			log.Printf("Failed to hash password: %v", err)
//			errors.HandleError(c, http.StatusInternalServerError, errors.ErrCreateUserFailed, "Failed to process password")
//			return
//		}
//		user.Password = string(hashedPassword)
//	}
//	if req.Role != "" {
//		if !models.IsValidRole(req.Role) {
//			log.Printf("Invalid role: %s", req.Role)
//			errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidRole, "Role must be 'user' or 'admin'")
//			return
//		}
//		user.Role = models.Role(req.Role)
//	}
//
//	if err := ctrl.DB.Save(&user).Error; err != nil {
//		log.Printf("Failed to update user: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrUpdateUserFailed, "Failed to update user")
//		return
//	}
//
//	var h codec.JsonHandle
//	var b []byte
//	enc := codec.NewEncoderBytes(&b, &h)
//
//	if err := enc.Encode(gin.H{
//		"id":    user.ID,
//		"name":  user.Name,
//		"email": user.Email,
//		"role":  user.Role,
//	}); err != nil {
//		log.Printf("Failed to encode update user response: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInvalidInput, "Failed to process response")
//		return
//	}
//
//	c.Data(http.StatusOK, "application/json; charset=utf-8", b)
//}
//
//// DeleteUser deletes a user
//func (ctrl *UserController) DeleteUser(c *gin.Context) {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		log.Printf("Invalid user ID: %v", err)
//		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid user ID")
//		return
//	}
//
//	var user models.User
//	if err := ctrl.DB.First(&user, id).Error; err == gorm.ErrRecordNotFound {
//		log.Printf("User not found: ID=%d", id)
//		errors.HandleError(c, http.StatusNotFound, errors.ErrUserNotFound, "User not found")
//		return
//	} else if err != nil {
//		log.Printf("Failed to fetch user: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrFetchUsersFailed, "Failed to fetch user")
//		return
//	}
//
//	if err := ctrl.DB.Delete(&user).Error; err != nil {
//		log.Printf("Failed to delete user: %v", err)
//		errors.HandleError(c, http.StatusInternalServerError, errors.ErrDeleteUserFailed, "Failed to delete user")
//		return
//	}
//
//	c.Status(http.StatusNoContent)
//}
//
//func generateJWT(user models.User) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":    user.ID,
//		"email": user.Email,
//		"role":  user.Role,
//		"exp":   time.Now().Add(time.Hour * 24).Unix(),
//	})
//
//	jwtSecret := os.Getenv("JWT_SECRET")
//	if jwtSecret == "" {
//		return "", fmt.Errorf("JWT_SECRET not set")
//	}
//
//	return token.SignedString([]byte(jwtSecret))
//}
