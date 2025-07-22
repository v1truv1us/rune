---
layout: ../../layouts/BaseLayout.astro
title: "Configuration Guide - Rune CLI"
description: "Complete guide to configuring Rune CLI"
---

# Configuration Guide

Rune uses a YAML configuration file located at `~/.rune/config.yaml`. This file controls all aspects of Rune's behavior.

## Configuration Structure

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m
  timezone: "America/New_York"

projects:
  - name: "main-app"
    detect: ["git:main-app", "dir:~/projects/main-app"]

rituals:
  start:
    global: []
    per_project: {}
  stop:
    global: []
    per_project: {}

integrations:
  git:
    enabled: true
  slack:
    enabled: false
  telemetry:
    enabled: true
```

## Settings Section

### Basic Settings

```yaml
settings:
  work_hours: 8.0 # Target work hours per day
  break_interval: 50m # Suggested break interval
  idle_threshold: 10m # Idle time before auto-pause
  timezone: "America/New_York" # Your timezone
  auto_start: false # Auto-start on first command
  auto_stop: true # Auto-stop at day end
```

### Notification Settings

```yaml
settings:
  notifications:
    enabled: true
    break_reminders: true
    session_alerts: true
    sound: true
```

### Focus Settings

```yaml
settings:
  focus:
    dnd_on_start: true # Enable Do Not Disturb on start
    dnd_on_break: false # Keep DND during breaks
    block_websites: [] # Websites to block during work
    allow_emergency: true # Allow emergency interruptions
```

## Projects Section

Projects define how Rune detects and categorizes your work.

### Project Detection

```yaml
projects:
  - name: "web-app"
    detect:
      - "git:web-app" # Git repository name
      - "dir:~/projects/web-app" # Directory path
      - "git-remote:github.com/me/app" # Git remote URL
    description: "Main web application"
    tags: ["frontend", "react"]

  - name: "api-service"
    detect:
      - "git:api-service"
      - "dir:~/work/api"
    description: "Backend API service"
    tags: ["backend", "go"]
```

### Detection Methods

- `git:name` - Match Git repository name
- `dir:path` - Match directory path (supports `~` expansion)
- `git-remote:url` - Match Git remote URL pattern
- `env:VAR=value` - Match environment variable

### Project Metadata

```yaml
projects:
  - name: "project-name"
    detect: ["git:project"]
    description: "Project description"
    tags: ["tag1", "tag2"]
    client: "Client Name"
    rate: 150.00 # Hourly rate for reporting
    budget: 40.0 # Hour budget
    deadline: "2024-12-31" # Project deadline
```

## Rituals Section

Rituals are automated command sequences that run during start/stop/break events.

### Basic Ritual Structure

```yaml
rituals:
  start:
    global:
      - name: "Update repositories"
        command: "git -C ~/projects pull --all"
        optional: false
        timeout: 30s

    per_project:
      web-app:
        - name: "Start dev server"
          command: "cd ~/projects/web-app && npm run dev"
          background: true

  stop:
    global:
      - name: "Commit work in progress"
        command: "git add -A && git commit -m 'WIP: End of day'"
        optional: true
```

### Ritual Command Options

```yaml
- name: "Command description"
  command: "shell command to execute"
  optional: true # Don't fail ritual if command fails
  background: false # Run in background
  timeout: 30s # Command timeout
  working_dir: "~/projects" # Working directory
  env: # Environment variables
    NODE_ENV: "development"
  when: # Conditional execution
    - "git_clean" # Only if git is clean
    - "weekday" # Only on weekdays
    - "time_after:09:00" # Only after 9 AM
