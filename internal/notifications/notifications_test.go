package notifications

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

// hasDisplay returns true if a graphical display is available
func hasDisplay() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}

// hasCommand returns true if the given command is available in PATH
func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func TestNotificationManager_Creation(t *testing.T) {
	nm := NewNotificationManager(true)
	if nm == nil {
		t.Fatal("NewNotificationManager returned nil")
	}
	if !nm.enabled {
		t.Error("Expected notifications to be enabled")
	}

	nmDisabled := NewNotificationManager(false)
	if nmDisabled.enabled {
		t.Error("Expected notifications to be disabled")
	}
}

func TestNotificationManager_DisabledSend(t *testing.T) {
	nm := NewNotificationManager(false)

	notification := Notification{
		Title:    "Test",
		Message:  "Test message",
		Type:     Custom,
		Priority: Normal,
		Sound:    false,
	}

	// Should not return an error when disabled
	err := nm.Send(notification)
	if err != nil {
		t.Errorf("Expected no error when notifications disabled, got: %v", err)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "30 seconds"},
		{90 * time.Second, "1 minutes"},
		{5 * time.Minute, "5 minutes"},
		{65 * time.Minute, "1 hours 5 minutes"},
		{2 * time.Hour, "2 hours"},
		{150 * time.Minute, "2 hours 30 minutes"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %q, expected %q", test.duration, result, test.expected)
		}
	}
}

func TestIsSupported(t *testing.T) {
	// This will depend on the platform the test is running on
	supported := IsSupported()

	// We can't assert a specific value since it depends on the platform,
	// but we can ensure the function doesn't panic
	t.Logf("Notifications supported on this platform: %v", supported)
}

// TestNotificationTypes tests that notifications can actually be sent.
// This test requires a graphical environment and appropriate notification tools.
func TestNotificationTypes(t *testing.T) {
	// Skip if no display available (headless environment)
	if !hasDisplay() {
		t.Skip("Skipping notification tests: no display available (headless environment)")
	}

	// Skip if notify-send not available on Linux
	if runtime.GOOS == "linux" {
		if !hasCommand("notify-send") {
			t.Skip("Skipping notification tests: notify-send not installed (install libnotify)")
		}
	}

	// Skip if terminal-notifier not available on macOS (osascript is fallback)
	if runtime.GOOS == "darwin" {
		// osascript is always available on macOS, so we can proceed
		// terminal-notifier is preferred but not required
		if !hasCommand("terminal-notifier") {
			t.Log("terminal-notifier not found, using osascript fallback")
		}
	}

	nm := NewNotificationManager(true)

	t.Run("BreakReminder", func(t *testing.T) {
		err := nm.SendBreakReminder(45 * time.Minute)
		if err != nil {
			t.Errorf("SendBreakReminder failed: %v", err)
		}
	})

	t.Run("EndOfDayReminder", func(t *testing.T) {
		err := nm.SendEndOfDayReminder(7*time.Hour+30*time.Minute, 8.0)
		if err != nil {
			t.Errorf("SendEndOfDayReminder failed: %v", err)
		}
	})

	t.Run("SessionComplete", func(t *testing.T) {
		err := nm.SendSessionComplete(2*time.Hour, "test-project")
		if err != nil {
			t.Errorf("SendSessionComplete failed: %v", err)
		}
	})

	t.Run("IdleDetected", func(t *testing.T) {
		err := nm.SendIdleDetected(10 * time.Minute)
		if err != nil {
			t.Errorf("SendIdleDetected failed: %v", err)
		}
	})
}

func TestGetSoundName(t *testing.T) {
	nm := NewNotificationManager(true)

	tests := []struct {
		notification Notification
		expected     string
	}{
		{
			Notification{Sound: false, Priority: Normal},
			"",
		},
		{
			Notification{Sound: true, Priority: Critical},
			"Basso",
		},
		{
			Notification{Sound: true, Priority: High},
			"Ping",
		},
		{
			Notification{Sound: true, Priority: Normal},
			"default",
		},
	}

	for _, test := range tests {
		result := nm.getSoundName(test.notification)
		if result != test.expected {
			t.Errorf("getSoundName(%+v) = %q, expected %q", test.notification, result, test.expected)
		}
	}
}

func TestGetUrgencyLevel(t *testing.T) {
	nm := NewNotificationManager(true)

	tests := []struct {
		priority Priority
		expected string
	}{
		{Critical, "critical"},
		{High, "normal"},
		{Normal, "normal"},
		{Low, "low"},
	}

	for _, test := range tests {
		result := nm.getUrgencyLevel(test.priority)
		if result != test.expected {
			t.Errorf("getUrgencyLevel(%v) = %q, expected %q", test.priority, result, test.expected)
		}
	}
}

// TestTestNotification tests the TestNotification method.
func TestTestNotification(t *testing.T) {
	// Skip if no display available (headless environment)
	if !hasDisplay() {
		t.Skip("Skipping test notification: no display available (headless environment)")
	}

	// Skip if notify-send not available on Linux
	if runtime.GOOS == "linux" {
		if !hasCommand("notify-send") {
			t.Skip("Skipping test notification: notify-send not installed")
		}
	}

	nm := NewNotificationManager(true)
	err := nm.TestNotification()
	if err != nil {
		t.Errorf("TestNotification failed: %v", err)
	}
}
