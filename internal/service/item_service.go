package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/config"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/repository"
)

// ItemService adalah interface untuk layanan barang
type ItemService interface {
	Create(ctx context.Context, item *domain.Item, userID uint) (*domain.ItemResponse, error)
	GetByID(ctx context.Context, id uint) (*domain.ItemResponse, error)
	GetAll(ctx context.Context, page, limit int, search, kategori, status string) ([]domain.ItemResponse, int, int64, error)
	GetByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.ItemResponse, int, int64, error)
	Update(ctx context.Context, id uint, item *domain.Item, userID uint) (*domain.ItemResponse, error)
	UpdateStatus(ctx context.Context, id uint, status domain.ItemStatus, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	// UploadImage mengembalikan string dalam format "fileID|fileName|viewURL"
	UploadImage(ctx *gin.Context, itemID uint, userID uint) (string, error)
}

// itemService adalah implementasi dari ItemService
type itemService struct {
	itemRepo repository.ItemRepository
	config   *config.Config
}

// NewItemService membuat instance baru dari ItemService
func NewItemService(itemRepo repository.ItemRepository, config *config.Config) ItemService {
	return &itemService{
		itemRepo: itemRepo,
		config:   config,
	}
}

// Create menambahkan barang baru
func (s *itemService) Create(ctx context.Context, item *domain.Item, userID uint) (*domain.ItemResponse, error) {
	// Set penjual ID
	item.PenjualID = userID

	// Set status default
	item.Status = domain.StatusTersedia

	// Buat barang baru
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	// Dapatkan barang yang baru dibuat dengan preload penjual
	createdItem, err := s.itemRepo.FindByID(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	// Kembalikan response
	response := createdItem.ToResponse(true)
	return &response, nil
}

// GetByID mendapatkan barang berdasarkan ID
func (s *itemService) GetByID(ctx context.Context, id uint) (*domain.ItemResponse, error) {
	item, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := item.ToResponse(true)
	return &response, nil
}

// GetAll mendapatkan semua barang dengan paginasi dan filter
func (s *itemService) GetAll(ctx context.Context, page, limit int, search, kategori, status string) ([]domain.ItemResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Dapatkan barang dari repository
	items, total, err := s.itemRepo.FindAll(ctx, page, limit, search, kategori, status)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var itemResponses []domain.ItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse(true))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return itemResponses, totalPages, total, nil
}

// GetByPenjualID mendapatkan barang berdasarkan ID penjual
func (s *itemService) GetByPenjualID(ctx context.Context, penjualID uint, page, limit int) ([]domain.ItemResponse, int, int64, error) {
	// Validasi input paginasi
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Dapatkan barang dari repository
	items, total, err := s.itemRepo.FindByPenjualID(ctx, penjualID, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Konversi ke format respons
	var itemResponses []domain.ItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse(false))
	}

	// Hitung total halaman
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return itemResponses, totalPages, total, nil
}

// Update memperbarui data barang
func (s *itemService) Update(ctx context.Context, id uint, itemData *domain.Item, userID uint) (*domain.ItemResponse, error) {
	// Dapatkan barang yang ada
	existingItem, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cek apakah pengguna adalah pemilik barang
	if existingItem.PenjualID != userID {
		return nil, errors.New("anda tidak memiliki izin untuk mengubah barang ini")
	}

	// Update data barang
	if itemData.NamaBarang != "" {
		existingItem.NamaBarang = itemData.NamaBarang
	}
	if itemData.Harga > 0 {
		existingItem.Harga = itemData.Harga
	}
	if itemData.Kategori != "" {
		existingItem.Kategori = itemData.Kategori
	}
	if itemData.Deskripsi != "" {
		existingItem.Deskripsi = itemData.Deskripsi
	}
	// Gambar tidak diupdate di sini, gunakan endpoint upload gambar

	// Simpan perubahan
	if err := s.itemRepo.Update(ctx, existingItem); err != nil {
		return nil, err
	}

	response := existingItem.ToResponse(true)
	return &response, nil
}

// UpdateStatus memperbarui status barang
func (s *itemService) UpdateStatus(ctx context.Context, id uint, status domain.ItemStatus, userID uint) error {
	// Dapatkan barang yang ada
	existingItem, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Cek apakah pengguna adalah pemilik barang
	if existingItem.PenjualID != userID {
		return errors.New("anda tidak memiliki izin untuk mengubah status barang ini")
	}

	// Update status
	return s.itemRepo.UpdateStatus(ctx, id, status)
}

