package rituals

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ferg-cod3s/rune/internal/config"
)

// filterEnvironment removes sensitive environment variables before passing to subprocesses
func filterEnvironment(env []string) []string {
	// Define sensitive prefixes and exact names to filter out
	sensitiveSubstrings := []string{
		"AWS_", "AZURE_", "GOOGLE_", "GCP_", "GCLOUD_", "DO_", "DIGITALOCEAN_",
		"SECRET", "TOKEN", "KEY", "PASSWORD", "PASS", "PRIVATE", "SSH_", "AUTH", "CREDENTIAL", "SESSION", "COOKIE", "BEARER",
	}
	sensitiveExact := map[string]struct{}{
		"RUNE_SENTRY_DSN":    {},
		"GITHUB_TOKEN":       {},
		"NPM_TOKEN":          {},
		"NETLIFY_AUTH_TOKEN": {},
		"VERCEL_TOKEN":       {},
	}

	filtered := make([]string, 0, len(env))
	for _, kv := range env {
		// kv is in the form KEY=VALUE
		key := kv
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			key = kv[:idx]
		}

		// Allowlist Rune-specific runtime flags that are safe
		if key == "RUNE_DEBUG" || key == "RUNE_ENV" {
			filtered = append(filtered, kv)
			continue
		}

		// Filter exact matches
		if _, ok := sensitiveExact[key]; ok {
			continue
		}

		upperKey := strings.ToUpper(key)
		blocked := false
		for _, sub := range sensitiveSubstrings {
			if strings.Contains(upperKey, sub) {
				blocked = true
				break
			}
		}
		if blocked {
			continue
		}

		filtered = append(filtered, kv)
	}
	return filtered
}

// Engine handles ritual execution
type Engine struct {
	config *config.Config
}

// NewEngine creates a new ritual engine
func NewEngine(cfg *config.Config) *Engine {
	return &Engine{
		config: cfg,
	}
}

// ExecuteStartRituals executes start rituals for the given project
func (e *Engine) ExecuteStartRituals(project string) error {
	fmt.Println("🔮 Executing start rituals...")

	// Execute global start rituals
	if err := e.executeCommands(e.config.Rituals.Start.Global, "global"); err != nil {
		return fmt.Errorf("failed to execute global start rituals: %w", err)
	}

	// Execute project-specific start rituals
	if projectCommands, exists := e.config.Rituals.Start.PerProject[project]; exists {
		if err := e.executeCommands(projectCommands, project); err != nil {
			return fmt.Errorf("failed to execute project start rituals: %w", err)
		}
	}

	return nil
}

// ExecuteStopRituals executes stop rituals for the given project
func (e *Engine) ExecuteStopRituals(project string) error {
	fmt.Println("🔮 Executing stop rituals...")

	// Execute project-specific stop rituals first
	if projectCommands, exists := e.config.Rituals.Stop.PerProject[project]; exists {
		if err := e.executeCommands(projectCommands, project); err != nil {
			return fmt.Errorf("failed to execute project stop rituals: %w", err)
		}
	}

	// Execute global stop rituals
	if err := e.executeCommands(e.config.Rituals.Stop.Global, "global"); err != nil {
		return fmt.Errorf("failed to execute global stop rituals: %w", err)
	}

	return nil
}

// executeCommands executes a list of commands
func (e *Engine) executeCommands(commands []config.Command, scope string) error {
	for _, cmd := range commands {
		if err := e.executeCommand(cmd, scope); err != nil {
			if cmd.Optional {
				fmt.Printf("⚠ Optional command failed: %s (%v)\n", cmd.Name, err)
				continue
			}
			return fmt.Errorf("command '%s' failed: %w", cmd.Name, err)
		}
	}
	return nil
}

// executeCommand executes a single command
func (e *Engine) executeCommand(cmd config.Command, _ string) error {
	fmt.Printf("  ⚡ %s...", cmd.Name)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse command and arguments
	parts := strings.Fields(cmd.Command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Create the command
	execCmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	// Set up environment with filtered variables to avoid leaking secrets
	execCmd.Env = filterEnvironment(os.Environ())

	if cmd.Background {
		// For background commands, just start them
		if err := execCmd.Start(); err != nil {
			fmt.Printf(" ❌\n")
			return err
		}
		fmt.Printf(" ✓ (background)\n")
		return nil
	}

	// For foreground commands, wait for completion
	output, err := execCmd.CombinedOutput()
	if err != nil {
		fmt.Printf(" ❌\n")
		if len(output) > 0 {
			fmt.Printf("    Output: %s\n", strings.TrimSpace(string(output)))
		}
		return err
	}

	fmt.Printf(" ✓\n")

	// Show output if verbose mode is enabled
	if len(output) > 0 && shouldShowOutput(string(output)) {
		fmt.Printf("    %s\n", strings.TrimSpace(string(output)))
	}

	return nil
}

// shouldShowOutput determines if command output should be displayed
func shouldShowOutput(output string) bool {
	// Don't show empty output or common uninteresting outputs
	output = strings.TrimSpace(output)
	if output == "" {
		return false
	}

	// Skip common Git outputs that aren't useful
	skipPatterns := []string{
		"Already up to date",
		"nothing to commit, working tree clean",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(output, pattern) {
			return false
		}
	}

	return true
}

// TestRitual tests a ritual without executing it
func (e *Engine) TestRitual(ritualType string, project string) error {
	fmt.Printf("🧪 Testing %s ritual for project: %s\n", ritualType, project)

	var commands []config.Command

	switch ritualType {
	case "start":
		commands = append(commands, e.config.Rituals.Start.Global...)
		if projectCommands, exists := e.config.Rituals.Start.PerProject[project]; exists {
			commands = append(commands, projectCommands...)
		}
	case "stop":
		if projectCommands, exists := e.config.Rituals.Stop.PerProject[project]; exists {
			commands = append(commands, projectCommands...)
		}
		commands = append(commands, e.config.Rituals.Stop.Global...)
	default:
		return fmt.Errorf("unknown ritual type: %s", ritualType)
	}

	if len(commands) == 0 {
		fmt.Println("  No commands configured for this ritual")
		return nil
	}

	fmt.Println("Commands that would be executed:")
	for i, cmd := range commands {
		optional := ""
		background := ""
		if cmd.Optional {
			optional = " (optional)"
		}
		if cmd.Background {
			background = " (background)"
		}
		fmt.Printf("  %d. %s: %s%s%s\n", i+1, cmd.Name, cmd.Command, optional, background)
	}

	return nil
}
