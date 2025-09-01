---
title: Interactive Rituals - Troubleshooting Guide
audience: mixed
version: Phase 3 Complete
date: 2025-08-31
feature: Interactive Rituals Troubleshooting
---

# Interactive Rituals - Troubleshooting Guide

This guide covers common issues, their root causes, and step-by-step resolution procedures for Interactive Rituals.

## Quick Diagnostics

### 1. Check Feature Availability

```bash
# Verify rune supports interactive rituals
rune --version
# Should show version with interactive support

# Check tmux availability (optional)
tmux -V
# Should show tmux version if installed

# Test basic functionality
rune config validate
# Should pass without interactive-related errors
```

### 2. Test Configuration Syntax

```bash
# Validate your config
rune config validate

# Show parsed configuration
rune config show

# Test ritual without execution
rune ritual test start [project]
```

### 3. Check System Requirements

```bash
# Verify PTY support (Unix systems)
ls -la /dev/ptmx
# Should exist and be readable

# Check tmux socket permissions (if using tmux)
ls -la /tmp/tmux-$(id -u)/
# Should show socket files if tmux is running
```

## Common Issues and Solutions

### Issue 1: Interactive Commands Don't Start

#### Symptoms
- Command appears to execute but no terminal interaction
- Process completes immediately without user input
- "Command finished" message appears instantly

#### Diagnosis
```bash
# Check if interactive flag is set
rune config show | grep -A5 -B5 interactive

# Test with explicit interactive command
rune ritual test start
```

#### Root Causes and Fixes

**Cause 1: Missing `interactive: true` flag**
```yaml
# ❌ Incorrect
- name: "Interactive Command"
  command: "vim file.txt"

# ✅ Correct  
- name: "Interactive Command"
  command: "vim file.txt"
  interactive: true
```

**Cause 2: Command doesn't require interaction**
```yaml
# ❌ Non-interactive command marked interactive
- name: "List Files"
  command: "ls -la"
  interactive: true

# ✅ Actually interactive command
- name: "Edit File"
  command: "vim config.yaml"
  interactive: true
```

**Cause 3: Running in non-TTY environment (CI/CD)**
- Interactive commands are automatically skipped in CI/CD
- This is expected behavior
- Test locally with proper terminal

### Issue 2: Tmux Session Creation Fails

#### Symptoms
- "tmux not available" error
- "failed to create session" error
- Session creation hangs or times out

#### Diagnosis
```bash
# Check tmux installation
which tmux
tmux -V

# Test basic tmux functionality
tmux new-session -d -s test-session
tmux list-sessions
tmux kill-session -t test-session

# Check socket permissions
ls -la /tmp/tmux-$(id -u)/
```

#### Root Causes and Fixes

**Cause 1: tmux not installed**
```bash
# Install tmux
# macOS:
brew install tmux

# Ubuntu/Debian:
sudo apt update && sudo apt install tmux

# CentOS/RHEL:
sudo yum install tmux

# Verify installation
tmux -V
```

**Cause 2: tmux server not running or corrupted**
```bash
# Kill tmux server and restart
tmux kill-server

# Test session creation
tmux new-session -d -s test
tmux list-sessions
```

**Cause 3: Socket permission issues**
```bash
# Check socket permissions
ls -la /tmp/tmux-$(id -u)/

# If permissions are wrong, fix them
sudo chown -R $(whoami) /tmp/tmux-$(id -u)/
sudo chmod 700 /tmp/tmux-$(id -u)/

# Or kill server to recreate
tmux kill-server
```

**Cause 4: Session name conflicts**
```bash
# List existing sessions
tmux list-sessions

# Kill conflicting session
tmux kill-session -t session-name

# Or use unique session names with timestamps
session_name: "dev-{{.Project}}-$(date +%H%M%S)"
```

### Issue 3: Template Processing Errors

#### Symptoms
- "template not found" error
- "failed to create session from template" error
- Template creates session but layout is wrong

#### Diagnosis
```bash
# Validate template configuration
rune config show | grep -A20 templates

# Test template processing
rune ritual test start [project]
```

