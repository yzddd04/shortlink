package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"link-shortener/internal/config"
	"link-shortener/internal/database"
	"link-shortener/internal/handlers"
	"link-shortener/internal/middleware"
	"link-shortener/internal/repository"
	"link-shortener/internal/services"
	"link-shortener/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database tables
	if err := db.InitTables(); err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	// Initialize JWT manager
	jwtMgr := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiry)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	linkRepo := repository.NewLinkRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtMgr)
	linkService := services.NewLinkService(linkRepo, fmt.Sprintf("http://localhost:%s", cfg.Server.Port))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	linkHandler := handlers.NewLinkHandler(linkService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtMgr)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute

	// Setup router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(rateLimiter.RateLimit())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Link Shortener API is running",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", authMiddleware.AuthRequired(), authHandler.GetProfile)
		}

		// Link routes (protected)
		links := api.Group("/links")
		links.Use(authMiddleware.AuthRequired())
		{
			links.POST("/", linkHandler.CreateLink)
			links.GET("/", linkHandler.GetLinks)
			links.GET("/stats", linkHandler.GetStats)
			links.GET("/:id", linkHandler.GetLink)
			links.PUT("/:id", linkHandler.UpdateLink)
			links.DELETE("/:id", linkHandler.DeleteLink)
		}
	}

	// Redirect route (public)
	router.GET("/r/:shortCode", linkHandler.Redirect)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
