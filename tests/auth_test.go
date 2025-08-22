package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"link-shortener/internal/config"
	"link-shortener/internal/database"
	"link-shortener/internal/handlers"
	"link-shortener/internal/middleware"
	"link-shortener/internal/models"
	"link-shortener/internal/repository"
	"link-shortener/internal/services"
	"link-shortener/internal/utils"
)

func setupTestRouter() (*gin.Engine, *database.Database) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load test config
	cfg, _ := config.Load()
	cfg.Database.Name = "link_shortener_test"

	// Initialize database
	db, _ := database.NewDatabase(cfg)
	db.InitTables()

	// Initialize dependencies
	jwtMgr := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiry)
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtMgr)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(jwtMgr)

	// Setup router
	router := gin.Default()
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", authMiddleware.AuthRequired(), authHandler.GetProfile)
		}
	}

	return router, db
}

func TestRegister(t *testing.T) {
	router, db := setupTestRouter()
	defer db.Close()

	tests := []struct {
		name           string
		requestBody    models.RegisterRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid registration",
			requestBody: models.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Invalid email",
			requestBody: models.RegisterRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Short password",
			requestBody: models.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectedError {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "data")
				assert.Contains(t, response["data"], "token")
			}
		})
	}
}

func TestLogin(t *testing.T) {
	router, db := setupTestRouter()
	defer db.Close()

	// First register a user
	registerBody := models.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	registerJSON, _ := json.Marshal(registerBody)
	registerReq, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
	registerReq.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReq)

	tests := []struct {
		name           string
		requestBody    models.LoginRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid login",
			requestBody: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Invalid password",
			requestBody: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "Non-existent user",
			requestBody: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectedError {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "data")
				assert.Contains(t, response["data"], "token")
			}
		})
	}
}
