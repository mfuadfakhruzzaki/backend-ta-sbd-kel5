package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// TransactionHandler menangani endpoint terkait transaksi
type TransactionHandler struct {
	transactionService service.TransactionService
}

// NewTransactionHandler membuat instance baru TransactionHandler
func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction membuat transaksi baru
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Bind data
	var transactionData struct {
		BarangID uint `json:"barang_id" binding:"required"`
	}
	
	// Binding JSON
	if err := c.ShouldBindJSON(&transactionData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}
	
	// Validasi manual
	if transactionData.BarangID == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak boleh kosong", nil)
		return
	}

	// Buat objek transaksi
	transaction := &domain.Transaction{
		BarangID:  transactionData.BarangID,
		PembeliID: userID.(uint),
	}

	// Buat transaksi baru
	createdTransaction, err := h.transactionService.Create(c.Request.Context(), transaction, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Transaksi berhasil dibuat", createdTransaction)
}

// GetTransaction mendapatkan data transaksi berdasarkan ID
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID transaksi tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Dapatkan data transaksi
	transaction, err := h.transactionService.GetByID(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data transaksi berhasil diambil", transaction)
}

// GetAllTransactions mendapatkan semua transaksi (admin only)
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	// Cek role dari context
	role, exists := c.Get("userRole")
	if !exists || role.(domain.Role) != domain.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: memerlukan hak admin", nil)
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

	// Dapatkan daftar transaksi
	transactions, totalPages, totalItems, err := h.transactionService.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar transaksi berhasil diambil", transactions, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// GetMyTransactionsAsPembeli mendapatkan daftar transaksi sebagai pembeli
func (h *TransactionHandler) GetMyTransactionsAsPembeli(c *gin.Context) {
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

	// Dapatkan daftar transaksi
	transactions, totalPages, totalItems, err := h.transactionService.GetByPembeliID(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar transaksi berhasil diambil", transactions, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// GetMyTransactionsAsPenjual mendapatkan daftar transaksi sebagai penjual
func (h *TransactionHandler) GetMyTransactionsAsPenjual(c *gin.Context) {
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

	// Dapatkan daftar transaksi
	transactions, totalPages, totalItems, err := h.transactionService.GetByPenjualID(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar transaksi berhasil diambil", transactions, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// UpdateTransactionStatus memperbarui status transaksi
func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID transaksi tidak valid", nil)
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
		Status domain.TransactionStatus `json:"status" binding:"required,oneof=Pending Selesai Dibatalkan"`
	}
	
	// Binding JSON
	if err := c.ShouldBindJSON(&statusData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}
	
	// Validasi manual
	if statusData.Status == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Status transaksi tidak boleh kosong", nil)
		return
	}

	// Update status transaksi
	if err := h.transactionService.UpdateStatus(c.Request.Context(), uint(id), statusData.Status, userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Status transaksi berhasil diperbarui", nil)
}

// DeleteTransaction menghapus transaksi
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID transaksi tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Hapus transaksi
	if err := h.transactionService.Delete(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaksi berhasil dihapus", nil)
}

// RegisterRoutes mendaftarkan route untuk TransactionHandler
func (h *TransactionHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	transactions := router.Group("/transactions")
	transactions.Use(authMiddleware)
	{
		transactions.POST("", h.CreateTransaction)
		transactions.GET("", h.GetAllTransactions)
		transactions.GET("/:id", h.GetTransaction)
		transactions.GET("/as-pembeli", h.GetMyTransactionsAsPembeli)
		transactions.GET("/as-penjual", h.GetMyTransactionsAsPenjual)
		transactions.PATCH("/:id/status", h.UpdateTransactionStatus)
		transactions.DELETE("/:id", h.DeleteTransaction)
	}
}