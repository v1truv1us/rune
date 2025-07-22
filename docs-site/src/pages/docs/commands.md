---
layout: ../../layouts/BaseLayout.astro
title: "Commands Reference - Rune CLI"
description: "Complete reference for all Rune CLI commands"
---

# Commands Reference

## Core Commands

### `rune start`

Start your workday and execute start rituals.

```bash
rune start [flags]
```

**Flags:**

- `--project <name>` - Start with specific project
- `--skip-rituals` - Skip start rituals
- `--dry-run` - Show what would be executed without running

**Examples:**

```bash
rune start
rune start --project my-app
rune start --skip-rituals
```

### `rune stop`

End your workday and execute stop rituals.

```bash
rune stop [flags]
```

**Flags:**

- `--skip-rituals` - Skip stop rituals
- `--force` - Force stop even if rituals fail
- `--dry-run` - Show what would be executed

**Examples:**

```bash
rune stop
rune stop --skip-rituals
rune stop --force
```

### `rune pause`

Pause the current work session.

```bash
rune pause [duration]
```

**Arguments:**

- `duration` - Optional pause duration (e.g., 15m, 1h)

**Examples:**

```bash
rune pause           # Indefinite pause
rune pause 15m       # Pause for 15 minutes
rune pause 1h30m     # Pause for 1 hour 30 minutes
```

### `rune resume`

Resume a paused work session.

```bash
rune resume
```

### `rune status`

Show current session status and statistics.

```bash
rune status [flags]
```

**Flags:**

- `--verbose` - Show detailed information
- `--json` - Output in JSON format

**Examples:**

```bash
rune status
rune status --verbose
rune status --json
```

## Time Tracking Commands

### `rune report`

Generate time tracking reports.

```bash
rune report [flags]
```

**Flags:**

- `--today` - Show today's report
- `--week` - Show this week's report
- `--month` - Show this month's report
- `--from <date>` - Start date (YYYY-MM-DD)
- `--to <date>` - End date (YYYY-MM-DD)
- `--project <name>` - Filter by project
- `--format <format>` - Output format (table, json, csv)

**Examples:**

```bash
rune report --today
rune report --week
rune report --from 2024-01-01 --to 2024-01-31
rune report --project my-app --format json
```

## Configuration Commands

### `rune config`

Manage configuration settings.

#### `rune config edit`

Open configuration file in default editor.

```bash
rune config edit
```

#### `rune config validate`

Validate configuration file syntax.

```bash
rune config validate [file]
```

#### `rune config show`

Display current configuration.

```bash
rune config show [flags]
```

**Flags:**

- `--section <name>` - Show specific section
- `--json` - Output in JSON format

#### `rune config migrate`

Migrate from other time tracking tools.

```bash
rune config migrate [flags]
```

**Flags:**

- `--from <tool>` - Source tool (watson, timewarrior)
- `--backup` - Create backup before migration

## Ritual Commands

### `rune ritual`

Manage and execute rituals.

#### `rune ritual list`

List all available rituals.

```bash
rune ritual list [flags]
```

**Flags:**

- `--type <type>` - Filter by type (start, stop, break)
- `--project <name>` - Filter by project

#### `rune ritual run`

Execute a specific ritual.

```bash
rune ritual run <name> [flags]
```

**Flags:**

- `--dry-run` - Show commands without executing
- `--verbose` - Show detailed output

#### `rune ritual test`

Test ritual configuration without execution.

```bash
rune ritual test <name>
```

## Utility Commands

### `rune init`

Initialize Rune configuration.

```bash
rune init [flags]
```

**Flags:**

- `--guided` - Interactive setup wizard
- `--template <name>` - Use configuration template
- `--force` - Overwrite existing configuration

**Templates:**

- `basic` - Minimal configuration
- `developer` - Developer workflow preset
- `freelancer` - Freelancer workflow preset

### `rune update`

Update Rune to the latest version.

```bash
rune update [flags]
```

**Flags:**

- `--check` - Check for updates without installing
- `--beta` - Include beta releases
- `--force` - Force update even if already latest

### `rune completion`

Generate shell completion scripts.

```bash
rune completion <shell>
```

**Supported shells:**

- `bash`
- `zsh`
- `fish`
- `powershell`

**Examples:**

```bash
# Bash
rune completion bash > /etc/bash_completion.d/rune

# Zsh
rune completion zsh > "${fpath[1]}/_rune"

# Fish
rune completion fish > ~/.config/fish/completions/rune.fish
```

## Global Flags

These flags are available for all commands:

- `--config <file>` - Specify config file (default: ~/.rune/config.yaml)
- `--verbose, -v` - Enable verbose output
- `--help, -h` - Show help for command
- `--version` - Show version information

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Configuration error
- `3` - Runtime error
- `4` - User interruption (Ctrl+C)