#### Root Causes and Fixes

**Cause 1: Template reference mismatch**
```yaml
# ❌ Template name mismatch
- name: "Setup Environment"
  tmux_template: "dev-env"    # Template name
  
templates:
  development:                # Different name!
    session_name: "dev"

# ✅ Correct template reference
- name: "Setup Environment"
  tmux_template: "development"
  
templates:
  development:
    session_name: "dev"
```

**Cause 2: Invalid template structure**
```yaml
# ❌ Missing required fields
templates:
  dev:
    windows:              # Missing session_name!
      - name: "main"

# ✅ Complete template structure
templates:
  dev:
    session_name: "dev-{{.Project}}"
    windows:
      - name: "main"
        panes:
          - "echo 'ready'"
```

**Cause 3: Invalid layout names**
```yaml
# ❌ Invalid layout
- name: "editor"
  layout: "custom-layout"     # Not a valid tmux layout

# ✅ Valid layouts
- name: "editor"
  layout: "main-horizontal"   # Valid tmux layout
```

Valid layouts:
- `main-horizontal`
- `main-vertical` 
- `even-horizontal`
- `even-vertical`
- `tiled`

**Cause 4: Variable expansion issues**
```yaml
# ❌ Invalid variable syntax
session_name: "{Project}"        # Wrong syntax
session_name: "{{Project}}"      # Missing dot

# ✅ Correct variable syntax
session_name: "{{.Project}}"     # Correct
```

### Issue 4: PTY/Terminal Issues

#### Symptoms
- Terminal appears corrupted or garbled
- Key presses don't register correctly
- Size/layout issues in terminal

#### Diagnosis
```bash
# Check terminal environment
echo $TERM
echo $SHELL

# Test PTY directly
python3 -c "import pty; print('PTY support available')"

# Check terminal size
tput cols
tput lines
```

#### Root Causes and Fixes

**Cause 1: Terminal compatibility issues**
- Try different terminal emulator
- Set TERM environment variable:
  ```bash
  export TERM=xterm-256color
  rune start
  ```

**Cause 2: PTY size issues**
```bash
# In tmux, resize manually
Ctrl-b : resize-window -A

# Or detach and reattach
Ctrl-b d
tmux attach -t session-name
```

**Cause 3: Raw mode issues**
- Exit and restart terminal
- Check for conflicting terminal settings
- Verify terminal supports raw mode

### Issue 5: Permission and Security Errors

#### Symptoms
- "permission denied" when creating sessions
- Commands fail with security errors
- Environment variables not available

#### Diagnosis
```bash
# Check user permissions
id
groups

# Check environment filtering
rune config show | grep -A10 -B5 filter

# Test with minimal environment
env -i PATH=$PATH rune start
```

#### Root Causes and Fixes

**Cause 1: Insufficient permissions**
```bash
# Check tmux socket permissions
ls -la /tmp/tmux-$(id -u)/

# Fix permissions if needed
sudo chown -R $(whoami) /tmp/tmux-$(id -u)/
```

**Cause 2: Environment variable filtering**
- Rune filters sensitive environment variables by design
- This is expected security behavior
- Use config file for settings instead of environment variables

**Cause 3: Conflicting process permissions**
```bash
# Kill any conflicting tmux processes
pkill -u $(whoami) tmux
tmux kill-server

# Restart fresh
rune start
```

## Platform-Specific Issues

### macOS Issues

#### Issue: Homebrew tmux not in PATH
```bash
# Check if tmux installed via Homebrew
brew list | grep tmux

# Add Homebrew to PATH
echo 'export PATH=/opt/homebrew/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
```

#### Issue: Terminal.app compatibility
- Some advanced features may not work in Terminal.app
- Try iTerm2 for better compatibility
- Or use basic PTY mode instead of tmux

### Linux Issues

#### Issue: Old tmux version
```bash
# Check version
tmux -V

# Update if too old (need 2.0+)
sudo apt update && sudo apt upgrade tmux
```