// Delete menghapus barang (soft delete)
func (s *itemService) Delete(ctx context.Context, id uint, userID uint) error {
	// Dapatkan barang yang ada
	existingItem, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Cek apakah pengguna adalah pemilik barang atau admin
	isAdmin, exists := ctx.Value("isAdmin").(bool)
	if !exists {
		isAdmin = false
	}
	
	if existingItem.PenjualID != userID && !isAdmin {
		return errors.New("anda tidak memiliki izin untuk menghapus barang ini")
	}

	// Hapus barang (soft delete)
	return s.itemRepo.Delete(ctx, id)
}

// UploadImage mengupload gambar untuk barang
func (s *itemService) UploadImage(ctx *gin.Context, itemID uint, userID uint) (string, error) {
	// Dapatkan barang yang ada
	existingItem, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return "", err
	}

	// Cek apakah pengguna adalah pemilik barang
	if existingItem.PenjualID != userID {
		return "", errors.New("anda tidak memiliki izin untuk mengupload gambar barang ini")
	}

	// Dapatkan file dari form
	file, err := ctx.FormFile("gambar")
	if err != nil {
		return "", err
	}

	// Cek ukuran file
	if file.Size > s.config.Upload.MaxSize {
		return "", fmt.Errorf("ukuran file terlalu besar (maksimal %d bytes)", s.config.Upload.MaxSize)
	}

	// Buka file untuk dibaca
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	// Baca file ke dalam byte array
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file: %v", err)
	}

	// Buat nama file unik
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%d%s", itemID, time.Now().Unix(), ext)
	fileID := fmt.Sprintf("item_%d_%d", itemID, time.Now().Unix())

	// Buat endpoint Appwrite
	appwriteEndpoint := s.config.Appwrite.Endpoint
	projectID := s.config.Appwrite.ProjectID
	bucketID := s.config.Appwrite.BucketID
	apiKey := s.config.Appwrite.APIKey
	
	// Tampilkan log untuk debug
	fmt.Printf("Appwrite Config:\n")
	fmt.Printf("  Endpoint: %s\n", appwriteEndpoint)
	fmt.Printf("  Project ID: %s\n", projectID)
	fmt.Printf("  Bucket ID: %s\n", bucketID)
	fmt.Printf("  File ID: %s\n", fileID)
	
	// Cek konfigurasi
	if projectID == "" {
		return "", fmt.Errorf("project ID tidak boleh kosong")
	}
	if bucketID == "" {
		return "", fmt.Errorf("bucket ID tidak boleh kosong")
	}
	if apiKey == "" {
		return "", fmt.Errorf("API key tidak boleh kosong")
	}
	
	// Simpan file ke sistem lokal sementara (sebagai backup)
	tempDir := "./temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return "", fmt.Errorf("gagal membuat direktori temp: %v", err)
		}
	}
	
	// Simpan file ke lokal sementara
	tempFilePath := filepath.Join(tempDir, fileName)
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		return "", fmt.Errorf("gagal menyimpan file sementara: %v", err)
	}
	defer os.Remove(tempFilePath) // Hapus file sementara setelah selesai
	
	// Persiapkan URL upload Appwrite
	uploadURL := fmt.Sprintf("%s/storage/buckets/%s/files", 
		appwriteEndpoint, 
		bucketID)
		
	// Buat body untuk request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Tambahkan fileId field
	if err := writer.WriteField("fileId", fileID); err != nil {
		return "", fmt.Errorf("gagal membuat field fileId: %v", err)
	}
	
	// Tambahkan file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("gagal membuat form file: %v", err)
	}
	
	// Copy file bytes ke form
	if _, err := io.Copy(part, bytes.NewReader(fileBytes)); err != nil {
		return "", fmt.Errorf("gagal menulis file ke form: %v", err)
	}
	
	// Tutup writer multipart
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("gagal menutup writer: %v", err)
	}
	
	// Buat request
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %v", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-Appwrite-Project", projectID)
	req.Header.Add("X-Appwrite-Key", apiKey)
	
	// Print semua header untuk debug
	fmt.Printf("Request Headers:\n")
	for name, values := range req.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}
	
	// Kirim request
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()
	
	// Baca response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca response: %v", err)
	}
	
	// Print response untuk debug
	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(respBody))
	
	// Cek response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("gagal upload file ke Appwrite: HTTP %d - %s", resp.StatusCode, string(respBody))
	}
	
	// Buat URL view untuk gambar
	viewURL := fmt.Sprintf("%s/storage/buckets/%s/files/%s/view", 
		appwriteEndpoint,
		bucketID,
		fileID)
	
	// Simpan URL gambar di database
	existingItem.Gambar = viewURL
	if err := s.itemRepo.Update(ctx, existingItem); err != nil {
		return "", err
	}
	
	// Return format fileID|fileName|viewURL untuk penggunaan di handler
	gambarInfo := fmt.Sprintf("%s|%s|%s", fileID, fileName, viewURL)
	return gambarInfo, nil
}