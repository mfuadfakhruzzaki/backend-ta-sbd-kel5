package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"gorm.io/gorm"
)

// ChatRepository adalah interface untuk operasi database chat
type ChatRepository interface {
	// Create membuat pesan chat baru
	Create(ctx context.Context, chat *domain.Chat) error
	
	// FindByID mencari chat berdasarkan ID
	FindByID(ctx context.Context, id uint) (*domain.Chat, error)
	
	// FindByBarangID mencari chat berdasarkan ID barang
	FindByBarangID(ctx context.Context, barangID uint, page, limit int) ([]domain.Chat, int64, error)
	
	// FindByUserIDs mencari chat antara dua pengguna untuk barang tertentu
	FindByUserIDs(ctx context.Context, pengirimID, penerimaID, barangID uint, page, limit int) ([]domain.Chat, int64, error)
	
	// FindChatPartners mencari semua partner chat untuk pengguna tertentu
	FindChatPartners(ctx context.Context, userID uint) ([]domain.User, error)
	
	// UpdateReadStatus memperbarui status dibaca untuk chat
	UpdateReadStatus(ctx context.Context, chatID uint, dibaca bool) error
	
	// Delete menghapus chat
	Delete(ctx context.Context, id uint) error
}

// chatRepositoryImpl adalah implementasi PostgreSQL dari ChatRepository
type chatRepositoryImpl struct {
	db *gorm.DB
}

// NewChatRepository membuat instance baru dari ChatRepository
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepositoryImpl{
		db: db,
	}
}

// Create membuat pesan chat baru
func (r *chatRepositoryImpl) Create(ctx context.Context, chat *domain.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

// FindByID mencari chat berdasarkan ID
func (r *chatRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Chat, error) {
	var chat domain.Chat
	if err := r.db.WithContext(ctx).Preload("Pengirim").Preload("Penerima").Preload("Barang").Where("id = ?", id).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat dengan ID %d tidak ditemukan", id)
		}
		return nil, err
	}
	return &chat, nil
}

// FindByBarangID mencari chat berdasarkan ID barang
func (r *chatRepositoryImpl) FindByBarangID(ctx context.Context, barangID uint, page, limit int) ([]domain.Chat, int64, error) {
	var chats []domain.Chat
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.Chat{}).Where("barang_id = ?", barangID)

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Preload("Pengirim").Preload("Penerima").Preload("Barang").
		Offset(offset).Limit(limit).Order("timestamp ASC").
		Find(&chats).Error; err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

// FindByUserIDs mencari chat antara dua pengguna untuk barang tertentu
func (r *chatRepositoryImpl) FindByUserIDs(ctx context.Context, pengirimID, penerimaID, barangID uint, page, limit int) ([]domain.Chat, int64, error) {
	var chats []domain.Chat
	var total int64

	// Hitung offset berdasarkan halaman dan batas
	offset := (page - 1) * limit

	// Buat query dasar
	query := r.db.WithContext(ctx).Model(&domain.Chat{}).
		Where(
			"(pengirim_id = ? AND penerima_id = ?) OR (pengirim_id = ? AND penerima_id = ?)",
			pengirimID, penerimaID, penerimaID, pengirimID,
		).
		Where("barang_id = ?", barangID)

	// Hitung total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jalankan query dengan paginasi
	if err := query.Preload("Pengirim").Preload("Penerima").Preload("Barang").
		Offset(offset).Limit(limit).Order("timestamp ASC").
		Find(&chats).Error; err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

// FindChatPartners mencari semua partner chat untuk pengguna tertentu
func (r *chatRepositoryImpl) FindChatPartners(ctx context.Context, userID uint) ([]domain.User, error) {
	var users []domain.User

	// Mencari semua pengguna yang pernah berinteraksi dengan pengguna ini
	// Buat subquery untuk pengirim
	query := r.db.WithContext(ctx).Table("pengguna").
		Where("id IN (?)",
			r.db.Table("chat").Select("DISTINCT penerima_id").Where("pengirim_id = ?", userID),
		).
		Or("id IN (?)",
			r.db.Table("chat").Select("DISTINCT pengirim_id").Where("penerima_id = ?", userID),
		)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateReadStatus memperbarui status dibaca untuk chat
func (r *chatRepositoryImpl) UpdateReadStatus(ctx context.Context, chatID uint, dibaca bool) error {
	return r.db.WithContext(ctx).Model(&domain.Chat{}).Where("id = ?", chatID).Update("dibaca", dibaca).Error
}

// Delete menghapus chat
func (r *chatRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Chat{}, id).Error
}