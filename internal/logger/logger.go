package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

var defaultLogger *slog.Logger

// LogLevel represents the log level
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Config holds logger configuration
type Config struct {
	Level  LogLevel `yaml:"level" mapstructure:"level"`
	Format string   `yaml:"format" mapstructure:"format"` // "json" or "text"
	Output string   `yaml:"output" mapstructure:"output"` // "stdout", "stderr", or file path
}

// Initialize sets up the global logger with the given configuration
func Initialize(cfg Config) error {
	// Determine log level
	var level slog.Level
	switch cfg.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Determine output destination
	var writer io.Writer
	switch cfg.Output {
	case "stdout", "":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// File output
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		writer = file
	}

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize time format
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	defaultLogger = slog.New(handler)
	return nil
}

// InitializeFromViper initializes logging from viper configuration
func InitializeFromViper() error {
	cfg := Config{
		Level:  LevelInfo,
		Format: "text",
		Output: "stderr",
	}

	// Check for debug mode from environment
	if os.Getenv("RUNE_DEBUG") == "true" {
		cfg.Level = LevelDebug
	}

	// Override from viper if available
	if viper.IsSet("logging.level") {
		cfg.Level = LogLevel(viper.GetString("logging.level"))
	}
	if viper.IsSet("logging.format") {
		cfg.Format = viper.GetString("logging.format")
	}
	if viper.IsSet("logging.output") {
		cfg.Output = viper.GetString("logging.output")
	}

	return Initialize(cfg)
}

// GetLogger returns the default logger instance
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		// Fallback to a basic logger if not initialized
		defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return defaultLogger
}

// Context-aware logging functions

// Debug logs a debug message with optional key-value pairs
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

// Info logs an info message with optional key-value pairs
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs a warning message with optional key-value pairs
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// With creates a logger with additional context
func With(args ...any) *slog.Logger {
	return GetLogger().With(args...)
}

// WithContext creates a logger with context
func WithContext(ctx context.Context) *slog.Logger {
	return GetLogger().With("request_id", ctx.Value("request_id"))
}

// Component-specific loggers

// TelemetryLogger returns a logger for telemetry operations
func TelemetryLogger() *slog.Logger {
	return GetLogger().With("component", "telemetry")
}

// CommandLogger returns a logger for command operations
func CommandLogger(command string) *slog.Logger {
	return GetLogger().With("component", "command", "command", command)
}

// RitualLogger returns a logger for ritual operations
func RitualLogger(project string) *slog.Logger {
	return GetLogger().With("component", "ritual", "project", project)
}

// TrackingLogger returns a logger for tracking operations
func TrackingLogger() *slog.Logger {
	return GetLogger().With("component", "tracking")
}

// ConfigLogger returns a logger for configuration operations
func ConfigLogger() *slog.Logger {
	return GetLogger().With("component", "config")
}

// Helper functions for common patterns

// LogError logs an error with additional context
func LogError(err error, msg string, args ...any) {
	if err != nil {
		allArgs := append([]any{"error", err.Error()}, args...)
		GetLogger().Error(msg, allArgs...)
	}
}

// LogDuration logs the duration of an operation
func LogDuration(start time.Time, operation string, args ...any) {
	duration := time.Since(start)
	allArgs := append([]any{"operation", operation, "duration_ms", duration.Milliseconds()}, args...)
	GetLogger().Info("operation completed", allArgs...)
}