# I Built a CLI That Automates My Entire Developer Workflow (And You Can Too!)

*Originally published on [Dev.to](https://dev.to)*

---

**TL;DR**: After getting frustrated with manually managing my development workflow, I built [Rune](https://github.com/ferg-cod3s/rune) - a CLI that automatically tracks time, manages focus mode, and runs custom rituals. It's now in beta and I'd love your feedback!

---

## The Problem: Context Switching is Killing My Flow

As developers, we're constantly switching between tasks, projects, and mental contexts. But how often do you actually track this? How much time do you spend on each project? When was your last real break?

I found myself constantly:
- **Forgetting to start/stop time tracking** (goodbye accurate client billing)
- **Getting distracted by notifications** during deep work
- **Manually running the same setup commands** for each project
- **Having no idea where my time actually went** at the end of the day

Sound familiar? 

## The Solution: Automate Everything

Instead of relying on willpower and memory, I decided to build a tool that would handle all of this automatically. Meet **Rune** - a developer-first CLI that automates your entire work ritual.

Here's what it does:

### ⏱️ **Intelligent Time Tracking**
```bash
# Just start working - Rune figures out the rest
rune start

# It automatically detects your project from Git, package.json, go.mod, etc.
# Session persists across restarts, handles idle time, tracks everything
```

### 🔕 **Cross-Platform Focus Mode**
```bash
# Automatically enables Do Not Disturb on macOS, Windows, and Linux
rune start  # DND turns on
rune stop   # DND turns off
```

No more manual notification management. Rune handles:
- **macOS**: Shortcuts integration with Focus modes
- **Windows**: Focus Assist automation
- **Linux**: GNOME, KDE, XFCE, i3, sway, and more

### 🤖 **Custom Ritual Automation**
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

Every time you start/stop work, Rune runs your custom commands. No more forgetting to start the dev server or clean up containers.

### 📊 **Beautiful Reporting**
```bash
rune report today
# ┌─────────────────────────────────────────────────────────────┐
# │                     📊 Daily Report                        │
# │                    July 16, 2025                           │
# ├─────────────────────────────────────────────────────────────┤
# │ Total Time: 7h 23m                                         │
# │ Projects:                                                   │
# │   • rune-cli: 4h 15m (57.5%)                              │
# │   • client-project: 2h 30m (33.8%)                        │
# │   • documentation: 38m (8.6%)                             │
# └─────────────────────────────────────────────────────────────┘
```

## The Journey: 6 Months of Development

Building Rune taught me a lot about developer productivity and cross-platform development. Here are some key insights:

### **Why Go?**
I chose Go for several reasons:
- **Cross-platform**: Single binary for macOS, Linux, Windows
- **Performance**: Sub-200ms startup time, <50MB memory usage
- **Reliability**: Strong typing and excellent error handling
- **Ecosystem**: Great CLI libraries (Cobra, Viper) and testing tools

### **The Hardest Part: Cross-Platform DND**
Getting Do Not Disturb working across all platforms was surprisingly complex:

**macOS**: Uses Shortcuts app integration with fallback to AppleScript
**Windows**: Registry manipulation for Focus Assist + PowerShell automation  
**Linux**: Supporting 6+ desktop environments (GNOME, KDE, XFCE, i3, sway, etc.)

Each platform has its own quirks and APIs. The key was building robust fallback mechanisms.

### **Database Choice: BBolt**
For local storage, I chose BBolt (embedded key-value store) because:
- **No dependencies**: Single file, no server required
- **ACID transactions**: Data integrity for time tracking
- **Performance**: Fast reads/writes for session data
- **Reliability**: Used by etcd and other production systems

### **Migration Tools: The Unexpected Challenge**
Many developers already use Watson or Timewarrior. Building migration tools turned out to be crucial for adoption:

```bash
# Import your Watson data
rune migrate watson ~/.config/watson/frames --dry-run

# Import Timewarrior with project mapping
rune migrate timewarrior ~/.local/share/timewarrior/data --project-map "old=new"
```

Parsing different time tracking formats and handling edge cases took weeks, but it's essential for user onboarding.

## Technical Architecture

Here's how Rune is structured:

```
rune/
├── cmd/rune/           # CLI entry point
├── internal/
│   ├── commands/       # Cobra command implementations
│   ├── tracking/       # Time tracking and session management
│   ├── dnd/           # Cross-platform Do Not Disturb
│   ├── rituals/       # Automation engine
│   ├── config/        # YAML configuration management
│   ├── notifications/ # System notifications
│   └── telemetry/     # Optional usage analytics
└── docs/              # Platform-specific documentation
```

### **Key Design Decisions**

**1. Local-First**: All data stays on your machine. No cloud dependencies.

**2. Configuration as Code**: YAML files that can be version controlled and shared.

**3. Progressive Complexity**: Works out of the box, but highly customizable.

**4. Extensive Testing**: Unit tests, integration tests, and cross-platform compatibility tests.

## What I Learned

### **Developer Productivity is Personal**
Everyone has different workflows. The key is building flexible primitives that can be composed in different ways.

### **Cross-Platform is Hard**
Supporting macOS, Linux, and Windows means dealing with:
- Different notification systems
- Various desktop environments  
- Platform-specific APIs and permissions
- Inconsistent command-line tools

### **Migration Matters**
Developers won't switch tools unless migration is seamless. Spending time on import tools pays off in adoption.

### **Documentation is Critical**
A CLI tool is only as good as its documentation. I spent significant time on:
- Platform-specific setup guides
- Troubleshooting documentation
- Example configurations
- Clear error messages

## The Beta Launch

After 6 months of development, Rune is now in beta! Here's what's included:

### ✅ **What's Working Great**
- Core time tracking functionality
- Cross-platform DND automation
- Ritual automation engine
- Watson/Timewarrior migration
- Shell completions (Bash, Zsh, Fish)
- Colored output with themes

### 🚧 **What's Coming in v1.0**
- Web-based reporting dashboard
- VS Code extension
- Plugin system
- Team collaboration features
- PowerShell completions

## Try It Yourself

Want to automate your workflow? Here's how to get started:

### **Installation**
```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash

# Or with Homebrew
brew install --cask ferg-cod3s/tap/rune
```

### **Quick Start**
```bash
# Interactive setup
rune init --guided

# Start your first session
rune start

# Check status
rune status

# End session and see report
rune stop
rune report today
```

### **Beta Feedback**
I'm actively looking for beta testers! If you try Rune, I'd love to hear:
- How it fits into your workflow
- What features are missing
- Any bugs or issues you encounter
- Ideas for improvements

**Feedback channels**:
- [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions)
- [GitHub Issues](https://github.com/ferg-cod3s/rune/issues)
- Email: beta@rune.dev

## What's Next?

The roadmap to v1.0 includes:

**Short-term (v0.9.x betas)**:
- PowerShell completion scripts
- Enhanced error messages
- Performance optimizations
- Bug fixes based on beta feedback

**v1.0 Release**:
- Web dashboard for advanced reporting
- VS Code extension for status display
- Full plugin system with SDK
- Team collaboration features

**Post-v1.0**:
- Mobile companion apps
- Additional IDE integrations
- External service integrations (Slack, Discord)
- AI-powered productivity insights

## Join the Beta!

If you're interested in automating your developer workflow, I'd love to have you try Rune. It's completely free and open source (MIT license).

**Links**:
- [GitHub Repository](https://github.com/ferg-cod3s/rune)
- [Beta Documentation](https://github.com/ferg-cod3s/rune/blob/main/BETA.md)
- [Installation Guide](https://github.com/ferg-cod3s/rune#installation)

## Questions?

Drop a comment below or reach out on:
- [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions)
- [Twitter](https://twitter.com/RuneCLI)
- Email: hello@rune.dev

Thanks for reading! I'm excited to see how Rune can help improve your development workflow. 🚀

---

*What tools do you use to manage your development workflow? Have you built any automation that saves you time? I'd love to hear about it in the comments!*

---

**Tags**: #productivity #cli #golang #automation #timetracking #developer #workflow #opensource