package interfaces

import (
	"context"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
)

// TransactionRepository adalah interface untuk operasi database transaksi
type TransactionRepository interface {
	// Create membuat transaksi baru
	Create(ctx context.Context, transaction *domain.Transaction) error
	
	// FindByID mencari transaksi berdasarkan ID
	FindByID(ctx context.Context, id uint) (*domain.Transaction, error)
	
	// FindAll mencari semua transaksi dengan paginasi
	FindAll(ctx context.Context, page, limit int) ([]domain.Transaction, int64, error)
	
	// FindByPembeliID mencari transaksi berdasarkan ID pembeli
	FindByPembeliID(ctx context.Context, pembeliID uint, page, limit int) ([]domain.Transaction, int64, error)
	
	// FindByPenjualID mencari transaksi berdasarkan ID penjual (melalui barang)
	FindByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.Transaction, int64, error)
	
	// FindByBarangID mencari transaksi berdasarkan ID barang
	FindByBarangID(ctx context.Context, barangID uint) (*domain.Transaction, error)
	
	// Update memperbarui data transaksi
	Update(ctx context.Context, transaction *domain.Transaction) error
	
	// UpdateStatus memperbarui status transaksi
	UpdateStatus(ctx context.Context, id uint, status domain.TransactionStatus) error
	
	// Delete menghapus transaksi
	Delete(ctx context.Context, id uint) error
}