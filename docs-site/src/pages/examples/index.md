---
layout: ../../layouts/BaseLayout.astro
title: "Examples - Rune CLI"
description: "Real-world examples and configurations for Rune CLI"
---

# Examples

## Workflow Examples

### Frontend Developer

Perfect for React, Vue, Svelte, or any frontend development workflow.

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 50m

projects:
  - name: "web-app"
    detect: ["git:web-app", "dir:~/projects/web-app"]

rituals:
  start:
    global:
      - name: "Pull latest changes"
        command: "git pull"
      - name: "Install dependencies"
        command: "pnpm install"
      - name: "Start dev server"
        command: "pnpm dev"
        background: true
      - name: "Open browser"
        command: "open http://localhost:3000"

  stop:
    global:
      - name: "Stop dev server"
        command: "pkill -f 'pnpm dev'"
      - name: "Commit WIP"
        command: "git add -A && git commit -m 'WIP: $(date)'"
        optional: true

integrations:
  git:
    enabled: true
    auto_detect_project: true
```

### Backend Developer

Ideal for API development with Docker services.

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 45m

projects:
  - name: "api-service"
    detect: ["git:api-service"]

rituals:
  start:
    global:
      - name: "Start Docker services"
        command: "docker-compose up -d postgres redis"
      - name: "Run migrations"
        command: "go run migrate.go"
      - name: "Start API server"
        command: "air" # Hot reload for Go
        background: true

  stop:
    global:
      - name: "Stop API server"
        command: "pkill -f air"
      - name: "Stop Docker services"
        command: "docker-compose down"
      - name: "Run tests"
        command: "go test ./..."
        optional: true
```

### DevOps Engineer

Perfect for infrastructure and deployment workflows.

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 60m

projects:
  - name: "infrastructure"
    detect: ["git:infra", "dir:~/infra"]

rituals:
  start:
    global:
      - name: "Check cluster status"
        command: "kubectl get nodes"
      - name: "Check monitoring"
        command: "open https://grafana.company.com"
      - name: "Update Terraform"
        command: "terraform plan"
        working_dir: "~/infra/terraform"

  stop:
    global:
      - name: "Security scan"
        command: "trivy fs ."
        optional: true
      - name: "Generate report"
        command: "./scripts/daily-report.sh"
```

### Freelancer/Consultant

Time tracking focused with client project separation.

```yaml
version: 1
settings:
  work_hours: 6.0
  break_interval: 45m

projects:
  - name: "client-acme"
    detect: ["dir:~/clients/acme"]
    client: "Acme Corp"
    rate: 150.00

  - name: "client-beta"
    detect: ["dir:~/clients/beta"]
    client: "Beta Inc"
    rate: 175.00

rituals:
  start:
    global:
      - name: "Time tracking reminder"
        command: "echo 'Starting work for {{project}} at ${{rate}}/hr'"

  stop:
    global:
      - name: "Log time entry"
        command: "rune report --today --project {{project}} --format csv >> ~/invoices/{{client}}-$(date +%Y-%m).csv"
      - name: "Backup work"
        command: "rsync -av ~/clients/ ~/backups/clients-$(date +%Y%m%d)/"
        optional: true
```

## Configuration Patterns

### Multi-Project Monorepo

```yaml
projects:
  - name: "frontend"
    detect: ["dir:~/monorepo/apps/web"]

  - name: "backend"
    detect: ["dir:~/monorepo/apps/api"]

  - name: "mobile"
    detect: ["dir:~/monorepo/apps/mobile"]

rituals:
  start:
    per_project:
      frontend:
        - name: "Start web app"
          command: "pnpm dev:web"
          background: true

      backend:
        - name: "Start API"
          command: "pnpm dev:api"
          background: true

      mobile:
        - name: "Start Metro"
          command: "pnpm dev:mobile"
          background: true
```

### Remote Work Setup

```yaml
settings:
  work_hours: 8.0
  focus:
    dnd_on_start: true
    block_websites:
      - "facebook.com"
      - "twitter.com"
      - "reddit.com"

rituals:
  start:
    global:
      - name: "Set Slack status"
        command: "slack-cli status '🔨 Deep work mode'"
      - name: "Block distractions"
        command: "sudo dscacheutil -flushcache"
      - name: "Start focus music"
        command: "open -a Spotify"

  break:
    global:
      - name: "Stretch reminder"
        command: "osascript -e 'display notification \"Time to stretch\!\" with title \"Break Time\"'"

  stop:
    global:
      - name: "Clear Slack status"
        command: "slack-cli status ''"
      - name: "End of day summary"
        command: "rune report --today"
