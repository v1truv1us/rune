// Package testutil provides shared test utilities for environment detection
// and test skip conditions across the Rune test suite.
package testutil

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// HasDisplay returns true if a graphical display is available.
// Checks both X11 (DISPLAY) and Wayland (WAYLAND_DISPLAY) environments.
func HasDisplay() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}

// IsHeadless returns true if running in a headless environment (no display).
func IsHeadless() bool {
	return !HasDisplay()
}

// HasCommand returns true if the given command is available in PATH.
func HasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// CanSendLinuxNotifications returns true if Linux notifications can be sent.
// Requires both a display and notify-send to be installed.
func CanSendLinuxNotifications() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	return HasDisplay() && HasCommand("notify-send")
}

// CanControlGNOMEDnD returns true if GNOME DnD can be controlled.
// Requires gsettings, a display, and access to the GNOME notification schema.
func CanControlGNOMEDnD() bool {
	if !HasCommand("gsettings") {
		return false
	}
	if IsHeadless() {
		return false
	}
	// Try to access the schema - this will fail without a proper dbus session
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.notifications", "show-banners")
	return cmd.Run() == nil
}

// CanControlKDEDnD returns true if KDE DnD can be controlled.
func CanControlKDEDnD() bool {
	hasKDE5 := HasCommand("kwriteconfig5")
	hasKDE6 := HasCommand("kwriteconfig6")
	if !hasKDE5 && !hasKDE6 {
		return false
	}
	return HasDisplay()
}

// CanControlXFCEDnD returns true if XFCE DnD can be controlled.
func CanControlXFCEDnD() bool {
	if !HasCommand("xfconf-query") {
		return false
	}
	return HasDisplay()
}

// CanControlTilingWMDnD returns true if a tiling WM notification daemon can be controlled.
func CanControlTilingWMDnD() bool {
	hasDunst := HasCommand("dunstctl")
	hasMako := HasCommand("makoctl")
	return hasDunst || hasMako
}

// GetDesktopEnvironment returns the detected desktop environment name.
// Returns "unknown" if detection fails.
func GetDesktopEnvironment() string {
	envVars := []string{
		"XDG_CURRENT_DESKTOP",
		"DESKTOP_SESSION",
		"XDG_SESSION_DESKTOP",
		"GDMSESSION",
	}

	for _, envVar := range envVars {
		value := strings.ToLower(os.Getenv(envVar))
		if value != "" {
			if strings.Contains(value, "gnome") || strings.Contains(value, "ubuntu") || strings.Contains(value, "pop") {
				return "gnome"
			}
			if strings.Contains(value, "kde") || strings.Contains(value, "plasma") {
				return "kde"
			}
			if strings.Contains(value, "xfce") {
				return "xfce"
			}
			if strings.Contains(value, "mate") {
				return "mate"
			}
			if strings.Contains(value, "cinnamon") {
				return "cinnamon"
			}
			if strings.Contains(value, "i3") {
				return "i3"
			}
			if strings.Contains(value, "sway") {
				return "sway"
			}
		}
	}

	return "unknown"
}

// IsRunningGNOME returns true if GNOME desktop environment is active.
func IsRunningGNOME() bool {
	de := GetDesktopEnvironment()
	return de == "gnome"
}

// IsRunningKDE returns true if KDE desktop environment is active.
func IsRunningKDE() bool {
	de := GetDesktopEnvironment()
	return de == "kde"
}

// SkipIfHeadless skips the test if running in a headless environment.
// Use this at the start of tests that require a display.
func SkipIfHeadless(t interface{ Skip(...interface{}) }) {
	if IsHeadless() {
		t.Skip("Skipping test: no display available (headless environment)")
	}
}

// SkipIfNoNotifySend skips the test if notify-send is not available.
func SkipIfNoNotifySend(t interface{ Skip(...interface{}) }) {
	if runtime.GOOS == "linux" && !HasCommand("notify-send") {
		t.Skip("Skipping test: notify-send not installed")
	}
}

// SkipIfCannotControlGNOMEDnD skips the test if GNOME DnD cannot be controlled.
func SkipIfCannotControlGNOMEDnD(t interface{ Skip(...interface{}) }) {
	if !CanControlGNOMEDnD() {
		t.Skip("Skipping test: cannot control GNOME DnD (no display, no gsettings, or no dbus session)")
	}
}
