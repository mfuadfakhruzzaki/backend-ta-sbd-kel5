package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"gorm.io/gorm"
)

// UserRepository adalah interface untuk operasi database pengguna
type UserRepository interface {
	// Create menambahkan pengguna baru
	Create(ctx context.Context, user *domain.User) error
	
	// FindByID mencari pengguna berdasarkan ID
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	
	// FindByEmail mencari pengguna berdasarkan email
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	
	// FindAll mencari semua pengguna dengan paginasi dan filter
	FindAll(ctx context.Context, page, limit int, search string) ([]domain.User, int64, error)
	
	// Update memperbarui data pengguna
	Update(ctx context.Context, user *domain.User) error
	
	// Delete menghapus pengguna (soft delete)
	Delete(ctx context.Context, id uint) error
	
	// HardDelete menghapus pengguna secara permanen
	HardDelete(ctx context.Context, id uint) error
}

// userRepositoryImpl adalah implementasi PostgreSQL dari UserRepository
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Create menambahkan pengguna baru
func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID mencari pengguna berdasarkan ID
func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengguna dengan ID %d tidak ditemukan", id)
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail mencari pengguna berdasarkan email
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pengguna dengan email %s tidak ditemukan", email)
		}
		return nil, err
	}
	return &user, nil
}

// FindAll mencari semua pengguna dengan paginasi dan filter
func (r *userRepositoryImpl) FindAll(ctx context.Context, page, limit int, search string) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.User{})

	// Tambahkan filter pencarian jika ada
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
	}

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update memperbarui data pengguna
func (r *userRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete menghapus pengguna (soft delete)
func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

// HardDelete menghapus pengguna secara permanen
func (r *userRepositoryImpl) HardDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&domain.User{}, id).Error
}