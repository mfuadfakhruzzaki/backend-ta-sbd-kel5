package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestIDKey is the context key for the request ID
const RequestIDKey = "request_id"

// LoggerMiddleware adalah middleware untuk logging dengan format terstruktur
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Create request-scoped logger
		reqLogger := log.With().Str("request_id", requestID).Logger()

		// Process request
		reqLogger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Request received")

		// Process the request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code and error if any
		status := c.Writer.Status()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log request completion
		logEvent := reqLogger.Info()
		if status >= 400 {
			logEvent = reqLogger.Error()
		}

		logEvent.
			Int("status", status).
			Dur("latency", latency).
			Int("bytes", c.Writer.Size())

		if errMsg != "" {
			logEvent.Str("error", errMsg)
		}

		logEvent.Msg("Request completed")
	}
}

// GetRequestID returns the request ID from the context
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get(RequestIDKey)
	if !exists {
		return "unknown"
	}
	return requestID.(string)
}

// GetRequestLogger returns the logger for the current request
func GetRequestLogger(c *gin.Context) zerolog.Logger {
	requestID := GetRequestID(c)
	
	logger := log.With().
		Str("request_id", requestID).
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method)
	
	// Add user ID if available
	userID, exists := c.Get("userID")
	if exists {
		logger = logger.Uint("user_id", userID.(uint))
	}
	
	return logger.Logger()
}