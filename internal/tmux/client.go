// Package tmux client provides a wrapper around the gotmux library
// for Rune's interactive ritual system with session management utilities.
package tmux

import (
	"fmt"
	"strings"
	"time"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/ferg-cod3s/rune/internal/config"
)

// Client wraps the gotmux.Tmux client and provides Rune-specific functionality
// for managing tmux sessions, windows, and panes during ritual execution.
type Client struct {
	tmux        *gotmux.Tmux
	persistence *SessionPersistence
}

// NewClient creates a new tmux client wrapper.
// Returns an error if tmux is not available on the system.
func NewClient() (*Client, error) {
	if !IsAvailable() {
		return nil, fmt.Errorf("tmux is not available on this system. Please install tmux to use interactive rituals")
	}

	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return nil, fmt.Errorf("failed to create tmux client: %w", err)
	}

	persistence, err := NewSessionPersistence()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize session persistence: %w", err)
	}

	return &Client{
		tmux:        tmux,
		persistence: persistence,
	}, nil
}

// SessionExists checks if a tmux session with the given name exists.
func (c *Client) SessionExists(sessionName string) bool {
	sessions, err := c.tmux.ListSessions()
	if err != nil {
		return false
	}

	for _, session := range sessions {
		if session.Name == sessionName {
			return true
		}
	}
	return false
}

// CreateSession creates a new tmux session with the given name.
// If the session already exists, returns an error.
func (c *Client) CreateSession(sessionName string) error {
	if c.SessionExists(sessionName) {
		return fmt.Errorf("session '%s' already exists", sessionName)
	}

	_, err := c.tmux.NewSession(&gotmux.SessionOptions{
		Name: sessionName,
	})
	if err != nil {
		return fmt.Errorf("failed to create session '%s': %w", sessionName, err)
	}

	return nil
}

// KillSession kills the tmux session with the given name.
// Returns an error if the session doesn't exist.
func (c *Client) KillSession(sessionName string) error {
	if !c.SessionExists(sessionName) {
		return fmt.Errorf("session '%s' does not exist", sessionName)
	}

	sessions, err := c.tmux.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, session := range sessions {
		if session.Name == sessionName {
			err = session.Kill()
			if err != nil {
				return fmt.Errorf("failed to kill session '%s': %w", sessionName, err)
			}
			return nil
		}
	}

	return fmt.Errorf("session '%s' not found", sessionName)
}

// AttachSession attaches to an existing tmux session.
// This will block until the session is detached.
func (c *Client) AttachSession(sessionName string) error {
	if !c.SessionExists(sessionName) {
		return fmt.Errorf("session '%s' does not exist", sessionName)
	}

	sessions, err := c.tmux.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, session := range sessions {
		if session.Name == sessionName {
			err = session.Attach()
			if err != nil {
				return fmt.Errorf("failed to attach to session '%s': %w", sessionName, err)
			}
			return nil
		}
	}

	return fmt.Errorf("session '%s' not found", sessionName)
}

// CreateFromTemplate creates a tmux session from a template configuration.
// This includes creating windows, panes, and running initial commands.
func (c *Client) CreateFromTemplate(template *config.TmuxTemplate, variables map[string]string) error {
	// Replace variables in session name
	sessionName := template.SessionName
	for key, value := range variables {
		sessionName = strings.ReplaceAll(sessionName, "{{."+key+"}}", value)
	}

	// Create the session
	if c.SessionExists(sessionName) {
		return fmt.Errorf("session '%s' already exists", sessionName)
	}

	session, err := c.tmux.NewSession(&gotmux.SessionOptions{
		Name: sessionName,
	})
	if err != nil {
		return fmt.Errorf("failed to create session from template: %w", err)
	}

	// Create windows and panes according to template
	for i, window := range template.Windows {
		var tmuxWindow *gotmux.Window

		if i == 0 {
			// Get the default window and rename it
			tmuxWindow, err = session.GetWindowByIndex(0)
			if err != nil {
				return fmt.Errorf("failed to get first window: %w", err)
			}
			err = tmuxWindow.Rename(window.Name)
			if err != nil {
				return fmt.Errorf("failed to rename first window: %w", err)
			}
		} else {
			// Create new window
			tmuxWindow, err = session.New()
			if err != nil {
				return fmt.Errorf("failed to create window '%s': %w", window.Name, err)
			}
			err = tmuxWindow.Rename(window.Name)
			if err != nil {
				return fmt.Errorf("failed to rename window '%s': %w", window.Name, err)
			}
		}

		// Set layout if specified
		if window.Layout != "" {
			// Convert string layout to gotmux layout type
			var layout gotmux.WindowLayout
			switch window.Layout {
			case "even-horizontal":
				layout = gotmux.WindowLayoutEvenHorizontal
			case "even-vertical":
				layout = gotmux.WindowLayoutEvenVertical
			case "main-horizontal":
				layout = gotmux.WindowLayoutMainVertical // Note: gotmux only has MainVertical
			case "main-vertical":
				layout = gotmux.WindowLayoutMainVertical
			case "tiled":
				layout = gotmux.WindowLayoutTiled
			default:
				return fmt.Errorf("unsupported layout '%s' for window '%s'", window.Layout, window.Name)
			}
			err = tmuxWindow.SelectLayout(layout)
			if err != nil {
				return fmt.Errorf("failed to set layout '%s' for window '%s': %w", window.Layout, window.Name, err)
			}
		}

		// Create panes and run commands
		if err := c.createPanes(tmuxWindow, window.Panes, variables); err != nil {
			return fmt.Errorf("failed to create panes for window '%s': %w", window.Name, err)
		}
	}

	return nil
}

