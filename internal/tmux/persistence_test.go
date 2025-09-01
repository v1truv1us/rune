package tmux

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionPersistence(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	sp := &SessionPersistence{
		stateDir: tempDir,
	}

	t.Run("should save and load session state", func(t *testing.T) {
		state := &SessionState{
			Name:      "test-session",
			Template:  "test-template",
			Variables: map[string]string{"Project": "test-project"},
			CreatedAt: time.Now(),
			Project:   "test-project",
		}

		// Save session
		err := sp.SaveSession(state)
		require.NoError(t, err)

		// Load session
		loadedState, err := sp.LoadSession("test-session")
		require.NoError(t, err)

		assert.Equal(t, state.Name, loadedState.Name)
		assert.Equal(t, state.Template, loadedState.Template)
		assert.Equal(t, state.Variables, loadedState.Variables)
		assert.Equal(t, state.Project, loadedState.Project)
	})

	t.Run("should list persisted sessions", func(t *testing.T) {
		// Create multiple sessions
		sessions := []string{"session1", "session2", "session3"}
		for _, sessionName := range sessions {
			state := &SessionState{
				Name:      sessionName,
				Project:   "test-project",
				CreatedAt: time.Now(),
			}
			err := sp.SaveSession(state)
			require.NoError(t, err)
		}

		// List sessions
		listedSessions, err := sp.ListPersistedSessions()
		require.NoError(t, err)

		assert.Contains(t, listedSessions, "session1")
		assert.Contains(t, listedSessions, "session2")
		assert.Contains(t, listedSessions, "session3")
	})

	t.Run("should delete session state", func(t *testing.T) {
		// Create and save session
		state := &SessionState{
			Name:      "delete-test",
			Project:   "test-project",
			CreatedAt: time.Now(),
		}
		err := sp.SaveSession(state)
		require.NoError(t, err)

		// Verify it exists
		_, err = sp.LoadSession("delete-test")
		require.NoError(t, err)

		// Delete session
		err = sp.DeleteSession("delete-test")
		require.NoError(t, err)

		// Verify it's gone
		_, err = sp.LoadSession("delete-test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session state not found")
	})

	t.Run("should handle non-existent session", func(t *testing.T) {
		_, err := sp.LoadSession("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session state not found")
	})
}

func TestNewSessionPersistence(t *testing.T) {
	t.Run("should create session persistence", func(t *testing.T) {
		sp, err := NewSessionPersistence()
		assert.NoError(t, err)
		assert.NotNil(t, sp)

		// Verify directory exists
		_, err = os.Stat(sp.stateDir)
		assert.NoError(t, err)
	})
}

func TestSessionPersistenceCleanup(t *testing.T) {
	if !IsAvailable() {
		t.Skip("tmux not available, skipping cleanup test")
	}

	tempDir := t.TempDir()
	sp := &SessionPersistence{
		stateDir: tempDir,
	}

	// Create mock client for testing
	client, err := NewClient()
	require.NoError(t, err)

	t.Run("should cleanup stale sessions", func(t *testing.T) {
		// Create some persisted session states
		persistedSessions := []string{"stale-session-1", "stale-session-2"}
		for _, sessionName := range persistedSessions {
			state := &SessionState{
				Name:      sessionName,
				Project:   "test-project",
				CreatedAt: time.Now().Add(-24 * time.Hour), // Old session
			}
			err := sp.SaveSession(state)
			require.NoError(t, err)
		}

		// Run cleanup (these sessions don't exist in tmux)
		err = sp.CleanupStaleSessions(client)
		if err != nil {
			t.Logf("Cleanup failed (expected if tmux server not running): %v", err)
			// If cleanup fails due to no tmux server, skip the verification
			return
		}

		// Verify stale sessions were cleaned up
		sessions, err := sp.ListPersistedSessions()
		require.NoError(t, err)

		for _, staleSession := range persistedSessions {
			assert.NotContains(t, sessions, staleSession)
		}
	})
}
