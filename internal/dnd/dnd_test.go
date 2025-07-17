package dnd

import (
	"strings"
	"testing"
	"time"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

func TestPackageExists(t *testing.T) {
	// Basic test to ensure package compiles
	// More comprehensive tests should be added as functionality is implemented
	t.Log("DND package test placeholder")
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

func TestDNDManagerCrossPlatform(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	t.Run("enable_disable_cycle", func(t *testing.T) {
		// Test enable/disable cycle - should not panic
		err := dndManager.Enable()
		if err != nil {
			t.Logf("Enable failed (expected on some platforms): %v", err)
		}

		enabled, err := dndManager.IsEnabled()
		if err != nil {
			t.Logf("IsEnabled failed (expected on some platforms): %v", err)
		} else {
			t.Logf("DND enabled status: %v", enabled)
		}

		err = dndManager.Disable()
		if err != nil {
			t.Logf("Disable failed (expected on some platforms): %v", err)
		}
	})

	t.Run("shortcuts_setup", func(t *testing.T) {
		available, err := dndManager.CheckShortcutsSetup()
		if err != nil {
			t.Logf("CheckShortcutsSetup failed (expected on some platforms): %v", err)
		} else {
			t.Logf("Shortcuts available: %v", available)
		}
	})

	t.Run("test_notifications", func(t *testing.T) {
		err := dndManager.TestNotifications()
		if err != nil {
			t.Logf("TestNotifications failed (expected with disabled notifications): %v", err)
		}
	})
}

func TestWindowsSpecificMethods(t *testing.T) {
	nm := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(nm)

	// Test Windows-specific methods (these should work on any platform for testing)
	t.Run("windows_focus_assist_methods_exist", func(t *testing.T) {
		// These methods should exist and be callable, even if they fail on non-Windows platforms
		err := dndManager.enableWindowsFocusAssist()
		if err != nil {
			t.Logf("enableWindowsFocusAssist failed (expected on non-Windows): %v", err)
		}

		err = dndManager.disableWindowsFocusAssist()
		if err != nil {
			t.Logf("disableWindowsFocusAssist failed (expected on non-Windows): %v", err)
		}

		_, err = dndManager.isEnabledWindowsFocusAssist()
		if err != nil {
			t.Logf("isEnabledWindowsFocusAssist failed (expected on non-Windows): %v", err)
		}
	})

	t.Run("windows_notification_methods_exist", func(t *testing.T) {
		err := dndManager.enableWindowsNotificationSettings()
		if err != nil {
			t.Logf("enableWindowsNotificationSettings failed (expected on non-Windows): %v", err)
		}

		err = dndManager.disableWindowsNotificationSettings()
		if err != nil {
			t.Logf("disableWindowsNotificationSettings failed (expected on non-Windows): %v", err)
		}

		_, err = dndManager.isEnabledWindowsNotifications()
		if err != nil {
			t.Logf("isEnabledWindowsNotifications failed (expected on non-Windows): %v", err)
		}
	})

	t.Run("windows_action_center_methods_exist", func(t *testing.T) {
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