// createPanes creates panes in a window and runs the specified commands.
func (c *Client) createPanes(window *gotmux.Window, panes []string, variables map[string]string) error {
	if len(panes) == 0 {
		return nil
	}

	// Get the first pane (should exist by default)
	windowPanes, err := window.ListPanes()
	if err != nil {
		return fmt.Errorf("failed to get window panes: %w", err)
	}

	if len(windowPanes) == 0 {
		return fmt.Errorf("no panes found in window")
	}

	// Run command in first pane if specified
	if len(panes) > 0 {
		command := c.replaceVariables(panes[0], variables)
		if command != "" {
			err = windowPanes[0].SendKeys(command)
			if err != nil {
				return fmt.Errorf("failed to send command to first pane: %w", err)
			}
		}
	}

	// Create additional panes and run commands
	for i := 1; i < len(panes); i++ {
		// Split the first pane to create a new pane
		// Alternate between horizontal and vertical splits
		var splitOptions *gotmux.SplitWindowOptions
		if i%2 == 0 {
			splitOptions = &gotmux.SplitWindowOptions{
				SplitDirection: gotmux.PaneSplitDirectionHorizontal,
			}
		} else {
			splitOptions = &gotmux.SplitWindowOptions{
				SplitDirection: gotmux.PaneSplitDirectionVertical,
			}
		}
		err := windowPanes[0].SplitWindow(splitOptions)
		if err != nil {
			return fmt.Errorf("failed to split pane: %w", err)
		}

		// Refresh pane list to get the new pane
		windowPanes, err = window.ListPanes()
		if err != nil {
			return fmt.Errorf("failed to refresh pane list: %w", err)
		}

		// Use the last pane (newest one) for the command
		if len(windowPanes) > i {
			newPane := windowPanes[len(windowPanes)-1]
			command := c.replaceVariables(panes[i], variables)
			if command != "" {
				err = newPane.SendKeys(command)
				if err != nil {
					return fmt.Errorf("failed to send command to pane: %w", err)
				}
			}
		}
	}

	return nil
}

// replaceVariables replaces template variables in a command string.
func (c *Client) replaceVariables(command string, variables map[string]string) string {
	result := command
	for key, value := range variables {
		result = strings.ReplaceAll(result, "{{."+key+"}}", value)
	}
	return result
}

// ListSessions returns a list of all active tmux sessions.
func (c *Client) ListSessions() ([]string, error) {
	sessions, err := c.tmux.ListSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	var sessionNames []string
	for _, session := range sessions {
		sessionNames = append(sessionNames, session.Name)
	}

	return sessionNames, nil
}

// SaveSessionState persists the current state of a session
func (c *Client) SaveSessionState(sessionName, template, project string, variables map[string]string) error {
	if c.persistence == nil {
		return fmt.Errorf("session persistence not initialized")
	}

	state := &SessionState{
		Name:      sessionName,
		Template:  template,
		Variables: variables,
		CreatedAt: time.Now(),
		Project:   project,
	}

	return c.persistence.SaveSession(state)
}

// RestoreSession attempts to restore a session from persisted state
func (c *Client) RestoreSession(sessionName string) error {
	if c.persistence == nil {
		return fmt.Errorf("session persistence not initialized")
	}

	return c.persistence.RestoreSession(c, sessionName)
}

// ListPersistedSessions returns all sessions with saved state
func (c *Client) ListPersistedSessions() ([]string, error) {
	if c.persistence == nil {
		return nil, fmt.Errorf("session persistence not initialized")
	}

	return c.persistence.ListPersistedSessions()
}

// CleanupStaleSessions removes saved state for non-existent sessions
func (c *Client) CleanupStaleSessions() error {
	if c.persistence == nil {
		return fmt.Errorf("session persistence not initialized")
	}

	return c.persistence.CleanupStaleSessions(c)
}
