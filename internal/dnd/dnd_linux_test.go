package dnd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

func TestLinuxDNDIntegration(t *testing.T) {
	// Skip if not on Linux
	if os.Getenv("GOOS") != "linux" && !isLinux() {
		t.Skip("Skipping Linux DND tests on non-Linux system")
	}

	notificationManager := notifications.NewNotificationManager(true)
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
	})

	t.Run("EnableDisableLinux", func(t *testing.T) {
		// Test enable
		err := dndManager.enableLinux()
		if err != nil {
			t.Logf("Enable failed (expected on some systems): %v", err)
		}

		// Test disable
		err = dndManager.disableLinux()
		if err != nil {
			t.Logf("Disable failed (expected on some systems): %v", err)
		}
	})

	t.Run("IsEnabledLinux", func(t *testing.T) {
		enabled, err := dndManager.isEnabledLinux()
		if err != nil {
			t.Logf("IsEnabled check failed (expected on some systems): %v", err)
		} else {
			t.Logf("DND enabled status: %v", enabled)
		}
	})
}

func TestGNOMEDND(t *testing.T) {
	if !hasCommand("gsettings") {
		t.Skip("gsettings not available, skipping GNOME DND tests")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableGNOME", func(t *testing.T) {
		err := dndManager.enableGNOME()
		if err != nil {
			t.Logf("GNOME enable failed (expected if not on GNOME): %v", err)
		}
	})

	t.Run("DisableGNOME", func(t *testing.T) {
		err := dndManager.disableGNOME()
		if err != nil {
			t.Logf("GNOME disable failed (expected if not on GNOME): %v", err)
		}
	})

	t.Run("IsEnabledGNOME", func(t *testing.T) {
		enabled, err := dndManager.isEnabledGNOME()
		if err != nil {
			t.Logf("GNOME status check failed (expected if not on GNOME): %v", err)
		} else {
			t.Logf("GNOME DND enabled: %v", enabled)
		}
	})
}

func TestKDEDND(t *testing.T) {
	hasKDE5 := hasCommand("kwriteconfig5")
	hasKDE6 := hasCommand("kwriteconfig6")

	if !hasKDE5 && !hasKDE6 {
		t.Skip("KDE config tools not available, skipping KDE DND tests")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableKDE", func(t *testing.T) {
		err := dndManager.enableKDE()
		if err != nil {
			t.Logf("KDE enable failed (expected if not on KDE): %v", err)
		}
	})

	t.Run("DisableKDE", func(t *testing.T) {
		err := dndManager.disableKDE()
		if err != nil {
			t.Logf("KDE disable failed (expected if not on KDE): %v", err)
		}
	})

	t.Run("IsEnabledKDE", func(t *testing.T) {
		enabled, err := dndManager.isEnabledKDE()
		if err != nil {
			t.Logf("KDE status check failed (expected if not on KDE): %v", err)
		} else {
			t.Logf("KDE DND enabled: %v", enabled)
		}
	})
}

func TestXFCEDND(t *testing.T) {
	if !hasCommand("xfconf-query") {
		t.Skip("xfconf-query not available, skipping XFCE DND tests")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableXFCE", func(t *testing.T) {
		err := dndManager.enableXFCE()
		if err != nil {
			t.Logf("XFCE enable failed (expected if not on XFCE): %v", err)
		}
	})

	t.Run("DisableXFCE", func(t *testing.T) {
		err := dndManager.disableXFCE()
		if err != nil {
			t.Logf("XFCE disable failed (expected if not on XFCE): %v", err)
		}
	})

	t.Run("IsEnabledXFCE", func(t *testing.T) {
		enabled, err := dndManager.isEnabledXFCE()
		if err != nil {
			t.Logf("XFCE status check failed (expected if not on XFCE): %v", err)
		} else {
			t.Logf("XFCE DND enabled: %v", enabled)
		}
	})
}

