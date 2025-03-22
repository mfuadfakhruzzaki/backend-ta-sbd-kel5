package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware adalah middleware untuk logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Waktu awal
		startTime := time.Now()

		// Proses request
		c.Next()

		// Waktu selesai
		endTime := time.Now()

		// Hitung durasi
		duration := endTime.Sub(startTime)

		// Log informasi
		fmt.Printf("[%s] %s %s %d %s\n",
			endTime.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}