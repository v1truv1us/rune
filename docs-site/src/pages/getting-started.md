---
layout: ../layouts/BaseLayout.astro
title: "Getting Started - Rune CLI"
description: "Get up and running with Rune CLI in minutes"
---

# Getting Started with Rune CLI Beta

> **🚀 Beta Status**: You're installing v0.9.0-beta.1 - the MVP feature-complete version. [Join the beta program](/beta) for updates and feedback channels.

## Installation

### Quick Install Script (Recommended)

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash
```

### Homebrew (macOS)

```bash
brew install --cask ferg-cod3s/tap/rune
```

### Go Install

```bash
go install github.com/ferg-cod3s/rune/cmd/rune@latest
```

### Manual Installation

1. Download the latest beta from [GitHub Releases](https://github.com/ferg-cod3s/rune/releases/tag/v0.9.0-beta.1)
2. Extract the binary and move to your `PATH`
3. Make executable: `chmod +x rune`

### Windows Installation

Download the Windows binary from [GitHub Releases](https://github.com/ferg-cod3s/rune/releases/tag/v0.9.0-beta.1) and add to your `PATH`.

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

## Beta Feedback

As a beta user, your feedback is invaluable! Please share:

- **GitHub Discussions**: [General feedback and questions](https://github.com/ferg-cod3s/rune/discussions)
- **GitHub Issues**: [Bug reports and feature requests](https://github.com/ferg-cod3s/rune/issues)
- **Email**: `beta@rune.dev` for private feedback

## Getting Help

- Run `rune --help` for command overview
- Run `rune <command> --help` for specific command help
- Check the [Beta Program page](/beta) for known issues
- Visit [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions) for community support
