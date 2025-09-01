// Package tmux provides utilities for detecting and managing tmux sessions
// for Rune's interactive ritual system.
package tmux

import (
	"os/exec"
	"runtime"
	"strings"
)

// IsAvailable checks if tmux is available on the system by attempting
// to execute 'tmux -V' command and checking if it succeeds.
// This is cross-platform and works on all systems where tmux might be installed.
func IsAvailable() bool {
	cmd := exec.Command("tmux", "-V")
	err := cmd.Run()
	return err == nil
}

// GetVersion returns the tmux version string if tmux is available.
// Returns an error if tmux is not found or fails to execute.
// The version string is typically in format "tmux X.Y" where X.Y is the version.
func GetVersion() (string, error) {
	cmd := exec.Command("tmux", "-V")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}

// GetDefaultInstallPath returns the expected installation path for tmux
// based on the current operating system. This can be used for installation
// guidance or troubleshooting when tmux is not found in PATH.
func GetDefaultInstallPath() string {
	switch runtime.GOOS {
	case "darwin":
		// macOS typically installs via Homebrew
		return "/opt/homebrew/bin/tmux or /usr/local/bin/tmux"
	case "linux":
		// Linux distributions typically install to /usr/bin
		return "/usr/bin/tmux"
	case "windows":
		// Windows might have tmux via WSL or other means
		return "Available via WSL or Windows package managers"
	default:
		return "Check your package manager or compile from source"
	}
}

// IsInTmuxSession checks if the current process is running inside a tmux session
// by checking for the TMUX environment variable that tmux sets.
func IsInTmuxSession() bool {
	cmd := exec.Command("tmux", "display-message", "-p", "#{session_name}")
	err := cmd.Run()
	return err == nil
}
