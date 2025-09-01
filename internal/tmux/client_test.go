package tmux

import (
	"testing"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("should create client when tmux available", func(t *testing.T) {
		if !IsAvailable() {
			t.Skip("tmux not available, skipping client creation test")
		}

		client, err := NewClient()
		require.NoError(t, err, "NewClient should not error when tmux is available")
		assert.NotNil(t, client, "Client should not be nil")
		assert.NotNil(t, client.tmux, "Internal tmux client should not be nil")
	})

	t.Run("should return error when tmux not available", func(t *testing.T) {
		if IsAvailable() {
			t.Skip("tmux is available, cannot test unavailable scenario")
		}

		client, err := NewClient()
		assert.Error(t, err, "NewClient should return error when tmux is not available")
		assert.Nil(t, client, "Client should be nil when creation fails")
		assert.Contains(t, err.Error(), "tmux is not available", "Error should mention tmux availability")
	})
}

func TestReplaceVariables(t *testing.T) {
	if !IsAvailable() {
		t.Skip("tmux not available, skipping client tests")
	}

	client, err := NewClient()
	require.NoError(t, err)

	tests := []struct {
		name      string
		command   string
		variables map[string]string
		expected  string
	}{
		{
			name:      "single variable replacement",
			command:   "cd {{.Project}}",
			variables: map[string]string{"Project": "myproject"},
			expected:  "cd myproject",
		},
		{
			name:      "multiple variable replacement",
			command:   "cd {{.Project}} && echo {{.User}}",
			variables: map[string]string{"Project": "myproject", "User": "john"},
			expected:  "cd myproject && echo john",
		},
		{
			name:      "no variables",
			command:   "ls -la",
			variables: map[string]string{},
			expected:  "ls -la",
		},
		{
			name:      "variable not in map",
			command:   "cd {{.Missing}}",
			variables: map[string]string{"Project": "myproject"},
			expected:  "cd {{.Missing}}",
		},
		{
			name:      "empty command",
			command:   "",
			variables: map[string]string{"Project": "myproject"},
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.replaceVariables(tt.command, tt.variables)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSessionManagement(t *testing.T) {
	if !IsAvailable() {
		t.Skip("tmux not available, skipping session management tests")
	}

	client, err := NewClient()
	require.NoError(t, err)

	testSessionName := "rune-test-session"

	// Cleanup any existing test session
	if client.SessionExists(testSessionName) {
		_ = client.KillSession(testSessionName)
	}

	t.Run("session should not exist initially", func(t *testing.T) {
		exists := client.SessionExists(testSessionName)
		assert.False(t, exists, "Test session should not exist initially")
	})

	t.Run("should create session successfully", func(t *testing.T) {
		err := client.CreateSession(testSessionName)
		assert.NoError(t, err, "CreateSession should not error")

		// Verify session exists
		exists := client.SessionExists(testSessionName)
		assert.True(t, exists, "Session should exist after creation")
	})

	t.Run("should not create duplicate session", func(t *testing.T) {
		err := client.CreateSession(testSessionName)
		assert.Error(t, err, "CreateSession should error when session already exists")
		assert.Contains(t, err.Error(), "already exists", "Error should mention session already exists")
	})

	t.Run("should list sessions including test session", func(t *testing.T) {
		sessions, err := client.ListSessions()
		assert.NoError(t, err, "ListSessions should not error")
		assert.Contains(t, sessions, testSessionName, "Sessions list should include test session")
	})

	t.Run("should kill session successfully", func(t *testing.T) {
		err := client.KillSession(testSessionName)
		assert.NoError(t, err, "KillSession should not error")

		// Verify session no longer exists
		exists := client.SessionExists(testSessionName)
		assert.False(t, exists, "Session should not exist after killing")
	})

	t.Run("should not kill non-existent session", func(t *testing.T) {
		err := client.KillSession("non-existent-session")
		assert.Error(t, err, "KillSession should error when session doesn't exist")
		assert.Contains(t, err.Error(), "does not exist", "Error should mention session doesn't exist")
	})
}

func TestCreateFromTemplate(t *testing.T) {
	if !IsAvailable() {
		t.Skip("tmux not available, skipping template tests")
	}

	client, err := NewClient()
	require.NoError(t, err)

	testSessionName := "rune-template-test"

	// Cleanup any existing test session
	if client.SessionExists(testSessionName) {
		_ = client.KillSession(testSessionName)
	}

	// Cleanup after test
	t.Cleanup(func() {
		if client.SessionExists(testSessionName) {
			_ = client.KillSession(testSessionName)
		}
	})

	t.Run("should create session from simple template", func(t *testing.T) {
		template := &config.TmuxTemplate{
			SessionName: testSessionName,
			Windows: []config.TmuxWindow{
				{
					Name:   "main",
					Layout: "main-horizontal",
					Panes:  []string{"echo 'Hello World'"},
				},
			},
		}

		variables := map[string]string{
			"Project": "testproject",
		}

		err := client.CreateFromTemplate(template, variables)
		assert.NoError(t, err, "CreateFromTemplate should not error")

		// Verify session exists
		exists := client.SessionExists(testSessionName)
		assert.True(t, exists, "Session should exist after template creation")
	})

	t.Run("should create session with variable substitution", func(t *testing.T) {
		// Kill previous session first
		if client.SessionExists(testSessionName) {
			_ = client.KillSession(testSessionName)
		}

		template := &config.TmuxTemplate{
			SessionName: "{{.Project}}-dev",
			Windows: []config.TmuxWindow{
				{
					Name:   "editor",
					Layout: "main-horizontal",
					Panes:  []string{"cd {{.Project}}", "echo {{.User}}"},
				},
			},
		}

		variables := map[string]string{
			"Project": "myapp",
			"User":    "developer",
		}

		expectedSessionName := "myapp-dev"

		// Cleanup the expected session name too
		if client.SessionExists(expectedSessionName) {
			_ = client.KillSession(expectedSessionName)
		}
		t.Cleanup(func() {
			if client.SessionExists(expectedSessionName) {
				_ = client.KillSession(expectedSessionName)
			}
		})

		err := client.CreateFromTemplate(template, variables)
		assert.NoError(t, err, "CreateFromTemplate should not error with variables")

		// Verify session exists with substituted name
		exists := client.SessionExists(expectedSessionName)
		assert.True(t, exists, "Session should exist with variable-substituted name")
	})

	t.Run("should not create duplicate template session", func(t *testing.T) {
		template := &config.TmuxTemplate{
			SessionName: testSessionName,
			Windows: []config.TmuxWindow{
				{
					Name:  "main",
					Panes: []string{"echo 'test'"},
				},
			},
		}

		// First creation should succeed
		err := client.CreateFromTemplate(template, map[string]string{})
		if !client.SessionExists(testSessionName) {
			require.NoError(t, err, "First CreateFromTemplate should not error")
		}

		// Second creation should fail
		err = client.CreateFromTemplate(template, map[string]string{})
		assert.Error(t, err, "CreateFromTemplate should error when session already exists")
		assert.Contains(t, err.Error(), "already exists", "Error should mention session already exists")
	})
}

// BenchmarkSessionExists benchmarks the SessionExists function
func BenchmarkSessionExists(b *testing.B) {
	if !IsAvailable() {
		b.Skip("tmux not available")
	}

	client, err := NewClient()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.SessionExists("non-existent-session")
	}
}
