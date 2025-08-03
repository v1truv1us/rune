package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		verify func(t *testing.T)
	}{
		{
			name: "debug level with text format",
			config: Config{
				Level:  LevelDebug,
				Format: "text",
				Output: "stderr",
			},
			verify: func(t *testing.T) {
				logger := GetLogger()
				assert.NotNil(t, logger)
				assert.True(t, logger.Enabled(nil, slog.LevelDebug))
			},
		},
		{
			name: "info level with json format",
			config: Config{
				Level:  LevelInfo,
				Format: "json",
				Output: "stdout",
			},
			verify: func(t *testing.T) {
				logger := GetLogger()
				assert.NotNil(t, logger)
				assert.True(t, logger.Enabled(nil, slog.LevelInfo))
				assert.False(t, logger.Enabled(nil, slog.LevelDebug))
			},
		},
		{
			name: "error level",
			config: Config{
				Level:  LevelError,
				Format: "text",
				Output: "stderr",
			},
			verify: func(t *testing.T) {
				logger := GetLogger()
				assert.NotNil(t, logger)
				assert.True(t, logger.Enabled(nil, slog.LevelError))
				assert.False(t, logger.Enabled(nil, slog.LevelWarn))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.config)
			require.NoError(t, err)
			tt.verify(t)
		})
	}
}

func TestInitializeFromViper(t *testing.T) {
	// Test with RUNE_DEBUG environment variable
	t.Run("debug mode from environment", func(t *testing.T) {
		os.Setenv("RUNE_DEBUG", "true")
		defer os.Unsetenv("RUNE_DEBUG")

		err := InitializeFromViper()
		require.NoError(t, err)

		logger := GetLogger()
		assert.True(t, logger.Enabled(nil, slog.LevelDebug))
	})

	t.Run("default configuration", func(t *testing.T) {
		os.Unsetenv("RUNE_DEBUG")

		err := InitializeFromViper()
		require.NoError(t, err)

		logger := GetLogger()
		assert.True(t, logger.Enabled(nil, slog.LevelInfo))
		assert.False(t, logger.Enabled(nil, slog.LevelDebug))
	})
}

func TestLoggingFunctions(t *testing.T) {
	// Capture output to verify logging
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time for consistent testing
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	defaultLogger = slog.New(handler)

	t.Run("debug logging", func(t *testing.T) {
		buf.Reset()
		Debug("test debug message", "key", "value")
		output := buf.String()
		assert.Contains(t, output, "test debug message")
		assert.Contains(t, output, "key=value")
		assert.Contains(t, output, "level=DEBUG")
	})

	t.Run("info logging", func(t *testing.T) {
		buf.Reset()
		Info("test info message", "user_id", "123")
		output := buf.String()
		assert.Contains(t, output, "test info message")
		assert.Contains(t, output, "user_id=123")
		assert.Contains(t, output, "level=INFO")
	})

	t.Run("warn logging", func(t *testing.T) {
		buf.Reset()
		Warn("test warn message", "component", "config")
		output := buf.String()
		assert.Contains(t, output, "test warn message")
		assert.Contains(t, output, "component=config")
		assert.Contains(t, output, "level=WARN")
	})

	t.Run("error logging", func(t *testing.T) {
		buf.Reset()
		Error("test error message", "error", "something went wrong")
		output := buf.String()
		assert.Contains(t, output, "test error message")
		assert.Contains(t, output, "error=\"something went wrong\"")
		assert.Contains(t, output, "level=ERROR")
	})
}

func TestComponentLoggers(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	defaultLogger = slog.New(handler)

	t.Run("telemetry logger", func(t *testing.T) {
		buf.Reset()
		logger := TelemetryLogger()
		logger.Info("telemetry event", "event", "session_started")
		output := buf.String()
		assert.Contains(t, output, "component=telemetry")
		assert.Contains(t, output, "event=session_started")
	})

	t.Run("command logger", func(t *testing.T) {
		buf.Reset()
		logger := CommandLogger("start")
		logger.Info("executing command")
		output := buf.String()
		assert.Contains(t, output, "component=command")
		assert.Contains(t, output, "command=start")
	})

	t.Run("ritual logger", func(t *testing.T) {
		buf.Reset()
		logger := RitualLogger("test-project")
		logger.Info("executing ritual")
		output := buf.String()
		assert.Contains(t, output, "component=ritual")
		assert.Contains(t, output, "project=test-project")
	})

	t.Run("tracking logger", func(t *testing.T) {
		buf.Reset()
		logger := TrackingLogger()
		logger.Info("session started")
		output := buf.String()
		assert.Contains(t, output, "component=tracking")
	})

	t.Run("config logger", func(t *testing.T) {
		buf.Reset()
		logger := ConfigLogger()
		logger.Info("config loaded")
		output := buf.String()
		assert.Contains(t, output, "component=config")
	})
}

func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	defaultLogger = slog.New(handler)

	t.Run("with error", func(t *testing.T) {
		buf.Reset()
		err := errors.New("test error")
		LogError(err, "operation failed", "operation", "test")
		output := buf.String()
		assert.Contains(t, output, "operation failed")
		assert.Contains(t, output, "error=\"test error\"")
		assert.Contains(t, output, "operation=test")
		assert.Contains(t, output, "level=ERROR")
	})

	t.Run("with nil error", func(t *testing.T) {
		buf.Reset()
		LogError(nil, "operation succeeded", "operation", "test")
		output := buf.String()
		assert.Empty(t, output) // Should not log anything for nil error
	})
}

func TestLogDuration(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	defaultLogger = slog.New(handler)

	start := time.Now().Add(-100 * time.Millisecond)
	LogDuration(start, "test_operation", "component", "test")
	
	output := buf.String()
	assert.Contains(t, output, "operation completed")
	assert.Contains(t, output, "operation=test_operation")
	assert.Contains(t, output, "component=test")
	assert.Contains(t, output, "duration_ms=")
}

func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, "2023-01-01T00:00:00Z")
			}
			return a
		},
	})
	defaultLogger = slog.New(handler)

	Info("test json message", "key", "value", "number", 42)
	
	output := strings.TrimSpace(buf.String())
	
	// Parse as JSON to verify structure
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(output), &parsed)
	require.NoError(t, err)
	
	assert.Equal(t, "INFO", parsed["level"])
	assert.Equal(t, "test json message", parsed["msg"])
	assert.Equal(t, "value", parsed["key"])
	assert.Equal(t, float64(42), parsed["number"]) // JSON numbers are float64
	assert.Equal(t, "2023-01-01T00:00:00Z", parsed["time"])
}

func TestWithContext(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	defaultLogger = slog.New(handler)

	logger := With("session_id", "abc123", "user", "testuser")
	logger.Info("test message")
	
	output := buf.String()
	assert.Contains(t, output, "session_id=abc123")
	assert.Contains(t, output, "user=testuser")
	assert.Contains(t, output, "test message")
}