# Rune CLI

```
 ______     __  __     __   __     ______   
/\  == \   /\ \/\ \   /\ "-.\ \   /\  ___\  
\ \  __<   \ \ \_\ \  \ \ \-.  \  \ \  __\  
 \ \_\ \_\  \ \_____\  \ \_\\"\_\  \ \_____\
  \/_/ /_/   \/_____/   \/_/ \/_/   \/_____/ 
```

**Ancient wisdom for modern workflows**

> ⚠️ **Alpha Release**: Rune is currently in alpha development. Features may change, and some functionality is still being implemented. Use with caution in production environments.

Rune is a developer-first CLI productivity platform that automates daily work rituals, enforces healthy work-life boundaries, and integrates seamlessly with existing developer workflows.

## Features

- **Ritual Automation**: Automate your daily start/stop workflows with custom commands
- **Time Tracking**: Intelligent time tracking with Git integration and project detection
- **Focus Protection**: OS-level Do Not Disturb automation and break reminders
- **Cross-Platform**: Works on macOS, Linux, and Windows with shell completions
- **YAML Configuration**: Simple, version-controlled configuration management
- **Security-First**: Sandboxed command execution with audit logging

## Quick Start

### Installation

**Homebrew (Recommended for macOS/Linux)**
```bash
brew install --cask ferg-cod3s/tap/rune
```

**Quick Install Script (Linux/macOS)**
```bash
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh
```

**Skip Homebrew suggestion in install script**
```bash
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --skip-homebrew
```

**Go Install**
```bash
go install github.com/ferg-cod3s/rune/cmd/rune@latest
```