func TestMATEDND(t *testing.T) {
	if !hasCommand("gsettings") {
		t.Skip("gsettings not available, skipping MATE DND tests")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableMATE", func(t *testing.T) {
		err := dndManager.enableMATE()
		if err != nil {
			t.Logf("MATE enable failed (expected if not on MATE): %v", err)
		}
	})

	t.Run("DisableMATE", func(t *testing.T) {
		err := dndManager.disableMATE()
		if err != nil {
			t.Logf("MATE disable failed (expected if not on MATE): %v", err)
		}
	})

	t.Run("IsEnabledMATE", func(t *testing.T) {
		enabled, err := dndManager.isEnabledMATE()
		if err != nil {
			t.Logf("MATE status check failed (expected if not on MATE): %v", err)
		} else {
			t.Logf("MATE DND enabled: %v", enabled)
		}
	})
}

func TestCinnamonDND(t *testing.T) {
	if !hasCommand("gsettings") {
		t.Skip("gsettings not available, skipping Cinnamon DND tests")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableCinnamon", func(t *testing.T) {
		err := dndManager.enableCinnamon()
		if err != nil {
			t.Logf("Cinnamon enable failed (expected if not on Cinnamon): %v", err)
		}
	})

	t.Run("DisableCinnamon", func(t *testing.T) {
		err := dndManager.disableCinnamon()
		if err != nil {
			t.Logf("Cinnamon disable failed (expected if not on Cinnamon): %v", err)
		}
	})

	t.Run("IsEnabledCinnamon", func(t *testing.T) {
		enabled, err := dndManager.isEnabledCinnamon()
		if err != nil {
			t.Logf("Cinnamon status check failed (expected if not on Cinnamon): %v", err)
		} else {
			t.Logf("Cinnamon DND enabled: %v", enabled)
		}
	})
}

func TestTilingWMDND(t *testing.T) {
	hasDunst := hasCommand("dunstctl")
	hasMako := hasCommand("makoctl")

	if !hasDunst && !hasMako {
		t.Skip("No supported notification daemon found, skipping tiling WM DND tests")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableTilingWM", func(t *testing.T) {
		err := dndManager.enableTilingWM()
		if err != nil {
			t.Logf("Tiling WM enable failed: %v", err)
		}
	})

	t.Run("DisableTilingWM", func(t *testing.T) {
		err := dndManager.disableTilingWM()
		if err != nil {
			t.Logf("Tiling WM disable failed: %v", err)
		}
	})

	t.Run("IsEnabledTilingWM", func(t *testing.T) {
		enabled, err := dndManager.isEnabledTilingWM()
		if err != nil {
			t.Logf("Tiling WM status check failed: %v", err)
		} else {
			t.Logf("Tiling WM DND enabled: %v", enabled)
		}
	})
}

func TestGenericLinuxDND(t *testing.T) {
	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	t.Run("EnableGenericLinux", func(t *testing.T) {
		err := dndManager.enableGenericLinux()
		if err != nil {
			t.Errorf("Generic Linux enable should not fail: %v", err)
		}
	})

	t.Run("DisableGenericLinux", func(t *testing.T) {
		err := dndManager.disableGenericLinux()
		if err != nil {
			t.Errorf("Generic Linux disable should not fail: %v", err)
		}
	})
}

func TestDesktopEnvironmentDetection(t *testing.T) {
	notificationManager := notifications.NewNotificationManager(true)
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
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "Linux"
}

func hasCommand(command string) bool {
	cmd := exec.Command("which", command)
	return cmd.Run() == nil
}

// Benchmark tests for performance

func BenchmarkDetectDesktopEnvironment(b *testing.B) {
	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dndManager.detectDesktopEnvironment()
	}
}

func BenchmarkEnableDisableLinux(b *testing.B) {
	if !isLinux() {
		b.Skip("Skipping Linux benchmark on non-Linux system")
	}

	notificationManager := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(notificationManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dndManager.enableLinux()
		_ = dndManager.disableLinux()
	}
}
