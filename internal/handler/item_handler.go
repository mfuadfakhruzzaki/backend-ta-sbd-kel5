package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/config"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// ItemHandler menangani endpoint terkait barang
type ItemHandler struct {
	itemService service.ItemService
}

// NewItemHandler membuat instance baru ItemHandler
func NewItemHandler(itemService service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// CreateItem menambahkan barang baru
func (h *ItemHandler) CreateItem(c *gin.Context) {
	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Cek apakah ada file gambar yang dikirim
	var hasImage bool
	_, err := c.FormFile("gambar")
	if err == nil {
		hasImage = true
	}

	var itemData domain.Item
	
	// Jika ada file, gunakan FormValue untuk membaca data lainnya
	if hasImage {
		// Binding dari form
		itemData.NamaBarang = c.PostForm("nama_barang")
		hargaStr := c.PostForm("harga")
		harga, err := strconv.ParseFloat(hargaStr, 64)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Harga harus berupa angka", nil)
			return
		}
		itemData.Harga = harga
		itemData.Kategori = domain.ItemCategory(c.PostForm("kategori"))
		itemData.Deskripsi = c.PostForm("deskripsi")
		
		// Validasi manual
		if itemData.NamaBarang == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Nama barang tidak boleh kosong", nil)
			return
		}
		
		if itemData.Harga <= 0 {
			utils.ErrorResponse(c, http.StatusBadRequest, "Harga harus lebih dari 0", nil)
			return
		}
		
		if itemData.Kategori == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Kategori tidak boleh kosong", nil)
			return
		}
	} else {
		// Binding JSON jika tidak ada gambar
		if err := c.ShouldBindJSON(&itemData); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
			return
		}
		
		// Validasi manual
		if itemData.NamaBarang == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Nama barang tidak boleh kosong", nil)
			return
		}
		
		if itemData.Harga <= 0 {
			utils.ErrorResponse(c, http.StatusBadRequest, "Harga harus lebih dari 0", nil)
			return
		}
		
		if itemData.Kategori == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Kategori tidak boleh kosong", nil)
			return
		}
	}

	// Set status default
	itemData.Status = domain.StatusTersedia

	// Buat barang baru
	newItem, err := h.itemService.Create(c.Request.Context(), &itemData, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Jika ada gambar, upload setelah item dibuat
	if hasImage {
		// Upload gambar dengan ID item yang baru dibuat
		fileInfo, err := h.itemService.UploadImage(c, newItem.ID, userID.(uint))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Barang berhasil dibuat tetapi gagal mengupload gambar: "+err.Error(), newItem)
			return
		}
		
		// Ambil URL gambar dari response
		parts := strings.Split(fileInfo, "|")
		if len(parts) > 2 {
			newItem.Gambar = parts[2] // URL gambar
		}
	}

	utils.SuccessResponse(c, http.StatusCreated, "Barang berhasil ditambahkan", newItem)
}

// GetItem mendapatkan data barang berdasarkan ID
func (h *ItemHandler) GetItem(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak valid", nil)
		return
	}

	// Dapatkan data barang
	item, err := h.itemService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data barang berhasil diambil", item)
}

