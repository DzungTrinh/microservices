package controllers

import (
	"bytes"
	"microservices/user-management/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	config.DB = gormDB
	return gormDB, mock
}

func TestLogin_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Mock database
	_, mock := setupTestDB(t)
	defer func() {
		sqlDB, err := config.DB.DB()
		assert.NoError(t, err)
		_ = sqlDB.Close()
	}()

	// Mock JWT_SECRET
	os.Setenv("JWT_SECRET", "my-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Mock user query
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow(1, "Admin", "admin@example.com", hashedPassword, "admin")
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\?").
		WithArgs("admin@example.com").
		WillReturnRows(rows)

	// Request
	body := `{"email":"admin@example.com","password":"admin123"}`
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Run
	r.POST("/login", Login)
	r.ServeHTTP(w, c.Request)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"token":"eyJhbGciOiJIUzI1NiIs`)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Mock database
	_, mock := setupTestDB(t)
	defer func() {
		sqlDB, err := config.DB.DB()
		assert.NoError(t, err)
		_ = sqlDB.Close()
	}()

	// Mock user query
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow(1, "Admin", "admin@example.com", hashedPassword, "admin")
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\?").
		WithArgs("admin@example.com").
		WillReturnRows(rows)

	// Request
	body := `{"email":"admin@example.com","password":"wrongpass"}`
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Run
	r.POST("/login", Login)
	r.ServeHTTP(w, c.Request)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	expected := `{"code":401,"error":"INCORRECT_PASSWORD","message":"Incorrect password"}`
	assert.JSONEq(t, expected, w.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}
