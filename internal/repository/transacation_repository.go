package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"gorm.io/gorm"
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

// transactionRepositoryImpl adalah implementasi PostgreSQL dari TransactionRepository
type transactionRepositoryImpl struct {
	db *gorm.DB
}

// NewTransactionRepository membuat instance baru dari TransactionRepository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
}

// Create membuat transaksi baru
func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// FindByID mencari transaksi berdasarkan ID
func (r *transactionRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.WithContext(ctx).Preload("Barang").Preload("Barang.Penjual").Preload("Pembeli").Where("id = ?", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaksi dengan ID %d tidak ditemukan", id)
		}
		return nil, err
	}
	return &transaction, nil
}

// FindAll mencari semua transaksi dengan paginasi
func (r *transactionRepositoryImpl) FindAll(ctx context.Context, page, limit int) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.Transaction{})

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Preload("Barang").Preload("Pembeli").Offset(offset).Limit(limit).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindByPembeliID mencari transaksi berdasarkan ID pembeli
func (r *transactionRepositoryImpl) FindByPembeliID(ctx context.Context, pembeliID uint, page, limit int) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.Transaction{}).Where("pembeli_id = ?", pembeliID)

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Preload("Barang").Preload("Barang.Penjual").Offset(offset).Limit(limit).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindByPenjualID mencari transaksi berdasarkan ID penjual (melalui barang)
func (r *transactionRepositoryImpl) FindByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar dengan join ke tabel barang
	query := r.db.WithContext(ctx).Table("transaksi").
		Joins("JOIN barang ON transaksi.barang_id = barang.id").
		Where("barang.penjual_id = ?", penjualID)

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := r.db.WithContext(ctx).
		Preload("Barang").
		Preload("Pembeli").
		Joins("JOIN barang ON transaksi.barang_id = barang.id").
		Where("barang.penjual_id = ?", penjualID).
		Offset(offset).
		Limit(limit).
		Order("transaksi.created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindByBarangID mencari transaksi berdasarkan ID barang
func (r *transactionRepositoryImpl) FindByBarangID(ctx context.Context, barangID uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Barang").
		Preload("Pembeli").
		Where("barang_id = ?", barangID).
		First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ada transaksi untuk barang ini
		}
		return nil, err
	}
	return &transaction, nil
}

// Update memperbarui data transaksi
func (r *transactionRepositoryImpl) Update(ctx context.Context, transaction *domain.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// UpdateStatus memperbarui status transaksi
func (r *transactionRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status domain.TransactionStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Transaction{}).Where("id = ?", id).Update("status_transaksi", status).Error
}

// Delete menghapus transaksi
func (r *transactionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Transaction{}, id).Error
}