```

### Conditional Execution

Available conditions:

- `git_clean` - Git working directory is clean
- `git_dirty` - Git working directory has changes
- `weekday` - Monday through Friday
- `weekend` - Saturday and Sunday
- `time_after:HH:MM` - After specific time
- `time_before:HH:MM` - Before specific time
- `env:VAR=value` - Environment variable matches

### Ritual Types

```yaml
rituals:
  start: # Run when starting work
  stop: # Run when stopping work
  break: # Run when taking breaks
  resume: # Run when resuming from break
  daily: # Run once per day
  weekly: # Run once per week
```

## Integrations Section

### Git Integration

```yaml
integrations:
  git:
    enabled: true
    auto_detect_project: true # Auto-detect project from Git
    commit_on_stop: false # Auto-commit on stop
    push_on_stop: false # Auto-push on stop
    branch_in_status: true # Show branch in status
```

### Slack Integration

```yaml
integrations:
  slack:
    enabled: true
    workspace: "myteam"
    token: "xoxp-your-token" # Use environment variable instead
    dnd_on_start: true
    status_on_start: "🔨 Working on {{project}}"
    status_on_break: "☕ On break"
    status_on_stop: ""
```

### Calendar Integration

```yaml
integrations:
  calendar:
    enabled: true
    provider: "google" # google, outlook, caldav
    block_calendar: true # Block time on calendar
    meeting_detection: true # Detect meetings for auto-pause
    credentials_file: "~/.rune/calendar-creds.json"
```

### Telemetry Settings

```yaml
integrations:
  telemetry:
    enabled: true
    segment_write_key: "" # Set via environment
    sentry_dsn: "" # Set via environment
    collect_errors: true
    collect_usage: true
    collect_performance: false
```

## Environment Variables

Sensitive values should be set via environment variables:

```bash
# Telemetry
export RUNE_SEGMENT_WRITE_KEY="your-key"
export RUNE_SENTRY_DSN="your-dsn"

# Integrations
export RUNE_SLACK_TOKEN="xoxp-your-token"
export RUNE_CALENDAR_CREDENTIALS="path/to/creds.json"

# General
export RUNE_CONFIG_FILE="~/.rune/config.yaml"
export RUNE_TELEMETRY_DISABLED="true"
export RUNE_DEBUG="true"
```

## Configuration Templates

### Developer Template

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 15m

projects:
  - name: "current-project"
    detect: ["git:{{git_repo_name}}"]

rituals:
  start:
    global:
      - name: "Update dependencies"
        command: "git pull && npm install"
      - name: "Start services"
        command: "docker-compose up -d"
        optional: true

  stop:
    global:
      - name: "Commit WIP"
        command: "git add -A && git commit -m 'WIP: $(date)'"
        optional: true

integrations:
  git:
    enabled: true
    auto_detect_project: true
```

### Freelancer Template

```yaml
version: 1
settings:
  work_hours: 6.0
  break_interval: 45m

projects:
  - name: "client-a"
    detect: ["dir:~/clients/client-a"]
    client: "Client A"
    rate: 150.00

  - name: "client-b"
    detect: ["dir:~/clients/client-b"]
    client: "Client B"
    rate: 175.00

rituals:
  start:
    global:
      - name: "Time tracking reminder"
        command: "echo 'Remember to track time accurately'"

  stop:
    global:
      - name: "Generate invoice data"
        command: "rune report --today --format csv >> ~/invoices/$(date +%Y-%m).csv"
```

## Validation

Validate your configuration:

```bash
rune config validate
```

Common validation errors:

- Invalid `YAML` syntax
- Unknown configuration keys
- Invalid duration formats
- Missing required fields
- Invalid regex patterns

## Migration

### From `Watson`

```bash
rune config migrate --from watson --backup
```

### From `Timewarrior`

```bash
rune config migrate --from timewarrior --backup
```

## Best Practices

1. **Version Control**: Keep your config in a dotfiles repository
2. **Environment Variables**: Use env vars for sensitive data
3. **Backup**: Regularly backup your configuration
4. **Testing**: Test rituals with `--dry-run` flag
5. **Documentation**: Document custom rituals and integrations
