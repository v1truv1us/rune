---
layout: ../layouts/DocsLayout.astro
title: "Install Rune CLI - Developer Time Tracking Tool | Watson Alternative Setup"
description: "Step-by-step guide to install Rune CLI, the modern Watson and Timewarrior alternative. Quick setup for developer time tracking, project detection, and productivity automation on macOS, Linux, and Windows."
image: "/og-getting-started.png"
category: "Getting Started"
lastModified: "2024-07-23T00:00:00Z"
---

> **🚀 Beta Status**: You're installing v0.2.0-beta.1 - the MVP feature-complete version. [Join the beta program](/beta) for updates and feedback channels.

## Installation

<div class="installation-cards">
  <div class="installation-card">
    <h4>🚀 Quick Install (Recommended)</h4>
    <p>One-line installer for macOS and Linux systems</p>
    <pre><code># macOS/Linux
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | bash</code></pre>
  </div>

  <div class="installation-card">
    <h4>🍺 Homebrew (macOS)</h4>
    <p>Install via Homebrew package manager</p>
    <pre><code>brew install --cask ferg-cod3s/tap/rune</code></pre>
  </div>

  <div class="installation-card">
    <h4>🐹 Go Install</h4>
    <p>Build from source with Go toolchain</p>
    <pre><code>go install github.com/ferg-cod3s/rune/cmd/rune@latest</code></pre>
  </div>

  <div class="installation-card">
    <h4>📦 Manual Installation</h4>
    <p>Download and install manually</p>
    <ol>
      <li>Download from <a href="https://github.com/ferg-cod3s/rune/releases/tag/v0.2.0-beta.1">GitHub Releases</a></li>
      <li>Extract binary and move to your <code>PATH</code></li>
      <li>Make executable: <code>chmod +x rune</code></li>
    </ol>
  </div>
</div>

<div class="installation-cards">
  <div class="installation-card">
    <h4>🪟 Windows Installation</h4>
    <p>Download the Windows binary from <a href="https://github.com/ferg-cod3s/rune/releases/tag/v0.2.0-beta.1">GitHub Releases</a> and add to your <code>PATH</code>.</p>
  </div>
</div>

## Initial Setup

<div class="installation-cards">
  <div class="installation-card">
    <h4>1️⃣ Initialize Configuration</h4>
    <p>Interactive setup to configure your work environment</p>
    <pre><code>rune init --guided</code></pre>
    <p><strong>This setup configures:</strong></p>
    <ul>
      <li>Work hours and break intervals</li>
      <li>Project detection rules</li>
      <li>Start/stop rituals</li>
      <li>Integration preferences</li>
    </ul>
  </div>

  <div class="installation-card">
    <h4>2️⃣ Verify Installation</h4>
    <p>Confirm everything is working correctly</p>
    <pre><code>rune --version
rune status</code></pre>
  </div>
</div>

## Basic Usage

### Daily Workflow

<div class="installation-cards">
  <div class="installation-card">
    <h4>🚀 Start Your Workday</h4>
    <pre><code># Start tracking and enable focus mode
rune start

# Check current status
rune status</code></pre>
  </div>

  <div class="installation-card">
    <h4>⏸️ Take Breaks</h4>
    <pre><code># Pause for a break
rune pause

# Resume work
rune resume</code></pre>
  </div>

  <div class="installation-card">
    <h4>🏁 End Your Day</h4>
    <pre><code># Stop tracking and generate report
rune stop</code></pre>
  </div>
</div>

### Time Tracking & Reporting

<div class="installation-cards">
  <div class="installation-card">
    <h4>📊 Daily Reports</h4>
    <pre><code># View today's report
rune report --today</code></pre>
  </div>

  <div class="installation-card">
    <h4>📈 Weekly Reports</h4>
    <pre><code># View this week's report
rune report --week</code></pre>
  </div>

  <div class="installation-card">
    <h4>📅 Custom Date Range</h4>
    <pre><code># View specific date range
rune report --from 2024-01-01 --to 2024-01-07</code></pre>
  </div>
</div>

### Configuration Management

<div class="installation-cards">
  <div class="installation-card">
    <h4>⚙️ Edit Configuration</h4>
    <pre><code># Edit configuration file
rune config edit</code></pre>
  </div>

  <div class="installation-card">
    <h4>✅ Validate Settings</h4>
    <pre><code># Validate configuration
rune config validate</code></pre>
  </div>

  <div class="installation-card">
    <h4>👀 View Current Config</h4>
    <pre><code># Show current configuration
rune config show</code></pre>
  </div>
</div>

## Next Steps

<div class="feature-grid">
  <div class="feature-card">
    <span class="emoji">🤖</span>
    <h4><a href="/docs/rituals">Configure Rituals</a></h4>
    <p>Set up custom commands that run when you start and stop work sessions.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">📁</span>
    <h4><a href="/docs/projects">Project Detection</a></h4>
    <p>Configure automatic project detection via Git repos and package files.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">⚡</span>
    <h4><a href="/docs/commands">Command Reference</a></h4>
    <p>Explore all available commands with detailed examples and options.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">💡</span>
    <h4><a href="/examples">Example Workflows</a></h4>
    <p>See real-world examples and configurations for different development setups.</p>
  </div>
</div>

## Beta Feedback

As a beta user, your feedback is invaluable! Please share:

- **GitHub Discussions**: [General feedback and questions](https://github.com/ferg-cod3s/rune/discussions)
- **GitHub Issues**: [Bug reports and feature requests](https://github.com/ferg-cod3s/rune/issues)
- **Email**: `beta@rune.dev` for private feedback

## Getting Help

- Run `rune --help` for command overview
- Run `rune <command> --help` for specific command help
- Check the [Beta Program page](/beta) for known issues
- Visit [GitHub Discussions](https://github.com/ferg-cod3s/rune/discussions) for community support
