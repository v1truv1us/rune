#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="rune-install"
if [ -z "$TMUX" ] && [ -z "$SKIP_AUTO_TMUX" ]; then
  # Check if tmux is available
  if command -v tmux >/dev/null 2>&1; then
    # If session exists, attach; otherwise create new session
    if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
      echo "Attaching to existing tmux session: $SESSION_NAME"
      exec tmux attach-session -t "$SESSION_NAME"
    else
      echo "Starting new tmux session: $SESSION_NAME"
      exec tmux new-session -s "$SESSION_NAME" "$0" "$@"
    fi
  else
    # Fallback: check if running in an interactive terminal
    if [ ! -t 0 ]; then
      echo "This script must be run in an interactive terminal (tmux not available)."
      exit 1
    fi
  fi
fi

# Rune CLI Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh
# Or with options: curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --skip-homebrew
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="ferg-cod3s/rune"
BINARY_NAME="rune"
INSTALL_DIR="/usr/local/bin"
SKIP_HOMEBREW=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-homebrew)
            SKIP_HOMEBREW=true
            shift
            ;;
        -h|--help)
            echo "Rune CLI Installation Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --skip-homebrew    Skip Homebrew suggestion even if available"
            echo "  -h, --help         Show this help message"
            echo ""
            echo "Examples:"
            echo "  curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh"
            echo "  curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --skip-homebrew"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case $OS in
    linux) OS="Linux" ;;
    darwin) OS="Darwin" ;;
    *) echo -e "${RED}Unsupported OS: $OS${NC}"; exit 1 ;;
esac

echo -e "${GREEN}Installing Rune CLI...${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"

# Check if Homebrew is available and suggest using it instead
if [ "$SKIP_HOMEBREW" = false ] && command -v brew >/dev/null 2>&1; then
    echo ""
    echo -e "${YELLOW}📦 Homebrew detected!${NC}"
    echo "For easier installation and updates, consider using:"
    echo "  brew install --cask ferg-cod3s/tap/rune"
    echo ""
    echo "Or continue with direct binary installation..."
    echo ""
fi

# Get latest release
echo -e "${YELLOW}Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}Failed to fetch latest release${NC}"
    exit 1
fi

echo "Latest version: $LATEST_RELEASE"

# Construct download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
echo -e "${YELLOW}Downloading $DOWNLOAD_URL...${NC}"
curl -sL "$DOWNLOAD_URL" | tar xz

# Check if binary exists
if [ ! -f "$BINARY_NAME" ]; then
    echo -e "${RED}Binary not found in archive${NC}"
    exit 1
fi

# Make binary executable
chmod +x "$BINARY_NAME"

# Install binary
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
fi

# Install shell completions if available
if [ -d "completions" ]; then
    echo -e "${YELLOW}Installing shell completions...${NC}"
    
    # Bash completion
    if [ -f "completions/rune.bash" ]; then
        if [ -d "/usr/share/bash-completion/completions" ]; then
            sudo cp "completions/rune.bash" "/usr/share/bash-completion/completions/rune" 2>/dev/null || true
        elif [ -d "/etc/bash_completion.d" ]; then
            sudo cp "completions/rune.bash" "/etc/bash_completion.d/rune" 2>/dev/null || true
        fi
    fi
    
    # Zsh completion
    if [ -f "completions/rune.zsh" ]; then
        if [ -d "/usr/share/zsh/site-functions" ]; then
            sudo cp "completions/rune.zsh" "/usr/share/zsh/site-functions/_rune" 2>/dev/null || true
        fi
    fi
    
    # Fish completion
    if [ -f "completions/rune.fish" ]; then
        if [ -d "/usr/share/fish/completions" ]; then
            sudo cp "completions/rune.fish" "/usr/share/fish/completions/" 2>/dev/null || true
        fi
    fi
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Rune CLI installed successfully!${NC}"
    echo ""
    echo "🚀 Getting Started:"
    echo "  rune --help              # Show available commands"
    echo "  rune init --guided       # Create your first configuration"
    echo "  rune update              # Update to latest version"
    echo ""
    if [ "$SKIP_HOMEBREW" = false ] && command -v brew >/dev/null 2>&1; then
        echo "💡 For future updates, you can also use:"
        echo "  brew install --cask ferg-cod3s/tap/rune  # Switch to Homebrew"
        echo "  brew upgrade rune                 # Update via Homebrew"
        echo ""
    fi
else
    echo -e "${RED}Installation failed. Please check that $INSTALL_DIR is in your PATH.${NC}"
    exit 1
fi