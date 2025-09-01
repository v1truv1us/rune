---
title: Interactive Rituals - User Guide
audience: user
version: Phase 3 Complete
date: 2025-08-31
feature: Interactive Rituals with Tmux Integration
---

# Interactive Rituals - User Guide

Interactive Rituals transform Rune from simple command automation to sophisticated development environment management. Create multi-pane tmux sessions, interactive terminals, and template-based development setups that integrate seamlessly with your daily workflow.

## Overview

Interactive Rituals provide:

- 🖥️ **Interactive terminal commands** with full TTY support
- 📺 **Automated tmux session management** with custom layouts
- 🎯 **Template-based environments** for complex development setups
- 🔄 **Variable expansion** for dynamic customization
- ↩️ **Graceful fallbacks** when advanced features are unavailable

## Prerequisites

### Required
- Rune v1.0+ with interactive ritual support
- Terminal emulator (Terminal.app, iTerm2, GNOME Terminal, etc.)

### Optional (for advanced features)
- **tmux 2.0+** for session management and multi-pane layouts
  - macOS: `brew install tmux`
  - Ubuntu/Debian: `sudo apt install tmux`
  - CentOS/RHEL: `sudo yum install tmux`

### Platform Support
- ✅ **macOS**: Full support
- ✅ **Linux**: Full support  
- ⚠️ **Windows**: Requires WSL2 for tmux features

## Getting Started

### 1. Basic Interactive Command

Add an interactive command to your ritual configuration:

```yaml
# ~/.rune/config.yaml
rituals:
  start:
    global:
      - name: "Interactive Shell"
        command: "zsh"
        interactive: true  # Enables interactive mode
        optional: true
```

**What happens**: Rune launches an interactive shell using PTY (pseudoterminal) allocation, giving you full terminal functionality.

### 2. Tmux Session Management

Create and manage tmux sessions automatically:

```yaml
rituals:
  start:
    per_project:
      myproject:
        - name: "Development Session"
          command: "echo 'Starting development for myproject...'"
          interactive: true
          tmux_session: "dev-{{.Project}}"  # Creates session named "dev-myproject"
          optional: true
```

**What happens**: 
1. Rune creates a tmux session named "dev-myproject" 
2. Runs the specified command in the session
3. Attaches you to the session for interactive use

### 3. Template-Based Development Environment

Create sophisticated multi-pane layouts using templates:

```yaml
rituals:
  start:
    global:
      - name: "Full Development Environment"
        command: "echo 'Setting up development environment...'"
        interactive: true
        tmux_template: "fullstack-dev"  # Uses template defined below
        optional: true

  # Template definitions
  templates:
    fullstack-dev:
      session_name: "fullstack-{{.Project}}"
      windows:
        - name: "editor"
          layout: "main-horizontal"
          panes:
            - "nvim ."          # Opens editor in top pane
            - "git status"      # Shows git status in bottom pane
        - name: "servers"
          panes:
            - "npm run dev"     # Starts dev server
            - "npm run test:watch"  # Runs tests in watch mode
```

**What happens**:
1. Creates tmux session "fullstack-myproject"
2. Creates "editor" window with horizontal split (editor top, git bottom)
3. Creates "servers" window with two panes (dev server and test runner)
4. Attaches you to the session

## Configuration Reference

### Interactive Command Fields

```yaml
- name: "Command Name"           # Required: Display name
  command: "command to run"      # Required: Command to execute
  interactive: true              # Required: Enables interactive mode
  tmux_session: "session-name"   # Optional: Creates/attaches to tmux session
  tmux_template: "template-name" # Optional: Uses predefined template
  optional: true                 # Recommended: Prevents blocking automation
  background: false              # Not used with interactive commands
```

### Template Configuration

```yaml
templates:
  template-name:
    session_name: "session-{{.Project}}"  # Session name with variable expansion
    windows:
      - name: "window-name"               # Window display name
        layout: "main-horizontal"         # Optional: Window layout
        panes:                           # List of commands for each pane
          - "command for pane 1"
          - "command for pane 2"
          - "command for pane 3"
```

