package service

import (
	"context"
	"errors"
	"time"

	"github.com/mfuadfakhruzzaki/jubel/internal/config"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/repository"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// AuthService adalah interface untuk layanan autentikasi
type AuthService interface {
	Register(ctx context.Context, user *domain.User) (uint, error)
	Login(ctx context.Context, email, password string) (string, *domain.UserResponse, error)
	ValidateToken(ctx context.Context, token string) (*utils.JWTClaims, error)
}

// authService adalah implementasi dari AuthService
type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register mendaftarkan pengguna baru
func (s *authService) Register(ctx context.Context, user *domain.User) (uint, error) {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return 0, errors.New("email sudah terdaftar")
	}

	// Set role default ke 'user'
	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return 0, err
	}

	// Simpan pengguna ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Login melakukan autentikasi pengguna
func (s *authService) Login(ctx context.Context, email, password string) (string, *domain.UserResponse, error) {
	// Cari pengguna berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("email atau password salah")
	}

	// Verifikasi password
	if !user.CheckPassword(password) {
		return "", nil, errors.New("email atau password salah")
	}

	// Generate token JWT
	token, err := utils.GenerateToken(user, s.config.JWT.Secret, s.config.JWT.ExpiryTime)
	if err != nil {
		return "", nil, err
	}

	// Kembalikan token dan data pengguna (tanpa password)
	userResponse := user.ToResponse()
	return token, &userResponse, nil
}

// ValidateToken memvalidasi token JWT
func (s *authService) ValidateToken(ctx context.Context, token string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateToken(token, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	// Cek apakah token sudah kedaluwarsa
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token kedaluwarsa")
	}

	// Cek apakah pengguna masih ada
	_, err = s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	return claims, nil
}