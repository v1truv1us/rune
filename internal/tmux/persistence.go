package tmux

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SessionState represents the state of a tmux session for persistence
type SessionState struct {
	Name         string            `json:"name"`
	Template     string            `json:"template,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	LastActivity time.Time         `json:"last_activity"`
	Project      string            `json:"project"`
	Windows      []WindowState     `json:"windows,omitempty"`
}

// WindowState represents the state of a tmux window
type WindowState struct {
	Name   string      `json:"name"`
	Layout string      `json:"layout,omitempty"`
	Panes  []PaneState `json:"panes,omitempty"`
}

// PaneState represents the state of a tmux pane
type PaneState struct {
	Command     string            `json:"command,omitempty"`
	WorkingDir  string            `json:"working_dir,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// SessionPersistence handles saving and loading tmux session states
type SessionPersistence struct {
	stateDir string
}

// NewSessionPersistence creates a new session persistence manager
func NewSessionPersistence() (*SessionPersistence, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	stateDir := filepath.Join(homeDir, ".rune", "sessions")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session state directory: %w", err)
	}

	return &SessionPersistence{
		stateDir: stateDir,
	}, nil
}

// SaveSession saves the state of a tmux session
func (sp *SessionPersistence) SaveSession(state *SessionState) error {
	state.LastActivity = time.Now()

	filename := fmt.Sprintf("%s.json", state.Name)
	filepath := filepath.Join(sp.stateDir, filename)

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write session state file: %w", err)
	}

	return nil
}

// LoadSession loads the state of a tmux session
func (sp *SessionPersistence) LoadSession(sessionName string) (*SessionState, error) {
	filename := fmt.Sprintf("%s.json", sessionName)
	filepath := filepath.Join(sp.stateDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session state not found: %s", sessionName)
		}
		return nil, fmt.Errorf("failed to read session state file: %w", err)
	}

	var state SessionState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session state: %w", err)
	}

	return &state, nil
}

// ListPersistedSessions returns a list of all persisted session names
func (sp *SessionPersistence) ListPersistedSessions() ([]string, error) {
	entries, err := os.ReadDir(sp.stateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read session state directory: %w", err)
	}

	var sessions []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			sessionName := entry.Name()[:len(entry.Name())-5] // Remove .json extension
			sessions = append(sessions, sessionName)
		}
	}

	return sessions, nil
}

// DeleteSession removes a persisted session state
func (sp *SessionPersistence) DeleteSession(sessionName string) error {
	filename := fmt.Sprintf("%s.json", sessionName)
	filepath := filepath.Join(sp.stateDir, filename)

	err := os.Remove(filepath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session state file: %w", err)
	}

	return nil
}

// CleanupStaleSession removes session states for sessions that no longer exist
func (sp *SessionPersistence) CleanupStaleSessions(client *Client) error {
	activeSessions, err := client.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list active sessions: %w", err)
	}

	persistedSessions, err := sp.ListPersistedSessions()
	if err != nil {
		return fmt.Errorf("failed to list persisted sessions: %w", err)
	}

	// Create a map for quick lookup
	activeMap := make(map[string]bool)
	for _, session := range activeSessions {
		activeMap[session] = true
	}

	// Delete stale session states
	for _, persistedSession := range persistedSessions {
		if !activeMap[persistedSession] {
			err := sp.DeleteSession(persistedSession)
			if err != nil {
				// Log error but continue cleanup
				fmt.Printf("Warning: failed to cleanup stale session state '%s': %v\n", persistedSession, err)
			}
		}
	}

	return nil
}

// RestoreSession attempts to restore a tmux session from persisted state
func (sp *SessionPersistence) RestoreSession(client *Client, sessionName string) error {
	state, err := sp.LoadSession(sessionName)
	if err != nil {
		return fmt.Errorf("failed to load session state: %w", err)
	}

	// Check if session already exists
	if client.SessionExists(sessionName) {
		return fmt.Errorf("session '%s' already exists", sessionName)
	}

	// If this was created from a template, try to recreate using the template
	if state.Template != "" {
		// This would require access to the config, so we'll implement this later
		// For now, we'll create a basic session
		return client.CreateSession(sessionName)
	}

	// Create basic session for now
	// A more sophisticated implementation would recreate the full window/pane structure
	return client.CreateSession(sessionName)
}
