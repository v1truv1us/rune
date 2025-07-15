---
layout: ../layouts/BaseLayout.astro
title: "Getting Started - Rune CLI"
description: "Get up and running with Rune CLI in minutes"
---

# Getting Started

## Installation

### Homebrew (Recommended for macOS/Linux)

```bash
brew install --cask ferg-cod3s/tap/rune
```

### Quick Install Script

```bash
# With Homebrew suggestion
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh

# Skip Homebrew suggestion
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --skip-homebrew
```

### Go Install

```bash
go install github.com/ferg-cod3s/rune/cmd/rune@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/ferg-cod3s/rune/releases) and place it in your PATH.

### Package Managers

- **Debian/Ubuntu**: Download `.deb` from releases
- **RHEL/CentOS**: Download `.rpm` from releases  
- **Arch Linux**: Available in AUR (coming soon)

## Initial Setup

### 1. Initialize Configuration

```bash
rune init --guided
```

This interactive setup will help you configure:
- Work hours and break intervals
- Project detection rules
- Start/stop rituals
- Integration preferences

### 2. Verify Installation

```bash
rune --version
rune status
```

## Basic Usage

### Daily Workflow

```bash
# Start your workday
rune start

# Check current status  
rune status

# Pause for a break
rune pause

# Resume work
rune resume

# End your workday
rune stop
```

### Time Tracking

```bash
# View today's report
rune report --today

# View this week's report
rune report --week

# View specific date range
rune report --from 2024-01-01 --to 2024-01-07
```

### Configuration Management

```bash
# Edit configuration
rune config edit

# Validate configuration
rune config validate

# Show current configuration
rune config show
```

## Next Steps

- [Configure your first rituals](/docs/rituals)
- [Set up project detection](/docs/projects)  
- [Explore command reference](/docs/commands)
- [See example workflows](/examples)

## Getting Help

- Run `rune --help` for command overview
- Run `rune <command> --help` for specific command help
- Visit our [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions)
- Report issues on [GitHub Issues](https://github.com/ferg-cod3s/rune/issues)
EOF < /dev/null