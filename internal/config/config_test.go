package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Projects: []Project{
					{Name: "test", Detect: []string{"git:test"}},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid version",
			config: Config{
				Version: 2,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "unsupported config version",
		},
		{
			name: "invalid work hours - zero",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "work_hours must be between 0 and 24",
		},
		{
			name: "invalid work hours - too high",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     25,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "work_hours must be between 0 and 24",
		},
		{
			name: "invalid break interval",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: -1 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "break_interval must be positive",
		},
		{
			name: "invalid idle threshold",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: -1 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "idle_threshold must be positive",
		},
		{
			name: "project with empty name",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Projects: []Project{
					{Name: "", Detect: []string{"git:test"}},
				},
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "project with empty detect patterns",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Projects: []Project{
					{Name: "test", Detect: []string{}},
				},
			},
			wantErr: true,
			errMsg:  "detect patterns cannot be empty",
		},
		{
			name: "template with empty session name",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Rituals: Rituals{
					Templates: map[string]TmuxTemplate{
						"test-template": {
							SessionName: "", // Empty session name
							Windows: []TmuxWindow{
								{Name: "main", Panes: []string{"vim"}},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "session_name cannot be empty",
		},
		{
			name: "template with empty window name",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Rituals: Rituals{
					Templates: map[string]TmuxTemplate{
						"test-template": {
							SessionName: "test-session",
							Windows: []TmuxWindow{
								{Name: "", Panes: []string{"vim"}}, // Empty window name
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "window[0] name cannot be empty",
		},
		{
			name: "command references undefined template",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Rituals: Rituals{
					Start: RitualSet{
						Global: []Command{
							{
								Name:         "test",
								Command:      "echo test",
								Interactive:  true,
								TmuxTemplate: "nonexistent-template", // References undefined template
							},
						},
					},
					Templates: map[string]TmuxTemplate{}, // Empty templates
				},
			},
			wantErr: true,
			errMsg:  "references undefined template",
		},
		{
			name: "interactive command without tmux configuration (PTY fallback)",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Rituals: Rituals{
					Start: RitualSet{
						Global: []Command{
							{
								Name:        "test",
								Command:     "echo test",
								Interactive: true,
								// Missing both TmuxTemplate and TmuxSession - should use PTY fallback
							},
						},
					},
				},
			},
			wantErr: false, // PTY fallback is now allowed
		},
		{
			name: "valid interactive configuration with template",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Rituals: Rituals{
					Start: RitualSet{
						Global: []Command{
							{
								Name:         "test",
								Command:      "echo test",
								Interactive:  true,
								TmuxTemplate: "test-template",
							},
						},
					},
					Templates: map[string]TmuxTemplate{
						"test-template": {
							SessionName: "test-session",
							Windows: []TmuxWindow{
								{Name: "main", Panes: []string{"vim ."}},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid interactive configuration with session",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Rituals: Rituals{
					Start: RitualSet{
						Global: []Command{
							{
								Name:        "test",
								Command:     "echo test",
								Interactive: true,
								TmuxSession: "custom-session",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".rune", "config.yaml")
	assert.Equal(t, expected, path)
}

func TestExists(t *testing.T) {
	// Test when config doesn't exist
	// We'll use a temporary directory to avoid interfering with real config
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	exists, err := Exists()
	require.NoError(t, err)
	assert.False(t, exists)

	// Create config file
	runeDir := filepath.Join(tempDir, ".rune")
	err = os.MkdirAll(runeDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(runeDir, "config.yaml")
	err = os.WriteFile(configPath, []byte("version: 1"), 0644)
	require.NoError(t, err)

	// Test when config exists
	exists, err = Exists()
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestLoadWithProfile_MissingProfile(t *testing.T) {
	// Set up temporary home directory
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create base config
	runeDir := filepath.Join(tempDir, ".rune")
	err := os.MkdirAll(runeDir, 0755)
	require.NoError(t, err)

	baseConfigPath := filepath.Join(runeDir, "config.yaml")
	baseConfig := `version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m
`
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	require.NoError(t, err)

	// Initialize global viper with base config
	viper.Reset()
	viper.SetConfigFile(baseConfigPath)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	require.NoError(t, err)

	// Test loading with non-existent profile
	_, err = LoadWithProfile("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "profile 'nonexistent' not found")
}

func TestLoadWithProfile_EmptyProfile(t *testing.T) {
	// Set up temporary home directory
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create base config
	runeDir := filepath.Join(tempDir, ".rune")
	err := os.MkdirAll(runeDir, 0755)
	require.NoError(t, err)

	baseConfigPath := filepath.Join(runeDir, "config.yaml")
	baseConfig := `version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m
projects:
  - name: base-project
    detect:
      - git:base
`
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	require.NoError(t, err)

	// Initialize global viper with base config
	viper.Reset()
	viper.SetConfigFile(baseConfigPath)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	require.NoError(t, err)

	// Test loading with empty profile (should return base config)
	cfg, err := LoadWithProfile("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 8.0, cfg.Settings.WorkHours)
	assert.Equal(t, 1, len(cfg.Projects))
	assert.Equal(t, "base-project", cfg.Projects[0].Name)
}

func TestLoadWithProfile_ValidProfile(t *testing.T) {
	// Set up temporary home directory
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create base config
	runeDir := filepath.Join(tempDir, ".rune")
	err := os.MkdirAll(runeDir, 0755)
	require.NoError(t, err)

	baseConfigPath := filepath.Join(runeDir, "config.yaml")
	baseConfig := `version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m
  notifications:
    enabled: true
    break_reminders: true
projects:
  - name: base-project
    detect:
      - git:base
rituals:
  start:
    global:
      - name: base-ritual
        command: echo base
  templates:
    base-template:
      session_name: base-session
      windows:
        - name: main
          panes:
            - vim
`
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	require.NoError(t, err)

	// Create profiles directory
	profilesDir := filepath.Join(runeDir, "profiles")
	err = os.MkdirAll(profilesDir, 0755)
	require.NoError(t, err)

	// Create work profile
	workProfilePath := filepath.Join(profilesDir, "work.yaml")
	workProfile := `version: 1
settings:
  work_hours: 10.0
  break_interval: 45m
  notifications:
    enabled: true
    break_reminders: false
projects:
  - name: work-project
    detect:
      - git:work
rituals:
  start:
    global:
      - name: work-ritual
        command: echo work
  templates:
    work-template:
      session_name: work-session
      windows:
        - name: editor
          panes:
            - nvim
`
	err = os.WriteFile(workProfilePath, []byte(workProfile), 0644)
	require.NoError(t, err)

	// Initialize global viper with base config
	viper.Reset()
	viper.SetConfigFile(baseConfigPath)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	require.NoError(t, err)

	// Test loading with work profile
	cfg, err := LoadWithProfile("work")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify merged settings
	assert.Equal(t, 10.0, cfg.Settings.WorkHours, "work_hours should be overridden by profile")
	assert.Equal(t, 45*time.Minute, cfg.Settings.BreakInterval, "break_interval should be overridden by profile")
	assert.Equal(t, 10*time.Minute, cfg.Settings.IdleThreshold, "idle_threshold should remain from base")

	// Verify notifications merge
	assert.True(t, cfg.Settings.Notifications.Enabled)
	assert.False(t, cfg.Settings.Notifications.BreakReminders, "notifications should be overridden by profile")

	// Verify projects are appended
	assert.Equal(t, 2, len(cfg.Projects), "projects should be appended")
	projectNames := []string{cfg.Projects[0].Name, cfg.Projects[1].Name}
	assert.Contains(t, projectNames, "base-project")
	assert.Contains(t, projectNames, "work-project")

	// Verify rituals are overridden
	assert.Equal(t, 1, len(cfg.Rituals.Start.Global))
	assert.Equal(t, "work-ritual", cfg.Rituals.Start.Global[0].Name, "rituals should be overridden by profile")

	// Verify templates are merged
	assert.Equal(t, 2, len(cfg.Rituals.Templates), "templates should be merged")
	assert.Contains(t, cfg.Rituals.Templates, "base-template")
	assert.Contains(t, cfg.Rituals.Templates, "work-template")
}

func TestMergeConfigs(t *testing.T) {
	base := &Config{
		Version: 1,
		UserID:  "base-user",
		Settings: Settings{
			WorkHours:     8.0,
			BreakInterval: 50 * time.Minute,
			IdleThreshold: 10 * time.Minute,
			Notifications: NotificationSettings{
				Enabled:        true,
				BreakReminders: true,
			},
		},
		Projects: []Project{
			{Name: "base-project", Detect: []string{"git:base"}},
		},
		Rituals: Rituals{
			Start: RitualSet{
				Global: []Command{{Name: "base-start", Command: "echo base"}},
			},
			Templates: map[string]TmuxTemplate{
				"base-template": {
					SessionName: "base-session",
					Windows: []TmuxWindow{
						{Name: "main", Panes: []string{"vim"}},
					},
				},
			},
		},
		Logging: Logging{
			Level:  "info",
			Format: "text",
		},
	}

	profile := &Config{
		Version: 1,
		UserID:  "profile-user",
		Settings: Settings{
			WorkHours:     10.0,
			BreakInterval: 45 * time.Minute,
			// IdleThreshold not set - should use base value
			Notifications: NotificationSettings{
				Enabled:        true,
				BreakReminders: false,
			},
		},
		Projects: []Project{
			{Name: "profile-project", Detect: []string{"git:profile"}},
		},
		Rituals: Rituals{
			Start: RitualSet{
				Global: []Command{{Name: "profile-start", Command: "echo profile"}},
			},
			Templates: map[string]TmuxTemplate{
				"profile-template": {
					SessionName: "profile-session",
					Windows: []TmuxWindow{
						{Name: "editor", Panes: []string{"nvim"}},
					},
				},
			},
		},
		Logging: Logging{
			Level: "debug",
			// Format not set - should use base value
		},
	}

	merged := mergeConfigs(base, profile)

	// Verify overrides
	assert.Equal(t, "profile-user", merged.UserID)
	assert.Equal(t, 10.0, merged.Settings.WorkHours)
	assert.Equal(t, 45*time.Minute, merged.Settings.BreakInterval)

	// Verify base values preserved when profile doesn't override
	assert.Equal(t, 10*time.Minute, merged.Settings.IdleThreshold)

	// Verify notifications merge
	assert.True(t, merged.Settings.Notifications.Enabled)
	assert.False(t, merged.Settings.Notifications.BreakReminders)

	// Verify projects are appended
	assert.Equal(t, 2, len(merged.Projects))
	assert.Equal(t, "base-project", merged.Projects[0].Name)
	assert.Equal(t, "profile-project", merged.Projects[1].Name)

	// Verify rituals are overridden
	assert.Equal(t, 1, len(merged.Rituals.Start.Global))
	assert.Equal(t, "profile-start", merged.Rituals.Start.Global[0].Name)

	// Verify templates are merged
	assert.Equal(t, 2, len(merged.Rituals.Templates))
	assert.Contains(t, merged.Rituals.Templates, "base-template")
	assert.Contains(t, merged.Rituals.Templates, "profile-template")

	// Verify logging merge
	assert.Equal(t, "debug", merged.Logging.Level)
	assert.Equal(t, "text", merged.Logging.Format)
}

func TestGetProfilePath(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	path, err := GetProfilePath("work")
	require.NoError(t, err)

	expected := filepath.Join(tempDir, ".rune", "profiles", "work.yaml")
	assert.Equal(t, expected, path)
}
