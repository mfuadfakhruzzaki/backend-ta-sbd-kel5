package main

import (
	"context"
	"fmt"
	"io"
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
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	// Load konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Initialize logger
	logLevel := zerolog.InfoLevel
	if cfg.Server.Env == "development" {
		logLevel = zerolog.DebugLevel
	}
	
	// Configure pretty output for development
	var logOutput io.Writer = os.Stdout
	if cfg.Server.Env == "development" {
		logOutput = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}
	
	// Create root logger
	zerolog.SetGlobalLevel(logLevel)
	rootLogger := zerolog.New(logOutput).With().Timestamp().Caller().Logger()
	
	// Replace standard logger with zerolog
	log.SetFlags(0)
	log.SetOutput(zerolog.NewConsoleWriter())
	
	// Log application startup
	rootLogger.Info().
		Str("environment", cfg.Server.Env).
		Str("port", cfg.Server.Port).
		Msg("Application starting")
	
	// Inisialisasi database
	db, err := initDatabase(cfg, rootLogger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Jalankan auto-migrate
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto-migration: %v", err)
	}

	// Inisialisasi router
	router := setupRouter(cfg, db, rootLogger)

	// Jalankan server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		rootLogger.Info().Str("port", cfg.Server.Port).Msg("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			rootLogger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server, rootLogger)
}

// initDatabase menginisialisasi koneksi database
func initDatabase(cfg *config.Config, logger zerolog.Logger) (*gorm.DB, error) {
	dbLogger := logger.With().Str("component", "database").Logger()
	
	// Konfigurasi logger GORM
	gormLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Output
		gormlogger.Config{
			SlowThreshold:             time.Second,   // Ambang batas query lambat
			LogLevel:                  gormlogger.Info,   // Level log
			IgnoreRecordNotFoundError: true,          // Abaikan error record not found
			Colorful:                  true,          // Gunakan warna
		},
	)

	// Log connection attempt
	dsn := cfg.Database.GetDSN()
	dbLogger.Info().
		Str("host", cfg.Database.Host).
		Str("database", cfg.Database.Name).
		Msg("Connecting to database")

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		dbLogger.Error().Err(err).Msg("Failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		dbLogger.Error().Err(err).Msg("Failed to get DB instance")
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		dbLogger.Error().Err(err).Msg("Failed to ping database")
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbLogger.Info().Msg("Database connection established successfully")
	return db, nil
}

// setupRouter menginisialisasi router dan endpoints
func setupRouter(cfg *config.Config, db *gorm.DB, logger zerolog.Logger) *gin.Engine {
	routerLogger := logger.With().Str("component", "router").Logger()
	
	// Mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New() // Use gin.New() instead of gin.Default() to avoid default middleware

	// Middleware global
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorMiddleware(routerLogger)) // Add error middleware

	// Root path handler
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Jubel API",
			"description": "Aplikasi Jual Beli Barang Bekas",
			"version":     "1.0.0",
			"status":      "running",
		})
	})

	// Inisialisasi repositories
	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	chatRepo := repository.NewChatRepository(db)

	routerLogger.Debug().Msg("Repositories initialized")

	// Inisialisasi services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	itemService := service.NewItemService(itemRepo, cfg)
	transactionService := service.NewTransactionService(transactionRepo, itemRepo)
	chatService := service.NewChatService(chatRepo, userRepo, itemRepo)
	
	routerLogger.Debug().Msg("Services initialized")

	// Inisialisasi middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Inisialisasi handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	itemHandler := handler.NewItemHandler(itemService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	chatHandler := handler.NewChatHandler(chatService)
	
	routerLogger.Debug().Msg("Handlers initialized")

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
	
	routerLogger.Info().Msg("Routes registered successfully")

	return router
}

// gracefulShutdown menangani shutdown server secara graceful
func gracefulShutdown(server *http.Server, logger zerolog.Logger) {
	// Channel for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-quit
	logger.Info().Str("signal", sig.String()).Msg("Shutting down server")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exiting")
}