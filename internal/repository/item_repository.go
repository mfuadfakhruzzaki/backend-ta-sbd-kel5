package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"gorm.io/gorm"
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

// itemRepositoryImpl adalah implementasi PostgreSQL dari ItemRepository
type itemRepositoryImpl struct {
	db *gorm.DB
}

// NewItemRepository membuat instance baru dari ItemRepository
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepositoryImpl{
		db: db,
	}
}

// Create menambahkan barang baru
func (r *itemRepositoryImpl) Create(ctx context.Context, item *domain.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// FindByID mencari barang berdasarkan ID
func (r *itemRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Item, error) {
	var item domain.Item
	if err := r.db.WithContext(ctx).Preload("Penjual").Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("barang dengan ID %d tidak ditemukan", id)
		}
		return nil, err
	}
	return &item, nil
}

// FindAll mencari semua barang dengan paginasi dan filter
func (r *itemRepositoryImpl) FindAll(ctx context.Context, page, limit int, search string, kategori string, status string) ([]domain.Item, int64, error) {
	var items []domain.Item
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.Item{}).Preload("Penjual")

	// Tambahkan filter pencarian jika ada
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("nama_barang ILIKE ? OR deskripsi ILIKE ?", searchQuery, searchQuery)
	}

	// Filter berdasarkan kategori
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	// Filter berdasarkan status
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// Default menampilkan barang yang tersedia
		query = query.Where("status = ?", domain.StatusTersedia)
	}

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByPenjualID mencari barang berdasarkan ID penjual
func (r *itemRepositoryImpl) FindByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.Item, int64, error) {
	var items []domain.Item
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.Item{}).Where("penjual_id = ?", penjualID)

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Update memperbarui data barang
func (r *itemRepositoryImpl) Update(ctx context.Context, item *domain.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// UpdateStatus memperbarui status barang
func (r *itemRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status domain.ItemStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Item{}).Where("id = ?", id).Update("status", status).Error
}

// Delete menghapus barang (soft delete)
func (r *itemRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Item{}, id).Error
}

// HardDelete menghapus barang secara permanen
func (r *itemRepositoryImpl) HardDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&domain.Item{}, id).Error
}