### Supported Layouts

- `main-horizontal`: Large pane on top, smaller panes below
- `main-vertical`: Large pane on left, smaller panes on right
- `even-horizontal`: All panes equal height
- `even-vertical`: All panes equal width
- `tiled`: Grid layout with equal-sized panes

### Variable Expansion

Use `{{.Variable}}` syntax in templates and session names:

- `{{.Project}}`: Current project name (auto-detected)
- More variables planned for future versions

## Step-by-Step Examples

### Example 1: Simple Interactive Development

**Goal**: Create an interactive development session for a web project.

**Step 1**: Add to your config:
```yaml
rituals:
  start:
    per_project:
      webapp:
        - name: "Start Development"
          command: "echo 'Starting webapp development...'"
          interactive: true
          tmux_session: "webapp-dev"
          optional: true
```

**Step 2**: Start your work session:
```bash
rune start
```

**Result**: Tmux session "webapp-dev" is created and you're attached to it.

### Example 2: Multi-Pane Development Environment

**Goal**: Create a full-stack development environment with code editor, dev server, and test runner.

**Step 1**: Add template to your config:
```yaml
rituals:
  start:
    global:
      - name: "Full Development Setup"
        interactive: true
        tmux_template: "fullstack"
        optional: true

  templates:
    fullstack:
      session_name: "dev-{{.Project}}"
      windows:
        - name: "code"
          layout: "main-vertical"
          panes:
            - "code ."                    # VS Code
            - "git log --oneline -10"     # Recent commits
        - name: "frontend" 
          layout: "even-horizontal"
          panes:
            - "cd frontend && npm start" # Dev server
            - "cd frontend && npm test"   # Test runner
        - name: "backend"
          panes:
            - "cd backend && go run main.go"  # Backend server
        - name: "monitoring"
          layout: "tiled"
          panes:
            - "htop"                      # System monitor
            - "docker stats"              # Container stats
            - "tail -f logs/app.log"      # Application logs
```

**Step 2**: Start development:
```bash
cd /path/to/your/project
rune start
```

**Result**: Creates session "dev-myproject" with 4 windows:
1. **code**: VS Code editor + git log
2. **frontend**: Dev server + test runner  
3. **backend**: Go server running
4. **monitoring**: System stats, Docker stats, and logs

**Step 3**: Navigate between windows:
- `Ctrl-b + 1` - Switch to code window
- `Ctrl-b + 2` - Switch to frontend window
- `Ctrl-b + 3` - Switch to backend window
- `Ctrl-b + 4` - Switch to monitoring window

### Example 3: Database Development Environment

**Goal**: Set up environment for database-heavy development.

**Configuration**:
```yaml
templates:
  database-dev:
    session_name: "db-dev-{{.Project}}"
    windows:
      - name: "editor"
        panes: ["vim ."]
      - name: "database"
        layout: "main-horizontal"
        panes:
          - "docker-compose up postgres"   # Start database
          - "psql -h localhost -U dev"     # Database client
      - name: "migration"
        panes:
          - "npm run db:migrate"           # Run migrations
          - "npm run db:seed"              # Seed data
```

## Troubleshooting

### Interactive Command Doesn't Start

**Symptom**: Command appears to run but no interactive terminal appears.

**Cause**: Missing `interactive: true` flag in command configuration.

**Fix**: Add `interactive: true` to your command:
```yaml
- name: "My Command"
  command: "vim"
  interactive: true  # Add this line
```

### Tmux Session Creation Fails

**Symptom**: Error message "tmux not available" or "failed to create session".

**Cause**: tmux not installed or not in PATH.

**Fix**:
1. Install tmux:
   - macOS: `brew install tmux`
   - Linux: `sudo apt install tmux` or `sudo yum install tmux`
2. Verify installation: `tmux -V`
3. Restart terminal and try again

### Session Already Exists Error

**Symptom**: "session 'name' already exists" error.

**Cause**: Previous session with same name still running.

