package dnd

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

// hasDisplay returns true if a graphical display is available
func hasDisplay() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}

// isHeadless returns true if running in a headless environment
func isHeadless() bool {
	return !hasDisplay()
}

func TestPackageExists(t *testing.T) {
	// Basic test to ensure package compiles
	t.Log("DND package compiles successfully")
}

func TestNewDNDManager(t *testing.T) {
	nm := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(nm)

	if dndManager == nil {
		t.Fatal("NewDNDManager returned nil")
	}

	if dndManager.notificationManager != nm {
		t.Error("DND manager should store the notification manager reference")
	}
}

func TestDNDManagerWithoutNotifications(t *testing.T) {
	dndManager := NewDNDManager(nil)

	if dndManager == nil {
		t.Fatal("NewDNDManager returned nil")
	}

	// These should not panic even with nil notification manager
	err := dndManager.SendBreakNotification(30 * time.Minute)
	if err != nil {
		t.Errorf("SendBreakNotification should not error with nil notification manager: %v", err)
	}

	err = dndManager.SendEndOfDayNotification(8*time.Hour, 8.0)
	if err != nil {
		t.Errorf("SendEndOfDayNotification should not error with nil notification manager: %v", err)
	}

	err = dndManager.SendSessionCompleteNotification(2*time.Hour, "test")
	if err != nil {
		t.Errorf("SendSessionCompleteNotification should not error with nil notification manager: %v", err)
	}

	err = dndManager.SendIdleNotification(10 * time.Minute)
	if err != nil {
		t.Errorf("SendIdleNotification should not error with nil notification manager: %v", err)
	}
}

func TestDNDManagerWithNotifications(t *testing.T) {
	nm := notifications.NewNotificationManager(false) // Disabled to avoid actual notifications in tests
	dndManager := NewDNDManager(nm)

	// These should not error (notifications are disabled)
	err := dndManager.SendBreakNotification(30 * time.Minute)
	if err != nil {
		t.Errorf("SendBreakNotification failed: %v", err)
	}

	err = dndManager.SendEndOfDayNotification(8*time.Hour, 8.0)
	if err != nil {
		t.Errorf("SendEndOfDayNotification failed: %v", err)
	}

	err = dndManager.SendSessionCompleteNotification(2*time.Hour, "test")
	if err != nil {
		t.Errorf("SendSessionCompleteNotification failed: %v", err)
	}

	err = dndManager.SendIdleNotification(10 * time.Minute)
	if err != nil {
		t.Errorf("SendIdleNotification failed: %v", err)
	}
}

// TestDNDManagerCrossPlatform tests the cross-platform DND functionality.
// This test requires a graphical environment to properly test enable/disable.
func TestDNDManagerCrossPlatform(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_disable_cycle", func(t *testing.T) {
		if isHeadless() {
			t.Skip("Skipping DnD enable/disable test: no display available (headless environment)")
		}

		// Get initial state for cleanup
		initialState, err := dndManager.IsEnabled()
		if err != nil {
			t.Skipf("Cannot determine initial DnD state, skipping: %v", err)
		}

		// Test Enable
		if err := dndManager.Enable(); err != nil {
			t.Fatalf("Enable() failed: %v", err)
		}

		// Verify Enable actually worked
		enabled, err := dndManager.IsEnabled()
		if err != nil {
			t.Errorf("IsEnabled() failed after Enable(): %v", err)
		} else if !enabled {
			t.Errorf("DnD should be enabled after Enable(), but IsEnabled() returned false")
		}

		// Test Disable
		if err := dndManager.Disable(); err != nil {
			t.Fatalf("Disable() failed: %v", err)
		}

		// Verify Disable actually worked
		enabled, err = dndManager.IsEnabled()
		if err != nil {
			t.Errorf("IsEnabled() failed after Disable(): %v", err)
		} else if enabled {
			t.Errorf("DnD should be disabled after Disable(), but IsEnabled() returned true")
		}

		// Cleanup: restore initial state
		if initialState {
			_ = dndManager.Enable()
		}
	})

	t.Run("shortcuts_setup", func(t *testing.T) {
		// This test can run on any platform - it just checks if shortcuts are available
		available, err := dndManager.CheckShortcutsSetup()
		if err != nil {
			// On non-macOS platforms, this may return an error which is expected
			if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
				t.Skipf("CheckShortcutsSetup not supported on %s: %v", runtime.GOOS, err)
			}
			t.Errorf("CheckShortcutsSetup failed: %v", err)
		}
		t.Logf("Shortcuts available on %s: %v", runtime.GOOS, available)
	})

	t.Run("test_notifications", func(t *testing.T) {
		if isHeadless() {
			t.Skip("Skipping notification test: no display available (headless environment)")
		}

		// Skip if notify-send not available on Linux
		if runtime.GOOS == "linux" {
			if _, err := exec.LookPath("notify-send"); err != nil {
				t.Skip("Skipping notification test: notify-send not installed")
			}
		}

		err := dndManager.TestNotifications()
		if err != nil {
			t.Errorf("TestNotifications failed: %v", err)
		}
	})
}

