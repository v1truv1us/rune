package dnd

import (
	"runtime"
	"strings"
	"testing"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

// TestWindowsFocusAssistIntegration tests Windows Focus Assist functionality.
// Requires Windows with Focus Assist support.
func TestWindowsFocusAssistIntegration(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows Focus Assist tests only run on Windows")
	}

	nm := notifications.NewNotificationManager(false) // Disabled to avoid actual notifications
	dndManager := NewDNDManager(nm)

	t.Run("enable_focus_assist", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledWindowsFocusAssist()

		err := dndManager.enableWindowsFocusAssist()
		if err != nil {
			t.Fatalf("Focus Assist enable failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Errorf("Failed to check Focus Assist status: %v", err)
		} else if !enabled {
			t.Errorf("Focus Assist should be enabled after enableWindowsFocusAssist()")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableWindowsFocusAssist()
		}
	})

	t.Run("check_focus_assist_status", func(t *testing.T) {
		enabled, err := dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Errorf("Focus Assist status check failed: %v", err)
		}
		t.Logf("Focus Assist status: %v", enabled)
	})

	t.Run("disable_focus_assist", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableWindowsFocusAssist()

		err := dndManager.disableWindowsFocusAssist()
		if err != nil {
			t.Fatalf("Focus Assist disable failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Errorf("Failed to check Focus Assist status after disable: %v", err)
		} else if enabled {
			t.Errorf("Focus Assist should be disabled after disableWindowsFocusAssist()")
		}
	})
}

// TestWindowsNotificationSettingsIntegration tests Windows notification settings.
func TestWindowsNotificationSettingsIntegration(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows notification tests only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_notification_dnd", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledWindowsNotifications()

		err := dndManager.enableWindowsNotificationSettings()
		if err != nil {
			t.Fatalf("Notification DND enable failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Errorf("Failed to check notification DND status: %v", err)
		} else if !enabled {
			t.Errorf("Notification DND should be enabled after enableWindowsNotificationSettings()")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableWindowsNotificationSettings()
		}
	})

	t.Run("check_notification_status", func(t *testing.T) {
		enabled, err := dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Errorf("Notification status check failed: %v", err)
		}
		t.Logf("Notification DND status: %v", enabled)
	})

	t.Run("disable_notification_dnd", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableWindowsNotificationSettings()

		err := dndManager.disableWindowsNotificationSettings()
		if err != nil {
			t.Fatalf("Notification DND disable failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Errorf("Failed to check notification DND status after disable: %v", err)
		} else if enabled {
			t.Errorf("Notification DND should be disabled after disableWindowsNotificationSettings()")
		}
	})
}

// TestWindowsDNDMethods tests the main Windows DND methods.
func TestWindowsDNDMethods(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows DND tests only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_windows_dnd", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledWindows()

		err := dndManager.enableWindows()
		if err != nil {
			// Check if the error message is informative
			if !strings.Contains(err.Error(), "Focus Assist") {
				t.Errorf("Expected informative error message about Focus Assist, got: %v", err)
			}
			t.Fatalf("Windows DND enable failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindows()
		if err != nil {
			t.Errorf("Failed to check Windows DND status: %v", err)
		} else if !enabled {
			t.Errorf("Windows DND should be enabled after enableWindows()")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableWindows()
		}
	})

	t.Run("check_windows_dnd_status", func(t *testing.T) {
		enabled, err := dndManager.isEnabledWindows()
		if err != nil {
			t.Errorf("Windows DND status check failed: %v", err)
		}
		t.Logf("Windows DND status: %v", enabled)
	})

	t.Run("disable_windows_dnd", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableWindows()

		err := dndManager.disableWindows()
		if err != nil {
			if !strings.Contains(err.Error(), "Focus Assist") {
				t.Errorf("Expected informative error message about Focus Assist, got: %v", err)
			}
			t.Fatalf("Windows DND disable failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledWindows()
		if err != nil {
			t.Errorf("Failed to check Windows DND status after disable: %v", err)
		} else if enabled {
			t.Errorf("Windows DND should be disabled after disableWindows()")
		}
	})
}

// TestWindowsDNDCrossPlatform tests cross-platform DND behavior on Windows.
func TestWindowsDNDCrossPlatform(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_dnd_cross_platform", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows cross-platform DND test on non-Windows platform")
		}

		// Get initial state for cleanup
		initialEnabled, _ := dndManager.IsEnabled()

		err := dndManager.Enable()
		if err != nil {
			t.Fatalf("DND enable failed on Windows: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.IsEnabled()
		if err != nil {
			t.Errorf("Failed to check DND status: %v", err)
		} else if !enabled {
			t.Errorf("DND should be enabled after Enable()")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.Disable()
		}
	})

	t.Run("check_dnd_status_cross_platform", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows cross-platform DND status test on non-Windows platform")
		}

		enabled, err := dndManager.IsEnabled()
		if err != nil {
			t.Errorf("DND status check failed on Windows: %v", err)
		}
		t.Logf("DND status on Windows: %v", enabled)
	})

	t.Run("disable_dnd_cross_platform", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows cross-platform DND disable test on non-Windows platform")
		}

		// First enable to ensure we have something to disable
		_ = dndManager.Enable()

		err := dndManager.Disable()
		if err != nil {
			t.Fatalf("DND disable failed on Windows: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.IsEnabled()
		if err != nil {
			t.Errorf("Failed to check DND status after disable: %v", err)
		} else if enabled {
			t.Errorf("DND should be disabled after Disable()")
		}
	})
}

