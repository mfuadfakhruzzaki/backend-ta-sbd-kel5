package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// AuthHandler menangani endpoint terkait otentikasi
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUser menangani pendaftaran pengguna baru
// @Summary      Register new user
// @Description  Mendaftarkan pengguna baru
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      domain.RegisterRequest  true  "User registration data"
// @Success      201      {object}  utils.StandardResponse{data=domain.RegisterResponse}
// @Failure      400      {object}  utils.StandardResponse
// @Failure      500      {object}  utils.StandardResponse
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// Gunakan struct pengguna sederhana tanpa validasi tag
	var userData struct {
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		Password string `json:"password"`
		NoHP     string `json:"no_hp"`
		Alamat   string `json:"alamat"`
	}
	
	// Binding JSON tanpa validasi dulu
	if err := c.ShouldBindJSON(&userData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data: "+err.Error(), nil)
		return
	}
	
	// Buat objek User dari data yang dibinding
	user := &domain.User{
		Nama:     userData.Nama,
		Email:    userData.Email,
		Password: userData.Password,
		NoHP:     userData.NoHP,
		Alamat:   userData.Alamat,
	}
	
	// Validasi manual
	if user.Nama == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nama tidak boleh kosong", nil)
		return
	}
	
	if user.Email == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email tidak boleh kosong", nil)
		return
	}
	
	if user.Password == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Password tidak boleh kosong", nil)
		return
	}
	
	if len(user.Password) < 8 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Password minimal 8 karakter", nil)
		return
	}
	
	if user.NoHP == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nomor HP tidak boleh kosong", nil)
		return
	}

	// Daftarkan pengguna
	userID, err := h.authService.Register(c.Request.Context(), user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pendaftaran berhasil", gin.H{
		"user_id": userID,
	})
}

// LoginUser menangani login pengguna
// @Summary      Login user
// @Description  Melakukan login dan mendapatkan token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      domain.LoginRequest  true  "User login credentials"
// @Success      200      {object}  utils.StandardResponse{data=domain.LoginResponse}
// @Failure      400      {object}  utils.StandardResponse
// @Failure      401      {object}  utils.StandardResponse
// @Failure      500      {object}  utils.StandardResponse
// @Router       /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	// Struct untuk data login
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Binding JSON
	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca data login: "+err.Error(), nil)
		return
	}

	// Validasi manual
	if loginData.Email == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email tidak boleh kosong", nil)
		return
	}
	
	if loginData.Password == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Password tidak boleh kosong", nil)
		return
	}

	// Login
	token, user, err := h.authService.Login(c.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

// RegisterRoutes mendaftarkan route untuk AuthHandler
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
	}
}