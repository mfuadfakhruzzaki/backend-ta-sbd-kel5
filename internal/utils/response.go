package utils

import (
	"github.com/gin-gonic/gin"
)

// Struktur respons API yang konsisten
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse mengembalikan respons sukses
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengembalikan respons error
func ErrorResponse(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

// PaginatedResponse adalah struktur untuk respons dengan paginasi
type PaginatedResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

// Meta berisi informasi meta untuk paginasi
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// SuccessPaginatedResponse mengembalikan respons sukses dengan paginasi
func SuccessPaginatedResponse(c *gin.Context, statusCode int, message string, data interface{}, meta Meta) {
	c.JSON(statusCode, PaginatedResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}