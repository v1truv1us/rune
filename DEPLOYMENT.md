# Rune CLI Deployment Guide

This document outlines how to make Rune installable and track its usage.

## Installation Methods

### 1. Go Install (Immediate)
Users can install directly from GitHub:
```bash
go install github.com/ferg-cod3s/rune/cmd/rune@latest
```

### 2. Install Script (Recommended)
Quick installation via curl:
```bash
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh
```

### 3. GitHub Releases (Automated)
- Binary releases for Linux, macOS, Windows
- Debian/RPM packages
- Homebrew formula (requires separate tap repository)

### 4. Package Managers
- **Homebrew**: `brew tap ferg-cod3s/tap && brew install --cask rune`
- **Debian/Ubuntu**: Download `.deb` from releases
- **RHEL/CentOS**: Download `.rpm` from releases

## Release Process

### Automated Releases
1. Tag a new version: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. GitHub Actions automatically:
   - Runs tests
   - Builds binaries for all platforms
   - Creates GitHub release
   - Generates packages (.deb, .rpm)
   - Updates Homebrew formula

### Manual Release
```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Create release
goreleaser release --clean
```

## Telemetry & Analytics

### What's Tracked
- Command usage patterns (which commands are run)
- Error occurrences (to identify bugs)
- Performance metrics (command execution times)
- System information (OS, architecture)

**No personal data, file contents, or command arguments are collected.**

### Implementation
- Anonymous user IDs stored in `~/.rune/config.yaml`
- Telemetry client in `internal/telemetry/`
- Middleware wraps commands automatically
- Graceful failure - telemetry errors don't affect CLI functionality

### Configuration
Configure telemetry endpoints:
```bash
# OpenTelemetry OTLP logs endpoint (optional)
export RUNE_OTLP_ENDPOINT="http://localhost:4318/v1/logs"
# Sentry DSN (optional)
export RUNE_SENTRY_DSN="https://public@sentry.io/project"
```

### Disable Telemetry
Users can opt out:
```bash
export RUNE_TELEMETRY_DISABLED=true
```

## Setting Up for Production

### 1. Create GitHub Repository
- Push code to GitHub
- Set up repository secrets if needed

### 2. Set Up Telemetry
- Set up an OpenTelemetry Collector (or use a managed observability platform) for logs
- Configure `RUNE_OTLP_ENDPOINT` to point to your Collector's logs endpoint
- Optionally configure Sentry DSN via `RUNE_SENTRY_DSN` for errors and performance

### 3. Create Homebrew Tap (Optional)
```bash
# Create a new repository: homebrew-tap
# GoReleaser will automatically update it
```

### 4. Set Up Domain (Optional)
- Point `get.rune.dev` to your install script
- Set up documentation at `docs.rune.dev`

## Monitoring Usage

### Key Metrics to Track
- **Installation Events**: Track when users install
- **Command Usage**: Most/least used commands
- **Error Rates**: Commands that fail frequently
- **User Retention**: Daily/weekly active users
- **Performance**: Command execution times

### Example Analytics Queries
```sql
-- Most popular commands
SELECT event_properties.command, COUNT(*) as usage_count
FROM events 
WHERE event = 'command_executed'
GROUP BY event_properties.command
ORDER BY usage_count DESC;

-- Error rates by command
SELECT 
  event_properties.command,
  COUNT(*) as total_executions,
  SUM(CASE WHEN event_properties.success = false THEN 1 ELSE 0 END) as errors,
  (SUM(CASE WHEN event_properties.success = false THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) as error_rate
FROM events 
WHERE event = 'command_executed'
GROUP BY event_properties.command;
```

## Security Considerations

### Telemetry
- All telemetry is anonymous
- No sensitive data is transmitted
- Users can easily opt out
- Telemetry failures don't affect CLI functionality

### Distribution
- Binaries are signed (configure in GoReleaser)
- Checksums provided for all releases
- HTTPS-only distribution
- Package manager verification

## Next Steps

1. **Push to GitHub**: Make repository public
2. **Set up telemetry**: Configure your analytics service
3. **Create first release**: Tag v1.0.0 and push
4. **Test installation**: Verify install script works
5. **Monitor usage**: Set up dashboards for key metrics

## Files Created

- `.github/workflows/release.yml` - Automated releases
- `.github/workflows/test.yml` - CI testing
- `.goreleaser.yaml` - Release configuration
- `install.sh` - Installation script
- `internal/telemetry/` - Telemetry system
- Updated `README.md` - Installation instructions
- Updated `internal/commands/start.go` - Example telemetry integration

Your CLI is now ready for distribution and usage tracking!