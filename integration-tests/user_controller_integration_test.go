package integration_test

//import (
//	"bytes"
//	"context"
//	"github.com/gin-gonic/gin"
//	"microservices/user-management/config"
//	"microservices/user-management/models"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//	"user-management/helpers/test"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/testcontainers/testcontainers-go"
//	"github.com/testcontainers/testcontainers-go/modules/mysql"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//func setupMySQLContainer(t *testing.T) (*mysql.MySQLContainer, string) {
//	ctx := context.Background()
//	mysqlContainer, err := mysql.RunContainer(ctx,
//		testcontainers.WithImage("mysql:8.0"),
//		mysql.WithDatabase("user_management"),
//		mysql.WithUsername("testuser"),
//		mysql.WithPassword("testpass"),
//	)
//	if err != nil {
//		t.Fatalf("Failed to start MySQL container: %v", err)
//	}
//
//	dsn, err := mysqlContainer.ConnectionString(ctx, "charset=utf8mb4&parseTime=True&loc=Local")
//	if err != nil {
//		t.Fatalf("Failed to get DSN: %v", err)
//	}
//
//	return mysqlContainer, dsn
//}
//
//func TestLoginAndCreateUser(t *testing.T) {
//	tests := []struct {
//		name           string
//		parallel       bool
//		endpoint       string
//		method         string
//		body           string
//		headers        map[string]string
//		setup          func(t *testing.T, setup *test.TestSetup) string // Returns token if needed
//		expectedStatus int
//		expectedBody   string
//	}{
//		{
//			name:           "Login Success",
//			parallel:       true,
//			endpoint:       "/login",
//			method:         "POST",
//			body:           `{"email":"admin@example.com","password":"admin123"}`,
//			headers:        map[string]string{"Content-Type": "application/json"},
//			setup:          func(t *testing.T, setup *test.TestSetup) string { return "" },
//			expectedStatus: http.StatusOK,
//			expectedBody:   `{"token":"eyJhbGciOiJIUzI1NiIs`,
//		},
//		{
//			name:           "Login Invalid Credentials",
//			parallel:       true,
//			endpoint:       "/login",
//			method:         "POST",
//			body:           `{"email":"admin@example.com","password":"wrongpass"}`,
//			headers:        map[string]string{"Content-Type": "application/json"},
//			setup:          func(t *testing.T, setup *test.TestSetup) string { return "" },
//			expectedStatus: http.StatusUnauthorized,
//			expectedBody:   `{"code":401,"error":"INCORRECT_PASSWORD","message":"Incorrect password"}`,
//		},
//		{
//			name:     "Create User Success",
//			parallel: false, // Sequential due to DB state
//			endpoint: "/users",
//			method:   "POST",
//			body:     `{"name":"Test User","email":"test@example.com","password":"secure123","role":"user"}`,
//			headers:  map[string]string{"Content-Type": "application/json"},
//			setup: func(t *testing.T, setup *test.TestSetup) string {
//				// Login to get token
//				w := httptest.NewRecorder()
//				body := `{"email":"admin@example.com","password":"admin123"}`
//				req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
//				req.Header.Set("Content-Type", "application/json")
//				setup.Router.ServeHTTP(w, req)
//				assert.Equal(t, http.StatusOK, w.Code)
//
//				var loginResp struct{ Token string }
//				assert.NoError(t, gin.H{"token": &loginResp.Token}.BindJSON(bytes.NewReader(w.Body.Bytes())))
//				return loginResp.Token
//			},
//			expectedStatus: http.StatusCreated,
//			expectedBody:   `{"id":2,"name":"Test User","email":"test@example.com","role":"user"}`,
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt // Capture range variable
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.parallel {
//				t.Parallel()
//			}
//
//			// Setup MySQL container
//			mysqlContainer, dsn := setupMySQLContainer(t)
//			defer func() {
//				if err := mysqlContainer.Terminate(context.Background()); err != nil {
//					t.Errorf("Failed to terminate container: %v", err)
//				}
//			}()
//
//			// Initialize GORM with container DSN
//			gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//			if err != nil {
//				t.Fatalf("Failed to open GORM DB: %v", err)
//			}
//			config.DB = gormDB
//			defer func() {
//				sqlDB, _ := gormDB.DB()
//				sqlDB.Close()
//			}()
//
//			// Migrate schema and seed admin
//			err = gormDB.AutoMigrate(&models.User{})
//			if err != nil {
//				t.Fatalf("Failed to migrate schema: %v", err)
//			}
//			os.Setenv("JWT_SECRET", "my-secret-key")
//			os.Setenv("ADMIN_PASSWORD", "admin123")
//			config.SeedAdmin(gormDB)
//
//			// Setup test
//			setup := test.SetupIntegrationTest(t)
//			setup.GormDB = gormDB // Override SQLite with MySQL
//
//			// Custom setup (e.g., get token)
//			token := tt.setup(t, setup)
//
//			// Request
//			req, err := http.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer([]byte(tt.body)))
//			assert.NoError(t, err)
//			for k, v := range tt.headers {
//				req.Header.Set(k, v)
//			}
//			if token != "" {
//				req.Header.Set("Authorization", "Bearer "+token)
//			}
//			setup.Context.Request = req
//
//			// Run
//			setup.Router.ServeHTTP(setup.Recorder, req)
//
//			// Assertions
//			assert.Equal(t, tt.expectedStatus, setup.Recorder.Code)
//			if tt.expectedBody[0] == '{' {
//				assert.JSONEq(t, tt.expectedBody, setup.Recorder.Body.String())
//			} else {
//				assert.Contains(t, setup.Recorder.Body.String(), tt.expectedBody)
//			}
//		})
//	}
//}
