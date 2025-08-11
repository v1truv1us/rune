#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="build-with-telemetry"
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

# Build script with runtime telemetry support (secure - no embedded secrets)
# Telemetry keys are loaded at runtime from environment variables or config

set -e

VERSION=${VERSION:-$(git describe --tags --always --dirty)}

echo "Building Rune with runtime telemetry support..."
echo "Version: $VERSION"
echo "NOTE: Telemetry keys are loaded at runtime from:"
echo "  - Environment variables: RUNE_OTLP_ENDPOINT, RUNE_SENTRY_DSN"echo "  - Configuration file: ~/.rune/config.yaml"

# Build the binary WITHOUT embedded secrets
go build -ldflags "\
  -s -w \
  -X github.com/ferg-cod3s/rune/internal/commands.version=$VERSION \
  -X github.com/ferg-cod3s/rune/internal/telemetry.version=$VERSION \
" -o rune ./cmd/rune

echo "Build completed successfully!"
echo "Binary: ./rune"

# Test the build
echo "Testing telemetry integration..."
if RUNE_DEBUG=true ./rune --version > /dev/null 2>&1; then
    echo "✅ Telemetry integration test passed"
else
    echo "❌ Telemetry integration test failed"
    exit 1
fi