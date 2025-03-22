package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// Bind data
	var itemData domain.Item
	valid, errors := utils.ValidateJSON(c, &itemData)
	if !valid {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", errors)
		return
	}

	// Buat barang baru
	item, err := h.itemService.Create(c.Request.Context(), &itemData, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Barang berhasil ditambahkan", item)
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
	var itemData domain.Item
	valid, errors := utils.ValidateJSON(c, &itemData)
	if !valid {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", errors)
		return
	}

	// Update barang
	updatedItem, err := h.itemService.Update(c.Request.Context(), uint(id), &itemData, userID.(uint))
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
		Status domain.ItemStatus `json:"status" validate:"required,oneof=Tersedia Terjual Dihapus"`
	}
	valid, errors := utils.ValidateJSON(c, &statusData)
	if !valid {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", errors)
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
	fileName, err := h.itemService.UploadImage(c, uint(id), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Gambar berhasil diupload", gin.H{
		"file_name": fileName,
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