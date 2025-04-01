package utils

import (
	"github.com/gin-gonic/gin"
)

// StandardResponse adalah struktur standar response API
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse adalah struktur response dengan paginasi
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    Meta        `json:"meta"`
}

// Meta untuk informasi paginasi
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// ErrorResponse mengirimkan response error
func ErrorResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Message: message,
		Data:    data,
	})
}

// SuccessResponse mengirimkan response sukses
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessPaginatedResponse mengirimkan response sukses dengan paginasi
func SuccessPaginatedResponse(c *gin.Context, statusCode int, message string, data interface{}, meta Meta) {
	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}