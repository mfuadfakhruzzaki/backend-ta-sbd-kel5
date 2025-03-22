package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationError adalah struktur untuk error validasi
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

var validate *validator.Validate

// init menginisialisasi validator
func init() {
	validate = validator.New()

	// Register fungsi untuk mendapatkan nilai dari tag json untuk nama field
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate memvalidasi struct berdasarkan tag validator
func Validate(data interface{}) (bool, []ValidationError) {
	err := validate.Struct(data)
	if err == nil {
		return true, nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		validationError := ValidationError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   fmt.Sprintf("%v", err.Value()),
			Message: getErrorMessage(err),
		}
		validationErrors = append(validationErrors, validationError)
	}

	return false, validationErrors
}

// ValidateJSON memvalidasi JSON dari request body
func ValidateJSON(c *gin.Context, data interface{}) (bool, []ValidationError) {
	if err := c.ShouldBindJSON(data); err != nil {
		// Cek apakah error dari validator
		if _, ok := err.(validator.ValidationErrors); ok {
			_, validationErrors := Validate(data)
			return false, validationErrors
		}
		
		// Error lain (format JSON tidak valid, dll)
		return false, []ValidationError{
			{
				Field:   "body",
				Tag:     "json",
				Value:   "",
				Message: "Format JSON tidak valid",
			},
		}
	}

	return Validate(data)
}

// getErrorMessage menerjemahkan error validator menjadi pesan yang mudah dibaca
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Bidang ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("Nilai minimal adalah %s", err.Param())
	case "max":
		return fmt.Sprintf("Nilai maksimal adalah %s", err.Param())
	case "gt":
		return fmt.Sprintf("Nilai harus lebih besar dari %s", err.Param())
	case "lt":
		return fmt.Sprintf("Nilai harus lebih kecil dari %s", err.Param())
	case "oneof":
		return fmt.Sprintf("Nilai harus salah satu dari [%s]", err.Param())
	default:
		return fmt.Sprintf("Validasi %s gagal", err.Tag())
	}
}