package telemetry

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOTLPIntegration placeholder for future OTLP log delivery tests using a mock collector
func TestOTLPIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Placeholder: In future, spin up a mock OTLP HTTP endpoint to assert log ingestion
	// For now, we'll test the client creation and method calls
	client := NewClient("")
	require.NotNil(t, client)

	// Test tracking an event
	client.Track("integration_test_event", map[string]interface{}{
		"test_property": "test_value",
		"timestamp":     time.Now().Unix(),
	})

	// Test tracking a command
	client.TrackCommand("test_command", 100*time.Millisecond, true)

	// Test tracking an error
	client.TrackError(fmt.Errorf("test error"), "test_command", map[string]interface{}{
		"context": "integration_test",
	})

	// Close the client to flush events
	client.Close()

	// In a real integration test, we would verify events were received
	// For now, we just ensure no panics occurred
}

// TestSentryIntegration tests Sentry integration with environment variables
func TestSentryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Only run if we have a test Sentry DSN
	testDSN := os.Getenv("RUNE_TEST_SENTRY_DSN")
	if testDSN == "" {
		t.Skip("RUNE_TEST_SENTRY_DSN not set, skipping Sentry integration test")
	}

	client := NewClient(testDSN)
	require.NotNil(t, client)
	assert.True(t, client.sentryEnabled)

	// Test capturing an exception
	client.CaptureException(
		fmt.Errorf("integration test error"),
		map[string]string{"test": "integration"},
		map[string]interface{}{"timestamp": time.Now().Unix()},
	)

	// Test starting a transaction
	transaction := client.StartTransaction("integration_test", "test")
	if transaction != nil {
		transaction.Finish()
	}

	// Test command tracking
	client.StartCommand("integration_test_command")
	client.EndCommand("integration_test_command", true, 50*time.Millisecond)

	// Close to flush events
	client.Close()
}

// TestTelemetryWithRealKeys tests with actual telemetry endpoints/keys if available
func TestTelemetryWithRealKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	otlpEndpoint := os.Getenv("RUNE_TEST_OTLP_ENDPOINT")
	sentryDSN := os.Getenv("RUNE_TEST_SENTRY_DSN")

	if otlpEndpoint == "" && sentryDSN == "" {
		t.Skip("No test telemetry configuration provided, skipping real integration test")
	}

	// If an OTLP test endpoint is provided, set it for this test run
	var restoreOTLP string
	if otlpEndpoint != "" {
		restoreOTLP = os.Getenv("RUNE_OTLP_ENDPOINT")
		_ = os.Setenv("RUNE_OTLP_ENDPOINT", otlpEndpoint)
		defer func() {
			if restoreOTLP == "" {
				_ = os.Unsetenv("RUNE_OTLP_ENDPOINT")
			} else {
				_ = os.Setenv("RUNE_OTLP_ENDPOINT", restoreOTLP)
			}
		}()
	}

	client := NewClient(sentryDSN)
	require.NotNil(t, client)

	// Send a test event
	client.Track("integration_test_real_keys", map[string]interface{}{
		"test_run":   true,
		"timestamp":  time.Now().Unix(),
		"go_version": "test",
	})

	// Test command execution tracking
	client.TrackCommand("integration_test_command", 25*time.Millisecond, true)

	// Test error tracking
	client.TrackError(
		fmt.Errorf("integration test error - this is expected"),
		"integration_test_command",
		map[string]interface{}{
			"test_context": "real_keys_test",
		},
	)

	// Close to ensure events are sent
	client.Close()

	t.Log("Integration test completed - check your OTLP collector and/or Sentry project for events")
}

// TestMiddlewareIntegration tests the middleware functions
func TestMiddlewareIntegration(t *testing.T) {
	// Initialize global client for testing
	Initialize("")
	defer Close()

	// Test global functions
	Track("middleware_test_event", map[string]interface{}{
		"source": "middleware_test",
	})

	TrackCommand("middleware_test_command", 10*time.Millisecond, true)

	TrackError(fmt.Errorf("middleware test error"), "test_command", map[string]interface{}{
		"middleware": true,
	})

	StartCommand("middleware_command")
	EndCommand("middleware_command", true, 5*time.Millisecond)
}

// TestEventProperties tests that events contain expected properties
func TestEventProperties(t *testing.T) {
	client := NewClient("")
	require.NotNil(t, client)

	// We can't easily test the actual properties sent to external services
	// without mocking, but we can test that the methods don't panic
	// and handle various input types correctly

	testCases := []struct {
		name       string
		event      string
		properties map[string]interface{}
	}{
		{
			name:  "string properties",
			event: "test_string_props",
			properties: map[string]interface{}{
				"string_prop": "test_value",
				"app_name":    "rune",
			},
		},
		{
			name:  "mixed properties",
			event: "test_mixed_props",
			properties: map[string]interface{}{
				"string_prop": "test",
				"int_prop":    42,
				"bool_prop":   true,
				"float_prop":  3.14,
				"time_prop":   time.Now(),
			},
		},
		{
			name:       "nil properties",
			event:      "test_nil_props",
			properties: nil,
		},
		{
			name:       "empty properties",
			event:      "test_empty_props",
			properties: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			client.Track(tc.event, tc.properties)
		})
	}
}
