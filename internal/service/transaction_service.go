package service

import (
	"context"
	"errors"
	"math"

	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/repository"
)

// TransactionService adalah interface untuk layanan transaksi
type TransactionService interface {
	Create(ctx context.Context, transaction *domain.Transaction, userID uint) (*domain.TransactionResponse, error)
	GetByID(ctx context.Context, id uint, userID uint) (*domain.TransactionResponse, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.TransactionResponse, int, int64, error)
	GetByPembeliID(ctx context.Context, pembeliID uint, page, limit int) ([]domain.TransactionResponse, int, int64, error)
	GetByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.TransactionResponse, int, int64, error)
	UpdateStatus(ctx context.Context, id uint, status domain.TransactionStatus, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
}

// transactionService adalah implementasi dari TransactionService
type transactionService struct {
	transactionRepo repository.TransactionRepository
	itemRepo        repository.ItemRepository
}

// NewTransactionService membuat instance baru dari TransactionService
func NewTransactionService(
	transactionRepo repository.TransactionRepository, 
	itemRepo repository.ItemRepository,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		itemRepo:        itemRepo,
	}
}

// Create membuat transaksi baru
func (s *transactionService) Create(ctx context.Context, transaction *domain.Transaction, userID uint) (*domain.TransactionResponse, error) {
	// Validasi pembeli
	if transaction.PembeliID != userID {
		return nil, errors.New("anda hanya dapat membuat transaksi untuk diri sendiri")
	}

	// Cek apakah barang ada dan tersedia
	item, err := s.itemRepo.FindByID(ctx, transaction.BarangID)
	if err != nil {
		return nil, err
	}

	// Cek status barang
	if item.Status != domain.StatusTersedia {
		return nil, errors.New("barang tidak tersedia untuk dibeli")
	}

	// Cek apakah pembeli adalah penjual
	if item.PenjualID == userID {
		return nil, errors.New("anda tidak dapat membeli barang anda sendiri")
	}

	// Cek apakah barang sudah ada transaksi pending
	existingTransaction, err := s.transactionRepo.FindByBarangID(ctx, transaction.BarangID)
	if err != nil {
		return nil, err
	}
	if existingTransaction != nil && existingTransaction.StatusTransaksi == domain.StatusPending {
		return nil, errors.New("barang ini sedang dalam proses transaksi")
	}

	// Set status default
	transaction.StatusTransaksi = domain.StatusPending

	// Buat transaksi
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Update status barang
	if err := s.itemRepo.UpdateStatus(ctx, item.ID, domain.StatusTerjual); err != nil {
		return nil, err
	}

	// Dapatkan transaksi yang baru dibuat dengan preload
	createdTransaction, err := s.transactionRepo.FindByID(ctx, transaction.ID)
	if err != nil {
		return nil, err
	}

	// Kembalikan response
	response := createdTransaction.ToResponse(true, true)
	return &response, nil
}

// GetByID mendapatkan transaksi berdasarkan ID
func (s *transactionService) GetByID(ctx context.Context, id uint, userID uint) (*domain.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cek apakah pengguna terlibat dalam transaksi
	if transaction.PembeliID != userID && transaction.Barang.PenjualID != userID {
		return nil, errors.New("anda tidak memiliki akses ke transaksi ini")
	}

	response := transaction.ToResponse(true, true)
	return &response, nil
}

// GetAll mendapatkan semua transaksi dengan paginasi
func (s *transactionService) GetAll(ctx context.Context, page, limit int) ([]domain.TransactionResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transactions, total, err := s.transactionRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var transactionResponses []domain.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, transaction.ToResponse(true, true))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return transactionResponses, totalPages, total, nil
}

// GetByPembeliID mendapatkan transaksi berdasarkan ID pembeli
func (s *transactionService) GetByPembeliID(ctx context.Context, pembeliID uint, page, limit int) ([]domain.TransactionResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transactions, total, err := s.transactionRepo.FindByPembeliID(ctx, pembeliID, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var transactionResponses []domain.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, transaction.ToResponse(true, false))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return transactionResponses, totalPages, total, nil
}

// GetByPenjualID mendapatkan transaksi berdasarkan ID penjual
func (s *transactionService) GetByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.TransactionResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transactions, total, err := s.transactionRepo.FindByPenjualID(ctx, penjualID, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var transactionResponses []domain.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, transaction.ToResponse(true, true))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return transactionResponses, totalPages, total, nil
}

// UpdateStatus memperbarui status transaksi
func (s *transactionService) UpdateStatus(ctx context.Context, id uint, status domain.TransactionStatus, userID uint) error {
	// Dapatkan transaksi yang ada
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Validasi pengguna
	if transaction.PembeliID != userID && transaction.Barang.PenjualID != userID {
		return errors.New("anda tidak memiliki akses ke transaksi ini")
	}

	// Validasi perubahan status
	if transaction.StatusTransaksi == domain.StatusSelesai || transaction.StatusTransaksi == domain.StatusDibatalkan {
		return errors.New("tidak dapat mengubah status transaksi yang sudah selesai atau dibatalkan")
	}

	// Update status transaksi
	if err := s.transactionRepo.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	// Jika transaksi dibatalkan, kembalikan status barang menjadi tersedia
	if status == domain.StatusDibatalkan {
		if err := s.itemRepo.UpdateStatus(ctx, transaction.BarangID, domain.StatusTersedia); err != nil {
			return err
		}
	}

	return nil
}

// Delete menghapus transaksi
func (s *transactionService) Delete(ctx context.Context, id uint, userID uint) error {
	// Dapatkan transaksi yang ada
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Validasi pengguna (hanya admin atau penjual yang dapat menghapus)
	if transaction.Barang.PenjualID != userID {
		return errors.New("anda tidak memiliki izin untuk menghapus transaksi ini")
	}

	// Hapus transaksi
	return s.transactionRepo.Delete(ctx, id)
}