**Download Binary**
Download the latest release from [GitHub Releases](https://github.com/ferg-cod3s/rune/releases) and place it in your PATH.

**Package Managers**
- **Debian/Ubuntu**: Download `.deb` from releases
- **RHEL/CentOS**: Download `.rpm` from releases
- **Arch Linux**: Available in AUR (coming soon)

### Initialize Your Configuration

```bash
rune init --guided
```

This will walk you through setting up your first rituals and work preferences.

### Basic Usage

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

# View time reports
rune report --today

# Update to latest version
rune update

# Check for updates without installing
rune update --check

# Update via Homebrew (if installed via brew)
brew upgrade rune
```

## Configuration

Rune uses a YAML configuration file at `~/.rune/config.yaml`:

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m
  
projects:
  - name: "main-app"
    detect: ["git:main-app", "dir:~/projects/main-app"]
    
rituals:
  start:
    global:
      - name: "Update repositories"
        command: "git -C ~/projects pull --all"
      - name: "Start Docker"
        command: "docker-compose up -d"
    per_project:
      main-app:
        - name: "Start dev server"
          command: "bun run dev"
          
  stop:
    global:
      - name: "Commit changes"
        command: "git add -A && git commit -m 'WIP: End of day'"
        optional: true
      - name: "Stop services"
        command: "docker-compose down"
        
integrations:
  git:
    enabled: true
    auto_detect_project: true
  slack:
    workspace: "myteam"
    dnd_on_start: true
  calendar:
    provider: "google"
    block_calendar: true
```

## Commands

### Core Commands

- `rune init` - Initialize configuration with guided setup
- `rune start` - Start workday and run start rituals
- `rune pause` - Pause current timer
- `rune resume` - Resume paused timer
- `rune status` - Show current session status
- `rune stop` - End workday and run stop rituals
- `rune report` - Generate time reports
- `rune update` - Update rune to the latest version

### Configuration Commands

- `rune config edit` - Edit configuration file
- `rune config validate` - Validate configuration
- `rune config migrate` - Migrate from Watson/Timewarrior

### Ritual Commands

- `rune ritual list` - List available rituals
- `rune ritual run <name>` - Run specific ritual
- `rune ritual test <name>` - Test ritual without execution

## Examples

### Frontend Developer Workflow

```yaml
rituals:
  start:
    global:
      - name: "Pull latest changes"
        command: "git pull"
      - name: "Install dependencies"
        command: "bun install"
      - name: "Start dev server"
        command: "bun run dev"
        background: true
      - name: "Open browser"
        command: "open http://localhost:3000"
        
  stop:
    global:
      - name: "Stop dev server"
        command: "pkill -f 'bun run dev'"
      - name: "Commit WIP"
        command: "git add -A && git commit -m 'WIP: $(date)'"
        optional: true
```

### DevOps Engineer Workflow

```yaml
rituals:
  start:
    global:
      - name: "Check cluster status"
        command: "kubectl get nodes"
      - name: "Update monitoring dashboard"
        command: "open https://grafana.company.com"
      - name: "Check CI/CD pipeline"
        command: "gh workflow list --repo company/main-app"
        
  stop:
    global:
      - name: "Generate daily report"
        command: "./scripts/daily-report.sh"
      - name: "Update team status"
        command: "slack-cli post '#devops' 'EOD: All systems green'"
```

## Development

### Prerequisites

- Go 1.21+
- Make
- Git

### Building from Source

```bash
git clone https://github.com/ferg-cod3s/rune.git
cd rune
make build
```

### Running Tests

```bash
make test
```

### Development Commands

```bash
make dev          # Run in development mode
make lint         # Run linter
make fmt          # Format code
make test-watch   # Run tests in watch mode
```

## Privacy & Telemetry

Rune collects anonymous usage analytics to help improve the tool. This includes:

- Command usage patterns (which commands are run)
- Error occurrences (to identify bugs)
- Performance metrics (command execution times)
- System information (OS, architecture)
- OpenTelemetry logs (if RUNE_OTLP_ENDPOINT is configured)

**No personal data, file contents, or command arguments are collected.**

### Telemetry Configuration

For telemetry to work, you can configure Sentry and/or an OTLP logs endpoint using environment variables:

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your keys:
# RUNE_OTLP_ENDPOINT=http://localhost:4318/v1/logs
# RUNE_SENTRY_DSN=https://your_sentry_dsn@sentry.io/project_id
```

**Important**: Never commit secrets to version control. Always use environment variables or the `.env` file.

### Disable Telemetry

You can disable telemetry in several ways:

```bash
# Environment variable (recommended)
export RUNE_TELEMETRY_DISABLED=true

# Or add to your shell profile
echo 'export RUNE_TELEMETRY_DISABLED=true' >> ~/.bashrc

# Or set in config file
# integrations:
#   telemetry:
#     enabled: false
```

## Security

Rune takes security seriously:

- All commands run with user privileges only
- No shell expansion without explicit consent
- Credentials stored in OS keychain
- Command audit logging available
- Optional sandboxing support

See [SECURITY.md](SECURITY.md) for detailed security information.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run tests: `make test`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [docs.rune.dev](https://docs.rune.dev)
- **Issues**: [GitHub Issues](https://github.com/ferg-cod3s/rune/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions)
- **Discord**: [Rune Community](https://discord.gg/rune)

## Roadmap

See [TODO.md](TODO.md) for current development priorities and [PRD.md](PRD.md) for the complete product roadmap.

## Acknowledgments

- Inspired by Watson, Timewarrior, and dijo
- Built with Go, Cobra, and Viper
- Logo design by [Designer Name]

---

**Cast your daily runes and master your workflow** ✨

## Codeflow Workflow - Claude Code

This project is set up for Claude Code with native slash commands.

### Available Commands

- `/research` - Comprehensive codebase and documentation analysis
- `/plan` - Create detailed implementation plans
- `/execute` - Implement plans with verification
- `/test` - Generate comprehensive test suites
- `/document` - Create user guides and API documentation
- `/commit` - Create structured git commits
- `/review` - Validate implementations against plans

### Usage

Simply use the slash commands directly in Claude Code:

```
/research "Analyze the authentication system for potential OAuth integration"

/plan "Create implementation plan based on the research findings"

/execute "Implement the OAuth integration following the plan"
```

Commands are located in `.claude/commands/` and can be customized for this project.
