# Interactive Rituals Guide

Rune supports interactive rituals that provide terminal-based development environments using tmux sessions or direct PTY allocation. This enables you to create sophisticated multi-pane development setups and interactive command execution.

## Overview

Interactive rituals extend Rune's automation capabilities by providing:

- 🖥️ **PTY-based interactive commands** with full terminal support
- 📺 **Tmux session management** with automatic creation and attachment
- 🎯 **Template-based multi-pane layouts** for complex development environments
- 🔄 **Variable expansion** for dynamic session and command customization
- ↩️ **Graceful fallbacks** when tmux is unavailable

## Configuration

### Basic Interactive Command

The simplest interactive command uses PTY for direct terminal access:

```yaml
rituals:
  start:
    global:
      - name: "Interactive Shell"
        command: "zsh"
        interactive: true  # Enables interactive mode
        optional: true
```

### Tmux Session Commands

Create and attach to a tmux session:

```yaml
rituals:
  start:
    global:
      - name: "Development Session"
        command: "echo 'Starting development...'"
        interactive: true
        tmux_session: "dev-{{.Project}}"  # Session name with variable expansion
        optional: true
```

### Tmux Template Commands

Use predefined templates for complex multi-pane environments:

```yaml
rituals:
  start:
    global:
      - name: "Full Development Environment"
        command: "echo 'Setting up full dev environment...'"
        interactive: true
        tmux_template: "fullstack-dev"  # References template below
        optional: true

  templates:
    fullstack-dev:
      session_name: "dev-{{.Project}}"
      windows:
        - name: "editor"
          layout: "main-horizontal"
          panes:
            - "vim ."
            - "git status"
        - name: "servers"
          panes:
            - "npm run dev"
            - "npm run test:watch"
```

## Template Configuration

Templates define multi-window, multi-pane tmux sessions:

### Template Structure

```yaml
templates:
  template-name:
    session_name: "session-{{.Project}}"  # Session name with variables
    windows:
      - name: "window1"           # Window name
        layout: "main-horizontal" # Optional layout
        panes:                    # Commands for each pane
          - "command1"
          - "command2"
      - name: "window2"
        panes:
          - "command3"
```

### Supported Layouts

- `main-horizontal` - Large pane on top, smaller panes below
- `main-vertical` - Large pane on left, smaller panes on right  
- `even-horizontal` - All panes same height
- `even-vertical` - All panes same width
- `tiled` - Grid layout with equal-sized panes

### Variable Expansion

Templates support variable substitution using `{{.Variable}}` syntax:

- `{{.Project}}` - Current project name
- `{{.Command}}` - Original command from ritual
- Custom variables can be added in future versions

## Execution Flow

1. **Interactive Detection**: Commands with `interactive: true` are processed differently
2. **Tmux Availability Check**: If tmux is available, proceed with tmux features
3. **Template Processing**: If `tmux_template` specified, load template and expand variables
4. **Session Management**: If `tmux_session` specified, create or attach to session
5. **PTY Fallback**: If no tmux configuration, use direct PTY execution
6. **Error Handling**: Graceful degradation with helpful error messages

## Examples

### Simple Development Session

```yaml
rituals:
  start:
    per_project:
      webapp:
        - name: "Start Dev Environment" 
          command: "echo 'Starting webapp development...'"
          interactive: true
          tmux_session: "webapp-dev"
          optional: true
```

### Complex Multi-Pane Setup

```yaml
rituals:
  start:
    global:
      - name: "Full Stack Development"
        command: "echo 'Setting up development environment...'"
        interactive: true
        tmux_template: "fullstack"
        optional: true

  templates:
    fullstack:
      session_name: "fullstack-{{.Project}}"
      windows:
        - name: "code"
          layout: "main-vertical"
          panes:
            - "nvim ."
            - "git log --oneline -10"
            
        - name: "frontend"
          layout: "even-horizontal"
          panes:
            - "cd frontend && npm run dev"
            - "cd frontend && npm run test:watch"
            - "cd frontend && npm run lint:watch"
            
        - name: "backend"  
          panes:
            - "cd backend && go run main.go"
            - "cd backend && go test ./... -watch"
            
        - name: "monitoring"
          layout: "tiled"
          panes:
            - "htop"
            - "docker stats"
            - "tail -f logs/app.log"
            - "netstat -tuln"
```