// TestWindowsSpecificMethods tests Windows-specific DND methods.
// These tests should only run on Windows.
func TestWindowsSpecificMethods(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("windows_focus_assist_methods_exist", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows Focus Assist test on non-Windows platform")
		}

		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledWindowsFocusAssist()

		err := dndManager.enableWindowsFocusAssist()
		if err != nil {
			t.Fatalf("enableWindowsFocusAssist failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Errorf("isEnabledWindowsFocusAssist failed: %v", err)
		} else if !enabled {
			t.Errorf("Focus Assist should be enabled after enableWindowsFocusAssist()")
		}

		err = dndManager.disableWindowsFocusAssist()
		if err != nil {
			t.Fatalf("disableWindowsFocusAssist failed: %v", err)
		}

		// Verify disable worked
		enabled, err = dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Errorf("isEnabledWindowsFocusAssist failed after disable: %v", err)
		} else if enabled {
			t.Errorf("Focus Assist should be disabled after disableWindowsFocusAssist()")
		}

		// Cleanup: restore initial state
		if initialEnabled {
			_ = dndManager.enableWindowsFocusAssist()
		}
	})

	t.Run("windows_notification_methods_exist", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows notification settings test on non-Windows platform")
		}

		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledWindowsNotifications()

		err := dndManager.enableWindowsNotificationSettings()
		if err != nil {
			t.Fatalf("enableWindowsNotificationSettings failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Errorf("isEnabledWindowsNotifications failed: %v", err)
		} else if !enabled {
			t.Errorf("Notification DnD should be enabled after enableWindowsNotificationSettings()")
		}

		err = dndManager.disableWindowsNotificationSettings()
		if err != nil {
			t.Fatalf("disableWindowsNotificationSettings failed: %v", err)
		}

		// Verify disable worked
		enabled, err = dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Errorf("isEnabledWindowsNotifications failed after disable: %v", err)
		} else if enabled {
			t.Errorf("Notification DnD should be disabled after disableWindowsNotificationSettings()")
		}

		// Cleanup: restore initial state
		if initialEnabled {
			_ = dndManager.enableWindowsNotificationSettings()
		}
	})

	t.Run("windows_action_center_methods_exist", func(t *testing.T) {
		// These methods are not yet implemented, so they should return specific errors
		err := dndManager.enableWindowsActionCenter()
		if err != nil {
			// This should always fail as it's not implemented yet
			if !strings.Contains(err.Error(), "not yet implemented") {
				t.Errorf("Expected 'not yet implemented' error, got: %v", err)
			}
		}

		err = dndManager.disableWindowsActionCenter()
		if err != nil {
			if !strings.Contains(err.Error(), "not yet implemented") {
				t.Errorf("Expected 'not yet implemented' error, got: %v", err)
			}
		}

		_, err = dndManager.isEnabledWindowsWinRT()
		if err != nil {
			if !strings.Contains(err.Error(), "not yet implemented") {
				t.Errorf("Expected 'not yet implemented' error, got: %v", err)
			}
		}
	})
}
