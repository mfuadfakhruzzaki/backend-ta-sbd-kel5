package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// ChatHandler menangani endpoint terkait chat
// @Summary      Chat handler
// @Description  Menangani endpoint terkait chat dan pesan
type ChatHandler struct {
	chatService service.ChatService
}

// NewChatHandler membuat instance baru ChatHandler
func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// SendMessage mengirim pesan chat baru
// @Summary      Send message
// @Description  Mengirim pesan chat baru
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        request  body      domain.SendMessageRequestSwagger  true  "Message data"
// @Security     BearerAuth
// @Success      201      {object}  utils.StandardResponse{data=domain.ChatResponseSwagger}
// @Failure      400      {object}  utils.StandardResponse
// @Failure      401      {object}  utils.StandardResponse
// @Failure      500      {object}  utils.StandardResponse
// @Router       /chats [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Bind data
	var chatData struct {
		PenerimaID uint   `json:"penerima_id" binding:"required"`
		BarangID   uint   `json:"barang_id" binding:"required"`
		Pesan      string `json:"pesan" binding:"required"`
	}
	
	// Binding JSON
	if err := c.ShouldBindJSON(&chatData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}
	
	// Validasi manual
	if chatData.PenerimaID == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID penerima tidak boleh kosong", nil)
		return
	}
	
	if chatData.BarangID == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID barang tidak boleh kosong", nil)
		return
	}
	
	if chatData.Pesan == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Pesan tidak boleh kosong", nil)
		return
	}

	// Buat objek chat
	chat := &domain.Chat{
		PengirimID: userID.(uint),
		PenerimaID: chatData.PenerimaID,
		BarangID:   chatData.BarangID,
		Pesan:      chatData.Pesan,
	}

	// Kirim pesan
	sentChat, err := h.chatService.SendMessage(c.Request.Context(), chat, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pesan berhasil dikirim", sentChat)
}

// GetChat mendapatkan data chat berdasarkan ID
// @Summary      Get chat by ID
// @Description  Mendapatkan detail chat berdasarkan ID
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Chat ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse{data=domain.ChatResponseSwagger}
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      404  {object}  utils.StandardResponse
// @Router       /chats/{id} [get]
func (h *ChatHandler) GetChat(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID chat tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Dapatkan data chat
	chat, err := h.chatService.GetByID(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data chat berhasil diambil", chat)
}

// GetChatsByBarang mendapatkan daftar chat berdasarkan ID barang
// @Summary      Get chats by item ID
// @Description  Mendapatkan daftar chat yang terkait dengan barang tertentu
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id     path      int  true   "Item ID"
// @Param        page   query     int  false  "Page number (default: 1)"
// @Param        limit  query     int  false  "Items per page (default: 10)"
// @Security     BearerAuth
// @Success      200    {object}  utils.PaginatedResponse{data=[]domain.ChatResponseSwagger}
// @Failure      400    {object}  utils.StandardResponse
// @Failure      401    {object}  utils.StandardResponse
// @Failure      500    {object}  utils.StandardResponse
// @Router       /chats/barang/{id} [get]
func (h *ChatHandler) GetChatsByBarang(c *gin.Context) {
	// Dapatkan ID barang dari URL
	barangIDStr := c.Param("id")
	barangID, err := strconv.ParseUint(barangIDStr, 10, 64)
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

	// Dapatkan daftar chat
	chats, totalPages, totalItems, err := h.chatService.GetByBarangID(c.Request.Context(), uint(barangID), page, limit, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar chat berhasil diambil", chats, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// GetConversation mendapatkan percakapan antara dua pengguna tentang suatu barang
// @Summary      Get conversation
// @Description  Mendapatkan percakapan antara dua pengguna tentang barang tertentu
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        partner_id  query     int  true   "Partner user ID"
// @Param        barang_id   query     int  true   "Item ID"
// @Param        page        query     int  false  "Page number (default: 1)"
// @Param        limit       query     int  false  "Items per page (default: 10)"
// @Security     BearerAuth
// @Success      200         {object}  utils.PaginatedResponse{data=[]domain.ChatResponseSwagger}
// @Failure      400         {object}  utils.StandardResponse
// @Failure      401         {object}  utils.StandardResponse
// @Failure      500         {object}  utils.StandardResponse
// @Router       /chats/conversation [get]
func (h *ChatHandler) GetConversation(c *gin.Context) {
	// Dapatkan ID pengirim dan penerima dari query
	penerimaIDStr := c.Query("penerima_id")
	barangIDStr := c.Query("barang_id")

	if penerimaIDStr == "" || barangIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Parameter penerima_id dan barang_id diperlukan", nil)
		return
	}

	penerimaID, err := strconv.ParseUint(penerimaIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID penerima tidak valid", nil)
		return
	}

	barangID, err := strconv.ParseUint(barangIDStr, 10, 64)
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

	// Dapatkan parameter paginasi
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	// Dapatkan percakapan
	chats, totalPages, totalItems, err := h.chatService.GetConversation(
		c.Request.Context(), userID.(uint), uint(penerimaID), uint(barangID), page, limit, userID.(uint),
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Percakapan berhasil diambil", chats, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// GetChatPartners mendapatkan daftar pengguna yang pernah chat dengan pengguna saat ini
// @Summary      Get chat partners
// @Description  Mendapatkan daftar pengguna yang pernah berkomunikasi dengan pengguna saat ini
// @Tags         chats
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse{data=[]domain.ChatPartnerResponseSwagger}
// @Failure      401  {object}  utils.StandardResponse
// @Failure      500  {object}  utils.StandardResponse
// @Router       /chats/partners [get]
func (h *ChatHandler) GetChatPartners(c *gin.Context) {
	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Dapatkan partner chat
	users, err := h.chatService.GetChatPartners(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar partner chat berhasil diambil", users)
}

// MarkAsRead menandai pesan chat sebagai telah dibaca
// @Summary      Mark message as read
// @Description  Menandai pesan chat sebagai telah dibaca
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Chat ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      404  {object}  utils.StandardResponse
// @Failure      500  {object}  utils.StandardResponse
// @Router       /chats/{id}/read [patch]
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID chat tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Tandai sebagai dibaca
	if err := h.chatService.MarkAsRead(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pesan ditandai sebagai dibaca", nil)
}

// DeleteChat menghapus pesan chat
// @Summary      Delete chat
// @Description  Menghapus pesan chat berdasarkan ID
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Chat ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      500  {object}  utils.StandardResponse
// @Router       /chats/{id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID chat tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Hapus chat
	if err := h.chatService.Delete(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Chat berhasil dihapus", nil)
}

// RegisterRoutes mendaftarkan route untuk ChatHandler
func (h *ChatHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	chats := router.Group("/chats")
	chats.Use(authMiddleware)
	{
		chats.POST("", h.SendMessage)
		chats.GET("/:id", h.GetChat)
		chats.GET("/barang/:id", h.GetChatsByBarang)
		chats.GET("/conversation", h.GetConversation)
		chats.GET("/partners", h.GetChatPartners)
		chats.PATCH("/:id/read", h.MarkAsRead)
		chats.DELETE("/:id", h.DeleteChat)
	}
}