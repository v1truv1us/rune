# 🚀 Rune CLI Beta Program

Welcome to the Rune CLI Beta! You're among the first to experience the next generation of developer productivity automation.

## 🎯 What is Rune?

Rune is a developer-first CLI productivity platform that automates your daily work rituals, enforces healthy work-life boundaries, and integrates seamlessly with your existing developer workflows. It's designed by developers, for developers who want to focus on code, not time management.

### ✨ Key Features in Beta

- **⏱️ Intelligent Time Tracking**: Automatic project detection via Git repositories, with persistent session management
- **🔕 Cross-Platform Focus Mode**: Native Do Not Disturb integration for macOS, Windows, and Linux
- **🤖 Ritual Automation**: Custom command execution for start/stop workflows (build, deploy, notifications)
- **📊 Smart Reporting**: Comprehensive daily/weekly summaries with detailed project breakdowns
- **🔄 Migration Tools**: Seamless data import from Watson, Timewarrior, and other time trackers
- **🎨 Beautiful CLI**: Rich colored output with customizable themes and progress indicators
- **🌍 Cross-Platform**: Native binaries for macOS (Intel/Apple Silicon), Linux (x64/ARM), and Windows
- **⚡ Performance**: Sub-200ms startup time, minimal memory footprint (<50MB)

## 🚧 Beta Status & What to Expect

**Current Version**: v0.2.0-beta.1  
**Release Cadence**: Weekly beta releases (Fridays)  
**Stability**: Production-ready core features, experimental advanced features

### ✅ Battle-Tested Features
- **Core time tracking**: Robust session management with automatic recovery
- **Cross-platform DND**: Native integration with macOS, Windows 10/11, and major Linux DEs
- **Ritual automation**: Reliable command execution with error handling
- **Data migration**: Proven Watson/Timewarrior import with data validation
- **Shell completions**: Full support for Bash, Zsh, and Fish
- **Configuration**: YAML-based config with validation and migration
- **Performance**: Optimized for developer workflows with minimal overhead

### 🔬 Experimental Features
- **Advanced reporting**: Enhanced analytics and insights (feedback welcome!)
- **Plugin architecture**: Early-stage extensibility framework
- **Telemetry integration**: Optional usage analytics for product improvement
- **Multi-project workflows**: Complex project detection and switching

### ⚠️ Current Limitations
- **PowerShell completions**: Windows PowerShell support coming in v0.3.0
- **Web dashboard**: Advanced reporting UI planned for v1.0
- **IDE extensions**: VS Code integration in development
- **Team features**: Collaboration tools planned for v1.0+

### 🐛 Known Issues & Workarounds
- **Migration tests**: String matching issues in tests (core functionality works perfectly)
- **Linux DND**: Some desktop environments need manual setup ([see docs](docs/linux-dnd.md))
- **Windows Focus Assist**: Requires Windows 10 version 1803+ or Windows 11
- **Git detection**: Very large repositories (>100k files) may have slower project detection

## 📥 Installation

### Quick Install (Recommended)
```bash
# macOS/Linux - One-line installer
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash

# Homebrew (macOS) - Cask installation
brew install --cask ferg-cod3s/tap/rune

# Verify installation
rune --version
```

### Platform-Specific Options

**macOS**
- Homebrew cask (recommended): `brew install --cask ferg-cod3s/tap/rune`
- Direct download: Intel and Apple Silicon binaries available
- Automatic PATH setup and shell completion installation

**Linux**
- Install script supports all major distributions
- Automatic detection of package manager (apt, yum, pacman, etc.)
- Shell completion setup for Bash, Zsh, Fish

**Windows**
- PowerShell installer coming soon
- Manual installation via GitHub releases
- Windows Terminal and PowerShell support

