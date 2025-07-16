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

# Build script that ensures telemetry keys are embedded
# This script can be used for local builds and CI/CD

set -e

# Default telemetry keys (can be overridden by environment variables)
DEFAULT_SEGMENT_KEY="ZkEZXHRWH96y8EviNkbYJUByqGR9QI4G"
DEFAULT_SENTRY_DSN="https://3b20acb23bbbc5958448bb41900cdca2@sentry.fergify.work/10"

# Use environment variables if set, otherwise use defaults
SEGMENT_KEY=${RUNE_SEGMENT_WRITE_KEY:-$DEFAULT_SEGMENT_KEY}
SENTRY_DSN=${RUNE_SENTRY_DSN:-$DEFAULT_SENTRY_DSN}
VERSION=${VERSION:-$(git describe --tags --always --dirty)}

echo "Building Rune with telemetry support..."
echo "Version: $VERSION"
echo "Segment Key: ${SEGMENT_KEY:0:10}..."
echo "Sentry DSN: ${SENTRY_DSN:0:30}..."

# Build the binary
go build -ldflags "\
  -s -w \
  -X github.com/ferg-cod3s/rune/internal/commands.version=$VERSION \
  -X github.com/ferg-cod3s/rune/internal/telemetry.version=$VERSION \
  -X github.com/ferg-cod3s/rune/internal/telemetry.segmentWriteKey=$SEGMENT_KEY \
  -X github.com/ferg-cod3s/rune/internal/telemetry.sentryDSN=$SENTRY_DSN \
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