package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/domain"
	"github.com/mfuadfakhruzzaki/jubel/internal/service"
	"github.com/mfuadfakhruzzaki/jubel/internal/utils"
)

// AuthMiddleware middleware untuk otentikasi JWT
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware membuat instance baru AuthMiddleware
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// VerifyToken digunakan untuk memverifikasi token JWT
func (m *AuthMiddleware) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dapatkan token dari header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token tidak ditemukan", nil)
			c.Abort()
			return
		}

		// Format token harus "Bearer [token]"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Format token tidak valid", nil)
			c.Abort()
			return
		}

		token := parts[1]

		// Validasi token
		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		// Set user info di context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireAdmin digunakan untuk memeriksa apakah pengguna adalah admin
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dapatkan role dari context
		role, exists := c.Get("userRole")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		// Cek apakah role adalah admin
		if role.(domain.Role) != domain.RoleAdmin {
			utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak: memerlukan hak admin", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}