```

### Open Source Contributor

```yaml
projects:
  - name: "oss-project-a"
    detect: ["git-remote:github.com/owner/repo-a"]

  - name: "oss-project-b"
    detect: ["git-remote:github.com/owner/repo-b"]

rituals:
  start:
    global:
      - name: "Check issues"
        command: "gh issue list --repo {{git_remote}}"
      - name: "Pull latest"
        command: "git pull upstream main"

  stop:
    global:
      - name: "Check if ready for PR"
        command: "git status --porcelain"
      - name: "Run tests"
        command: "npm test"
        optional: true
```

## Integration Examples

### Slack Integration

```yaml
integrations:
  slack:
    enabled: true
    workspace: "myteam"
    dnd_on_start: true
    status_on_start: "🔨 Working on {{project}}"
    status_on_break: "☕ Taking a break"
    status_on_stop: ""

rituals:
  start:
    global:
      - name: "Notify team"
        command: "slack-cli post '#dev' 'Starting work on {{project}}'"
        optional: true
```

### Calendar Integration

```yaml
integrations:
  calendar:
    enabled: true
    provider: "google"
    block_calendar: true
    meeting_detection: true

rituals:
  start:
    global:
      - name: "Block focus time"
        command: "gcal-cli create 'Deep Work - {{project}}' --duration 4h"
        optional: true
```

### Git Automation

```yaml
integrations:
  git:
    enabled: true
    auto_detect_project: true
    branch_in_status: true

rituals:
  start:
    global:
      - name: "Create feature branch"
        command: "git checkout -b feature/daily-work-$(date +%Y%m%d)"
        optional: true
        when: ["git_clean"]

  stop:
    global:
      - name: "Auto-commit progress"
        command: "git add -A && git commit -m 'Daily progress: $(date)'"
        optional: true
        when: ["git_dirty"]
```

## Ritual Patterns

### Development Environment Setup

```yaml
rituals:
  start:
    global:
      - name: "Check system resources"
        command: "df -h && free -h"
      - name: "Start Docker"
        command: "docker-compose up -d database cache"
      - name: "Update packages"
        command: "npm outdated"
        optional: true
      - name: "Start development tools"
        command: "code . && open -a Terminal"
```

### End of Day Cleanup

```yaml
rituals:
  stop:
    global:
      - name: "Save session"
        command: "tmux capture-session -p > ~/logs/session-$(date +%Y%m%d).log"
        optional: true
      - name: "Backup important files"
        command: "rsync -av ~/important/ ~/backups/"
        optional: true
      - name: "Clean temporary files"
        command: "rm -rf /tmp/dev-*"
        optional: true
      - name: "Generate daily summary"
        command: "rune report --today --format json > ~/reports/$(date +%Y%m%d).json"
```

### Break Automation

```yaml
rituals:
  break:
    global:
      - name: "Lock screen"
        command: "osascript -e 'tell application \"System Events\" to keystroke \"q\" using {control down, command down}'"
      - name: "Pause music"
        command: "osascript -e 'tell application \"Spotify\" to pause'"
        optional: true

  resume:
    global:
      - name: "Resume music"
        command: "osascript -e 'tell application \"Spotify\" to play'"
        optional: true
      - name: "Check notifications"
        command: "echo 'Welcome back\! Checking for important updates...'"
```

## Platform-Specific Examples

### macOS

```yaml
rituals:
  start:
    global:
      - name: "Enable Do Not Disturb"
        command: "shortcuts run 'Enable DND'"
      - name: "Set desktop wallpaper"
        command: 'osascript -e ''tell application "System Events" to set picture of every desktop to "/path/to/work-wallpaper.jpg"'''

  stop:
    global:
      - name: "Disable Do Not Disturb"
        command: "shortcuts run 'Disable DND'"
```

### Linux

```yaml
rituals:
  start:
    global:
      - name: "Set work profile"
        command: "gsettings set org.gnome.desktop.background picture-uri file:///home/user/wallpapers/work.jpg"
      - name: "Enable focus mode"
        command: "dunstctl set-paused true"

  stop:
    global:
      - name: "Disable focus mode"
        command: "dunstctl set-paused false"
```

### Windows

```yaml
rituals:
  start:
    global:
      - name: "Set focus assist"
        command: 'powershell -Command "& {Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait(''^{F1}'')}"'

  stop:
    global:
      - name: "Disable focus assist"
        command: 'powershell -Command "& {Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait(''^{F1}'')}"'
```

EOF < /dev/null
