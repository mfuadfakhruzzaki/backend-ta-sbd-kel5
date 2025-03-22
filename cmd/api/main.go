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
	"github.com/mfuadfakhruzzaki/jubel/internal/config"
	"github.com/mfuadfakhruzzaki/jubel/internal/database"
	"github.com/mfuadfakhruzzaki/jubel/internal/handler"
	"github.com/mfuadfakhruzzaki/jubel/internal/middleware"
	"github.com/mfuadfakhruzzaki/jubel/internal/repository"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inisialisasi database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Jalankan auto-migrate
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto-migration: %v", err)
	}

	// Inisialisasi router
	router := setupRouter(cfg, db)

	// Jalankan server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

// initDatabase menginisialisasi koneksi database
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Konfigurasi logger GORM
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Output
		logger.Config{
			SlowThreshold:             time.Second,   // Ambang batas query lambat
			LogLevel:                  logger.Info,   // Level log
			IgnoreRecordNotFoundError: true,          // Abaikan error record not found
			Colorful:                  true,          // Gunakan warna
		},
	)

	// Open connection
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// setupRouter menginisialisasi router dan endpoints
func setupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// Mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Middleware global
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Inisialisasi repositories
	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	chatRepo := repository.NewChatRepository(db)

	// Inisialisasi services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	itemService := service.NewItemService(itemRepo, cfg)
	transactionService := service.NewTransactionService(transactionRepo, itemRepo)
	chatService := service.NewChatService(chatRepo, userRepo, itemRepo)

	// Inisialisasi middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Inisialisasi handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	itemHandler := handler.NewItemHandler(itemService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	chatHandler := handler.NewChatHandler(chatService)

	// Setup static file serving
	router.Static("/uploads", cfg.Upload.Dir)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "Server is running",
		})
	})

	// API version
	v1 := router.Group("/api/v1")
	{
		// Register routes
		authHandler.RegisterRoutes(v1)
		userHandler.RegisterRoutes(v1, authMiddleware.VerifyToken(), authMiddleware.RequireAdmin())
		itemHandler.RegisterRoutes(v1, authMiddleware.VerifyToken(), authMiddleware.RequireAdmin())
		transactionHandler.RegisterRoutes(v1, authMiddleware.VerifyToken())
		chatHandler.RegisterRoutes(v1, authMiddleware.VerifyToken())
	}

	return router
}

// gracefulShutdown menangani shutdown server secara graceful
func gracefulShutdown(server *http.Server) {
	// Channel for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}