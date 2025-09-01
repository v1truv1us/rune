package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the main configuration structure
type Config struct {
	Version      int          `yaml:"version" mapstructure:"version"`
	UserID       string       `yaml:"user_id" mapstructure:"user_id"`
	Settings     Settings     `yaml:"settings" mapstructure:"settings"`
	Projects     []Project    `yaml:"projects" mapstructure:"projects"`
	Rituals      Rituals      `yaml:"rituals" mapstructure:"rituals"`
	Integrations Integrations `yaml:"integrations" mapstructure:"integrations"`
	Logging      Logging      `yaml:"logging" mapstructure:"logging"`
}

// Settings contains global application settings
type Settings struct {
	WorkHours     float64              `yaml:"work_hours" mapstructure:"work_hours"`
	BreakInterval time.Duration        `yaml:"break_interval" mapstructure:"break_interval"`
	IdleThreshold time.Duration        `yaml:"idle_threshold" mapstructure:"idle_threshold"`
	Notifications NotificationSettings `yaml:"notifications" mapstructure:"notifications"`
}

// NotificationSettings contains notification preferences
type NotificationSettings struct {
	Enabled           bool `yaml:"enabled" mapstructure:"enabled"`
	BreakReminders    bool `yaml:"break_reminders" mapstructure:"break_reminders"`
	EndOfDayReminders bool `yaml:"end_of_day_reminders" mapstructure:"end_of_day_reminders"`
	SessionComplete   bool `yaml:"session_complete" mapstructure:"session_complete"`
	IdleDetection     bool `yaml:"idle_detection" mapstructure:"idle_detection"`
	Sound             bool `yaml:"sound" mapstructure:"sound"`
}

// Project represents a project configuration
type Project struct {
	Name   string   `yaml:"name" mapstructure:"name"`
	Detect []string `yaml:"detect" mapstructure:"detect"`
}

// Rituals contains start and stop ritual configurations
type Rituals struct {
	Start     RitualSet               `yaml:"start" mapstructure:"start"`
	Stop      RitualSet               `yaml:"stop" mapstructure:"stop"`
	Templates map[string]TmuxTemplate `yaml:"templates,omitempty" mapstructure:"templates"`
}

// RitualSet contains global and per-project rituals
type RitualSet struct {
	Global     []Command            `yaml:"global" mapstructure:"global"`
	PerProject map[string][]Command `yaml:"per_project" mapstructure:"per_project"`
}

// Command represents a ritual command
type Command struct {
	Name         string `yaml:"name" mapstructure:"name"`
	Command      string `yaml:"command" mapstructure:"command"`
	Optional     bool   `yaml:"optional" mapstructure:"optional"`
	Background   bool   `yaml:"background" mapstructure:"background"`
	Interactive  bool   `yaml:"interactive" mapstructure:"interactive"`
	TmuxSession  string `yaml:"tmux_session,omitempty" mapstructure:"tmux_session"`
	TmuxTemplate string `yaml:"tmux_template,omitempty" mapstructure:"tmux_template"`
}

// TmuxTemplate represents a tmux session template configuration
type TmuxTemplate struct {
	SessionName string       `yaml:"session_name" mapstructure:"session_name"`
	Windows     []TmuxWindow `yaml:"windows" mapstructure:"windows"`
}

// TmuxWindow represents a tmux window configuration
type TmuxWindow struct {
	Name   string   `yaml:"name" mapstructure:"name"`
	Layout string   `yaml:"layout,omitempty" mapstructure:"layout"`
	Panes  []string `yaml:"panes" mapstructure:"panes"`
}

// Integrations contains external service integrations
type Integrations struct {
	Git       GitIntegration       `yaml:"git" mapstructure:"git"`
	Slack     SlackIntegration     `yaml:"slack" mapstructure:"slack"`
	Calendar  CalendarIntegration  `yaml:"calendar" mapstructure:"calendar"`
	Telemetry TelemetryIntegration `yaml:"telemetry" mapstructure:"telemetry"`
}

