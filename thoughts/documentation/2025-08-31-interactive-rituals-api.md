---
title: Interactive Rituals - API Reference
audience: api
version: Phase 3 Complete  
date: 2025-08-31
feature: Interactive Rituals with Tmux Integration
---

# Interactive Rituals - API Reference

This document provides comprehensive API reference for Interactive Rituals, including CLI commands, configuration schema, and programmatic interfaces.

## CLI Commands

### Basic Ritual Commands

All existing ritual commands support interactive features transparently.

#### `rune ritual test <start|stop> [project]`

Test interactive rituals without execution.

**Usage:**
```bash
rune ritual test start [project]
rune ritual test stop [project]
```

**Options:**
- `--config string`: Custom config file path
- `--verbose`: Show detailed command information

**Example Output:**
```bash
$ rune ritual test start myproject
🧪 Testing start ritual for project: myproject
Commands that would be executed:
  1. Interactive Development: echo 'Starting...' (optional)
  2. Setup Environment: echo 'Environment ready' (optional)
```

#### `rune start [project]`

Execute start rituals, including interactive commands.

**Usage:**
```bash
rune start [project]
```

**Interactive Behavior:**
- Non-interactive commands execute normally
- Interactive commands launch with TTY allocation
- Tmux sessions are created and attached automatically
- PTY fallback used when tmux unavailable

**Example:**
```bash
$ rune start webapp
🔮 Executing start rituals...
  ⚡ Check Dependencies... ✓
  ⚡ Interactive Development... 📺 (created and attaching to session 'webapp-dev')
# User is now attached to tmux session
```

## Configuration Schema

### Command Configuration

Interactive commands extend the standard `Command` schema:

```yaml
# Standard Command Schema
- name: string                    # Required: Display name
  command: string                 # Required: Command to execute  
  optional: bool                  # Optional: Continue on failure (default: false)
  background: bool                # Optional: Run in background (default: false)
  
  # Interactive Extensions
  interactive: bool               # Optional: Enable interactive mode (default: false)
  tmux_session: string           # Optional: Tmux session name (supports {{.Variable}})
  tmux_template: string          # Optional: Template name to use
```

### Template Schema

Templates define multi-window tmux session layouts:

```yaml
templates:
  template-name:                  # Template identifier
    session_name: string          # Session name (supports {{.Variable}})
    windows:                      # Array of window definitions
      - name: string              # Window name
        layout: string            # Optional: tmux layout name
        panes: string[]           # Array of commands for each pane
```

### Complete Configuration Example

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 30m
  idle_threshold: 10m

projects:
  - name: "webapp"
    detect: ["package.json"]

rituals:
  start:
    global:
      # PTY interactive command
      - name: "Interactive Shell"
        command: "zsh"
        interactive: true
        optional: true
      
      # Tmux session command
      - name: "Development Session"
        command: "echo 'Starting development...'"
        interactive: true
        tmux_session: "dev-{{.Project}}"
        optional: true
        
      # Tmux template command  
      - name: "Full Environment"
        command: "echo 'Setting up environment...'"
        interactive: true
        tmux_template: "fullstack"
        optional: true

  templates:
    fullstack:
      session_name: "fullstack-{{.Project}}"
      windows:
        - name: "editor"
          layout: "main-horizontal"
          panes:
            - "vim ."
            - "git status"
        - name: "servers"
          layout: "even-horizontal"
          panes:
            - "npm run dev"
            - "npm run test:watch"
```

### Validation Rules

1. **Interactive Commands:**
   - Must have `interactive: true`
   - Can specify `tmux_session`, `tmux_template`, or neither (PTY fallback)
   - Cannot specify both `tmux_session` and `tmux_template`
   - Should be marked `optional: true` for non-blocking automation

2. **Tmux Templates:**
   - Template name must exist in `templates` section if referenced
   - `session_name` is required
   - `windows` array must have at least one window
   - Each window must have `name` and `panes` array
   - `panes` array must have at least one command

3. **Variable Expansion:**
   - `{{.Project}}` is supported in `session_name`, `tmux_session`, and pane commands
   - Variables must use exact syntax: `{{.Variable}}`
   - Unknown variables are left unexpanded

## Go API Reference

### Core Types

#### `config.Command`

```go
type Command struct {
    Name         string `yaml:"name" mapstructure:"name"`
    Command      string `yaml:"command" mapstructure:"command"`
    Optional     bool   `yaml:"optional,omitempty" mapstructure:"optional"`
    Background   bool   `yaml:"background,omitempty" mapstructure:"background"`
    Interactive  bool   `yaml:"interactive,omitempty" mapstructure:"interactive"`
    TmuxSession  string `yaml:"tmux_session,omitempty" mapstructure:"tmux_session"`
    TmuxTemplate string `yaml:"tmux_template,omitempty" mapstructure:"tmux_template"`
}
```

#### `config.TmuxTemplate`

```go
type TmuxTemplate struct {
    SessionName string       `yaml:"session_name" mapstructure:"session_name"`
    Windows     []TmuxWindow `yaml:"windows" mapstructure:"windows"`
}
```

#### `config.TmuxWindow`

```go
type TmuxWindow struct {
    Name   string   `yaml:"name" mapstructure:"name"`
    Layout string   `yaml:"layout,omitempty" mapstructure:"layout"`
    Panes  []string `yaml:"panes" mapstructure:"panes"`
}
```

### Ritual Engine API

#### `rituals.Engine`

```go
type Engine struct {
    config        *config.Config
    tmuxClient    *tmux.Client
    ptySupport    bool
    activeSession map[string]*tmux.Session
}
```

**Methods:**

```go
// NewEngine creates a new ritual engine with tmux support
func NewEngine(cfg *config.Config) *Engine