### PTY-Only Interactive Command

```yaml
rituals:
  start:
    global:
      - name: "Interactive Git Status"
        command: "git status && git log --oneline -5"
        interactive: true  # Will use PTY since no tmux config
        optional: true
```

## Best Practices

### 1. Always Use Optional Flag
Interactive commands should be optional to prevent blocking automation:

```yaml
- name: "Development Environment"
  interactive: true
  tmux_template: "dev"
  optional: true  # Important!
```

### 2. Meaningful Session Names
Use descriptive session names with project variables:

```yaml
session_name: "{{.Project}}-development"  # Good
session_name: "dev"                       # Less useful
```

### 3. Logical Window Organization
Group related tasks in windows:

```yaml
windows:
  - name: "editor"     # Code editing
  - name: "servers"    # Running services  
  - name: "testing"    # Tests and quality
  - name: "monitoring" # System monitoring
```

### 4. Progressive Enhancement
Start simple and add complexity:

```yaml
# Start with basic session
tmux_session: "dev-{{.Project}}"

# Later, upgrade to template
tmux_template: "advanced-dev"
```

## Common Patterns

### Development Server Setup

```yaml
templates:
  webdev:
    session_name: "webdev-{{.Project}}"
    windows:
      - name: "editor"
        panes: ["code ."]
      - name: "server"
        panes: 
          - "npm run dev"
          - "npm run build:watch"
      - name: "testing"
        panes:
          - "npm run test:watch"
          - "npm run e2e:watch"
```

### Full-Stack Development

```yaml
templates:
  fullstack:
    session_name: "fullstack-{{.Project}}"
    windows:
      - name: "frontend"
        layout: "main-horizontal" 
        panes:
          - "cd frontend && npm start"
          - "cd frontend && npm test"
      - name: "backend"
        panes:
          - "cd backend && go run main.go"
          - "cd backend && go test ./..."
      - name: "database"
        panes:
          - "docker-compose up postgres"
          - "psql -h localhost -U dev"
```

### Monitoring Dashboard

```yaml
templates:
  monitor:
    session_name: "monitor-{{.Project}}"
    windows:
      - name: "system"
        layout: "tiled"
        panes:
          - "htop"
          - "iotop" 
          - "nethogs"
          - "df -h"
      - name: "logs"
        layout: "even-horizontal"
        panes:
          - "tail -f /var/log/syslog"
          - "tail -f logs/app.log"
          - "journalctl -f"
```

## Integration with Rune Workflow

Interactive rituals integrate seamlessly with Rune's session management:

1. **Start Ritual**: Creates interactive environment
2. **Work Session**: Use the interactive environment  
3. **Stop Ritual**: Clean up and finalize work

Example workflow:
```bash
rune start           # Creates tmux session with dev environment
# Work in the created tmux session
rune stop            # Runs cleanup commands, can preserve session
```

## Platform Support

- ✅ **macOS**: Full support with Homebrew tmux
- ✅ **Linux**: Full support with system tmux
- ⚠️ **Windows**: Limited support via WSL2
- ✅ **PTY Fallback**: Works on all platforms

## Security Considerations

Interactive rituals maintain Rune's security model:

- ✅ **Environment filtering** applied to all interactive commands
- ✅ **No secret exposure** in tmux sessions  
- ✅ **Optional by design** - won't break automation
- ✅ **User control** - sessions can be detached/killed manually

Interactive rituals provide a powerful way to create sophisticated development environments while maintaining the simplicity and security that makes Rune effective for daily workflow automation.