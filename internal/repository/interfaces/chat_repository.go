package interfaces

import (
	"context"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
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