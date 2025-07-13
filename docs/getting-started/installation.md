# Installation Guide

This guide covers all the ways to install Rune on your system.

## Quick Install (Recommended)

### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh
```

### Windows (PowerShell)
```powershell
# Coming soon - use Go install for now
```

## Package Managers

### Homebrew (macOS/Linux)
```bash
brew tap ferg-cod3s/tap
brew install --cask rune
```

### Go Install
```bash
go install github.com/ferg-cod3s/rune/cmd/rune@latest
```

## Manual Installation

### Download Binary
1. Go to [GitHub Releases](https://github.com/ferg-cod3s/rune/releases)
2. Download the appropriate binary for your platform
3. Extract and place in your PATH

### Build from Source
```bash
git clone https://github.com/ferg-cod3s/rune.git
cd rune
make build
sudo cp bin/rune /usr/local/bin/
```

## Platform-Specific Packages

### Debian/Ubuntu
```bash
# Download .deb from releases
wget https://github.com/ferg-cod3s/rune/releases/latest/download/rune_linux_amd64.deb
sudo dpkg -i rune_linux_amd64.deb
```

### RHEL/CentOS/Fedora
```bash
# Download .rpm from releases
wget https://github.com/ferg-cod3s/rune/releases/latest/download/rune_linux_amd64.rpm
sudo rpm -i rune_linux_amd64.rpm
```

### Arch Linux
```bash
# Available in AUR (coming soon)
yay -S rune-cli
```

## Verify Installation

```bash
rune --version
```

You should see:
```
 ______     __  __     __   __     ______   
/\  == \   /\ \/\ \   /\ "-.\ \   /\  ___\  
\ \  __<   \ \ \_\ \  \ \ \-.  \  \ \  __\  
 \ \_\ \_\  \ \_____\  \ \_\\"\_\  \ \_____\
  \/_/ /_/   \/_____/   \/_/ \/_/   \/_____/ 

version x.x.x
```

## Shell Completions

### Bash
```bash
rune completion bash > /etc/bash_completion.d/rune
```

### Zsh
```bash
rune completion zsh > "${fpath[1]}/_rune"
```

### Fish
```bash
rune completion fish > ~/.config/fish/completions/rune.fish
```

## Next Steps

- [Quick Start Tutorial](./quickstart.md)
- [Configuration Setup](../configuration/setup.md)
- [First Ritual](./first-ritual.md)

## Troubleshooting

See [Troubleshooting Guide](./troubleshooting.md) for common installation issues.