### Manual Installation
1. Download the latest beta from [GitHub Releases](https://github.com/ferg-cod3s/rune/releases/latest)
2. Extract the binary for your platform
3. Move to a directory in your PATH (e.g., `/usr/local/bin`, `~/.local/bin`)
4. Make executable: `chmod +x rune` (macOS/Linux)
5. Run `rune init --guided` to complete setup

### Verification
```bash
# Check installation
rune --version
rune config validate

# Test core functionality
rune init --guided
rune start
rune status
rune stop
```

## 🚀 Quick Start

```bash
# Initialize Rune with guided setup
rune init --guided

# Start your first work session
rune start

# Check your status
rune status

# Take a break
rune pause

# Resume work
rune resume

# End your session
rune stop

# View your daily report
rune report today
```

## 💬 Beta Feedback & Community

Your feedback directly shapes Rune's development. We're building this tool for developers, by developers.

### 🎯 How to Provide Feedback

**GitHub Discussions** (Preferred)
- [💡 Ideas & Feature Requests](https://github.com/ferg-cod3s/rune/discussions/categories/ideas)
- [❓ Q&A and Support](https://github.com/ferg-cod3s/rune/discussions/categories/q-a)
- [📢 Show and Tell](https://github.com/ferg-cod3s/rune/discussions/categories/show-and-tell) - Share your workflows!

**GitHub Issues**
- [🐛 Bug Reports](https://github.com/ferg-cod3s/rune/issues/new?template=bug_report.md)
- [🚀 Feature Requests](https://github.com/ferg-cod3s/rune/issues/new?template=feature_request.md)

**Direct Contact**
- **Email**: beta@rune.dev (for sensitive feedback or private discussions)
- **Twitter**: [@RuneCLI](https://twitter.com/RuneCLI) for quick questions

### 📋 High-Priority Feedback Areas

**Critical for v1.0**
- **Daily workflow integration**: How does Rune fit into your routine?
- **Performance bottlenecks**: Startup time, memory usage, responsiveness
- **Cross-platform compatibility**: Platform-specific issues or improvements
- **Configuration complexity**: What's confusing or could be simplified?

**Feature Development**
- **Missing integrations**: What tools/services should Rune connect with?
- **Ritual automation**: What commands/workflows would you automate?
- **Reporting needs**: What insights would help your productivity?
- **Team collaboration**: How would you use Rune in team environments?

### 🏆 Beta Contributor Recognition

**Hall of Fame**
- Top contributors featured in v1.0 release announcement
- Special recognition badge in GitHub Discussions
- Early access to premium features and enterprise tools

**Community Perks**
- Direct access to development team via GitHub Discussions
- Influence on roadmap priorities and feature development
- Beta merchandise for significant contributors (stickers, shirts, etc.)
- Invitation to exclusive beta community events and AMAs

**Contribution Types We Value**
- Detailed bug reports with reproduction steps
- Feature requests with clear use cases and examples
- Documentation improvements and workflow examples
- Community support and helping other beta testers
- Code contributions (bug fixes, features, tests)

## 📊 Beta Success Metrics

Help us measure success by:
- **Using Rune daily** for at least 2 weeks
- **Providing feedback** on GitHub Discussions
- **Reporting bugs** you encounter
- **Sharing your workflow** configurations
- **Inviting colleagues** who might benefit

## 🛠️ Troubleshooting

### Common Issues

**Installation Problems**
```bash
# Check if Rune is in your PATH
which rune

# Verify installation
rune --version

# Reinstall if needed
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash
```

**Configuration Issues**
```bash
# Reset configuration
rune config reset

# Validate current config
rune config validate

# Re-run guided setup
rune init --guided
```

**DND Not Working**
- **macOS**: Create Shortcuts named "Turn On Do Not Disturb" and "Turn Off Do Not Disturb"
- **Windows**: Ensure Windows 10/11 with Focus Assist enabled
- **Linux**: Check [Linux DND documentation](docs/linux-dnd.md) for your desktop environment

### Getting Help
1. Check the [documentation](docs/)
2. Search [GitHub Issues](https://github.com/ferg-cod3s/rune/issues)
3. Ask in [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions)
4. Email beta@rune.dev for urgent issues

## 🗺️ Development Roadmap

### Next Release (v0.3.0) - Target: August 2025
- **PowerShell completions** for Windows users
- **Enhanced error messages** with actionable suggestions
- **Configuration validation** improvements
- **Performance optimizations** for large Git repositories
- **Ritual condition system** (time-based, project-based triggers)

### v0.5.0 - Target: September 2025
- **Plugin architecture** foundation with example plugins
- **Advanced reporting** with exportable formats (JSON, CSV, PDF)
- **Multi-project session** support
- **Configuration encryption** for sensitive ritual commands
- **Improved Linux DND** support for more desktop environments

### v1.0 - Target: Q4 2025
- **Web dashboard** for advanced analytics and team insights
- **VS Code extension** with status bar integration and commands
- **Team collaboration** features (shared projects, team reports)
- **API endpoints** for external integrations
- **Comprehensive plugin system** with marketplace

### v1.1+ - 2026 and Beyond
- **Mobile companion apps** (iOS/Android) for time tracking on-the-go
- **IDE integrations** (JetBrains, Vim, Emacs, Sublime Text)
- **Service integrations** (Slack, Discord, Calendar, Jira, Linear)
- **AI-powered insights** for productivity optimization
- **Enterprise features** (SSO, audit logs, compliance reporting)

### Community-Driven Features
Features prioritized based on beta feedback and GitHub discussions:
- **Custom themes** and CLI appearance options
- **Backup/sync** solutions for configuration and data
- **Advanced ritual scripting** with conditional logic
- **Integration templates** for popular development workflows

## 🤝 Contributing

Beta testers can contribute in many ways:

### Non-Code Contributions
- **Documentation improvements**
- **Workflow examples** and configurations
- **Bug reports** with detailed reproduction steps
- **Feature requests** with use case descriptions
- **Community support** in discussions

### Code Contributions
- **Bug fixes** for reported issues
- **Feature implementations** from the roadmap
- **Test coverage** improvements
- **Performance optimizations**

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## 📜 License & Privacy

**Open Source Commitment**
- **License**: MIT License - fully open source, no restrictions
- **Source Code**: All code available on [GitHub](https://github.com/ferg-cod3s/rune)
- **Transparency**: No hidden functionality, everything is auditable

**Privacy & Data Handling**
- **Local-First**: All your data stays on your machine by default
- **Optional Telemetry**: Usage analytics are opt-in and can be disabled anytime
- **No Cloud Dependencies**: Rune works completely offline
- **Data Ownership**: You own your time tracking data, export anytime
- **No Personal Data**: We don't collect personal information without explicit consent

**Telemetry Details** (When Enabled)
- Anonymous usage statistics (commands used, feature adoption)
- Performance metrics (startup time, memory usage)
- Error reporting for debugging (no personal data included)
- Disable anytime: `rune config set telemetry.enabled false`

## 🤝 Contributing to Beta

**Non-Technical Contributions**
- **Documentation**: Improve setup guides, add workflow examples
- **Community Support**: Help other beta testers in GitHub Discussions
- **Testing**: Try Rune on different platforms and report compatibility
- **Feedback**: Share detailed use cases and feature requests

**Technical Contributions**
- **Bug Fixes**: Tackle issues from our [GitHub Issues](https://github.com/ferg-cod3s/rune/issues)
- **Feature Development**: Implement features from the roadmap
- **Testing**: Add test coverage for new functionality
- **Performance**: Profile and optimize critical paths

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding standards, and contribution guidelines.

## 🙏 Thank You, Beta Testers!

You're not just testing software - you're helping shape the future of developer productivity tools. Every bug report, feature request, and piece of feedback makes Rune better for the entire developer community.

**Special thanks to our early beta contributors** who've provided invaluable feedback and helped us identify critical improvements.

**Ready to boost your productivity?** 🚀

---

**Beta Program**: v0.2.0-beta.1  
**Last Updated**: July 18, 2025  
**Next Beta Release**: July 25, 2025 (Weekly Fridays)  
**Estimated v1.0**: Q4 2025

**Stay Connected**
- 🐦 Follow [@RuneCLI](https://twitter.com/RuneCLI) for updates
- ⭐ Star the [GitHub repository](https://github.com/ferg-cod3s/rune)
- 💬 Join [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions)
- 📧 Email beta@rune.dev for direct feedback