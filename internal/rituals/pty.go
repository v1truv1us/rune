package rituals

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/ferg-cod3s/rune/internal/config"
	"golang.org/x/term"
)

// executePTYCommand executes a command with direct PTY allocation for interactive terminal access
func (e *Engine) executePTYCommand(cmd config.Command) error {
	// Parse command and arguments
	parts := strings.Fields(cmd.Command)
	if len(parts) == 0 {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("empty command")
	}

	// Create the command
	execCmd := exec.Command(parts[0], parts[1:]...)

	// Set up environment with filtered variables to avoid leaking secrets
	execCmd.Env = filterEnvironment(os.Environ())

	// Start the command with a pty
	ptmx, err := pty.Start(execCmd)
	if err != nil {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("failed to start command with PTY: %w", err)
	}
	defer func() {
		_ = ptmx.Close()
	}()

	fmt.Printf(" 🖥️ (interactive)\n")

	// Handle pty size changes
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				// Log error but don't fail - this is just for terminal resize
				fmt.Printf("Error resizing pty: %v\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize

	// Set stdin in raw mode for proper terminal interaction
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("Warning: failed to set raw mode: %v\n", err)
	} else {
		defer func() {
			_ = term.Restore(fd, oldState)
		}()
	}

	// Copy stdin/stdout/stderr to/from the pty
	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	// Copy output from pty to stdout
	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	// Wait for the command to finish
	err = execCmd.Wait()
	if err != nil {
		// For interactive commands, exit codes might be expected (e.g., user quits)
		// So we don't always treat this as a hard error
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				fmt.Printf("Interactive command exited with code: %d\n", exitError.ExitCode())
			}
		}
		return err
	}

	return nil
}

// ptySupported checks if PTY functionality is supported on the current platform
// Currently unused but kept for future platform-specific checks
func ptySupported() bool { //nolint:unused
	// PTY support is available on Unix-like systems (macOS, Linux)
	// Windows support is limited but may work through WSL
	return true // creack/pty handles platform differences
}
