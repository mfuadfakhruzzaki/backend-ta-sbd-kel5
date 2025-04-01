package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/jubel/internal/errors"
	"github.com/rs/zerolog"
)

// ErrorMiddleware returns a middleware that handles errors and returns appropriate HTTP responses
func ErrorMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			logErr := logger.With().
				Str("requestID", RequestIDFromContext(c)).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Logger()

			// Try to convert to StandardError
			if stdErr, ok := errors.AsStandardError(err); ok {
				// Log with context from the standard error
				errContext := logErr.With().
					Str("errorType", stdErr.Type).
					Int("statusCode", stdErr.Code)

				if stdErr.Metadata != nil {
					errContext = errContext.Fields(stdErr.Metadata)
				}

				// Create the logger from context and log the message
				errorLogger := errContext.Logger()
				if stdErr.Err != nil {
					errorLogger.Error().Err(stdErr.Err).Msg(stdErr.Message)
				} else {
					errorLogger.Error().Msg(stdErr.Message)
				}

				// Respond with JSON
				c.JSON(stdErr.Code, gin.H{
					"error": gin.H{
						"type":    stdErr.Type,
						"message": stdErr.Message,
					},
				})
				return
			}

			// If not a StandardError, treat as internal error
			logErr.Error().Err(err).Msg("Internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"type":    "internal_error",
					"message": "An internal server error occurred",
				},
			})
		}
	}
}

// RequestIDFromContext retrieves the request ID from gin context
func RequestIDFromContext(c *gin.Context) string {
	if reqID, exists := c.Get("RequestID"); exists {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return "unknown"
} 