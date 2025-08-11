package otellogger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// OtelLoggerConfig holds configuration for OpenTelemetry logging
type OtelLoggerConfig struct {
	SentryDSN       string            // Sentry DSN for error tracking
	OTLPEndpoint    string            // Optional OTLP endpoint URL (http://localhost:4318/v1/logs by default)
	ServiceName     string            // Service name for telemetry
	ServiceVersion  string            // Service version
	Environment     string            // Environment (production, development, etc.)
	ExtraAttributes map[string]string // Additional resource attributes
	DisableSentry   bool              // Disable Sentry integration
	DisableOTLP     bool              // Disable OTLP export
	DebugMode       bool              // Enable debug logging
}

// OtelLogger wraps OpenTelemetry logging with Sentry integration
type OtelLogger struct {
	config       *OtelLoggerConfig
	logger       *slog.Logger
	logProvider  *sdklog.LoggerProvider
	sentryClient *sentry.Client
	enabled      bool
}

var globalLogger *OtelLogger

// Initialize creates and configures the global OpenTelemetry logger with Sentry integration
func Initialize(config *OtelLoggerConfig) (*OtelLogger, error) {
	if config == nil {
		config = &OtelLoggerConfig{}
	}

	// Set defaults
	if config.ServiceName == "" {
		config.ServiceName = "rune"
	}
	if config.ServiceVersion == "" {
		config.ServiceVersion = getVersion()
	}
	if config.Environment == "" {
		config.Environment = getEnvironment()
	}
	if config.OTLPEndpoint == "" {
		config.OTLPEndpoint = "http://localhost:4318/v1/logs"
	}

	// Check if telemetry is globally disabled
	enabled := os.Getenv("RUNE_TELEMETRY_DISABLED") != "true"
	if !enabled {
		if config.DebugMode {
			slog.Default().Debug("OpenTelemetry disabled via env", "env", "RUNE_TELEMETRY_DISABLED")
		}
		return &OtelLogger{config: config, enabled: false}, nil
	}

	logger := &OtelLogger{
		config:  config,
		enabled: enabled,
	}

	if err := logger.initializeOpenTelemetry(); err != nil {
		return nil, fmt.Errorf("failed to initialize OpenTelemetry: %w", err)
	}

	globalLogger = logger
	return logger, nil
}

// initializeOpenTelemetry sets up the OpenTelemetry logging infrastructure
func (ol *OtelLogger) initializeOpenTelemetry() error {
	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(ol.config.ServiceName),
			semconv.ServiceVersionKey.String(ol.config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(ol.config.Environment),
			semconv.OSNameKey.String(runtime.GOOS),
			semconv.OSVersionKey.String(getOSVersion()),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Add extra attributes if provided
	if len(ol.config.ExtraAttributes) > 0 {
		attrs := make([]attribute.KeyValue, 0, len(ol.config.ExtraAttributes))
		for key, value := range ol.config.ExtraAttributes {
			attrs = append(attrs, attribute.String(key, value))
		}
		res, err = resource.Merge(res, resource.NewWithAttributes("", attrs...))
		if err != nil {
			return fmt.Errorf("failed to merge extra attributes: %w", err)
		}
	}

	// Create exporters
	var exporters []sdklog.Exporter

	// Initialize Sentry integration if enabled and DSN provided
	if !ol.config.DisableSentry && ol.config.SentryDSN != "" {
		if ol.config.DebugMode {
			slog.Default().Debug("Initializing Sentry OpenTelemetry integration")
		}

		// Initialize Sentry with OpenTelemetry integration
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              ol.config.SentryDSN,
			Environment:      ol.config.Environment,
			Release:          ol.config.ServiceVersion,
			AttachStacktrace: true,
			EnableTracing:    true,
			TracesSampleRate: 1.0,
			EnableLogs:       true,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				// Add additional context
				event.Contexts["app"] = map[string]interface{}{
					"name":    ol.config.ServiceName,
					"version": ol.config.ServiceVersion,
				}
				return event
			},
		})
		if err != nil {
			if ol.config.DebugMode {
				slog.Default().Debug("Sentry initialization failed", "error", err)
			}
			slog.Default().Warn("Sentry integration failed - continuing with local logging only")
		} else {
			ol.sentryClient = sentry.CurrentHub().Client()
			if ol.config.DebugMode {
				slog.Default().Debug("Sentry OpenTelemetry integration initialized successfully")
			}
		}
	}

	// Initialize OTLP exporter if enabled and endpoint provided
	if !ol.config.DisableOTLP && ol.config.OTLPEndpoint != "" {
		if ol.config.DebugMode {
			slog.Default().Debug("Initializing OTLP log exporter", "endpoint", ol.config.OTLPEndpoint)
		}

		// Choose the correct option based on the value provided.
		// - If a full URL (http/https) is provided, use WithEndpointURL
		// - Otherwise, treat it as host[:port] and use WithEndpoint
		// Note: WithURLPath defaults to /v1/logs when not specified.
		opts := []otlploghttp.Option{}

		ep := ol.config.OTLPEndpoint
		isURL := strings.HasPrefix(ep, "http://") || strings.HasPrefix(ep, "https://")
		if isURL {
			opts = append(opts, otlploghttp.WithEndpointURL(ep))
			if strings.HasPrefix(ep, "http://") {
				// Insecure only for http
				opts = append(opts, otlploghttp.WithInsecure())
			}
		} else {
			opts = append(opts, otlploghttp.WithEndpoint(ep))
		}

		otlpExporter, err := otlploghttp.New(context.Background(), opts...)
		if err != nil {
			if ol.config.DebugMode {
				slog.Default().Debug("OTLP exporter initialization failed", "error", err)
			}
			slog.Default().Info("OTLP log export not available", "endpoint", ol.config.OTLPEndpoint)
		} else {
			exporters = append(exporters, otlpExporter)
			if ol.config.DebugMode {
				slog.Default().Debug("OTLP log exporter initialized successfully")
			}
		}
	}
	// Create the logger provider with all exporters
	opts := []sdklog.LoggerProviderOption{
		sdklog.WithResource(res),
	}

	// Add processors for each exporter
	for _, exporter := range exporters {
		opts = append(opts, sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)))
	}

	ol.logProvider = sdklog.NewLoggerProvider(opts...)

	// Set as global logger provider
	global.SetLoggerProvider(ol.logProvider)

	// Create slog logger with OpenTelemetry bridge
	ol.logger = otelslog.NewLogger("rune", otelslog.WithLoggerProvider(ol.logProvider))

	if ol.config.DebugMode {
		slog.Default().Debug("OpenTelemetry logging initialized", "exporters", len(exporters))
	}

	// Set up OpenTelemetry trace provider with Sentry if enabled
	if ol.sentryClient != nil {
		// Create a tracer provider with Sentry span processor
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())
		if ol.config.DebugMode {
			slog.Default().Debug("Sentry OpenTelemetry tracer provider set")
		}
	}

	return nil
}

