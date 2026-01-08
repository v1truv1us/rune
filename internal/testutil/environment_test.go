package testutil

import (
	"os"
	"runtime"
	"testing"
)

func TestHasDisplay(t *testing.T) {
	// Save original values
	origDisplay := os.Getenv("DISPLAY")
	origWayland := os.Getenv("WAYLAND_DISPLAY")
	defer func() {
		os.Setenv("DISPLAY", origDisplay)
		os.Setenv("WAYLAND_DISPLAY", origWayland)
	}()

	t.Run("no display", func(t *testing.T) {
		os.Setenv("DISPLAY", "")
		os.Setenv("WAYLAND_DISPLAY", "")
		if HasDisplay() {
			t.Error("HasDisplay() should return false when no display is set")
		}
	})

	t.Run("X11 display", func(t *testing.T) {
		os.Setenv("DISPLAY", ":0")
		os.Setenv("WAYLAND_DISPLAY", "")
		if !HasDisplay() {
			t.Error("HasDisplay() should return true when DISPLAY is set")
		}
	})

	t.Run("Wayland display", func(t *testing.T) {
		os.Setenv("DISPLAY", "")
		os.Setenv("WAYLAND_DISPLAY", "wayland-0")
		if !HasDisplay() {
			t.Error("HasDisplay() should return true when WAYLAND_DISPLAY is set")
		}
	})
}

func TestIsHeadless(t *testing.T) {
	// Save original values
	origDisplay := os.Getenv("DISPLAY")
	origWayland := os.Getenv("WAYLAND_DISPLAY")
	defer func() {
		os.Setenv("DISPLAY", origDisplay)
		os.Setenv("WAYLAND_DISPLAY", origWayland)
	}()

	os.Setenv("DISPLAY", "")
	os.Setenv("WAYLAND_DISPLAY", "")
	if !IsHeadless() {
		t.Error("IsHeadless() should return true when no display is set")
	}

	os.Setenv("DISPLAY", ":0")
	if IsHeadless() {
		t.Error("IsHeadless() should return false when DISPLAY is set")
	}
}

func TestHasCommand(t *testing.T) {
	// Test with a command that should exist on all systems
	var existingCmd string
	switch runtime.GOOS {
	case "windows":
		existingCmd = "cmd"
	default:
		existingCmd = "sh"
	}

	if !HasCommand(existingCmd) {
		t.Errorf("HasCommand(%q) should return true", existingCmd)
	}

	if HasCommand("this-command-definitely-does-not-exist-12345") {
		t.Error("HasCommand() should return false for non-existent command")
	}
}

func TestGetDesktopEnvironment(t *testing.T) {
	// Save original value
	origDE := os.Getenv("XDG_CURRENT_DESKTOP")
	defer os.Setenv("XDG_CURRENT_DESKTOP", origDE)

	testCases := []struct {
		envValue string
		expected string
	}{
		{"GNOME", "gnome"},
		{"ubuntu:GNOME", "gnome"},
		{"KDE", "kde"},
		{"plasma", "kde"},
		{"XFCE", "xfce"},
		{"MATE", "mate"},
		{"X-Cinnamon", "cinnamon"},
		{"i3", "i3"},
		{"sway", "sway"},
		{"", "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.envValue, func(t *testing.T) {
			os.Setenv("XDG_CURRENT_DESKTOP", tc.envValue)
			result := GetDesktopEnvironment()
			if result != tc.expected {
				t.Errorf("GetDesktopEnvironment() = %q, want %q", result, tc.expected)
			}
		})
	}
}
