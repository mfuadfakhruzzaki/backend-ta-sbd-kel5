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
// @Summary      Transactions handler
// @Description  Menangani endpoint terkait transaksi
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
// @Summary      Create a transaction
// @Description  Membuat transaksi baru untuk pembelian barang
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CreateTransactionRequest  true  "Transaction data"
// @Security     BearerAuth
// @Success      201      {object}  utils.StandardResponse{data=domain.TransactionResponse}
// @Failure      400      {object}  utils.StandardResponse
// @Failure      401      {object}  utils.StandardResponse
// @Failure      500      {object}  utils.StandardResponse
// @Router       /transactions [post]
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
// @Summary      Get transaction by ID
// @Description  Mendapatkan detail transaksi berdasarkan ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Transaction ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse{data=domain.TransactionResponse}
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      404  {object}  utils.StandardResponse
// @Router       /transactions/{id} [get]
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
// @Summary      List all transactions
// @Description  Mendapatkan daftar semua transaksi (admin only)
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        page   query    int  false  "Page number (default: 1)"
// @Param        limit  query    int  false  "Items per page (default: 10)"
// @Security     BearerAuth
// @Success      200    {object}  utils.PaginatedResponse{data=[]domain.TransactionResponse}
// @Failure      401    {object}  utils.StandardResponse
// @Failure      403    {object}  utils.StandardResponse
// @Failure      500    {object}  utils.StandardResponse
// @Router       /transactions [get]
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
// @Summary      Get my transactions as buyer
// @Description  Mendapatkan daftar transaksi dimana user adalah pembeli
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        page   query    int  false  "Page number (default: 1)"
// @Param        limit  query    int  false  "Items per page (default: 10)"
// @Security     BearerAuth
// @Success      200    {object}  utils.PaginatedResponse{data=[]domain.TransactionResponse}
// @Failure      401    {object}  utils.StandardResponse
// @Failure      500    {object}  utils.StandardResponse
// @Router       /transactions/as-pembeli [get]
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
// @Summary      Get my transactions as seller
// @Description  Mendapatkan daftar transaksi dimana user adalah penjual
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        page   query    int  false  "Page number (default: 1)"
// @Param        limit  query    int  false  "Items per page (default: 10)"
// @Security     BearerAuth
// @Success      200    {object}  utils.PaginatedResponse{data=[]domain.TransactionResponse}
// @Failure      401    {object}  utils.StandardResponse
// @Failure      500    {object}  utils.StandardResponse
// @Router       /transactions/as-penjual [get]
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
// @Summary      Update transaction status
// @Description  Memperbarui status transaksi (Pending, Selesai, Dibatalkan)
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id       path      int                              true  "Transaction ID"
// @Param        request  body      domain.UpdateTransactionStatusRequest  true  "Transaction status data"
// @Security     BearerAuth
// @Success      200      {object}  utils.StandardResponse
// @Failure      400      {object}  utils.StandardResponse
// @Failure      401      {object}  utils.StandardResponse
// @Failure      403      {object}  utils.StandardResponse
// @Failure      404      {object}  utils.StandardResponse
// @Failure      500      {object}  utils.StandardResponse
// @Router       /transactions/{id}/status [patch]
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
// @Summary      Delete transaction
// @Description  Menghapus transaksi berdasarkan ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Transaction ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      403  {object}  utils.StandardResponse
// @Failure      404  {object}  utils.StandardResponse
// @Failure      500  {object}  utils.StandardResponse
// @Router       /transactions/{id} [delete]
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