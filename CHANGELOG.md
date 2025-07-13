# Changelog

All notable changes to the Rune CLI project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Nothing yet

### Changed
- Nothing yet

### Deprecated
- Nothing yet

### Removed
- Nothing yet

### Fixed
- Nothing yet

### Security
- Nothing yet

## [0.1.0-alpha.4] - 2025-01-09

### Added
- **Homebrew Tap**: Created homebrew-tap repository for proper Homebrew distribution
- **Homebrew Formula**: Enabled Homebrew formula generation in release workflow

### Fixed
- **Release Workflow**: Fixed Homebrew tap integration to enable `brew install --cask ferg-cod3s/tap/rune`

## [0.1.0-alpha.3] - 2025-01-09

### Added
- **Comprehensive Test Coverage**: Added test files for commands, dnd, rituals, and telemetry modules
- **Update Command**: Implemented self-update functionality with version checking and binary replacement
- **Build Automation**: Added Sentry release automation script
- **Uninstall Script**: Added clean removal script for complete uninstallation
- **Enhanced Installation**: Improved install.sh with better error handling and validation

### Changed
- **Build Configuration**: Updated GoReleaser and GitHub Actions workflow
- **Documentation**: Enhanced README with updated installation and usage instructions
- **Telemetry**: Improved telemetry middleware and functionality

### Fixed
- **Build Process**: Resolved build configuration issues and ensured all tests pass
- **Module Dependencies**: Updated Go module dependencies for better compatibility

## [0.1.0-alpha.1] - 2025-01-09

### Added
- **Core CLI Structure**: Complete command structure with Cobra framework
- **Dual Logo System**: 
  - SubZero ASCII art for daily CLI use (highly readable)
  - Runic ceremonial logo for special initialization moments
- **Security-First Configuration**: Environment variable support for API keys
- **Comprehensive Documentation**: Installation guides, quick start, and API docs
- **Cross-Platform Support**: macOS, Linux, Windows binaries via GoReleaser
- **Shell Completions**: Bash, Zsh, Fish completion scripts
- **Commands Implemented**:
  - `rune init --guided` - Interactive configuration setup with ceremonial logo
  - `rune --help` - Command reference with SubZero logo
  - `rune --version` - Version display with branding
  - `rune config` - Configuration management (framework)
  - `rune start/stop/pause/resume` - Time tracking (framework)
  - `rune status` - Session status display (framework)
  - `rune report` - Time reporting (framework)
  - `rune ritual` - Ritual management (framework)

### Security
- Environment variable configuration for sensitive data (.env support)
- No hardcoded secrets in codebase
- Secure .gitignore preventing credential commits
- Template .env.example for safe local development

### Documentation
- Complete installation guide for all platforms
- Quick start tutorial with examples
- Command reference documentation
- Security best practices guide
- Structured docs/ directory ready for docs site

### Infrastructure
- GitHub Actions CI/CD pipeline
- GoReleaser configuration for multi-platform builds
- Automated testing on push and PR
- Release workflow with binary artifacts

### Notes
- **Alpha Release**: Core functionality is framework-only
- Time tracking, ritual execution, and integrations not yet implemented
- Configuration parsing and validation are in place
- Focus on CLI experience, branding, and developer setup

## [0.1.0] - TBD (Planned Beta)

### Added
- Basic CLI structure with Cobra framework
- YAML configuration management with Viper
- Time tracking core functionality
- Ritual automation engine
- Cross-platform Do Not Disturb integration
- Shell completions for Bash, Zsh, Fish, PowerShell
- Basic reporting system
- Git integration for project detection

### Security
- Command sandboxing implementation
- Secure credential storage using OS keychain
- Input validation and injection prevention
- Audit logging for executed commands

## Development Milestones

### Phase 1: Core MVP (Months 1-3)
- [ ] Basic time tracking with start/stop/pause
- [ ] YAML configuration with validation
- [ ] Simple ritual execution
- [ ] Cross-platform DND automation
- [ ] Basic reporting (daily/weekly)
- [ ] Shell completions

### Phase 2: Integration & Intelligence (Months 4-6)
- [ ] Git hooks for automatic project detection
- [ ] IDE plugins (VS Code, JetBrains)
- [ ] Slack/Discord status integration
- [ ] Calendar blocking
- [ ] Smart break reminders
- [ ] Advanced reporting with visualizations

### Phase 3: Collaboration & Extensibility (Months 7-9)
- [ ] Team features with shared configurations
- [ ] Plugin marketplace
- [ ] Web dashboard for analytics
- [ ] Mobile companion app
- [ ] AI-powered ritual suggestions
- [ ] Export to time tracking services

### Phase 4: Enterprise & Scale (Months 10-12)
- [ ] SSO and enterprise authentication
- [ ] Compliance reporting (SOC2, GDPR)
- [ ] Advanced analytics and insights
- [ ] Custom integrations API
- [ ] White-label options

## Version History

### Version Numbering
- **Major** (X.0.0): Breaking changes, major feature additions
- **Minor** (0.X.0): New features, backwards compatible
- **Patch** (0.0.X): Bug fixes, security patches

### Release Schedule
- **Alpha**: Internal testing and development
- **Beta**: Limited public testing (100 developers)
- **RC**: Release candidate for final testing
- **Stable**: Public release with full support

## Breaking Changes

### Future Breaking Changes
- Configuration schema changes will be documented here
- Command interface modifications will be noted
- API changes for plugin developers will be tracked

## Migration Guides

### From Watson/Timewarrior
- Migration tools will be provided in v0.1.0
- Automatic data import functionality
- Configuration conversion utilities

### Between Rune Versions
- Detailed migration instructions for each major version
- Automated migration tools where possible
- Backwards compatibility notes

## Contributors

### Core Team
- John Ferguson (@ferg-cod3s) - Project Lead & Primary Developer

### Community Contributors
- Contributors will be listed here as the project grows
- Special recognition for significant contributions
- Acknowledgment of bug reports and feature suggestions

## Release Notes Template

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features and capabilities

### Changed
- Changes to existing functionality

### Deprecated
- Features marked for removal in future versions

### Removed
- Features removed in this version

### Fixed
- Bug fixes and corrections

### Security
- Security improvements and vulnerability fixes
```

---

**Note**: This changelog will be updated with each release. For the most current development status, see [TODO.md](TODO.md).