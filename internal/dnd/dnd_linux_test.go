package dnd

import (
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

// canAccessGNOMESettings returns true if GNOME settings are accessible
func canAccessGNOMESettings() bool {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.notifications", "show-banners")
	return cmd.Run() == nil
}

// canAccessMATESettings returns true if MATE settings are accessible
func canAccessMATESettings() bool {
	cmd := exec.Command("gsettings", "get", "org.mate.NotificationDaemon", "popup-location")
	return cmd.Run() == nil
}

// canAccessCinnamonSettings returns true if Cinnamon settings are accessible
func canAccessCinnamonSettings() bool {
	cmd := exec.Command("gsettings", "get", "org.cinnamon.desktop.notifications", "display-notifications")
	return cmd.Run() == nil
}

// TestLinuxDNDIntegration tests Linux DND functionality.
// Requires a graphical environment with a detectable desktop environment.
func TestLinuxDNDIntegration(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux DND tests on non-Linux system")
	}

	if isHeadless() {
		t.Skip("Skipping Linux DND tests: no display available (headless environment)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("DetectDesktopEnvironment", func(t *testing.T) {
		de := dndManager.detectDesktopEnvironment()
		t.Logf("Detected desktop environment: %s", de)

		// Should return a valid desktop environment or "unknown"
		validDEs := []string{"gnome", "kde", "xfce", "mate", "cinnamon", "i3", "sway", "dwm", "awesome", "bspwm", "unknown"}
		found := false
		for _, validDE := range validDEs {
			if de == validDE {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("detectDesktopEnvironment returned invalid DE: %s", de)
		}

		// If we have a display but DE is unknown, that's a detection failure
		if de == "unknown" {
			t.Errorf("Desktop environment detection failed - returned 'unknown'. " +
				"Ensure XDG_CURRENT_DESKTOP or similar env vars are set, or a DE process is running.")
		}
	})

	t.Run("EnableDisableLinux", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, err := dndManager.isEnabledLinux()
		if err != nil {
			t.Skipf("Cannot determine initial DnD state: %v", err)
		}

		// Test enable
		if err := dndManager.enableLinux(); err != nil {
			t.Fatalf("enableLinux() failed: %v", err)
		}

		// Verify enable worked
		enabled, err := dndManager.isEnabledLinux()
		if err != nil {
			t.Errorf("isEnabledLinux() failed after enable: %v", err)
		} else if !enabled {
			t.Errorf("DnD should be enabled after enableLinux(), but isEnabledLinux() returned false")
		}

		// Test disable
		if err := dndManager.disableLinux(); err != nil {
			t.Fatalf("disableLinux() failed: %v", err)
		}

		// Verify disable worked
		enabled, err = dndManager.isEnabledLinux()
		if err != nil {
			t.Errorf("isEnabledLinux() failed after disable: %v", err)
		} else if enabled {
			t.Errorf("DnD should be disabled after disableLinux(), but isEnabledLinux() returned true")
		}

		// Cleanup: restore initial state
		if initialEnabled {
			_ = dndManager.enableLinux()
		}
	})

	t.Run("IsEnabledLinux", func(t *testing.T) {
		enabled, err := dndManager.isEnabledLinux()
		if err != nil {
			t.Errorf("isEnabledLinux() failed: %v", err)
		}
		t.Logf("DND enabled status: %v", enabled)
	})
}