#### Issue: Missing dependencies
```bash
# Install development dependencies
sudo apt install build-essential

# Or use package manager
sudo apt install tmux screen
```

### Windows/WSL Issues

#### Issue: tmux not working in Windows
- tmux requires Unix-like environment
- Install WSL2: Windows Subsystem for Linux
- Run rune inside WSL2, not native Windows

#### Issue: WSL terminal issues
```bash
# In WSL, install tmux
sudo apt update && sudo apt install tmux

# Set proper terminal
export TERM=screen-256color

# Or use Windows Terminal instead of command prompt
```

## Advanced Troubleshooting

### Debug Mode

```bash
# Enable debug logging
export RUNE_DEBUG=true
rune start

# Check logs
rune --log
```

### Verbose Output

```bash
# Run with verbose output
rune --verbose start

# Show detailed ritual test
rune --verbose ritual test start
```

### Manual Session Testing

```bash
# Test tmux manually
tmux new-session -d -s debug-session
tmux send-keys -t debug-session 'echo "test"' Enter
tmux list-sessions
tmux attach -t debug-session
```

### Configuration Debugging

```bash
# Show effective configuration
rune config show

# Test specific project
rune ritual test start myproject

# Validate against schema
rune config validate --verbose
```

## Performance Troubleshooting

### Slow Session Creation

**Symptoms**: Sessions take >5 seconds to create

**Diagnosis**:
```bash
# Time session creation
time rune start

# Check system load
top
htop
```

**Solutions**:
- Reduce template complexity
- Check for hanging commands in template
- Verify system resources available

### Memory Issues

**Symptoms**: High memory usage, system slowdown

**Diagnosis**:
```bash
# Check tmux session memory usage
ps aux | grep tmux

# List active sessions
tmux list-sessions
```

**Solutions**:
- Limit number of concurrent sessions
- Kill unused sessions: `tmux kill-session -t name`
- Use simpler templates with fewer panes

### I/O Performance Issues

**Symptoms**: Slow terminal response, input lag

**Solutions**:
- Reduce output volume in template commands
- Use background commands for high-output processes
- Check terminal emulator performance settings

## Recovery Procedures

### Complete Reset

If everything is broken, perform complete reset:

```bash
# 1. Kill all tmux processes
tmux kill-server
pkill -u $(whoami) tmux

# 2. Clear tmux socket directory
rm -rf /tmp/tmux-$(id -u)/

# 3. Reset rune configuration (backup first!)
cp ~/.rune/config.yaml ~/.rune/config.yaml.backup
rune config edit  # Fix configuration issues

# 4. Test basic functionality
rune config validate
rune ritual test start

# 5. Start fresh
rune start
```

### Session Recovery

If sessions are corrupted but working:

```bash
# List sessions
tmux list-sessions

# Kill problematic session
tmux kill-session -t session-name

# Start fresh ritual
rune start
```

### Configuration Recovery

If configuration is broken:

```bash
# Use example configuration
cp examples/config-developer.yaml ~/.rune/config.yaml

# Or start minimal
rune config edit
# Add minimal working configuration
```

## Getting Help

### Diagnostic Information

When reporting issues, include:

```bash
# System information
rune --version
tmux -V  # if using tmux
uname -a
echo $TERM

# Configuration (sanitized)
rune config show

# Error output with debug
export RUNE_DEBUG=true
rune start 2>&1 | tee error.log
```

### Common Support Channels

1. **Documentation**: Check user guide and API reference
2. **Issues**: GitHub issues for bug reports
3. **Discussions**: Community discussions for usage questions
4. **Examples**: Reference example configurations

### Self-Help Checklist

Before reporting issues:

- [ ] Checked this troubleshooting guide
- [ ] Validated configuration with `rune config validate`
- [ ] Tested with minimal example configuration
- [ ] Verified system requirements (tmux version, etc.)
- [ ] Tried with different terminal emulator
- [ ] Checked debug logs with `RUNE_DEBUG=true`

This troubleshooting guide should resolve most common issues with Interactive Rituals. For persistent problems, provide diagnostic information when seeking help.