// ExecuteStartRituals executes start rituals, handling interactive commands
func (e *Engine) ExecuteStartRituals(project string) error

// ExecuteStopRituals executes stop rituals  
func (e *Engine) ExecuteStopRituals(project string) error

// TestRitual shows what commands would be executed
func (e *Engine) TestRitual(ritualType string, project string) error
```

### Tmux Client API

#### `tmux.Client`

```go
type Client struct {
    tmux *gotmux.Tmux
}
```

**Methods:**

```go
// NewClient creates a new tmux client
func NewClient() (*Client, error)

// SessionExists checks if a session exists
func (c *Client) SessionExists(sessionName string) bool

// CreateSession creates a new tmux session
func (c *Client) CreateSession(sessionName string) error

// AttachSession attaches to an existing session
func (c *Client) AttachSession(sessionName string) error

// KillSession terminates a session
func (c *Client) KillSession(sessionName string) error

// ListSessions returns all active session names
func (c *Client) ListSessions() ([]string, error)

// CreateFromTemplate creates session from template
func (c *Client) CreateFromTemplate(template *config.TmuxTemplate, variables map[string]string) error
```

**Utility Functions:**

```go
// IsAvailable checks if tmux is available on the system
func IsAvailable() bool

// GetVersion returns tmux version string
func GetVersion() (string, error)

// GetDefaultInstallPath returns platform-specific install path
func GetDefaultInstallPath() string

// IsInTmuxSession checks if currently running inside tmux
func IsInTmuxSession() bool
```

## Error Codes and Messages

### Common Error Patterns

#### Configuration Errors

```
Error: failed to load config: config validation failed: <specific issue>
```

**Causes:**
- Invalid YAML syntax
- Missing required fields
- Invalid template references

**Resolution:** Fix configuration file and validate with `rune config validate`

#### Tmux Errors

```
Error: tmux not available on this system. Please install tmux to use interactive rituals
```

**Cause:** tmux not installed or not in PATH

**Resolution:** Install tmux using system package manager

```
Error: failed to create session '<name>': session already exists
```

**Cause:** Session with same name already running

**Resolution:** Use unique session name or kill existing session

#### PTY Errors

```
Error: failed to start command with PTY: <system error>
```

**Cause:** PTY allocation failed (system-level issue)

**Resolution:** Check system permissions and available PTYs

### Error Handling Flow

1. **Configuration Phase:**
   - Validate interactive command syntax
   - Check template references
   - Warn about missing optional dependencies

2. **Execution Phase:**
   - Detect tmux availability
   - Attempt session creation/attachment
   - Fall back to PTY if tmux unavailable
   - Fall back to standard execution if PTY fails

3. **Cleanup Phase:**
   - Restore terminal state
   - Clean up orphaned processes
   - Log errors for debugging

## Examples

### Creating Interactive Commands Programmatically

```go
package main

import (
    "github.com/ferg-cod3s/rune/internal/config"
    "github.com/ferg-cod3s/rune/internal/rituals"
)

func main() {
    cfg := &config.Config{
        Rituals: config.Rituals{
            Start: config.RitualSet{
                Global: []config.Command{
                    {
                        Name:        "Development Environment",
                        Command:     "echo 'Starting dev environment...'",
                        Interactive: true,
                        TmuxTemplate: "fullstack",
                        Optional:    true,
                    },
                },
            },
            Templates: map[string]config.TmuxTemplate{
                "fullstack": {
                    SessionName: "dev-{{.Project}}",
                    Windows: []config.TmuxWindow{
                        {
                            Name: "editor",
                            Layout: "main-horizontal",
                            Panes: []string{"vim .", "git status"},
                        },
                    },
                },
            },
        },
    }
    
    engine := rituals.NewEngine(cfg)
    err := engine.ExecuteStartRituals("myproject")
    if err != nil {
        panic(err)
    }
}
```

### Testing Templates

```bash
# Test what commands would execute
rune ritual test start myproject

# Validate configuration
rune config validate

# Show current configuration
rune config show
```

### Session Management

```bash
# List all tmux sessions
tmux list-sessions

# Attach to specific session
tmux attach -t session-name

# Kill session
tmux kill-session -t session-name

# Kill all rune sessions (example pattern)
tmux list-sessions | grep "rune-" | cut -d: -f1 | xargs -I {} tmux kill-session -t {}
```

## Integration Examples

### CI/CD Integration

Interactive rituals are automatically skipped in non-TTY environments:

```yaml
# .github/workflows/test.yml
- name: Test Rituals
  run: |
    rune ritual test start myproject  # Works in CI
    # rune start would skip interactive commands in CI
```

### Docker Integration

```dockerfile
# Install tmux for interactive development
RUN apt-get update && apt-get install -y tmux

# Copy rune config
COPY .rune/ /app/.rune/

# Set up development environment
CMD ["rune", "start", "webapp"]
```

### Shell Integration

```bash
# ~/.bashrc or ~/.zshrc
alias dev="rune start"
alias devstop="rune stop"

# Function to attach to project session
dev-attach() {
    local project="${1:-$(basename $(pwd))}"
    tmux attach -t "dev-$project" 2>/dev/null || echo "No session for $project"
}
```

This API reference provides comprehensive coverage of Interactive Rituals functionality, enabling developers to integrate the feature effectively in their workflows and applications.