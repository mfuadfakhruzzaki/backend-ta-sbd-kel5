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