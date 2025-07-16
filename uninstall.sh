#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="rune-uninstall"
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

# Rune CLI Uninstall Script
# Usage: curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/uninstall.sh | sh
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="rune"
INSTALL_DIR="/usr/local/bin"

echo -e "${YELLOW}Uninstalling Rune CLI...${NC}"

# Function to remove file if it exists
remove_file() {
    local file="$1"
    local description="$2"
    
    if [ -f "$file" ]; then
        echo "Removing $description..."
        if [ -w "$(dirname "$file")" ]; then
            rm "$file"
        else
            sudo rm "$file"
        fi
        echo -e "${GREEN}✓ Removed $file${NC}"
    fi
}

# Remove main binary
remove_file "$INSTALL_DIR/$BINARY_NAME" "main binary"

# Remove shell completions
echo -e "${YELLOW}Removing shell completions...${NC}"

# Bash completions
remove_file "/usr/share/bash-completion/completions/rune" "bash completion"
remove_file "/etc/bash_completion.d/rune" "bash completion (alternative location)"

# Zsh completions
remove_file "/usr/share/zsh/site-functions/_rune" "zsh completion"

# Fish completions
remove_file "/usr/share/fish/completions/rune.fish" "fish completion"

# Check if binary is still accessible
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Warning: 'rune' command is still accessible${NC}"
    echo "This might be because:"
    echo "  - It's installed via Homebrew (run: brew uninstall rune)"
    echo "  - It's installed in a different location"
    echo "  - Your shell cache needs to be refreshed (run: hash -r)"
    echo ""
    echo "Current location: $(which rune)"
else
    echo -e "${GREEN}✓ Rune CLI uninstalled successfully!${NC}"
fi

echo ""
echo "Note: This script only removes files installed by the curl-based installer."
echo "If you installed Rune via Homebrew, run: brew uninstall rune"
echo ""
echo "Your configuration files in ~/.config/rune/ were left untouched."
echo "Remove them manually if desired: rm -rf ~/.config/rune/"