---
layout: ../layouts/DocsLayout.astro
title: "Beta Program - Rune CLI | Join Early Access Testing"
description: "Join the Rune CLI beta program and help shape the future of developer productivity automation. Get early access to features, provide feedback, and connect with the developer community."
image: "/og-beta.png"
category: "Beta Program"
lastModified: "2024-07-23T00:00:00Z"
---

Welcome to the Rune CLI Beta! You're among the first to experience the next generation of developer productivity automation.

## What is Rune?

Rune is a developer-first CLI that automates your entire work ritual - from time tracking and focus mode to project-specific automation. Think of it as your personal productivity assistant that understands developer workflows.

## Key Features in Beta

<div class="feature-grid">
  <div class="feature-card">
    <span class="emoji">⏱️</span>
    <h4>Intelligent Time Tracking</h4>
    <p>Automatic project detection via Git, package.json, go.mod with session persistence across restarts.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🔕</span>
    <h4>Cross-Platform Focus Mode</h4>
    <p>Do Not Disturb automation for macOS, Windows, and Linux with 6+ desktop environment support.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🤖</span>
    <h4>Ritual Automation</h4>
    <p>Custom commands that run when you start/stop work with YAML configuration and error handling.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">📊</span>
    <h4>Smart Reporting</h4>
    <p>Daily/weekly summaries with project breakdowns and colored CLI output with themes.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🔄</span>
    <h4>Migration Tools</h4>
    <p>Import your data from Watson and Timewarrior with project mapping and dry-run preview.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🌍</span>
    <h4>Cross-Platform</h4>
    <p>Native support for macOS, Linux, and Windows with platform-specific optimizations.</p>
  </div>
</div>

## Installation

<div class="installation-cards">
  <div class="installation-card">
    <h4>🚀 Quick Install (Recommended)</h4>
    <pre><code># macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash

# Or with Homebrew (macOS)
brew install --cask ferg-cod3s/tap/rune</code></pre>
  </div>

  <div class="installation-card">
    <h4>📦 Manual Installation</h4>
    <ol>
      <li>Download the latest beta from <a href="https://github.com/ferg-cod3s/rune/releases">GitHub Releases</a></li>
      <li>Extract and move to your PATH</li>
      <li>Run <code>rune init --guided</code> to get started</li>
    </ol>
  </div>
</div>

## Quick Start

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

## Beta Feedback Channels

We want to hear from you! Your feedback shapes Rune's future.

<div class="feature-grid">
  <div class="feature-card">
    <span class="emoji">💬</span>
    <h4><a href="https://github.com/ferg-cod3s/rune/discussions">GitHub Discussions</a></h4>
    <p>Share ideas, ask questions, and connect with other beta testers in our community forum.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🐛</span>
    <h4><a href="https://github.com/ferg-cod3s/rune/issues">GitHub Issues</a></h4>
    <p>Report bugs and request new features with our issue tracking system.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">📧</span>
    <h4>Email Feedback</h4>
    <p>Send private feedback directly to <strong>beta@rune.dev</strong> for sensitive topics.</p>
  </div>
</div>

### What We're Looking For

- **Workflow integration**: How does Rune fit into your daily routine?
- **Feature requests**: What's missing from your ideal productivity tool?
- **Bug reports**: What's not working as expected?
- **Performance feedback**: Startup time, memory usage, responsiveness
- **Documentation gaps**: What's confusing or missing?
- **Platform-specific issues**: macOS, Linux, Windows compatibility

## 🚧 Beta Status & Limitations

**Current Version**: v0.2.0-beta.1

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

**Beta Program**: v0.2.0-beta.1  
**Last Updated**: July 16, 2025  
**Next Beta Release**: Weekly (Fridays)

For the latest updates, follow [@RuneCLI](https://twitter.com/RuneCLI) or watch the [GitHub repository](https://github.com/ferg-cod3s/rune).