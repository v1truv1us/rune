package telemetry

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name            string
		sentryDSN       string
		envDisabled     string
		expectedEnabled bool
	}{
		{
			name:            "valid sentry dsn",
			sentryDSN:       "https://test@sentry.io/123",
			expectedEnabled: true,
		},
		{
			name:            "empty dsn",
			sentryDSN:       "",
			expectedEnabled: true, // Still enabled, just no external services
		},
		{
			name:            "telemetry disabled via env",
			sentryDSN:       "https://test@sentry.io/123",
			envDisabled:     "true",
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if specified
			if tt.envDisabled != "" {
				os.Setenv("RUNE_TELEMETRY_DISABLED", tt.envDisabled)
				defer os.Unsetenv("RUNE_TELEMETRY_DISABLED")
			}

			client := NewClient(tt.sentryDSN)

			require.NotNil(t, client)
			assert.Equal(t, tt.expectedEnabled, client.enabled)
			assert.NotEmpty(t, client.userID)
			assert.NotEmpty(t, client.sessionID)
		})
	}
}

func TestClientTrack(t *testing.T) {
	tests := []struct {
		name       string
		enabled    bool
		event      string
		properties map[string]interface{}
	}{
		{
			name:    "enabled client with properties",
			enabled: true,
			event:   "test_event",
			properties: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		},
		{
			name:       "enabled client without properties",
			enabled:    true,
			event:      "test_event_no_props",
			properties: nil,
		},
		{
			name:    "disabled client",
			enabled: false,
			event:   "test_event",
			properties: map[string]interface{}{
				"key1": "value1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				enabled:   tt.enabled,
				userID:    "test_user",
				sessionID: "test_session",
			}

			// This should not panic
			client.Track(tt.event, tt.properties)
		})
	}
}

func TestClientTrackError(t *testing.T) {
	client := &Client{
		enabled:   true,
		userID:    "test_user",
		sessionID: "test_session",
	}

	err := assert.AnError
	command := "test_command"
	properties := map[string]interface{}{
		"extra": "data",
	}

	// This should not panic
	client.TrackError(err, command, properties)
}

func TestClientTrackCommand(t *testing.T) {
	client := &Client{
		enabled:   true,
		userID:    "test_user",
		sessionID: "test_session",
	}

	command := "test_command"
	duration := 100 * time.Millisecond
	success := true

	// This should not panic
	client.TrackCommand(command, duration, success)
}

func TestClientClose(t *testing.T) {
	client := &Client{
		enabled:   true,
		userID:    "test_user",
		sessionID: "test_session",
	}

	// This should not panic
	client.Close()
}

func TestGenerateAnonymousID(t *testing.T) {
	id1 := generateAnonymousID()
	time.Sleep(1 * time.Second) // Ensure different Unix timestamps
	id2 := generateAnonymousID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2) // Should be different due to timestamp
	assert.Contains(t, id1, "anon_")
}

func TestGenerateSessionID(t *testing.T) {
	id1 := generateSessionID()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	id2 := generateSessionID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2) // Should be different due to timestamp
	assert.Contains(t, id1, "session_")
}

func TestGetVersion(t *testing.T) {
	// Test default version
	version = ""
	assert.Equal(t, "dev", getVersion())

	// Test set version
	version = "1.0.0"
	assert.Equal(t, "1.0.0", getVersion())
}

func TestGetEnvironment(t *testing.T) {
	// Test default environment
	os.Unsetenv("RUNE_ENV")
	assert.Equal(t, "production", getEnvironment())

	// Test custom environment
	os.Setenv("RUNE_ENV", "development")
	defer os.Unsetenv("RUNE_ENV")
	assert.Equal(t, "development", getEnvironment())
}

func TestClientStartEndCommand(t *testing.T) {
	client := &Client{
		enabled:       true,
		sentryEnabled: false, // Disable Sentry for unit tests
		userID:        "test_user",
		sessionID:     "test_session",
	}

	command := "test_command"
	duration := 50 * time.Millisecond

	// Test start command
	client.StartCommand(command)

	// Test end command - success
	client.EndCommand(command, true, duration)

	// Test end command - failure
	client.EndCommand(command, false, duration)
}

func TestClientCaptureException(t *testing.T) {
	client := &Client{
		enabled:       true,
		sentryEnabled: false, // Disable Sentry for unit tests
		userID:        "test_user",
		sessionID:     "test_session",
	}

	err := assert.AnError
	tags := map[string]string{
		"component": "test",
	}
	extra := map[string]interface{}{
		"data": "test_data",
	}

	// This should not panic
	client.CaptureException(err, tags, extra)
}

func TestClientStartTransaction(t *testing.T) {
	client := &Client{
		enabled:       true,
		sentryEnabled: false, // Disable Sentry for unit tests
		userID:        "test_user",
		sessionID:     "test_session",
	}

	// Should return nil when Sentry is disabled
	transaction := client.StartTransaction("test_transaction", "test_operation")
	assert.Nil(t, transaction)
}

func TestMaskKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "[not provided]",
		},
		{
			name:     "Very short key (1 char)",
			input:    "a",
			expected: "[masked]",
		},
		{
			name:     "Short key (7 chars)",
			input:    "1234567",
			expected: "[masked]",
		},
		{
			name:     "Exactly 8 chars",
			input:    "12345678",
			expected: "1234****5678",
		},
		{
			name:     "Long key",
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: "abcd****wxyz",
		},
		{
			name:     "API key format",
			input:    "test_live_abc123def456ghi789jkl012mno345pqr678stu901vwx234",
			expected: "test****x234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