// TestGNOMEDND tests GNOME-specific DND functionality.
func TestGNOMEDND(t *testing.T) {
	if !hasCommand("gsettings") {
		t.Skip("gsettings not available, skipping GNOME DND tests")
	}

	if isHeadless() {
		t.Skip("Skipping GNOME DND tests: no display available (headless environment)")
	}

	if !canAccessGNOMESettings() {
		t.Skip("GNOME notification schema not accessible (not running GNOME or no dbus session)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableGNOME", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledGNOME()

		if err := dndManager.enableGNOME(); err != nil {
			t.Fatalf("enableGNOME() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledGNOME()
		if err != nil {
			t.Errorf("isEnabledGNOME() failed: %v", err)
		} else if !enabled {
			t.Errorf("GNOME DnD should be enabled after enableGNOME(), but isEnabledGNOME() returned false")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableGNOME()
		}
	})

	t.Run("DisableGNOME", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableGNOME()

		if err := dndManager.disableGNOME(); err != nil {
			t.Fatalf("disableGNOME() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledGNOME()
		if err != nil {
			t.Errorf("isEnabledGNOME() failed: %v", err)
		} else if enabled {
			t.Errorf("GNOME DnD should be disabled after disableGNOME(), but isEnabledGNOME() returned true")
		}
	})

	t.Run("IsEnabledGNOME", func(t *testing.T) {
		enabled, err := dndManager.isEnabledGNOME()
		if err != nil {
			t.Errorf("isEnabledGNOME() failed: %v", err)
		}
		t.Logf("GNOME DND enabled: %v", enabled)
	})
}

// TestKDEDND tests KDE-specific DND functionality.
func TestKDEDND(t *testing.T) {
	hasKDE5 := hasCommand("kwriteconfig5")
	hasKDE6 := hasCommand("kwriteconfig6")

	if !hasKDE5 && !hasKDE6 {
		t.Skip("KDE config tools not available, skipping KDE DND tests")
	}

	if isHeadless() {
		t.Skip("Skipping KDE DND tests: no display available (headless environment)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableKDE", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledKDE()

		if err := dndManager.enableKDE(); err != nil {
			t.Fatalf("enableKDE() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledKDE()
		if err != nil {
			t.Errorf("isEnabledKDE() failed: %v", err)
		} else if !enabled {
			t.Errorf("KDE DnD should be enabled after enableKDE(), but isEnabledKDE() returned false")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableKDE()
		}
	})

	t.Run("DisableKDE", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableKDE()

		if err := dndManager.disableKDE(); err != nil {
			t.Fatalf("disableKDE() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledKDE()
		if err != nil {
			t.Errorf("isEnabledKDE() failed: %v", err)
		} else if enabled {
			t.Errorf("KDE DnD should be disabled after disableKDE(), but isEnabledKDE() returned true")
		}
	})

	t.Run("IsEnabledKDE", func(t *testing.T) {
		enabled, err := dndManager.isEnabledKDE()
		if err != nil {
			t.Errorf("isEnabledKDE() failed: %v", err)
		}
		t.Logf("KDE DND enabled: %v", enabled)
	})
}

// TestXFCEDND tests XFCE-specific DND functionality.
func TestXFCEDND(t *testing.T) {
	if !hasCommand("xfconf-query") {
		t.Skip("xfconf-query not available, skipping XFCE DND tests")
	}

	if isHeadless() {
		t.Skip("Skipping XFCE DND tests: no display available (headless environment)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableXFCE", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledXFCE()

		if err := dndManager.enableXFCE(); err != nil {
			t.Fatalf("enableXFCE() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledXFCE()
		if err != nil {
			t.Errorf("isEnabledXFCE() failed: %v", err)
		} else if !enabled {
			t.Errorf("XFCE DnD should be enabled after enableXFCE(), but isEnabledXFCE() returned false")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableXFCE()
		}
	})

	t.Run("DisableXFCE", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableXFCE()

		if err := dndManager.disableXFCE(); err != nil {
			t.Fatalf("disableXFCE() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledXFCE()
		if err != nil {
			t.Errorf("isEnabledXFCE() failed: %v", err)
		} else if enabled {
			t.Errorf("XFCE DnD should be disabled after disableXFCE(), but isEnabledXFCE() returned true")
		}
	})

	t.Run("IsEnabledXFCE", func(t *testing.T) {
		enabled, err := dndManager.isEnabledXFCE()
		if err != nil {
			t.Errorf("isEnabledXFCE() failed: %v", err)
		}
		t.Logf("XFCE DND enabled: %v", enabled)
	})
}

// TestMATEDND tests MATE-specific DND functionality.
func TestMATEDND(t *testing.T) {
	if !hasCommand("gsettings") {
		t.Skip("gsettings not available, skipping MATE DND tests")
	}

	if isHeadless() {
		t.Skip("Skipping MATE DND tests: no display available (headless environment)")
	}

	if !canAccessMATESettings() {
		t.Skip("MATE notification schema not accessible (not running MATE or no dbus session)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableMATE", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledMATE()

		if err := dndManager.enableMATE(); err != nil {
			t.Fatalf("enableMATE() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledMATE()
		if err != nil {
			t.Errorf("isEnabledMATE() failed: %v", err)
		} else if !enabled {
			t.Errorf("MATE DnD should be enabled after enableMATE(), but isEnabledMATE() returned false")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableMATE()
		}
	})

	t.Run("DisableMATE", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableMATE()

		if err := dndManager.disableMATE(); err != nil {
			t.Fatalf("disableMATE() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledMATE()
		if err != nil {
			t.Errorf("isEnabledMATE() failed: %v", err)
		} else if enabled {
			t.Errorf("MATE DnD should be disabled after disableMATE(), but isEnabledMATE() returned true")
		}
	})

	t.Run("IsEnabledMATE", func(t *testing.T) {
		enabled, err := dndManager.isEnabledMATE()
		if err != nil {
			t.Errorf("isEnabledMATE() failed: %v", err)
		}
		t.Logf("MATE DND enabled: %v", enabled)
	})
}

// TestCinnamonDND tests Cinnamon-specific DND functionality.
func TestCinnamonDND(t *testing.T) {
	if !hasCommand("gsettings") {
		t.Skip("gsettings not available, skipping Cinnamon DND tests")
	}

	if isHeadless() {
		t.Skip("Skipping Cinnamon DND tests: no display available (headless environment)")
	}

	if !canAccessCinnamonSettings() {
		t.Skip("Cinnamon notification schema not accessible (not running Cinnamon or no dbus session)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableCinnamon", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledCinnamon()

		if err := dndManager.enableCinnamon(); err != nil {
			t.Fatalf("enableCinnamon() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledCinnamon()
		if err != nil {
			t.Errorf("isEnabledCinnamon() failed: %v", err)
		} else if !enabled {
			t.Errorf("Cinnamon DnD should be enabled after enableCinnamon(), but isEnabledCinnamon() returned false")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableCinnamon()
		}
	})

	t.Run("DisableCinnamon", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableCinnamon()

		if err := dndManager.disableCinnamon(); err != nil {
			t.Fatalf("disableCinnamon() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledCinnamon()
		if err != nil {
			t.Errorf("isEnabledCinnamon() failed: %v", err)
		} else if enabled {
			t.Errorf("Cinnamon DnD should be disabled after disableCinnamon(), but isEnabledCinnamon() returned true")
		}
	})

	t.Run("IsEnabledCinnamon", func(t *testing.T) {
		enabled, err := dndManager.isEnabledCinnamon()
		if err != nil {
			t.Errorf("isEnabledCinnamon() failed: %v", err)
		}
		t.Logf("Cinnamon DND enabled: %v", enabled)
	})
}

// TestTilingWMDND tests tiling window manager DND functionality.
func TestTilingWMDND(t *testing.T) {
	hasDunst := hasCommand("dunstctl")
	hasMako := hasCommand("makoctl")

	if !hasDunst && !hasMako {
		t.Skip("No supported notification daemon found (dunstctl or makoctl), skipping tiling WM DND tests")
	}

	if isHeadless() {
		t.Skip("Skipping tiling WM DND tests: no display available (headless environment)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableTilingWM", func(t *testing.T) {
		// Get initial state for cleanup
		initialEnabled, _ := dndManager.isEnabledTilingWM()

		if err := dndManager.enableTilingWM(); err != nil {
			t.Fatalf("enableTilingWM() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledTilingWM()
		if err != nil {
			t.Errorf("isEnabledTilingWM() failed: %v", err)
		} else if !enabled {
			t.Errorf("Tiling WM DnD should be enabled after enableTilingWM(), but isEnabledTilingWM() returned false")
		}

		// Cleanup: restore initial state
		if !initialEnabled {
			_ = dndManager.disableTilingWM()
		}
	})

	t.Run("DisableTilingWM", func(t *testing.T) {
		// First enable to ensure we have something to disable
		_ = dndManager.enableTilingWM()

		if err := dndManager.disableTilingWM(); err != nil {
			t.Fatalf("disableTilingWM() failed: %v", err)
		}

		// Verify it worked
		enabled, err := dndManager.isEnabledTilingWM()
		if err != nil {
			t.Errorf("isEnabledTilingWM() failed: %v", err)
		} else if enabled {
			t.Errorf("Tiling WM DnD should be disabled after disableTilingWM(), but isEnabledTilingWM() returned true")
		}
	})

	t.Run("IsEnabledTilingWM", func(t *testing.T) {
		enabled, err := dndManager.isEnabledTilingWM()
		if err != nil {
			t.Errorf("isEnabledTilingWM() failed: %v", err)
		}
		t.Logf("Tiling WM DND enabled: %v", enabled)
	})
}

// TestGenericLinuxDND tests the generic Linux DND fallback.
// This should always work as it's a best-effort approach.
func TestGenericLinuxDND(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping generic Linux DND tests on non-Linux system")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableGenericLinux", func(t *testing.T) {
		// Generic Linux enable should not fail (it's best-effort)
		err := dndManager.enableGenericLinux()
		if err != nil {
			t.Errorf("Generic Linux enable should not fail: %v", err)
		}
	})

	t.Run("DisableGenericLinux", func(t *testing.T) {
		// Generic Linux disable should not fail (it's best-effort)
		err := dndManager.disableGenericLinux()
		if err != nil {
			t.Errorf("Generic Linux disable should not fail: %v", err)
		}
	})
}

// TestDesktopEnvironmentDetection tests the desktop environment detection logic.
func TestDesktopEnvironmentDetection(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping desktop environment detection tests on non-Linux system")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	testCases := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name: "GNOME",
			envVars: map[string]string{
				"XDG_CURRENT_DESKTOP": "GNOME",
			},
			expected: "gnome",
		},
		{
			name: "KDE",
			envVars: map[string]string{
				"XDG_CURRENT_DESKTOP": "KDE",
			},
			expected: "kde",
		},
		{
			name: "XFCE",
			envVars: map[string]string{
				"XDG_CURRENT_DESKTOP": "XFCE",
			},
			expected: "xfce",
		},
		{
			name: "i3",
			envVars: map[string]string{
				"XDG_CURRENT_DESKTOP": "i3",
			},
			expected: "i3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original environment
			originalEnv := make(map[string]string)
			for key := range tc.envVars {
				originalEnv[key] = os.Getenv(key)
			}

			// Set test environment
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			// Test detection
			detected := dndManager.detectDesktopEnvironment()

			// Restore original environment
			for key, value := range originalEnv {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}

			if detected != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, detected)
			}
		})
	}
}

// Helper functions

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func hasCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Benchmark tests for performance

func BenchmarkDetectDesktopEnvironment(b *testing.B) {
	if runtime.GOOS != "linux" {
		b.Skip("Skipping Linux benchmark on non-Linux system")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dndManager.detectDesktopEnvironment()
	}
}

func BenchmarkEnableDisableLinux(b *testing.B) {
	if runtime.GOOS != "linux" {
		b.Skip("Skipping Linux benchmark on non-Linux system")
	}

	if isHeadless() {
		b.Skip("Skipping Linux benchmark: no display available (headless environment)")
	}

	notificationManager := notifications.NewNotificationManager(false)
	dndManager := NewDNDManager(notificationManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dndManager.enableLinux()
		_ = dndManager.disableLinux()
	}
}
