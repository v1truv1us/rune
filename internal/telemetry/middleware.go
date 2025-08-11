package telemetry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

var globalClient *Client

// Initialize sets up the global telemetry client
func Initialize(sentryDSN string) {
	globalClient = NewClient(sentryDSN)
}

// TrackCommand tracks command execution
func TrackCommand(command string, duration time.Duration, success bool) {
	if globalClient != nil {
		globalClient.TrackCommand(command, duration, success)
	}
}

// TrackError tracks errors
func TrackError(err error, command string, properties map[string]interface{}) {
	if globalClient != nil {
		globalClient.TrackError(err, command, properties)
	}
}

// Track tracks custom events
func Track(event string, properties map[string]interface{}) {
	if globalClient != nil {
		globalClient.Track(event, properties)
	}
}

// Close closes the global telemetry client
func Close() {
	if globalClient != nil {
		globalClient.Close()
		globalClient = nil
	}
}

// WrapCommand wraps a cobra command with telemetry tracking
func WrapCommand(cmd *cobra.Command, originalRun func(cmd *cobra.Command, args []string) error) {
	if originalRun == nil {
		return
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		commandName := cmd.CommandPath()

		// Start command tracking for release health
		StartCommand(commandName)

		err := originalRun(cmd, args)

		duration := time.Since(start)
		success := err == nil

		// End command tracking for release health
		EndCommand(commandName, success, duration)

		TrackCommand(commandName, duration, success)

		if err != nil {
			TrackError(err, commandName, map[string]interface{}{
				"args": args,
			})
		}

		return err
	}
}

// WrapCommandNoError wraps a cobra command that doesn't return an error
func WrapCommandNoError(cmd *cobra.Command, originalRun func(cmd *cobra.Command, args []string)) {
	if originalRun == nil {
		return
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		start := time.Now()
		commandName := cmd.CommandPath()

		originalRun(cmd, args)

		duration := time.Since(start)
		TrackCommand(commandName, duration, true)
	}
}

// StartTransaction starts a Sentry transaction for performance monitoring
func StartTransaction(name, operation string) *sentry.Span {
	if globalClient != nil {
		return globalClient.StartTransaction(name, operation)
	}
	return nil
}

// CaptureException captures an exception with additional context
func CaptureException(err error, tags map[string]string, extra map[string]interface{}) {
	if globalClient != nil {
		globalClient.CaptureException(err, tags, extra)
	}
}

// CaptureMessage captures a message with additional context
func CaptureMessage(message string, level sentry.Level, tags map[string]string) {
	if globalClient != nil {
		globalClient.CaptureMessage(message, level, tags)
	}
}

// StartCommand starts tracking a command for release health
func StartCommand(command string) {
	if globalClient != nil {
		globalClient.StartCommand(command)
	}
}

// EndCommand ends tracking a command for release health
func EndCommand(command string, success bool, duration time.Duration) {
	if globalClient != nil {
		globalClient.EndCommand(command, success, duration)
	}
}
