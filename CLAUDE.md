# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Rune is a developer-first CLI productivity platform written in Go that automates daily work rituals, enforces healthy work-life boundaries, and integrates seamlessly with existing developer workflows. It's built using Cobra for CLI functionality, Viper for configuration management, and includes telemetry via Segment and Sentry.

## Architecture

### Core Components

- **Commands** (`internal/commands/`): CLI command implementations using Cobra framework
  - Root command with global configuration and telemetry initialization
  - Subcommands: start, stop, pause, resume, status, report, ritual, config, init, update
- **Configuration** (`internal/config/`): YAML-based configuration management using Viper
- **Tracking** (`internal/tracking/`): Time tracking, session management, and project detection
- **Rituals** (`internal/rituals/`): Automation engine for executing custom commands/workflows
- **Telemetry** (`internal/telemetry/`): Analytics and error reporting integration
- **Notifications** (`internal/notifications/`): System notification handling
- **DND** (`internal/dnd/`): Do Not Disturb functionality for focus management

### Key Technologies

- **Go 1.23+** with toolchain 1.24.5
- **Cobra** for CLI framework
- **Viper** for configuration management
- **BBolt** for local database storage
- **Sentry** for error tracking
- **Segment** for analytics

### Entry Point

- Main entry: `cmd/rune/main.go` → `internal/commands/Execute()`
- Commands are defined in `internal/commands/` with each command in its own file

## Development Commands

### Building
- `make build` - Build the binary
- `make dev` - Build with race detection for development
- `make build-telemetry` - Build with telemetry support (embeds API keys)

### Testing
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make test-coverage-detailed` - Detailed coverage with 70% threshold
- `make test-watch` - Run tests in watch mode (requires `entr`)

### Code Quality
- `make lint` - Run golangci-lint
- `make fmt` - Format code with go fmt and gofmt
- `make vet` - Run go vet
- `make pre-commit` - Run fmt, vet, lint, and test

### Security
- `make security` - Basic vulnerability check with govulncheck
- `make security-all` - Comprehensive security checks (deps, vulns, static analysis, secrets)
- `make security-build` - Check built binary for embedded secrets

### Running
- `make run` - Run the application (use `ARGS="..."` for arguments)
- `./bin/rune` - Run built binary directly

### Release Validation
- `goreleaser check` - Validate .goreleaser.yaml configuration
- `goreleaser build --snapshot --clean` - Test build without release
- **IMPORTANT**: Always run `goreleaser check` before pushing new tags to ensure release workflow will succeed

## Configuration

Rune uses YAML configuration at `~/.rune/config.yaml` with the following structure:
- **settings**: work hours, break intervals, idle thresholds
- **projects**: project detection rules (git repos, directories)
- **rituals**: start/stop automation commands (global and per-project)
- **integrations**: git, slack, calendar, telemetry settings

## Testing Approach

- Unit tests alongside source files (`*_test.go`)
- Integration tests in telemetry package
- Benchmark tests for performance-critical components
- Coverage reporting with HTML output
- Test utilities in `internal/commands/utils.go`

## Telemetry Integration

Rune includes optional telemetry for usage analytics and error reporting:
- **Segment** for usage analytics (embedded key: starts with `ZkEZXHRWH96y8EviNkbYJUByqGR9QI4G`)
- **Sentry** for error tracking (DSN: `https://3b20acb23bbbc5958448bb41900cdca2@sentry.fergify.work/10`)
- Telemetry can be disabled via `RUNE_TELEMETRY_DISABLED=true`
- Debug mode: `RUNE_DEBUG=true`

## Code Signing

Releases are signed using cosign for verification and security:
- **cosign** signs checksums and binaries using keyless signing with GitHub OIDC
- Requires `COSIGN_EXPERIMENTAL=1` environment variable in CI
- Users can verify signatures with: `cosign verify-blob --certificate <artifact>.pem --signature <artifact>.sig <artifact>`
- No additional secrets required - uses GitHub's OIDC token for signing

## Important Files

- `Makefile` - Comprehensive build and development commands
- `go.mod` - Go module dependencies
- `internal/commands/root.go` - Main CLI setup and initialization
- `internal/config/config.go` - Configuration structure and loading
- `internal/tracking/session.go` - Core time tracking logic
- `internal/rituals/engine.go` - Automation execution engine