package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the global logger instance
	Logger zerolog.Logger
)

// Config contains logger configuration
type Config struct {
	// LogLevel is the logging level (debug, info, warn, error)
	LogLevel string
	// Pretty enables pretty console output
	Pretty bool
	// WithCaller adds caller information to log entries
	WithCaller bool
}

// Initialize sets up the logger with the given configuration
func Initialize(config Config) {
	// Set global time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Set global level
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output
	var output io.Writer = os.Stdout
	if config.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("@%s", i)
			},
		}
	}

	// Create logger
	logger := zerolog.New(output).With().Timestamp()
	if config.WithCaller {
		logger = logger.Caller()
	}
	Logger = logger.Logger()

	// Set as global logger
	log.Logger = Logger
}

// WithRequestID returns a logger with the request ID field set
func WithRequestID(requestID string) zerolog.Logger {
	return Logger.With().Str("request_id", requestID).Logger()
}

// WithUserID returns a logger with the user ID field set
func WithUserID(userID uint) zerolog.Logger {
	return Logger.With().Uint("user_id", userID).Logger()
}

// Fields returns a map of structured fields that can be used in logging
type Fields map[string]interface{}

// WithFields returns a logger with the given fields
func WithFields(fields Fields) zerolog.Logger {
	ctx := Logger.With()
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			ctx = ctx.Str(k, val)
		case int:
			ctx = ctx.Int(k, val)
		case uint:
			ctx = ctx.Uint(k, val)
		case float64:
			ctx = ctx.Float64(k, val)
		case bool:
			ctx = ctx.Bool(k, val)
		case error:
			ctx = ctx.Err(val)
		default:
			ctx = ctx.Interface(k, val)
		}
	}
	return ctx.Logger()
} 