// GetAllItems mendapatkan daftar barang
func (h *ItemHandler) GetAllItems(c *gin.Context) {
	// Dapatkan parameter paginasi dan filter
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")
	kategori := c.DefaultQuery("kategori", "")
	status := c.DefaultQuery("status", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Dapatkan daftar barang
	items, totalPages, totalItems, err := h.itemService.GetAll(c.Request.Context(), page, limit, search, kategori, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar barang berhasil diambil", items, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// GetItemsByPenjual mendapatkan daftar barang berdasarkan penjual
func (h *ItemHandler) GetItemsByPenjual(c *gin.Context) {
	// Dapatkan ID penjual dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID penjual tidak valid", nil)
		return
	}

	// Dapatkan parameter paginasi
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Dapatkan daftar barang
	items, totalPages, totalItems, err := h.itemService.GetByPenjualID(c.Request.Context(), uint(id), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar barang berhasil diambil", items, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// GetMyItems mendapatkan daftar barang milik pengguna yang login
func (h *ItemHandler) GetMyItems(c *gin.Context) {
	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Dapatkan parameter paginasi
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Dapatkan daftar barang
	items, totalPages, totalItems, err := h.itemService.GetByPenjualID(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar barang berhasil diambil", items, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// UpdateItem memperbarui data barang
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Bind data
	var itemData struct {
		NamaBarang string               `json:"nama_barang"`
		Harga      float64              `json:"harga" binding:"omitempty,gt=0"`
		Kategori   domain.ItemCategory  `json:"kategori" binding:"omitempty,oneof=Buku Elektronik Perabotan Kos-kosan Lainnya"`
		Deskripsi  string               `json:"deskripsi"`
	}
	
	// Binding JSON
	if err := c.ShouldBindJSON(&itemData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}

	// Validasi harga jika tidak kosong
	if itemData.Harga < 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Harga tidak boleh negatif", nil)
		return
	}

	// Buat objek item
	item := &domain.Item{
		NamaBarang: itemData.NamaBarang,
		Harga:      itemData.Harga,
		Kategori:   itemData.Kategori,
		Deskripsi:  itemData.Deskripsi,
	}

	// Update barang
	updatedItem, err := h.itemService.Update(c.Request.Context(), uint(id), item, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data barang berhasil diperbarui", updatedItem)
}

// UpdateItemStatus memperbarui status barang
func (h *ItemHandler) UpdateItemStatus(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Bind data
	var statusData struct {
		Status domain.ItemStatus `json:"status" binding:"required,oneof=Tersedia Terjual Dihapus"`
	}
	
	// Binding JSON
	if err := c.ShouldBindJSON(&statusData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}

	// Update status barang
	if err := h.itemService.UpdateStatus(c.Request.Context(), uint(id), statusData.Status, userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Status barang berhasil diperbarui", nil)
}

// DeleteItem menghapus barang
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Cek role
	role, exists := c.Get("userRole")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Jika admin, bisa hapus langsung
	if role.(domain.Role) == domain.RoleAdmin {
		if err := h.itemService.Delete(c.Request.Context(), uint(id), userID.(uint)); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utils.SuccessResponse(c, http.StatusOK, "Barang berhasil dihapus", nil)
		return
	}

	// Jika bukan admin, cek kepemilikan
	if err := h.itemService.Delete(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Barang berhasil dihapus", nil)
}

// UploadItemImage mengupload gambar barang
func (h *ItemHandler) UploadItemImage(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Upload gambar
	fileInfo, err := h.itemService.UploadImage(c, uint(id), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Parse fileInfo to get fileID, fileName, and viewURL
	parts := strings.Split(fileInfo, "|")
	fileID := parts[0]
	fileName := parts[0]
	var viewURL string
	
	if len(parts) > 1 {
		fileName = parts[1]
	}
	
	if len(parts) > 2 {
		viewURL = parts[2]
	} else {
		// Get Appwrite config for creating URLs if not included in response
		appwriteConfig, err := config.LoadConfig()
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memuat konfigurasi", nil)
			return
		}

		// Create file view URL
		appwriteEndpoint := appwriteConfig.Appwrite.Endpoint
		projectID := appwriteConfig.Appwrite.ProjectID
		bucketID := appwriteConfig.Appwrite.BucketID
		
		// Debug logs for troubleshooting
		fmt.Printf("Appwrite View URL: %s\n", appwriteEndpoint)
		fmt.Printf("Project ID: %s\n", projectID)
		fmt.Printf("Bucket ID: %s\n", bucketID)
		fmt.Printf("File ID: %s\n", fileID)
		
		viewURL = fmt.Sprintf("%s/storage/buckets/%s/files/%s/view?project=%s", 
			appwriteEndpoint,
			bucketID,
			fileID,
			projectID)
	}

	utils.SuccessResponse(c, http.StatusOK, "Gambar berhasil diupload", gin.H{
		"file_name": fileName,
		"file_id": fileID,
		"view_url": viewURL,
	})
}

// RegisterRoutes mendaftarkan route untuk ItemHandler
func (h *ItemHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	items := router.Group("/items")
	{
		items.GET("", h.GetAllItems)
		items.GET("/:id", h.GetItem)
		items.GET("/penjual/:id", h.GetItemsByPenjual)
		items.GET("/my", authMiddleware, h.GetMyItems)
		items.POST("", authMiddleware, h.CreateItem)
		items.PATCH("/:id", authMiddleware, h.UpdateItem)
		items.PATCH("/:id/status", authMiddleware, h.UpdateItemStatus)
		items.DELETE("/:id", authMiddleware, h.DeleteItem)
		items.POST("/:id/upload", authMiddleware, h.UploadItemImage)
	}

	// Admin routes
	admin := router.Group("/admin")
	{
		admin.DELETE("/items/:id", authMiddleware, adminMiddleware, h.DeleteItem)
	}
}