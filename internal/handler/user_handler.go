package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// UserHandler menangani endpoint terkait pengguna
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser mendapatkan data pengguna berdasarkan ID
// @Summary      Get user by ID
// @Description  Mendapatkan detail pengguna berdasarkan ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse{data=domain.UserResponse}
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      404  {object}  utils.StandardResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", nil)
		return
	}

	// Dapatkan data pengguna
	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data pengguna berhasil diambil", user)
}

// GetCurrentUser mendapatkan data pengguna yang sedang login
// @Summary      Get current user
// @Description  Mendapatkan detail pengguna yang sedang login
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse{data=domain.UserResponse}
// @Failure      401  {object}  utils.StandardResponse
// @Failure      404  {object}  utils.StandardResponse
// @Router       /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Dapatkan data pengguna
	user, err := h.userService.GetByID(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data pengguna berhasil diambil", user)
}

// GetAllUsers mendapatkan daftar pengguna
// @Summary      List all users
// @Description  Mendapatkan daftar pengguna dengan paginasi (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page    query    int     false  "Page number (default: 1)"
// @Param        limit   query    int     false  "Items per page (default: 10)"
// @Param        search  query    string  false  "Search query"
// @Security     BearerAuth
// @Success      200     {object}  utils.PaginatedResponse{data=[]domain.UserResponse}
// @Failure      401     {object}  utils.StandardResponse
// @Failure      403     {object}  utils.StandardResponse
// @Failure      500     {object}  utils.StandardResponse
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Dapatkan parameter paginasi
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Dapatkan daftar pengguna
	users, totalPages, totalItems, err := h.userService.GetAll(c.Request.Context(), page, limit, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Buat response paginasi
	utils.SuccessPaginatedResponse(c, http.StatusOK, "Daftar pengguna berhasil diambil", users, utils.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})
}

// UpdateUser memperbarui data pengguna
// @Summary      Update user
// @Description  Memperbarui data pengguna
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "User ID"
// @Param        request  body      domain.UpdateUserRequest  true  "User update data"
// @Security     BearerAuth
// @Success      200      {object}  utils.StandardResponse{data=domain.UserResponse}
// @Failure      400      {object}  utils.StandardResponse
// @Failure      401      {object}  utils.StandardResponse
// @Failure      403      {object}  utils.StandardResponse
// @Failure      404      {object}  utils.StandardResponse
// @Failure      500      {object}  utils.StandardResponse
// @Router       /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Validasi akses
	if uint(id) != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "Anda hanya dapat mengubah data diri sendiri", nil)
		return
	}

	// Bind data menggunakan struct khusus untuk update
	var userData struct {
		Nama     string `json:"nama"`
		NoHP     string `json:"no_hp"`
		Alamat   string `json:"alamat"`
		Password string `json:"password"`
	}
	
	// Binding JSON
	if err := c.ShouldBindJSON(&userData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}
	
	// Buat objek user untuk update
	user := &domain.User{
		Nama:     userData.Nama,
		NoHP:     userData.NoHP,
		Alamat:   userData.Alamat,
		Password: userData.Password,
	}

	// Update pengguna
	updatedUser, err := h.userService.Update(c.Request.Context(), uint(id), user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data pengguna berhasil diperbarui", updatedUser)
}

// DeleteUser menghapus pengguna
// @Summary      Delete user
// @Description  Menghapus pengguna berdasarkan ID (self atau admin)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  utils.StandardResponse
// @Failure      400  {object}  utils.StandardResponse
// @Failure      401  {object}  utils.StandardResponse
// @Failure      403  {object}  utils.StandardResponse
// @Failure      500  {object}  utils.StandardResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Dapatkan ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", nil)
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Dapatkan role dari context
	role, exists := c.Get("userRole")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Validasi akses
	if uint(id) != userID.(uint) && role.(domain.Role) != domain.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk menghapus pengguna ini", nil)
		return
	}

	// Hapus pengguna
	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pengguna berhasil dihapus", nil)
}

// RegisterRoutes mendaftarkan route untuk UserHandler
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	{
		users.GET("", authMiddleware, adminMiddleware, h.GetAllUsers)
		users.GET("/me", authMiddleware, h.GetCurrentUser)
		users.GET("/:id", authMiddleware, h.GetUser)
		users.PATCH("/:id", authMiddleware, h.UpdateUser)
		users.DELETE("/:id", authMiddleware, h.DeleteUser)
	}
}