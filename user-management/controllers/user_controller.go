package controllers

import (
	"log"
	"net/http"
	"os"
	"time"
	"user-management/config"
	"user-management/models"

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

// Login - Đăng nhập và tạo JWT
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Lỗi khi bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Xác thực dữ liệu đầu vào
	if err := validate.Struct(&input); err != nil {
		log.Printf("Lỗi validation: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Email không tồn tại: %s", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email không tồn tại"})
		return
	}

	// In giá trị user.Password từ database để debug
	log.Printf("Mật khẩu đã mã hóa trong database cho %s: %s", input.Email, user.Password)
	// In mật khẩu đầu vào để xác nhận
	log.Printf("Mật khẩu đầu vào cho %s: %s", input.Email, input.Password)

	// Kiểm tra password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("Mật khẩu không đúng cho email %s, hash trong DB: %s, mật khẩu nhập: %s, lỗi: %v", input.Email, user.Password, input.Password, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mật khẩu không đúng"})
		return
	}

	// Tạo JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Lỗi khi tạo token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	log.Printf("Đăng nhập thành công cho email: %s", input.Email)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// CreateUser - Tạo người dùng mới (chỉ admin đã đăng nhập được phép gọi, yêu cầu JWT)
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Lỗi khi bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Xác thực dữ liệu
	if err := validate.Struct(&user); err != nil {
		log.Printf("Lỗi validation: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("Lỗi khi tạo người dùng: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo người dùng"})
		return
	}

	// Không trả về password trong phản hồi
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// GetUsers - Lấy danh sách tất cả người dùng
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		log.Printf("Lỗi khi lấy danh sách người dùng: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách người dùng"})
		return
	}

	// Không trả về password trong phản hồi
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID - Lấy thông tin người dùng theo ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		log.Printf("Lỗi khi tìm người dùng ID %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}

	// Không trả về password
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUser - Cập nhật thông tin người dùng
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		log.Printf("Lỗi khi tìm người dùng ID %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Lỗi khi bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Xác thực dữ liệu (không validate password nếu không thay đổi)
	if err := validate.StructPartial(&input, "Name", "Email"); err != nil {
		log.Printf("Lỗi validation: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cập nhật các trường
	config.DB.Model(&user).Updates(map[string]interface{}{
		"Name":  input.Name,
		"Email": input.Email,
	})

	// Không trả về password
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser - Xóa người dùng
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		log.Printf("Lỗi khi tìm người dùng ID %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		log.Printf("Lỗi khi xóa người dùng ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa người dùng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa người dùng thành công"})
}
