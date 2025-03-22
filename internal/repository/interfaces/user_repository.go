package interfaces

import (
	"context"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
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