// GitIntegration contains Git-related settings
type GitIntegration struct {
	Enabled           bool `yaml:"enabled" mapstructure:"enabled"`
	AutoDetectProject bool `yaml:"auto_detect_project" mapstructure:"auto_detect_project"`
}

// SlackIntegration contains Slack-related settings
type SlackIntegration struct {
	Workspace  string `yaml:"workspace" mapstructure:"workspace"`
	DNDOnStart bool   `yaml:"dnd_on_start" mapstructure:"dnd_on_start"`
}

// CalendarIntegration contains calendar-related settings
type CalendarIntegration struct {
	Provider      string `yaml:"provider" mapstructure:"provider"`
	BlockCalendar bool   `yaml:"block_calendar" mapstructure:"block_calendar"`
}

// TelemetryIntegration contains telemetry-related settings
type TelemetryIntegration struct {
	Enabled   bool   `yaml:"enabled" mapstructure:"enabled"`
	SentryDSN string `yaml:"sentry_dsn" mapstructure:"sentry_dsn"`
}

// Logging contains logging configuration
type Logging struct {
	Level     string `yaml:"level" mapstructure:"level"`           // debug, info, warn, error
	Format    string `yaml:"format" mapstructure:"format"`         // text, json
	Output    string `yaml:"output" mapstructure:"output"`         // stdout, stderr, or file path
	ErrorFile string `yaml:"error_file" mapstructure:"error_file"` // JSON file for structured error logging
}

// Load loads the configuration from the default location or specified file
func Load() (*Config, error) {
	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version: %d (expected: 1)", c.Version)
	}

	if c.Settings.WorkHours <= 0 || c.Settings.WorkHours > 24 {
		return fmt.Errorf("work_hours must be between 0 and 24, got: %f", c.Settings.WorkHours)
	}

	if c.Settings.BreakInterval <= 0 {
		return fmt.Errorf("break_interval must be positive, got: %v", c.Settings.BreakInterval)
	}

	if c.Settings.IdleThreshold <= 0 {
		return fmt.Errorf("idle_threshold must be positive, got: %v", c.Settings.IdleThreshold)
	}

	// Validate projects
	for i, project := range c.Projects {
		if project.Name == "" {
			return fmt.Errorf("project[%d]: name cannot be empty", i)
		}
		if len(project.Detect) == 0 {
			return fmt.Errorf("project[%d]: detect patterns cannot be empty", i)
		}
	}

	// Validate templates
	for templateName, template := range c.Rituals.Templates {
		if template.SessionName == "" {
			return fmt.Errorf("template '%s': session_name cannot be empty", templateName)
		}
		for i, window := range template.Windows {
			if window.Name == "" {
				return fmt.Errorf("template '%s': window[%d] name cannot be empty", templateName, i)
			}
		}
	}

	// Validate template references in commands
	allCommands := [][]Command{c.Rituals.Start.Global, c.Rituals.Stop.Global}
	for _, commands := range c.Rituals.Start.PerProject {
		allCommands = append(allCommands, commands)
	}
	for _, commands := range c.Rituals.Stop.PerProject {
		allCommands = append(allCommands, commands)
	}

	for _, commands := range allCommands {
		for _, cmd := range commands {
			// Validate template references
			if cmd.TmuxTemplate != "" {
				if _, exists := c.Rituals.Templates[cmd.TmuxTemplate]; !exists {
					return fmt.Errorf("command '%s' references undefined template '%s'", cmd.Name, cmd.TmuxTemplate)
				}
			}

			// Interactive commands can use tmux_template, tmux_session, or fallback to PTY
			// No validation needed - all combinations are valid
		}
	}

	return nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".rune", "config.yaml"), nil
}

// Exists checks if the configuration file exists
func Exists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check config file: %w", err)
	}

	return true, nil
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	return Load()
}

// SaveConfig saves the configuration to file
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.Set("version", cfg.Version)
	viper.Set("user_id", cfg.UserID)
	viper.Set("settings", cfg.Settings)
	viper.Set("projects", cfg.Projects)
	viper.Set("rituals", cfg.Rituals)
	viper.Set("integrations", cfg.Integrations)

	return viper.WriteConfigAs(configPath)
}
