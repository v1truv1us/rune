---
layout: ../layouts/BaseLayout.astro
title: "Beta Program - Rune CLI"
---

# 🚀 Rune CLI Beta Program

Welcome to the Rune CLI Beta! You're among the first to experience the next generation of developer productivity automation.

## 🎯 What is Rune?

Rune is a developer-first CLI that automates your entire work ritual - from time tracking and focus mode to project-specific automation. Think of it as your personal productivity assistant that understands developer workflows.

### ✨ Key Features in Beta

- **⏱️ Intelligent Time Tracking**: Automatic project detection via Git, with session persistence
- **🔕 Cross-Platform Focus Mode**: Do Not Disturb automation for macOS, Windows, and Linux
- **🤖 Ritual Automation**: Custom commands that run when you start/stop work
- **📊 Smart Reporting**: Daily/weekly summaries with project breakdowns
- **🔄 Migration Tools**: Import your data from Watson and Timewarrior
- **🎨 Beautiful CLI**: Colored output with multiple themes
- **🌍 Cross-Platform**: Native support for macOS, Linux, and Windows

## 📥 Installation

### Quick Install (Recommended)
```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash

# Or with Homebrew (macOS)
brew install --cask ferg-cod3s/tap/rune
```

### Manual Installation
1. Download the latest beta from [GitHub Releases](https://github.com/ferg-cod3s/rune/releases)
2. Extract and move to your PATH
3. Run `rune init --guided` to get started

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

## 💬 Beta Feedback Channels

We want to hear from you! Your feedback shapes Rune's future.

### 🎯 Primary Feedback Channels
- **GitHub Discussions**: [Share ideas, ask questions, report bugs](https://github.com/ferg-cod3s/rune/discussions)
- **GitHub Issues**: [Report bugs and request features](https://github.com/ferg-cod3s/rune/issues)
- **Email**: beta@rune.dev (for private feedback)

### 📋 What We're Looking For
- **Workflow integration**: How does Rune fit into your daily routine?
- **Feature requests**: What's missing from your ideal productivity tool?
- **Bug reports**: What's not working as expected?
- **Performance feedback**: Startup time, memory usage, responsiveness
- **Documentation gaps**: What's confusing or missing?
- **Platform-specific issues**: macOS, Linux, Windows compatibility

## 🚧 Beta Status & Limitations

**Current Version**: v0.9.0-beta.1

### ✅ What's Working Great
- Core time tracking functionality
- macOS, Windows, and Linux DND integration
- Basic ritual automation
- Watson/Timewarrior migration
- Shell completions (Bash, Zsh, Fish)
- Configuration management

### ⚠️ Known Limitations
- **PowerShell completions**: Not yet implemented
- **Advanced reporting**: Web dashboard planned for v1.0
- **Plugin system**: Foundation exists, full system coming in v1.0
- **IDE integrations**: VS Code extension planned for v1.0
- **Mobile companion**: iOS/Android apps planned for v1.1

### 🐛 Known Issues
- Migration integration tests have string matching issues (functionality works)
- Some Linux desktop environments may need manual DND setup
- Windows Focus Assist requires Windows 10/11

## 🗺️ Roadmap to v1.0

### Coming Soon (v0.9.x)
- **PowerShell completions**
- **Enhanced error messages**
- **Configuration encryption**
- **Performance optimizations**

### v1.0 Features
- **Web dashboard** for advanced reporting
- **VS Code extension** for status display
- **Plugin system** for extensibility
- **Team collaboration** features
- **Advanced ritual conditions**

### Post-v1.0
- **Mobile companion apps**
- **IDE integrations** (JetBrains, Vim, Emacs)
- **External service integrations** (Slack, Discord, Calendar)
- **AI-powered insights**

## 🏆 Beta Tester Recognition
- **Contributors**: Beta feedback contributors get recognition in v1.0 release notes
- **Early access**: First access to new features and premium capabilities
- **Swag**: Top contributors receive Rune CLI merchandise
- **Community**: Join our private beta Discord for direct developer access

## 🙏 Thank You

Thank you for being part of the Rune CLI beta! Your feedback and support make this project possible.

**Happy coding!** 🚀

---

**Beta Program**: v0.9.0-beta.1  
**Last Updated**: July 16, 2025  
**Next Beta Release**: Weekly (Fridays)

For the latest updates, follow [@RuneCLI](https://twitter.com/RuneCLI) or watch the [GitHub repository](https://github.com/ferg-cod3s/rune).