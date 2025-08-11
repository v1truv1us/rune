# 🚀 Rune CLI v0.2.0-beta.1 - MVP Feature Complete!

**Welcome to the Rune CLI Beta Program!** 

After months of development, we're excited to share the first beta release of Rune - a developer-first CLI that automates your entire work ritual. This release marks **MVP feature completion** with all core functionality implemented and tested.

## 🎯 What is Rune?

Rune transforms how developers manage their daily workflow by automating the repetitive tasks that break your flow:
- **Automatic time tracking** with intelligent project detection
- **Cross-platform focus mode** that silences distractions
- **Custom ritual automation** for project-specific workflows
- **Beautiful reporting** with insights into your productivity patterns

## ✨ Major Features in Beta

### ⏱️ **Intelligent Time Tracking**
- **Automatic project detection** via Git repositories, package.json, go.mod, Cargo.toml
- **Session persistence** across restarts using BBolt database
- **Idle detection** with configurable thresholds
- **Pause/resume functionality** for breaks and interruptions
- **Historical tracking** with daily, weekly, and project-based summaries

### 🔕 **Cross-Platform Do Not Disturb**
- **macOS**: Shortcuts integration with Focus modes and Control Center
- **Windows**: Focus Assist automation via registry and PowerShell
- **Linux**: Comprehensive desktop environment support:
  - GNOME/Ubuntu/Pop!_OS (gsettings)
  - KDE Plasma 5 & 6 (kwriteconfig/kreadconfig)
  - XFCE (xfconf-query)
  - MATE (gsettings/dconf)
  - Cinnamon (gsettings/dconf)
  - Tiling WMs (i3, sway, dwm, awesome, bspwm) via dunst/mako
  - Smart desktop environment detection with fallback strategies

### 🤖 **Ritual Automation Engine**
- **YAML-based configuration** with schema validation
- **Conditional execution** based on day, time, project, or custom conditions
- **Command chaining** with error handling and rollback
- **Progress indicators** for long-running operations
- **Global and project-specific rituals**

### 🔄 **Migration Tools**
- **Watson migration**: Import JSON frames with project mapping
- **Timewarrior migration**: Parse data files with tag handling
- **Dry-run preview** to validate imports before execution
- **Project mapping** to rename/reorganize during migration
- **Comprehensive statistics** and import summaries

### 🎨 **Enhanced User Experience**
- **Colored output** with 8 built-in themes plus custom theme support
- **Relative time formatting**: "2h ago", "yesterday", "last week"
- **Shell completions** for Bash, Zsh, and Fish
- **Interactive setup** with `rune init --guided`
- **Comprehensive help** and error messages

### 📊 **Smart Reporting**
- **Daily/weekly summaries** with time breakdowns
- **Project-based analytics** showing time allocation
- **Export capabilities** to CSV and JSON formats
- **Terminal-based visualization** with colored output
- **Session history** with detailed activity logs

## 🛠️ Installation

### Quick Install
```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash

# Homebrew (macOS)
brew install --cask ferg-cod3s/tap/rune

# Manual download
# Download from GitHub Releases and add to PATH
```

### First Run
```bash
# Interactive setup
rune init --guided

# Start tracking
rune start

# Check status
rune status
```

## 📋 Example Workflows

### Basic Time Tracking
```bash
# Start work session (auto-detects project)
rune start

# Take a break
rune pause

# Resume work
rune resume

# End session
rune stop

# View today's summary
rune report today
```

### Migration from Watson
```bash
# Preview Watson import
rune migrate watson ~/.config/watson/frames --dry-run

# Import with project mapping
rune migrate watson ~/.config/watson/frames --project-map "old-name=new-name"
```

### Custom Rituals
```yaml
# ~/.rune/config.yaml
rituals:
  start:
    global:
      - name: "Setup Development Environment"
        commands:
          - "tmux new-session -d -s work"
          - "code ."
          - "npm run dev"
  stop:
    global:
      - name: "Cleanup"
        commands:
          - "tmux kill-session -t work"
          - "docker-compose down"
```

## 🌟 What Makes Rune Special

### **Developer-First Design**
- Built by developers, for developers
- Understands Git workflows and project structures
- Integrates seamlessly with existing tools
- Respects developer privacy (all data local)

### **Cross-Platform Excellence**
- Native support for macOS, Linux, and Windows
- Platform-specific optimizations (Shortcuts, Focus Assist, desktop environments)
- Consistent experience across all platforms
- Comprehensive fallback mechanisms

### **Automation Without Complexity**
- Simple YAML configuration
- Intelligent defaults that work out of the box
- Progressive complexity (simple start, powerful when needed)
- Extensive documentation and examples

