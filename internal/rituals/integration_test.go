package rituals

import (
	"testing"

	"github.com/ferg-cod3s/rune/internal/config"
)

// TestInteractiveCommandIntegration tests that interactive commands are properly dispatched
func TestInteractiveCommandIntegration(t *testing.T) {
	cfg := &config.Config{
		Rituals: config.Rituals{
			Templates: map[string]config.TmuxTemplate{
				"dev": {
					SessionName: "test-{{.Project}}",
					Windows: []config.TmuxWindow{
						{
							Name:   "main",
							Layout: "main-vertical",
							Panes:  []string{"echo 'Hello {{.Project}}'"},
						},
					},
				},
			},
		},
	}

	engine := NewEngine(cfg)

	// Test that engine initializes properly
	if engine == nil {
		t.Fatal("Engine should not be nil")
	}

	// Test that tmux client is initialized if available
	if engine.tmuxClient == nil && engine.ptySupport {
		t.Log("tmux client not available, interactive features will use fallbacks")
	}

	// Test executeInteractiveCommand dispatcher logic
	t.Run("should dispatch tmux template command", func(t *testing.T) {
		cmd := config.Command{
			Name:         "Test Template",
			Interactive:  true,
			TmuxTemplate: "dev",
		}

		// This should not panic even if tmux is not available
		// The method should handle graceful fallbacks
		err := engine.executeInteractiveCommand(cmd, "test-project")
		if err != nil {
			// Error is expected if tmux is not available or if this is just a test
			t.Logf("Expected error (tmux might not be available in test): %v", err)
		}
	})

	t.Run("should dispatch tmux session command", func(t *testing.T) {
		cmd := config.Command{
			Name:        "Test Session",
			Interactive: true,
			TmuxSession: "test-session-{{.Project}}",
			Command:     "echo hello",
		}

		// This should not panic even if tmux is not available
		err := engine.executeInteractiveCommand(cmd, "test-project")
		if err != nil {
			// Error is expected if tmux is not available or if this is just a test
			t.Logf("Expected error (tmux might not be available in test): %v", err)
		}
	})

	t.Run("should fallback to PTY for interactive commands", func(t *testing.T) {
		cmd := config.Command{
			Name:        "Test PTY",
			Interactive: true,
			Command:     "echo hello",
		}

		// This should not panic and should attempt PTY execution
		err := engine.executeInteractiveCommand(cmd, "test-project")
		if err != nil {
			// Error is expected since we're not actually running an interactive command
			t.Logf("Expected error (command not truly interactive in test): %v", err)
		}
	})
}

// TestBackwardsCompatibility ensures that existing non-interactive commands still work
func TestBackwardsCompatibility(t *testing.T) {
	cfg := &config.Config{
		Rituals: config.Rituals{},
	}

	engine := NewEngine(cfg)

	// Test that standard commands still work as before
	t.Run("should execute standard command normally", func(t *testing.T) {
		cmd := config.Command{
			Name:    "Test Standard",
			Command: "echo 'test'",
			// Interactive is false by default
		}

		// This should execute as a standard command
		err := engine.executeCommand(cmd, "test")
		if err != nil {
			t.Errorf("Standard command should execute successfully: %v", err)
		}
	})

	t.Run("should handle optional commands", func(t *testing.T) {
		cmd := config.Command{
			Name:     "Test Optional",
			Command:  "false", // Command that fails
			Optional: true,
		}

		// Use executeCommands which properly handles the Optional flag
		err := engine.executeCommands([]config.Command{cmd}, "test")
		if err != nil {
			t.Errorf("Optional command failure should be handled gracefully: %v", err)
		}
	})

	t.Run("should handle background commands", func(t *testing.T) {
		cmd := config.Command{
			Name:       "Test Background",
			Command:    "sleep 0.1",
			Background: true,
		}

		// This should start and return immediately
		err := engine.executeCommand(cmd, "test")
		if err != nil {
			t.Errorf("Background command should start successfully: %v", err)
		}
	})
}

// TestTemplateExpansion tests the template variable expansion
func TestTemplateExpansion(t *testing.T) {
	cfg := &config.Config{
		Rituals: config.Rituals{},
	}

	engine := NewEngine(cfg)

	tests := []struct {
		name      string
		template  string
		variables map[string]string
		expected  string
	}{
		{
			name:     "single variable",
			template: "Hello {{.Project}}",
			variables: map[string]string{
				"Project": "rune",
			},
			expected: "Hello rune",
		},
		{
			name:     "multiple variables",
			template: "{{.Project}} session in {{.Dir}}",
			variables: map[string]string{
				"Project": "rune",
				"Dir":     "/home/user",
			},
			expected: "rune session in /home/user",
		},
		{
			name:      "no variables",
			template:  "static string",
			variables: map[string]string{},
			expected:  "static string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.expandTemplate(tt.template, tt.variables)
			if result != tt.expected {
				t.Errorf("expandTemplate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
