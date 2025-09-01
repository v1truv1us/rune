package tmux

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAvailable(t *testing.T) {
	t.Run("should return consistent result with manual check", func(t *testing.T) {
		// Manual check using exec.Command
		cmd := exec.Command("tmux", "-V")
		manualErr := cmd.Run()
		manualAvailable := manualErr == nil

		// Our function result
		functionResult := IsAvailable()

		assert.Equal(t, manualAvailable, functionResult, "IsAvailable() should match manual tmux check")
	})
}

func TestGetVersion(t *testing.T) {
	if !IsAvailable() {
		t.Skip("tmux not available, skipping version test")
	}

	t.Run("should return valid version string", func(t *testing.T) {
		version, err := GetVersion()
		require.NoError(t, err, "GetVersion should not error when tmux is available")
		assert.NotEmpty(t, version, "Version string should not be empty")
		assert.Contains(t, version, "tmux", "Version string should contain 'tmux'")
	})
}

func TestGetVersionWhenNotAvailable(t *testing.T) {
	if IsAvailable() {
		t.Skip("tmux is available, cannot test unavailable scenario")
	}

	t.Run("should return error when tmux not available", func(t *testing.T) {
		_, err := GetVersion()
		assert.Error(t, err, "GetVersion should return error when tmux is not available")
	})
}

func TestGetDefaultInstallPath(t *testing.T) {
	t.Run("should return platform-specific path", func(t *testing.T) {
		path := GetDefaultInstallPath()
		assert.NotEmpty(t, path, "Install path should not be empty")

		switch runtime.GOOS {
		case "darwin":
			assert.Contains(t, path, "homebrew", "macOS path should mention homebrew")
		case "linux":
			assert.Contains(t, path, "/usr/bin", "Linux path should mention /usr/bin")
		case "windows":
			assert.Contains(t, path, "WSL", "Windows path should mention WSL")
		default:
			assert.Contains(t, path, "package manager", "Unknown OS should mention package manager")
		}
	})
}

func TestIsInTmuxSession(t *testing.T) {
	t.Run("should not panic when tmux unavailable", func(t *testing.T) {
		// This should not panic regardless of tmux availability
		assert.NotPanics(t, func() {
			IsInTmuxSession()
		})
	})

	t.Run("should return boolean result", func(t *testing.T) {
		result := IsInTmuxSession()
		// Result can be true or false, just ensure it's a valid boolean
		assert.IsType(t, true, result, "Should return boolean value")
	})
}

// TestTmuxIntegration tests the basic tmux integration
// This test only runs if tmux is available
func TestTmuxIntegration(t *testing.T) {
	if !IsAvailable() {
		t.Skip("tmux not available, skipping integration tests")
	}

	t.Run("basic tmux commands work", func(t *testing.T) {
		// Test that we can get version
		version, err := GetVersion()
		require.NoError(t, err)
		assert.NotEmpty(t, version)

		// Test that IsInTmuxSession doesn't crash
		assert.NotPanics(t, func() {
			IsInTmuxSession()
		})
	})
}

// BenchmarkIsAvailable benchmarks the IsAvailable function
func BenchmarkIsAvailable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsAvailable()
	}
}

// BenchmarkGetVersion benchmarks the GetVersion function
func BenchmarkGetVersion(b *testing.B) {
	if !IsAvailable() {
		b.Skip("tmux not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetVersion()
	}
}