// TestWindowsShortcutsSetup tests the shortcuts setup check on Windows.
func TestWindowsShortcutsSetup(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("check_shortcuts_setup", func(t *testing.T) {
		available, err := dndManager.CheckShortcutsSetup()
		if err != nil {
			t.Errorf("CheckShortcutsSetup failed: %v", err)
		}

		if runtime.GOOS == "windows" {
			// Windows should always return true since it doesn't use shortcuts
			if !available {
				t.Error("Windows should always report shortcuts as available")
			}
		}

		t.Logf("Shortcuts available on %s: %v", runtime.GOOS, available)
	})
}

// TestWindowsDNDErrorHandling tests error handling and edge cases.
func TestWindowsDNDErrorHandling(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("test_with_nil_notification_manager", func(t *testing.T) {
		dndManagerNil := NewDNDManager(nil)

		// These should not panic
		err := dndManagerNil.SendBreakNotification(0)
		if err != nil {
			t.Errorf("SendBreakNotification should not error with nil manager: %v", err)
		}

		err = dndManagerNil.SendEndOfDayNotification(0, 0)
		if err != nil {
			t.Errorf("SendEndOfDayNotification should not error with nil manager: %v", err)
		}

		err = dndManagerNil.SendSessionCompleteNotification(0, "test")
		if err != nil {
			t.Errorf("SendSessionCompleteNotification should not error with nil manager: %v", err)
		}

		err = dndManagerNil.SendIdleNotification(0)
		if err != nil {
			t.Errorf("SendIdleNotification should not error with nil manager: %v", err)
		}

		err = dndManagerNil.TestNotifications()
		if err == nil {
			t.Error("TestNotifications should error with nil manager")
		}
	})

	t.Run("test_notification_methods", func(t *testing.T) {
		// Test that notification methods work properly with disabled notifications
		err := dndManager.SendBreakNotification(30 * 60 * 1000000000) // 30 minutes in nanoseconds
		if err != nil {
			t.Errorf("SendBreakNotification failed: %v", err)
		}

		err = dndManager.SendEndOfDayNotification(8*60*60*1000000000, 8.0) // 8 hours
		if err != nil {
			t.Errorf("SendEndOfDayNotification failed: %v", err)
		}

		err = dndManager.SendSessionCompleteNotification(2*60*60*1000000000, "test-project") // 2 hours
		if err != nil {
			t.Errorf("SendSessionCompleteNotification failed: %v", err)
		}

		err = dndManager.SendIdleNotification(10 * 60 * 1000000000) // 10 minutes
		if err != nil {
			t.Errorf("SendIdleNotification failed: %v", err)
		}
	})
}

// Benchmark tests for Windows DND operations

func BenchmarkWindowsFocusAssistEnable(b *testing.B) {
	if runtime.GOOS != "windows" {
		b.Skip("Windows Focus Assist benchmarks only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dndManager.enableWindowsFocusAssist()
	}
}

func BenchmarkWindowsFocusAssistStatus(b *testing.B) {
	if runtime.GOOS != "windows" {
		b.Skip("Windows Focus Assist benchmarks only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = dndManager.isEnabledWindowsFocusAssist()
	}
}

func BenchmarkWindowsDNDCrossPlatform(b *testing.B) {
	if runtime.GOOS != "windows" {
		b.Skip("Windows DND benchmarks only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = dndManager.IsEnabled()
	}
}
