package rituals

import (
	"os"
	"testing"

	"github.com/ferg-cod3s/rune/internal/config"
)

func TestFilterEnvironment(t *testing.T) {
	// Prepare environment slice with sensitive and non-sensitive variables
	env := []string{
		"PATH=/usr/bin",
		"HOME=/home/user",
		"RUNE_DEBUG=true",
		"RUNE_ENV=development",
		"RUNE_SENTRY_DSN=https://public@sentry.io/1",
		"GITHUB_TOKEN=abc",
		"AWS_SECRET_ACCESS_KEY=shhh",
		"MY_APP_PASSWORD=secret",
		"CUSTOM_TOKEN=tok",
		"SAFE_VAR=value",
	}

	filtered := filterEnvironment(env)

	// Helper to check presence
	contains := func(key string) bool {
		for _, kv := range filtered {
			if len(kv) >= len(key)+1 && kv[:len(key)+1] == key+"=" {
				return true
			}
		}
		return false
	}

	// Should keep safe and allowed vars
	if !contains("PATH") || !contains("HOME") || !contains("SAFE_VAR") {
		t.Fatalf("expected safe vars to be present: %v", filtered)
	}
	if !contains("RUNE_DEBUG") || !contains("RUNE_ENV") {
		t.Fatalf("expected allowed Rune vars to be present: %v", filtered)
	}

	// Should remove sensitive ones
	sensitiveKeys := []string{
		"RUNE_SENTRY_DSN", "GITHUB_TOKEN", "AWS_SECRET_ACCESS_KEY", "MY_APP_PASSWORD", "CUSTOM_TOKEN",
	}
	for _, k := range sensitiveKeys {
		if contains(k) {
			t.Fatalf("expected %s to be filtered out", k)
		}
	}
}

func TestFilterEnvironmentFromOS(t *testing.T) {
	// Set some env vars in process and ensure they are filtered
	_ = os.Setenv("RUNE_SENTRY_DSN", "dsn")
	_ = os.Setenv("NPM_TOKEN", "tok")
	_ = os.Setenv("SAFE", "1")
	defer os.Unsetenv("RUNE_SENTRY_DSN")
	defer os.Unsetenv("NPM_TOKEN")
	defer os.Unsetenv("SAFE")

	filtered := filterEnvironment(os.Environ())
	for _, disallowed := range []string{"RUNE_SENTRY_DSN", "NPM_TOKEN"} {
		for _, kv := range filtered {
			if len(kv) >= len(disallowed)+1 && kv[:len(disallowed)+1] == disallowed+"=" {
				t.Fatalf("%s should have been filtered", disallowed)
			}
		}
	}
}

func TestNewEngine(t *testing.T) {
	cfg := &config.Config{}
	engine := NewEngine(cfg)

	if engine == nil {
		t.Fatal("NewEngine should not return nil")
	}

	if engine.config != cfg {
		t.Error("Engine config should match provided config")
	}

	if engine.activeSessions == nil {
		t.Error("Engine activeSessions should be initialized")
	}

	if !engine.ptySupport {
		t.Error("Engine ptySupport should be true by default")
	}
}

func TestExpandTemplate(t *testing.T) {
	cfg := &config.Config{}
	engine := NewEngine(cfg)

	tests := []struct {
		name      string
		template  string
		variables map[string]string
		expected  string
	}{
		{
			name:      "simple variable replacement",
			template:  "hello-{{.Project}}",
			variables: map[string]string{"Project": "myproject"},
			expected:  "hello-myproject",
		},
		{
			name:      "multiple variables",
			template:  "{{.User}}/{{.Project}}",
			variables: map[string]string{"User": "john", "Project": "myapp"},
			expected:  "john/myapp",
		},
		{
			name:      "no variables",
			template:  "static-text",
			variables: map[string]string{},
			expected:  "static-text",
		},
		{
			name:      "missing variable unchanged",
			template:  "{{.Missing}}-text",
			variables: map[string]string{},
			expected:  "{{.Missing}}-text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.expandTemplate(tt.template, tt.variables)
			if result != tt.expected {
				t.Errorf("expandTemplate(%q) = %q, want %q", tt.template, result, tt.expected)
			}
		})
	}
}

func TestShouldShowOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected bool
	}{
		{
			name:     "empty output",
			output:   "",
			expected: false,
		},
		{
			name:     "whitespace only",
			output:   "   \n\t  ",
			expected: false,
		},
		{
			name:     "already up to date git message",
			output:   "Already up to date.",
			expected: false,
		},
		{
			name:     "nothing to commit message",
			output:   "nothing to commit, working tree clean",
			expected: false,
		},
		{
			name:     "meaningful output",
			output:   "Build completed successfully",
			expected: true,
		},
		{
			name:     "error output",
			output:   "Error: command failed",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldShowOutput(tt.output)
			if result != tt.expected {
				t.Errorf("shouldShowOutput(%q) = %v, want %v", tt.output, result, tt.expected)
			}
		})
	}
}

func TestTestRitual(t *testing.T) {
	cfg := &config.Config{
		Rituals: config.Rituals{
			Start: config.RitualSet{
				Global: []config.Command{
					{Name: "global start", Command: "echo start"},
				},
				PerProject: map[string][]config.Command{
					"testproject": {
						{Name: "project start", Command: "echo project"},
					},
				},
			},
			Stop: config.RitualSet{
				Global: []config.Command{
					{Name: "global stop", Command: "echo stop"},
				},
			},
		},
	}
	engine := NewEngine(cfg)

	t.Run("test start ritual with project", func(t *testing.T) {
		err := engine.TestRitual("start", "testproject")
		if err != nil {
			t.Errorf("TestRitual should not error: %v", err)
		}
	})

	t.Run("test stop ritual", func(t *testing.T) {
		err := engine.TestRitual("stop", "testproject")
		if err != nil {
			t.Errorf("TestRitual should not error: %v", err)
		}
	})

	t.Run("test unknown ritual type", func(t *testing.T) {
		err := engine.TestRitual("unknown", "testproject")
		if err == nil {
			t.Error("TestRitual should error for unknown type")
		}
	})

	t.Run("test ritual with no commands", func(t *testing.T) {
		emptyCfg := &config.Config{}
		emptyEngine := NewEngine(emptyCfg)
		err := emptyEngine.TestRitual("start", "nonexistent")
		if err != nil {
			t.Errorf("TestRitual should not error for empty config: %v", err)
		}
	})
}

func TestFilterEnvironmentEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		env      []string
		expected map[string]bool // true = should be present
	}{
		{
			name: "azure credentials filtered",
			env: []string{
				"AZURE_CLIENT_SECRET=secret",
				"SAFE=value",
			},
			expected: map[string]bool{
				"AZURE_CLIENT_SECRET": false,
				"SAFE":                true,
			},
		},
		{
			name: "google cloud credentials filtered",
			env: []string{
				"GOOGLE_APPLICATION_CREDENTIALS=/path/to/key",
				"GCP_PROJECT=myproject",
				"NORMAL=value",
			},
			expected: map[string]bool{
				"GOOGLE_APPLICATION_CREDENTIALS": false,
				"GCP_PROJECT":                    false,
				"NORMAL":                         true,
			},
		},
		{
			name: "ssh keys filtered",
			env: []string{
				"SSH_PRIVATE_KEY=key",
				"SSH_AUTH_SOCK=/tmp/socket",
				"NORMAL=value",
			},
			expected: map[string]bool{
				"SSH_PRIVATE_KEY": false,
				"SSH_AUTH_SOCK":   false,
				"NORMAL":          true,
			},
		},
		{
			name: "bearer tokens filtered",
			env: []string{
				"BEARER_TOKEN=token",
				"MY_BEARER=value",
				"NORMAL=value",
			},
			expected: map[string]bool{
				"BEARER_TOKEN": false,
				"MY_BEARER":    false,
				"NORMAL":       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterEnvironment(tt.env)

			contains := func(key string) bool {
				for _, kv := range filtered {
					if len(kv) >= len(key)+1 && kv[:len(key)+1] == key+"=" {
						return true
					}
				}
				return false
			}

			for key, shouldBePresent := range tt.expected {
				isPresent := contains(key)
				if isPresent != shouldBePresent {
					if shouldBePresent {
						t.Errorf("expected %s to be present in filtered env", key)
					} else {
						t.Errorf("expected %s to be filtered out", key)
					}
				}
			}
		})
	}
}
