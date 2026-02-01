# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Rune is a developer-first CLI productivity platform written in Go that automates daily work rituals, enforces healthy work-life boundaries, and integrates seamlessly with existing developer workflows. It's built using Cobra for CLI functionality, Viper for configuration management, with telemetry via OpenTelemetry (OTLP HTTP logs) and Sentry.

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
- **OpenTelemetry (OTLP logs)** for usage analytics

### Entry Point

- Main entry: `cmd/rune/main.go` → `internal/commands/Execute()`
- Commands are defined in `internal/commands/` with each command in its own file

## Development Commands

### Building
- `make build` - Build the binary
- `make dev` - Build with race detection for development
- `make build-telemetry` - Build with runtime telemetry support (no embedded keys)

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
- **OpenTelemetry (OTLP logs)** for usage analytics
- **Sentry** for error tracking
- **No keys embedded by default** - telemetry endpoints/keys must be provided at runtime
- Keys loaded from:
  1. Environment variables: `RUNE_OTLP_ENDPOINT`, `RUNE_SENTRY_DSN`
  2. Config file: `~/.rune/config.yaml` under `integrations.telemetry`
- Telemetry can be disabled via `RUNE_TELEMETRY_DISABLED=true`
- Debug mode: `RUNE_DEBUG=true`
- **Beta users get clean experience** - no telemetry sent without explicit key configuration

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

## Recent Learnings

### 2026-02-01: Interactive Environment Management
- **Learning**: Rune successfully implemented sophisticated interactive tmux automation (bf724c87) - moving beyond basic ritual commands to full development environment orchestration with PTY integration and session persistence. This represents significant maturity in the ritual engine architecture.
- **Pattern**: Feature verification is critical before documentation - recent commits (5948e322, b604afb2, 073e31df) discovered project is ~85% ready for release when actual code review contradicted initial TODOs. Always verify implementation status against actual codebase.
- **Gotcha**: Documentation debt can obscure true project state - 4 consecutive docs commits were needed to establish correct understanding of what's actually implemented vs what's planned.

### 2026-02-01: Test Coverage Gap Identified
- **Learning**: Current test coverage at 62.9% (e71e5057) with specific gaps in telemetry masking and command integration testing. Coverage targets set at 80% overall, 85% for critical modules - clear path defined in a5bf9dd6 and 073e31df.
- **Pattern**: Test improvement should follow strategic plan (a5bf9dd6 spec) rather than ad-hoc additions - spec defines 4-phase approach with user stories, acceptance criteria, and effort estimates for each phase.

### 2026-02-01: Release Readiness Criteria
- **Learning**: Project MVP is complete with all core features functional. Release blockers are quality-focused: test coverage gaps, documentation site deployment, and linting tool setup. Version management critical - beta.7 designation (26aa5bac) maintains clear pre-release status.
- **Pattern**: Quality gates should precede documentation vs ship-now decisions - ensures users interact with tested, well-documented system rather than discovering gaps post-release.