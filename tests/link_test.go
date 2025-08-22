package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"link-shortener/internal/config"
	"link-shortener/internal/database"
	"link-shortener/internal/handlers"
	"link-shortener/internal/middleware"
	"link-shortener/internal/models"
	"link-shortener/internal/repository"
	"link-shortener/internal/services"
)

func setupLinkTestRouter() (*gin.Engine, *handlers.LinkHandler, uuid.UUID) {
	// Load test config
	cfg, _ := config.Load()
	
	// Initialize test database (you might want to use a test database)
	db, _ := database.NewDatabase(cfg)
	
	// Initialize repositories
	linkRepo := repository.NewLinkRepository(db)
	
	// Initialize services
	linkService := services.NewLinkService(linkRepo, "http://localhost:8080")
	
	// Initialize handlers
	linkHandler := handlers.NewLinkHandler(linkService)
	
	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add middleware
	router.Use(middleware.CORS())
	
	// Setup routes
	api := router.Group("/api")
	links := api.Group("/links")
	{
		links.POST("/", linkHandler.CreateLink)
		links.GET("/", linkHandler.GetLinks)
		links.GET("/:id", linkHandler.GetLink)
		links.PUT("/:id", linkHandler.UpdateLink)
		links.DELETE("/:id", linkHandler.DeleteLink)
		links.GET("/stats", linkHandler.GetStats)
	}
	
	// Redirect route
	router.GET("/r/:shortCode", linkHandler.Redirect)
	
	// Create a test user ID
	testUserID := uuid.New()
	
	return router, linkHandler, testUserID
}

func TestCreateLink(t *testing.T) {
	router, _, testUserID := setupLinkTestRouter()
	
	// Test data
	createLinkReq := models.CreateLinkRequest{
		OriginalURL: "https://example.com/very-long-url",
		CustomAlias: "test-link",
		Title:       "Test Link",
	}
	
	reqBody, _ := json.Marshal(createLinkReq)
	
	// Create request
	req, _ := http.NewRequest("POST", "/api/links", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	
	// Set user ID in context (simulating auth middleware)
	req.Header.Set("X-Test-User-ID", testUserID.String())
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "Link created successfully", response["message"])
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "https://example.com/very-long-url", data["original_url"])
	assert.Equal(t, "test-link", data["short_code"])
}

func TestGetLinks(t *testing.T) {
	router, _, testUserID := setupLinkTestRouter()
	
	// Create request
	req, _ := http.NewRequest("GET", "/api/links", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-Test-User-ID", testUserID.String())
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Should have data array
	assert.NotNil(t, response["data"])
}

func TestUpdateLink(t *testing.T) {
	router, _, testUserID := setupLinkTestRouter()
	
	// Test data
	updateLinkReq := models.UpdateLinkRequest{
		OriginalURL: "https://updated-example.com",
		Title:       "Updated Test Link",
	}
	
	reqBody, _ := json.Marshal(updateLinkReq)
	
	// Create request with a test link ID
	testLinkID := uuid.New()
	req, _ := http.NewRequest("PUT", "/api/links/"+testLinkID.String(), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-Test-User-ID", testUserID.String())
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Should return 400 since link doesn't exist
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLink(t *testing.T) {
	router, _, testUserID := setupLinkTestRouter()
	
	// Create request with a test link ID
	testLinkID := uuid.New()
	req, _ := http.NewRequest("DELETE", "/api/links/"+testLinkID.String(), nil)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-Test-User-ID", testUserID.String())
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Should return 400 since link doesn't exist
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetStats(t *testing.T) {
	router, _, testUserID := setupLinkTestRouter()
	
	// Create request
	req, _ := http.NewRequest("GET", "/api/links/stats", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-Test-User-ID", testUserID.String())
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Should have data
	assert.NotNil(t, response["data"])
}

func TestRedirect(t *testing.T) {
	router, _, _ := setupLinkTestRouter()
	
	// Create request
	req, _ := http.NewRequest("GET", "/r/test-short-code", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Should return 404 since short code doesn't exist
	assert.Equal(t, http.StatusNotFound, w.Code)
}
