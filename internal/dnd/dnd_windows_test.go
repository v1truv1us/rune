package dnd

import (
	"runtime"
	"strings"
	"testing"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

func TestWindowsFocusAssistIntegration(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows Focus Assist tests only run on Windows")
	}

	nm := notifications.NewNotificationManager(false) // Disabled to avoid actual notifications
	dndManager := NewDNDManager(nm)

	t.Run("enable_focus_assist", func(t *testing.T) {
		err := dndManager.enableWindowsFocusAssist()
		if err != nil {
			// This might fail if we don't have proper permissions or if Focus Assist isn't available
			t.Logf("Focus Assist enable failed (expected on some systems): %v", err)
		} else {
			t.Log("Focus Assist enabled successfully")
		}
	})

	t.Run("check_focus_assist_status", func(t *testing.T) {
		enabled, err := dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Logf("Focus Assist status check failed (expected on some systems): %v", err)
		} else {
			t.Logf("Focus Assist status: %v", enabled)
		}
	})

	t.Run("disable_focus_assist", func(t *testing.T) {
		err := dndManager.disableWindowsFocusAssist()
		if err != nil {
			t.Logf("Focus Assist disable failed (expected on some systems): %v", err)
		} else {
			t.Log("Focus Assist disabled successfully")
		}
	})
}

func TestWindowsNotificationSettingsIntegration(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows notification tests only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_notification_dnd", func(t *testing.T) {
		err := dndManager.enableWindowsNotificationSettings()
		if err != nil {
			t.Logf("Notification DND enable failed (expected on some systems): %v", err)
		} else {
			t.Log("Notification DND enabled successfully")
		}
	})

	t.Run("check_notification_status", func(t *testing.T) {
		enabled, err := dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Logf("Notification status check failed (expected on some systems): %v", err)
		} else {
			t.Logf("Notification DND status: %v", enabled)
		}
	})

	t.Run("disable_notification_dnd", func(t *testing.T) {
		err := dndManager.disableWindowsNotificationSettings()
		if err != nil {
			t.Logf("Notification DND disable failed (expected on some systems): %v", err)
		} else {
			t.Log("Notification DND disabled successfully")
		}
	})
}

func TestWindowsDNDMethods(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows DND tests only run on Windows")
	}

	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_windows_dnd", func(t *testing.T) {
		err := dndManager.enableWindows()
		// This should try multiple methods and may succeed with at least one
		if err != nil {
			// Check if the error message is informative
			if !strings.Contains(err.Error(), "Focus Assist") {
				t.Errorf("Expected informative error message about Focus Assist, got: %v", err)
			}
			t.Logf("Windows DND enable failed (expected on some systems): %v", err)
		} else {
			t.Log("Windows DND enabled successfully")
		}
	})

	t.Run("check_windows_dnd_status", func(t *testing.T) {
		enabled, err := dndManager.isEnabledWindows()
		if err != nil {
			t.Logf("Windows DND status check failed: %v", err)
		} else {
			t.Logf("Windows DND status: %v", enabled)
		}
	})

	t.Run("disable_windows_dnd", func(t *testing.T) {
		err := dndManager.disableWindows()
		if err != nil {
			if !strings.Contains(err.Error(), "Focus Assist") {
				t.Errorf("Expected informative error message about Focus Assist, got: %v", err)
			}
			t.Logf("Windows DND disable failed (expected on some systems): %v", err)
		} else {
			t.Log("Windows DND disabled successfully")
		}
	})
}

func TestWindowsDNDCrossPlatform(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_dnd_cross_platform", func(t *testing.T) {
		err := dndManager.Enable()
		if runtime.GOOS == "windows" {
			// On Windows, this should attempt to enable Focus Assist
			if err != nil {
				t.Logf("DND enable failed on Windows (expected on some systems): %v", err)
			} else {
				t.Log("DND enabled successfully on Windows")
			}
		} else {
			// On other platforms, test behavior may vary
			t.Logf("DND enable on %s: %v", runtime.GOOS, err)
		}
	})

	t.Run("check_dnd_status_cross_platform", func(t *testing.T) {
		enabled, err := dndManager.IsEnabled()
		if err != nil {
			t.Logf("DND status check failed on %s: %v", runtime.GOOS, err)
		} else {
			t.Logf("DND status on %s: %v", runtime.GOOS, enabled)
		}
	})

	t.Run("disable_dnd_cross_platform", func(t *testing.T) {
		err := dndManager.Disable()
		if runtime.GOOS == "windows" {
			if err != nil {
				t.Logf("DND disable failed on Windows (expected on some systems): %v", err)
			} else {
				t.Log("DND disabled successfully on Windows")
			}
		} else {
			t.Logf("DND disable on %s: %v", runtime.GOOS, err)
		}
	})
}

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

// Test error handling and edge cases
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
		// Test that notification methods work properly
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
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = dndManager.IsEnabled()
	}
}
