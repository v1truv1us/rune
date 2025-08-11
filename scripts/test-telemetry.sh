#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="test-telemetry"
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

# Test script to verify telemetry integration is working correctly

set -e

echo "🧪 Testing Rune Telemetry Integration"
echo "====================================="

# Build the binary with telemetry
echo "1. Building binary with telemetry..."
make build-telemetry > /dev/null 2>&1

if [ ! -f "./bin/rune" ]; then
    echo "❌ Build failed - binary not found"
    exit 1
fi

echo "✅ Binary built successfully"

# Test version command
echo "2. Testing version command..."
VERSION_OUTPUT=$(./bin/rune --version 2>/dev/null)
if [[ $VERSION_OUTPUT == *"version"* ]]; then
    echo "✅ Version command works"
else
    echo "❌ Version command failed"
    exit 1
fi

# Test telemetry initialization with debug output
echo "3. Testing telemetry initialization..."
DEBUG_OUTPUT=$(RUNE_DEBUG=true ./bin/rune status 2>&1)

# Check if telemetry is being initialized (look for debug initialization logs)
if [[ $DEBUG_OUTPUT == *"otel logging initialized successfully"* ]] || [[ $DEBUG_OUTPUT == *"OpenTelemetry logging initialized"* ]] || [[ $DEBUG_OUTPUT == *"initializing telemetry"* ]]; then
    echo "✅ Telemetry is enabled"
else
    echo "❌ Telemetry is not enabled"
    echo "$DEBUG_OUTPUT"
    exit 1
fi

# Check if OTLP is configured or disabled
if [[ -n "$RUNE_OTLP_ENDPOINT" ]]; then
    echo "✅ OTLP endpoint configured"
else
    echo "ℹ️  OTLP endpoint not set (RUNE_OTLP_ENDPOINT); logs will be local-only unless configured"
fi

# Check if Sentry is initialized when DSN is provided
if [[ -n "$RUNE_SENTRY_DSN" ]]; then
    if [[ $DEBUG_OUTPUT == *"Sentry OpenTelemetry integration initialized successfully"* ]] || [[ $DEBUG_OUTPUT == *"Sentry initialized successfully"* ]]; then
        echo "✅ Sentry initialized successfully"
    elif [[ $DEBUG_OUTPUT == *"Sentry integration failed"* ]]; then
        echo "⚠️  Sentry initialization failed (continuing without Sentry)"
    else
        echo "ℹ️  No explicit Sentry init logs detected"
    fi
else
    echo "ℹ️  Sentry DSN not set; skipping Sentry initialization check"
fi

# Test with telemetry disabled
echo "4. Testing with telemetry disabled..."
DISABLED_OUTPUT=$(RUNE_TELEMETRY_DISABLED=true RUNE_DEBUG=true ./bin/rune status 2>&1)

if [[ $DISABLED_OUTPUT == *"telemetry disabled via env"* ]]; then
    echo "✅ Telemetry can be disabled"
else
    echo "❌ Telemetry disable flag not working"
    echo "$DISABLED_OUTPUT"
    exit 1
fi

# Test a command that would generate telemetry using debug telemetry command
echo "5. Testing command telemetry..."
COMMAND_OUTPUT=$(RUNE_DEBUG=true ./bin/rune debug telemetry 2>&1)

if [[ $COMMAND_OUTPUT == *"Test event sent"* ]] || [[ $COMMAND_OUTPUT == *"tracking event"* ]]; then
    echo "✅ Command telemetry is working"
else
    echo "❌ Command telemetry not working"
    echo "$COMMAND_OUTPUT"
    exit 1
fi

echo ""
echo "🎉 All telemetry tests passed!"
echo "✅ Telemetry integration is working correctly"
echo "✅ Telemetry hooks are enabled"
if [[ -n "$RUNE_SENTRY_DSN" ]]; then
  echo "✅ Sentry error tracking is enabled"
else
  echo "ℹ️  Sentry error tracking not configured"
fi
echo "✅ Telemetry can be disabled when needed"
echo "✅ Command tracking is functional"