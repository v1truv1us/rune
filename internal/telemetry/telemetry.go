package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/otellogger"
	"github.com/getsentry/sentry-go"
)

type Client struct {
	sentryEnabled bool
	enabled       bool
	userID        string
	sessionID     string
	otelLogger    *otellogger.OtelLogger
}

func NewClient(userSentryDSN string) *Client {
	// Check if telemetry is disabled via environment variable or config
	enabled := os.Getenv("RUNE_TELEMETRY_DISABLED") != "true"

	// Generate or load user ID (anonymous)
	userID := getUserID()

	debugMode := os.Getenv("RUNE_DEBUG") == "true"

	if debugMode {
		// Use structured logging for debug output
		log := slog.Default()
		log.Debug("telemetry init",
			"enabled", enabled,
			"user_sentry_dsn", maskKey(userSentryDSN),
			"user_id", userID,
		)
	}

	client := &Client{
		enabled:       enabled,
		userID:        userID,
		sentryEnabled: userSentryDSN != "",
		sessionID:     generateSessionID(),
	}

	if !enabled {
		if debugMode {
			slog.Default().Info("telemetry disabled via env", "env", "RUNE_TELEMETRY_DISABLED")
		}
		return client
	}

	// Initialize OpenTelemetry logging with Sentry integration
	otelConfig := &otellogger.OtelLoggerConfig{
		SentryDSN:      userSentryDSN,
		ServiceName:    "rune",
		ServiceVersion: getVersion(),
		Environment:    getEnvironment(),
		DebugMode:      debugMode,
		OTLPEndpoint:   os.Getenv("RUNE_OTLP_ENDPOINT"), // Allow override via env var
	}

	// Check for additional OTLP configuration
	if otelConfig.OTLPEndpoint == "" {
		// Default to localhost for development, but can be overridden
		if otelConfig.Environment == "development" {
			otelConfig.OTLPEndpoint = "http://localhost:4318/v1/logs"
		} else {
			// Disable OTLP by default in production unless explicitly configured
			otelConfig.DisableOTLP = true
		}
	}

	otelLogger, err := otellogger.Initialize(otelConfig)
	if err != nil {
		if debugMode {
			slog.Default().Debug("otel init failed", "error", err)
		}
		slog.Default().Warn("otel logging init failed - using local logging only")
		// Continue without OpenTelemetry but keep basic telemetry working
		client.otelLogger = nil
	} else {
		client.otelLogger = otelLogger
		if debugMode {
			slog.Default().Debug("otel logging initialized successfully")
		}
	}

	return client
}

func (c *Client) Track(event string, properties map[string]interface{}) {
	if !c.enabled {
		if os.Getenv("RUNE_DEBUG") == "true" {
			slog.Default().Debug("telemetry disabled - not tracking event", "event", event)
		}
		return
	}

	if os.Getenv("RUNE_DEBUG") == "true" {
		slog.Default().Debug("tracking event", "event", event)
	}

	// Convert properties to slog attributes
	attrs := []interface{}{"event_name", event}
	for key, value := range properties {
		attrs = append(attrs, key, value)
	}

	// Add system context
	attrs = append(attrs,
		"app_name", "rune",
		"app_version", getVersion(),
		"os_name", runtime.GOOS,
		"os_version", getOSVersion(),
		"user_id", c.userID,
		"session_id", c.sessionID,
	)

	// Log with OpenTelemetry
	if c.otelLogger != nil {
		c.otelLogger.LogEvent(slog.LevelInfo, event, attrs...)
	}
}

func (c *Client) TrackError(err error, command string, properties map[string]interface{}) {
	if !c.enabled {
		return
	}

	// Convert properties to slog attributes
	attrs := []interface{}{"command", command}
	for key, value := range properties {
		attrs = append(attrs, key, value)
	}

	// Add system context
	attrs = append(attrs,
		"user_id", c.userID,
		"session_id", c.sessionID,
	)

	// Log with OpenTelemetry
	if c.otelLogger != nil {
		message := fmt.Sprintf("Command failed: %s", command)
		c.otelLogger.LogError(err, message, attrs...)
	}
}