### **Privacy & Performance**
- All data stored locally (no cloud dependencies)
- Sub-200ms startup time
- Minimal memory footprint (<50MB)
- Optional telemetry (fully transparent and disableable)

## 🧪 Beta Testing Focus Areas

We're particularly interested in feedback on:

### **Core Functionality**
- Time tracking accuracy and reliability
- DND integration across different platforms/environments
- Ritual automation effectiveness
- Migration tool completeness

### **User Experience**
- Onboarding and setup process
- Command discoverability and help
- Error messages and troubleshooting
- Documentation clarity

### **Platform Compatibility**
- macOS versions (Big Sur, Monterey, Ventura, Sonoma)
- Linux distributions (Ubuntu, Fedora, Arch, etc.)
- Windows versions (10, 11)
- Desktop environments and window managers

### **Performance & Reliability**
- Startup time and responsiveness
- Memory usage during long sessions
- Database reliability and corruption resistance
- Error recovery and graceful degradation

## 🐛 Known Issues & Limitations

### **Current Limitations**
- PowerShell completions not yet implemented
- Some migration integration tests have string matching issues (functionality works)
- Advanced reporting dashboard planned for v1.0
- Plugin system foundation exists but full system coming in v1.0

### **Platform-Specific Notes**
- **macOS**: Requires Shortcuts app for DND (available on macOS 12+)
- **Windows**: Focus Assist requires Windows 10/11
- **Linux**: Some desktop environments may need manual DND configuration

### **Performance Notes**
- First startup may be slower due to database initialization
- Large Watson/Timewarrior imports may take several seconds
- Git repository detection adds ~10ms per project scan

## 🗺️ Roadmap to v1.0

### **v0.9.x Beta Releases** (Weekly)
- PowerShell completion scripts
- Enhanced error messages and recovery
- Configuration file encryption
- Performance optimizations
- Bug fixes based on beta feedback

### **v1.0 Release** (Target: 4-6 weeks)
- Web-based reporting dashboard
- VS Code extension for status display
- Full plugin system with SDK
- Team collaboration features
- Advanced ritual conditions and triggers

### **Post-v1.0**
- Mobile companion apps (iOS/Android)
- Additional IDE integrations (JetBrains, Vim, Emacs)
- External service integrations (Slack, Discord, Calendar)
- AI-powered productivity insights

## 💬 Beta Feedback

Your feedback shapes Rune's future! Please share:

### **GitHub Discussions** (Preferred)
- General feedback and feature requests
- Workflow sharing and tips
- Community support and questions

### **GitHub Issues**
- Bug reports with reproduction steps
- Specific feature requests
- Documentation improvements

### **Direct Contact**
- Email: beta@rune.dev
- Critical issues or private feedback

## 🏆 Beta Tester Benefits

- **Early access** to new features and releases
- **Recognition** in v1.0 release notes for contributors
- **Direct influence** on product direction and priorities
- **Community access** to private beta Discord
- **Swag** for top contributors (stickers, t-shirts, etc.)

## 📊 Technical Details

### **Architecture**
- **Language**: Go 1.23+ for performance and cross-platform support
- **CLI Framework**: Cobra for command structure and help
- **Configuration**: Viper for YAML parsing and validation
- **Database**: BBolt for local session storage
- **Testing**: Comprehensive unit and integration tests

### **Dependencies**
- Minimal external dependencies for security and reliability
- No cloud services or external APIs required
- Optional telemetry via OpenTelemetry (OTLP logs) and Sentry (can be disabled)
- Platform-specific integrations use native APIs

### **Security**
- No sensitive data collection
- All configuration and data stored locally
- Optional telemetry is transparent and disableable
- Regular security audits and dependency updates

## 🙏 Acknowledgments

Special thanks to:
- **Early testers** who provided invaluable feedback during development
- **Open source community** for inspiration and tools
- **Watson and Timewarrior** projects for pioneering CLI time tracking
- **Go community** for excellent tooling and libraries

## 📜 License & Legal

- **License**: MIT (fully open source)
- **Privacy**: No personal data collection without explicit consent
- **Telemetry**: Optional, transparent, and easily disabled
- **Data**: All user data remains local and under user control

---

## 🚀 Get Started

Ready to transform your development workflow?

1. **Install Rune**: Follow the installation instructions above
2. **Join the beta**: Read [BETA.md](BETA.md) for full beta program details
3. **Get support**: Use GitHub Discussions for questions and feedback
4. **Share feedback**: Help us make Rune even better

**Welcome to the future of developer productivity!** 🎉

---

**Release**: v0.2.0-beta.1  
**Date**: July 16, 2025  
**Compatibility**: macOS 11+, Linux (most distributions), Windows 10/11  
**Download**: [GitHub Releases](https://github.com/ferg-cod3s/rune/releases/tag/v0.2.0-beta.1)