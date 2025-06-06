package controllers

import (
	"log"
	"microservices/user-management/config"
	"microservices/user-management/models"
	"microservices/user-management/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Login - Authenticate user and generate JWT
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrInvalidInput, "Invalid input format")
		return
	}

	// Validate input
	if err := validate.Struct(&input); err != nil {
		log.Printf("Validation error: %v", err)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrValidationFailed, "Validation failed: "+err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Email not found: %s", input.Email)
		utils.HandleError(c, http.StatusUnauthorized, utils.ErrEmailNotFound, "Email not found")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("Password mismatch for email %s", input.Email)
		utils.HandleError(c, http.StatusUnauthorized, utils.ErrIncorrectPassword, "Incorrect password")
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrTokenGeneration, "Failed to generate token")
		return
	}

	log.Printf("Login successful for user: %s", input.Email)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// CreateUser - Create a new user (admin only, requires JWT)
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrInvalidInput, "Invalid input format")
		return
	}

	// Validate input
	if err := validate.Struct(&user); err != nil {
		log.Printf("Validation error: %v", err)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrValidationFailed, "Validation failed: "+err.Error())
		return
	}

	// Validate role
	if user.Role != models.RoleUser && user.Role != models.RoleAdmin {
		log.Printf("Invalid role: %s", user.Role)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrInvalidRole, "Role must be 'user' or 'admin'")
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrCreateUserFailed, "Failed to create user")
		return
	}

	// Omit password in response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// GetUsers - List all users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrFetchUsersFailed, "Failed to fetch users")
		return
	}

	// Omit passwords in response
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID - Get user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		log.Printf("Error finding user ID %s: %v", id, err)
		utils.HandleError(c, http.StatusNotFound, utils.ErrUserNotFound, "User not found")
		return
	}

	// Omit password in response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUser - Update user information
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		log.Printf("Error finding user ID %s: %v", id, err)
		utils.HandleError(c, http.StatusNotFound, utils.ErrUserNotFound, "User not found")
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrInvalidInput, "Invalid input format")
		return
	}

	// Validate input (partial)
	if err := validate.StructPartial(&input, "Name", "Email"); err != nil {
		log.Printf("Validation error: %v", err)
		utils.HandleError(c, http.StatusBadRequest, utils.ErrValidationFailed, "Validation failed: "+err.Error())
		return
	}

	// Update fields
	config.DB.Model(&user).Updates(map[string]interface{}{
		"Name":  input.Name,
		"Email": input.Email,
	})

	// Omit password in response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser - Delete a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		log.Printf("Error finding user ID %s: %v", id, err)
		utils.HandleError(c, http.StatusNotFound, utils.ErrUserNotFound, "User not found")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		log.Printf("Error deleting user ID %s: %v", id, err)
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrDeleteUserFailed, "Failed to delete user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