// GetLogger returns the underlying slog.Logger with OpenTelemetry integration
func (ol *OtelLogger) GetLogger() *slog.Logger {
	if ol.logger == nil {
		// Fallback to default slog if not initialized
		return slog.Default()
	}
	return ol.logger
}

// Close shuts down the OpenTelemetry logger and flushes pending logs
func (ol *OtelLogger) Close() error {
	if !ol.enabled {
		return nil
	}

	var err error

	// Shutdown log provider
	if ol.logProvider != nil {
		if shutdownErr := ol.logProvider.Shutdown(context.Background()); shutdownErr != nil {
			err = fmt.Errorf("failed to shutdown log provider: %w", shutdownErr)
		}
	}

	// Flush Sentry
	if ol.sentryClient != nil {
		sentry.Flush(5)
	}

	return err
}

// GetGlobalLogger returns the global OpenTelemetry logger instance
func GetGlobalLogger() *OtelLogger {
	return globalLogger
}

// CloseGlobal closes the global OpenTelemetry logger
func CloseGlobal() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

// LogError logs an error with OpenTelemetry attributes and sends to Sentry if configured
func (ol *OtelLogger) LogError(err error, message string, attrs ...interface{}) {
	if !ol.enabled {
		return
	}

	// Add error to attributes
	allAttrs := append([]interface{}{"error", err.Error(), "error_type", fmt.Sprintf("%T", err)}, attrs...)

	// Log with slog
	ol.logger.Error(message, allAttrs...)

	// Send to Sentry if configured
	if ol.sentryClient != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)

			// Add attributes as extra data
			for i := 0; i < len(allAttrs)-1; i += 2 {
				if key, ok := allAttrs[i].(string); ok {
					scope.SetExtra(key, allAttrs[i+1])
				}
			}

			sentry.CaptureException(err)
		})
	}
}

// LogEvent logs a structured event with OpenTelemetry
func (ol *OtelLogger) LogEvent(level slog.Level, message string, attrs ...interface{}) {
	if !ol.enabled {
		return
	}

	ol.logger.Log(context.Background(), level, message, attrs...)

	// Add as breadcrumb to Sentry if configured
	if ol.sentryClient != nil && level >= slog.LevelWarn {
		var sentryLevel sentry.Level
		switch level {
		case slog.LevelWarn:
			sentryLevel = sentry.LevelWarning
		case slog.LevelError:
			sentryLevel = sentry.LevelError
		default:
			sentryLevel = sentry.LevelInfo
		}

		data := make(map[string]interface{})
		for i := 0; i < len(attrs)-1; i += 2 {
			if key, ok := attrs[i].(string); ok {
				data[key] = attrs[i+1]
			}
		}

		sentry.AddBreadcrumb(&sentry.Breadcrumb{
			Message:  message,
			Category: "otel.log",
			Level:    sentryLevel,
			Data:     data,
		})
	}
}

// Helper functions
func getVersion() string {
	if version := os.Getenv("RUNE_VERSION"); version != "" {
		return version
	}
	return "dev"
}

func getEnvironment() string {
	if env := os.Getenv("RUNE_ENV"); env != "" {
		return env
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		return env
	}
	return "production"
}

func getOSVersion() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}
