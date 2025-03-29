package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
	Appwrite AppwriteConfig
}

// ServerConfig menyimpan konfigurasi server
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig menyimpan konfigurasi database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig menyimpan konfigurasi JWT
type JWTConfig struct {
	Secret     string
	ExpiryTime time.Duration
}

// UploadConfig menyimpan konfigurasi upload file
type UploadConfig struct {
	Dir          string
	MaxSize      int64
	AllowedTypes []string
}

// AppwriteConfig menyimpan konfigurasi Appwrite
type AppwriteConfig struct {
	Endpoint   string
	ProjectID  string
	APIKey     string
	BucketID   string
}

// LoadConfig memuat konfigurasi dari file .env
func LoadConfig() (*Config, error) {
	// Coba membaca dari file .env terlebih dahulu
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file tidak ditemukan, mencoba dari environment variables")
	}

	// Konfigurasi server
	port := getEnv("APP_PORT", "8080")
	env := getEnv("APP_ENV", "development")

	// Konfigurasi database
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "jual_beli_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Konfigurasi JWT
	jwtSecret := getEnv("JWT_SECRET", "rahasia_jwt_sangat_aman")
	jwtExpiryStr := getEnv("JWT_EXPIRY", "24h")
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("gagal parse JWT_EXPIRY: %v", err)
	}

	// Konfigurasi upload
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")
	maxUploadSizeStr := getEnv("MAX_UPLOAD_SIZE", "5242880") // Default 5MB
	maxUploadSize, err := strconv.ParseInt(maxUploadSizeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("gagal parse MAX_UPLOAD_SIZE: %v", err)
	}

	// Konfigurasi Appwrite
	appwriteEndpoint := getEnv("APPWRITE_ENDPOINT", "http://tasbd-appwrite-38c346-34-101-66-9.traefik.me/v1")
	appwriteProjectID := getEnv("APPWRITE_PROJECT_ID", "67e7bbfb003b2a88a380")
	appwriteAPIKey := getEnv("APPWRITE_API_KEY", "standard_cb368bd976b49d276ae32ee22b541ab37799d406326f72aab5499c9b9a4adebaefd92489e80fab8fa66f255b80ed9c0f8b0829d52b1e67779f584e42c2b0ffff4280abf621e17c588ab1523659d36515d2fc741fba1f2b8d24e4383b3319b69e125dfb38689dbeb90722e86f37a4ce3ec25a305a88d9a18c260ad128404cd27f")
	appwriteBucketID := getEnv("APPWRITE_BUCKET_ID", "67e7bc05000d4dde5eb1")

	// Pastikan direktori upload ada
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat direktori upload: %v", err)
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Env:  env,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
		JWT: JWTConfig{
			Secret:     jwtSecret,
			ExpiryTime: jwtExpiry,
		},
		Upload: UploadConfig{
			Dir:     uploadDir,
			MaxSize: maxUploadSize,
			AllowedTypes: []string{
				"image/jpeg",
				"image/png",
				"image/gif",
			},
		},
		Appwrite: AppwriteConfig{
			Endpoint:   appwriteEndpoint,
			ProjectID:  appwriteProjectID,
			APIKey:     appwriteAPIKey,
			BucketID:   appwriteBucketID,
		},
	}, nil
}

// GetDSN mengembalikan connection string untuk PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// Helper function untuk mendapatkan environment variable dengan nilai default
func getEnv(key, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}