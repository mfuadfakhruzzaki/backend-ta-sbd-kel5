package service

import (
	"context"
	"math"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/repository"
)

// UserService adalah interface untuk layanan pengguna
type UserService interface {
	GetByID(ctx context.Context, id uint) (*domain.UserResponse, error)
	GetAll(ctx context.Context, page, limit int, search string) ([]domain.UserResponse, int, int64, error)
	Update(ctx context.Context, id uint, user *domain.User) (*domain.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	HardDelete(ctx context.Context, id uint) error
}

// userService adalah implementasi dari UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService membuat instance baru dari UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetByID mendapatkan pengguna berdasarkan ID
func (s *userService) GetByID(ctx context.Context, id uint) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userResponse := user.ToResponse()
	return &userResponse, nil
}

// GetAll mendapatkan semua pengguna dengan paginasi dan filter
func (s *userService) GetAll(ctx context.Context, page, limit int, search string) ([]domain.UserResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Dapatkan pengguna dari repository
	users, total, err := s.userRepo.FindAll(ctx, page, limit, search)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var userResponses []domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return userResponses, totalPages, total, nil
}

// Update memperbarui data pengguna
func (s *userService) Update(ctx context.Context, id uint, userData *domain.User) (*domain.UserResponse, error) {
	// Dapatkan pengguna yang ada
	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update data pengguna
	if userData.Nama != "" {
		existingUser.Nama = userData.Nama
	}
	if userData.NoHP != "" {
		existingUser.NoHP = userData.NoHP
	}
	if userData.Alamat != "" {
		existingUser.Alamat = userData.Alamat
	}
	if userData.Password != "" {
		existingUser.Password = userData.Password // Simpan password yang belum di-hash
		if err := existingUser.HashPassword(); err != nil {
			return nil, err
		}
	}

	// Role tidak bisa diubah melalui endpoint ini
	// Gunakan endpoint admin untuk mengubah role

	// Simpan perubahan
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	userResponse := existingUser.ToResponse()
	return &userResponse, nil
}

// Delete menghapus pengguna (soft delete)
func (s *userService) Delete(ctx context.Context, id uint) error {
	// Cek apakah pengguna ada
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, id)
}

// HardDelete menghapus pengguna secara permanen
func (s *userService) HardDelete(ctx context.Context, id uint) error {
	// Cek apakah pengguna ada
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepo.HardDelete(ctx, id)
}