package rituals

import (
	"fmt"
	"strings"

	"github.com/ferg-cod3s/rune/internal/config"
)

// executeTmuxCommand executes a command using tmux session management
func (e *Engine) executeTmuxCommand(cmd config.Command, scope string) error {
	if e.tmuxClient == nil {
		fmt.Printf(" ⚠ (tmux not available, falling back to standard execution)\n")
		return e.executeStandardCommand(cmd)
	}

	// Prepare template variables
	variables := map[string]string{
		"Project": scope,
		"Command": cmd.Command,
	}

	// Handle template-based session creation
	if cmd.TmuxTemplate != "" {
		return e.executeTemplateCommand(cmd, variables)
	}

	// Handle direct session management
	if cmd.TmuxSession != "" {
		return e.executeSessionCommand(cmd, variables)
	}

	// This shouldn't happen due to validation, but fallback to PTY
	fmt.Printf(" ⚠ (no tmux configuration specified, falling back to PTY)\n")
	return e.executePTYCommand(cmd)
}

// executeTemplateCommand creates a tmux session from a template
func (e *Engine) executeTemplateCommand(cmd config.Command, variables map[string]string) error {
	// Find the template in configuration
	template, exists := e.config.Rituals.Templates[cmd.TmuxTemplate]
	if !exists {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("template '%s' not found in configuration", cmd.TmuxTemplate)
	}

	// Create session from template
	err := e.tmuxClient.CreateFromTemplate(&template, variables)
	if err != nil {
		// If session already exists, try to attach to it instead
		if strings.Contains(err.Error(), "already exists") {
			sessionName := e.expandTemplate(template.SessionName, variables)
			fmt.Printf(" 📺 (attaching to existing session '%s')\n", sessionName)
			return e.tmuxClient.AttachSession(sessionName)
		}
		fmt.Printf(" ❌\n")
		return fmt.Errorf("failed to create session from template: %w", err)
	}

	// Save session state for persistence
	sessionName := e.expandTemplate(template.SessionName, variables)
	err = e.tmuxClient.SaveSessionState(sessionName, cmd.TmuxTemplate, variables["Project"], variables)
	if err != nil {
		// Log but don't fail - persistence is optional
		fmt.Printf(" ⚠ (failed to save session state: %v)\n", err)
	}

	// Attach to the newly created session
	fmt.Printf(" 📺 (created and attaching to session '%s')\n", sessionName)
	return e.tmuxClient.AttachSession(sessionName)
}

// executeSessionCommand manages a direct tmux session
func (e *Engine) executeSessionCommand(cmd config.Command, variables map[string]string) error {
	sessionName := e.expandTemplate(cmd.TmuxSession, variables)

	// Check if session exists
	if e.tmuxClient.SessionExists(sessionName) {
		fmt.Printf(" 📺 (attaching to existing session '%s')\n", sessionName)
		return e.tmuxClient.AttachSession(sessionName)
	}

	// Create new session
	err := e.tmuxClient.CreateSession(sessionName)
	if err != nil {
		fmt.Printf(" ❌\n")
		return fmt.Errorf("failed to create session '%s': %w", sessionName, err)
	}

	// Run the command in the new session if specified
	if cmd.Command != "" {
		// Parse and send the command to the session
		_ = e.expandTemplate(cmd.Command, variables) // TODO: Send command to session

		// Get the session and send the command
		// Note: This is a simplified approach. The gotmux library might need
		// more sophisticated command sending for complex commands.
		sessions, err := e.tmuxClient.ListSessions()
		if err != nil {
			fmt.Printf(" ⚠ (session created but failed to send command: %v)\n", err)
		} else {
			// Find our session and send the command
			for _, session := range sessions {
				if session == sessionName {
					// For now, we'll create the session and let the user manually run commands
					// A more sophisticated implementation would send the command to the session
					break
				}
			}
		}
	}

	fmt.Printf(" 📺 (created and attaching to session '%s')\n", sessionName)
	return e.tmuxClient.AttachSession(sessionName)
}

// cleanupSessions cleans up orphaned tmux sessions created by Rune
// Currently unused but kept for future cleanup functionality
func (e *Engine) cleanupSessions() error { //nolint:unused
	if e.tmuxClient == nil {
		return nil
	}

	sessions, err := e.tmuxClient.ListSessions()
	if err != nil {
		return err
	}

	// Clean up sessions that match Rune naming patterns
	// This is a basic implementation - could be more sophisticated
	for _, sessionName := range sessions {
		if strings.HasPrefix(sessionName, "rune-") {
			// Only kill sessions that appear to be orphaned
			// This is a placeholder for more sophisticated cleanup logic
			_ = e.tmuxClient.KillSession(sessionName)
		}
	}

	return nil
}
