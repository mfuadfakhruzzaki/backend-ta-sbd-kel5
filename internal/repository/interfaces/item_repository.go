package interfaces

import (
	"context"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
)

// ItemRepository adalah interface untuk operasi database barang
type ItemRepository interface {
	// Create menambahkan barang baru
	Create(ctx context.Context, item *domain.Item) error
	
	// FindByID mencari barang berdasarkan ID
	FindByID(ctx context.Context, id uint) (*domain.Item, error)
	
	// FindAll mencari semua barang dengan paginasi dan filter
	FindAll(ctx context.Context, page, limit int, search string, kategori string, status string) ([]domain.Item, int64, error)
	
	// FindByPenjualID mencari barang berdasarkan ID penjual
	FindByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.Item, int64, error)
	
	// Update memperbarui data barang
	Update(ctx context.Context, item *domain.Item) error
	
	// UpdateStatus memperbarui status barang
	UpdateStatus(ctx context.Context, id uint, status domain.ItemStatus) error
	
	// Delete menghapus barang (soft delete)
	Delete(ctx context.Context, id uint) error
	
	// HardDelete menghapus barang secara permanen
	HardDelete(ctx context.Context, id uint) error
}