func (c *Client) TrackCommand(command string, duration time.Duration, success bool) {
	if !c.enabled {
		return
	}

	// Debug console log for tests and local visibility
	if os.Getenv("RUNE_DEBUG") == "true" {
		slog.Default().Debug("tracking event", "event", "command_executed", "command", command, "success", success, "duration_ms", duration.Milliseconds())
	}

	// Emit an OpenTelemetry log event for command execution
	if c.otelLogger != nil {
		c.otelLogger.LogEvent(slog.LevelInfo, "command_executed",
			"command", command,
			"duration_ms", duration.Milliseconds(),
			"success", success,
		)
	}

	// Add performance monitoring to Sentry
	if c.sentryEnabled {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("command", command)
			scope.SetTag("success", fmt.Sprintf("%t", success))
			scope.SetExtra("duration_ms", duration.Milliseconds())

			// Create a transaction for performance monitoring
			ctx := sentry.SetHubOnContext(context.Background(), sentry.CurrentHub())
			transaction := sentry.StartTransaction(ctx, fmt.Sprintf("command.%s", command))
			transaction.SetTag("command", command)
			transaction.SetTag("success", fmt.Sprintf("%t", success))
			transaction.SetData("duration_ms", duration.Milliseconds())

			if !success {
				transaction.Status = sentry.SpanStatusInternalError
				sentry.CaptureMessage(fmt.Sprintf("Command failed: %s", command))
			} else {
				transaction.Status = sentry.SpanStatusOK
			}
			transaction.Finish()
		})
	}
}

func (c *Client) Close() {
	// Flush Sentry events
	if c.sentryEnabled {
		if os.Getenv("RUNE_DEBUG") == "true" {
			slog.Default().Debug("flushing Sentry events before close")
		}
		sentry.Flush(5 * time.Second) // Increased timeout for better reliability
	}

	// Shutdown OpenTelemetry logger provider if initialized
	if c.otelLogger != nil {
		if os.Getenv("RUNE_DEBUG") == "true" {
			slog.Default().Debug("shutting down OpenTelemetry logger provider")
		}
		_ = c.otelLogger.Close()
		c.otelLogger = nil
	}
}

// StartTransaction starts a Sentry transaction for performance monitoring
func (c *Client) StartTransaction(name, operation string) *sentry.Span {
	if !c.sentryEnabled {
		return nil
	}

	ctx := sentry.SetHubOnContext(context.Background(), sentry.CurrentHub())
	return sentry.StartTransaction(ctx, name)
}

// CaptureException captures an exception with additional context
func (c *Client) CaptureException(err error, tags map[string]string, extra map[string]interface{}) {
	if !c.sentryEnabled {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		for key, value := range extra {
			scope.SetExtra(key, value)
		}
		sentry.CaptureException(err)
	})
}

// CaptureMessage captures a message with additional context
func (c *Client) CaptureMessage(message string, level sentry.Level, tags map[string]string) {
	if !c.sentryEnabled {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		sentry.CaptureMessage(message)
	})
}

func getUserID() string {
	// Try to get from config first
	cfg, err := config.Load()
	if err == nil && cfg.UserID != "" {
		return cfg.UserID
	}

	// Generate a new anonymous ID
	userID := generateAnonymousID()

	// Try to save it to config
	if cfg != nil {
		cfg.UserID = userID
		_ = config.SaveConfig(cfg) // Ignore errors
	}

	return userID
}

func generateAnonymousID() string {
	// Simple anonymous ID generation
	hostname, _ := os.Hostname()
	return fmt.Sprintf("anon_%s_%d", hostname, time.Now().Unix())
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// StartCommand starts tracking a command execution for release health
func (c *Client) StartCommand(command string) {
	if !c.sentryEnabled {
		return
	}

	// Add breadcrumb for command start
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message:   fmt.Sprintf("Command started: %s", command),
		Category:  "command",
		Level:     sentry.LevelInfo,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"command": command,
			"action":  "start",
		},
	})
}

// EndCommand ends tracking a command execution
func (c *Client) EndCommand(command string, success bool, duration time.Duration) {
	if !c.sentryEnabled {
		return
	}

	// Add breadcrumb for command end
	level := sentry.LevelInfo
	if !success {
		level = sentry.LevelError
	}

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message:   fmt.Sprintf("Command %s: %s", map[bool]string{true: "completed", false: "failed"}[success], command),
		Category:  "command",
		Level:     level,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"command":     command,
			"action":      "end",
			"success":     success,
			"duration_ms": duration.Milliseconds(),
		},
	})
}

// maskKey masks sensitive keys for logging purposes
func maskKey(key string) string {
	if key == "" {
		return "[not provided]"
	}
	if len(key) < 8 {
		return "[masked]"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// Build-time version variable set via ldflags
var version string

func getVersion() string {
	if version != "" {
		return version
	}
	return "dev"
}

func getEnvironment() string {
	if env := os.Getenv("RUNE_ENV"); env != "" {
		return env
	}
	return "production"
}

func getOSVersion() string {
	switch runtime.GOOS {
	case "darwin":
		return getMacOSVersion()
	case "linux":
		return getLinuxVersion()
	case "windows":
		return getWindowsVersion()
	default:
		return "unknown"
	}
}

func getMacOSVersion() string {
	// Simple implementation - you might want to use a more robust method
	return "unknown"
}

func getLinuxVersion() string {
	// Simple implementation - you might want to use a more robust method
	return "unknown"
}

func getWindowsVersion() string {
	// Simple implementation - you might want to use a more robust method
	return "unknown"
}
