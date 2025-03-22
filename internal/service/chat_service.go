package service

import (
	"context"
	"errors"
	"math"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/repository"
)

// ChatService adalah interface untuk layanan chat
type ChatService interface {
	SendMessage(ctx context.Context, chat *domain.Chat, userID uint) (*domain.ChatResponse, error)
	GetByID(ctx context.Context, id uint, userID uint) (*domain.ChatResponse, error)
	GetByBarangID(ctx context.Context, barangID uint, page, limit int, userID uint) ([]domain.ChatResponse, int, int64, error)
	GetConversation(ctx context.Context, pengirimID, penerimaID, barangID uint, page, limit int, userID uint) ([]domain.ChatResponse, int, int64, error)
	GetChatPartners(ctx context.Context, userID uint) ([]domain.UserResponse, error)
	MarkAsRead(ctx context.Context, chatID uint, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
}

// chatService adalah implementasi dari ChatService
type chatService struct {
	chatRepo repository.ChatRepository
	userRepo repository.UserRepository
	itemRepo repository.ItemRepository
}

// NewChatService membuat instance baru dari ChatService
func NewChatService(
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	itemRepo repository.ItemRepository,
) ChatService {
	return &chatService{
		chatRepo: chatRepo,
		userRepo: userRepo,
		itemRepo: itemRepo,
	}
}

// SendMessage mengirim pesan chat baru
func (s *chatService) SendMessage(ctx context.Context, chat *domain.Chat, userID uint) (*domain.ChatResponse, error) {
	// Validasi pengirim
	if chat.PengirimID != userID {
		return nil, errors.New("anda hanya dapat mengirim pesan sebagai diri sendiri")
	}

	// Cek apakah penerima ada
	_, err := s.userRepo.FindByID(ctx, chat.PenerimaID)
	if err != nil {
		return nil, errors.New("penerima tidak ditemukan")
	}

	// Cek apakah barang ada
	item, err := s.itemRepo.FindByID(ctx, chat.BarangID)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	// Validasi chat terkait barang:
	// 1. Penjual dapat chat ke siapapun terkait barangnya
	// 2. Pembeli hanya dapat chat ke penjual
	if item.PenjualID != userID && item.PenjualID != chat.PenerimaID {
		return nil, errors.New("anda hanya dapat mengirim pesan kepada penjual barang ini")
	}

	// Buat chat baru
	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	// Dapatkan chat yang baru dibuat dengan preload
	createdChat, err := s.chatRepo.FindByID(ctx, chat.ID)
	if err != nil {
		return nil, err
	}

	// Kembalikan response
	response := createdChat.ToResponse(true, true, true)
	return &response, nil
}

// GetByID mendapatkan chat berdasarkan ID
func (s *chatService) GetByID(ctx context.Context, id uint, userID uint) (*domain.ChatResponse, error) {
	chat, err := s.chatRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cek apakah pengguna terlibat dalam chat
	if chat.PengirimID != userID && chat.PenerimaID != userID {
		return nil, errors.New("anda tidak memiliki akses ke chat ini")
	}

	response := chat.ToResponse(true, true, true)
	return &response, nil
}

// GetByBarangID mendapatkan chat berdasarkan ID barang
func (s *chatService) GetByBarangID(ctx context.Context, barangID uint, page, limit int, userID uint) ([]domain.ChatResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Dapatkan barang untuk validasi
	item, err := s.itemRepo.FindByID(ctx, barangID)
	if err != nil {
		return nil, 0, 0, errors.New("barang tidak ditemukan")
	}

	// Validasi akses (hanya penjual yang bisa melihat semua chat terkait barangnya)
	if item.PenjualID != userID {
		return nil, 0, 0, errors.New("anda tidak memiliki akses ke chat barang ini")
	}

	// Dapatkan chat dari repository
	chats, total, err := s.chatRepo.FindByBarangID(ctx, barangID, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var chatResponses []domain.ChatResponse
	for _, chat := range chats {
		chatResponses = append(chatResponses, chat.ToResponse(true, true, true))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return chatResponses, totalPages, total, nil
}

// GetConversation mendapatkan percakapan antara dua pengguna untuk barang tertentu
func (s *chatService) GetConversation(ctx context.Context, pengirimID, penerimaID, barangID uint, page, limit int, userID uint) ([]domain.ChatResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20 // Default lebih besar untuk percakapan
	}

	// Validasi pengguna terlibat dalam percakapan
	if pengirimID != userID && penerimaID != userID {
		return nil, 0, 0, errors.New("anda tidak memiliki akses ke percakapan ini")
	}

	// Dapatkan percakapan dari repository
	chats, total, err := s.chatRepo.FindByUserIDs(ctx, pengirimID, penerimaID, barangID, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Mark pesan sebagai dibaca jika pengguna adalah penerima
	for _, chat := range chats {
		if chat.PenerimaID == userID && !chat.Dibaca {
			s.chatRepo.UpdateReadStatus(ctx, chat.ID, true)
		}
	}

	// Konversi ke format respons
	var chatResponses []domain.ChatResponse
	for _, chat := range chats {
		chatResponses = append(chatResponses, chat.ToResponse(true, true, true))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return chatResponses, totalPages, total, nil
}

// GetChatPartners mendapatkan daftar partner chat
func (s *chatService) GetChatPartners(ctx context.Context, userID uint) ([]domain.UserResponse, error) {
	users, err := s.chatRepo.FindChatPartners(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Konversi ke format respons
	var userResponses []domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, nil
}

// MarkAsRead menandai pesan sebagai telah dibaca
func (s *chatService) MarkAsRead(ctx context.Context, chatID uint, userID uint) error {
	// Dapatkan chat yang ada
	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return err
	}

	// Validasi penerima
	if chat.PenerimaID != userID {
		return errors.New("anda tidak dapat menandai pesan ini sebagai dibaca")
	}

	// Update status dibaca
	return s.chatRepo.UpdateReadStatus(ctx, chatID, true)
}

// Delete menghapus chat
func (s *chatService) Delete(ctx context.Context, id uint, userID uint) error {
	// Dapatkan chat yang ada
	chat, err := s.chatRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Validasi pengguna (hanya pengirim yang dapat menghapus)
	if chat.PengirimID != userID {
		return errors.New("anda tidak memiliki izin untuk menghapus pesan ini")
	}

	// Hapus chat
	return s.chatRepo.Delete(ctx, id)
}