**Fix**: Choose one option:
1. **Kill existing session**: `tmux kill-session -t session-name`
2. **Attach to existing session**: `tmux attach -t session-name`  
3. **Use unique session names**: Add timestamp or unique suffix

### Pane Commands Don't Execute

**Symptom**: Tmux session creates but pane commands don't run.

**Cause**: Commands may have syntax errors or missing dependencies.

**Fix**:
1. Test commands manually: `cd /path/to/project && your-command`
2. Check for typos in template configuration
3. Ensure all required tools are installed
4. Use absolute paths if relative paths don't work

### Terminal Size Issues

**Symptom**: Terminal appears corrupted or incorrectly sized.

**Cause**: Terminal size not properly detected.

**Fix**: 
1. Resize terminal window manually
2. In tmux: `Ctrl-b + :` then type `resize-window -A`
3. Detach and reattach: `Ctrl-b + d` then `tmux attach`

### Permission Denied Errors

**Symptom**: "permission denied" when creating tmux sessions.

**Cause**: tmux socket permissions or user access issues.

**Fix**:
1. Check tmux socket: `ls -la /tmp/tmux-*`
2. Kill tmux server: `tmux kill-server`
3. Restart tmux: `tmux new -s test`
4. If persistent, check user groups and permissions

### Windows/WSL Issues

**Symptom**: tmux features don't work on Windows.

**Cause**: tmux requires Unix-like environment.

**Fix**:
1. **Use WSL2**: Install Windows Subsystem for Linux 2
2. **Install tmux in WSL**: `sudo apt install tmux`
3. **Run Rune in WSL**: Execute rune commands from WSL terminal
4. **Alternative**: Use PTY-only interactive commands (work everywhere)

## Best Practices

### 1. Always Use Optional Flag
Interactive commands can block automation, so always mark them optional:

```yaml
- name: "Interactive Command"
  interactive: true
  optional: true  # Critical for non-blocking automation
```

### 2. Meaningful Session Names
Use descriptive names with project variables:

```yaml
# Good: Clear project association
session_name: "{{.Project}}-development"

# Less useful: Generic name
session_name: "dev"
```

### 3. Progressive Enhancement
Start simple and add complexity as needed:

```yaml
# Step 1: Basic interactive command
- name: "Development Shell"
  command: "zsh"
  interactive: true
  
# Step 2: Add tmux session  
- name: "Development Shell"
  interactive: true
  tmux_session: "dev-{{.Project}}"

# Step 3: Use template for complex setup
- name: "Development Environment"
  interactive: true
  tmux_template: "fullstack-dev"
```

### 4. Logical Window Organization
Group related functionality in windows:

```yaml
windows:
  - name: "editor"     # Code editing and file management
  - name: "servers"    # Development servers and services
  - name: "testing"    # Test runners and quality checks  
  - name: "monitoring" # System monitoring and logs
```

### 5. Error Handling
Include commands that can handle failures gracefully:

```yaml
panes:
  - "npm install && npm start"           # Install then start
  - "docker-compose up || echo 'Docker not available'"  # Fallback message
```

## Integration with Rune Workflow

Interactive rituals integrate seamlessly with your daily workflow:

### Morning Startup
```bash
rune start           # Creates development environment
# Work in the interactive tmux session
rune pause           # Pause tracking (keeps session)  
# Continue working...
rune resume          # Resume tracking
```

### End of Day
```bash
rune stop            # Run cleanup commands
# Tmux session persists - you can detach/reattach
```

### Session Management
```bash
# List active tmux sessions
tmux list-sessions

# Attach to specific session  
tmux attach -t session-name

# Detach from session (keeps it running)
# Press: Ctrl-b then d
```

## Next Steps

1. **Start Simple**: Add one interactive command to test the feature
2. **Experiment**: Try tmux session management with your current project  
3. **Create Templates**: Build templates for your common development patterns
4. **Customize**: Adapt layouts and commands to your workflow
5. **Share**: Export working templates for team consistency

Interactive Rituals transform Rune from automation tool to development environment orchestrator. Start with basic interactive commands and gradually build sophisticated multi-pane development environments tailored to your projects and workflow.