#!/bin/bash

# Test script to verify telemetry keys are properly embedded in builds
# This simulates what happens in GitHub Actions

set -e

echo "=== Testing Telemetry Build Process ==="

# Set test environment variables (using the same keys as build-with-telemetry.sh)
export RUNE_SEGMENT_WRITE_KEY="ZkEZXHRWH96y8EviNkbYJUByqGR9QI4G"
export RUNE_SENTRY_DSN="https://3b20acb23bbbc5958448bb41900cdca2@sentry.fergify.work/10"
export VERSION="test-build"

echo "Environment variables set:"
echo "  RUNE_SEGMENT_WRITE_KEY: ${RUNE_SEGMENT_WRITE_KEY:0:10}..."
echo "  RUNE_SENTRY_DSN: ${RUNE_SENTRY_DSN:0:30}..."
echo "  VERSION: $VERSION"

# Test GoReleaser build (snapshot mode)
echo ""
echo "=== Testing GoReleaser Build ==="
goreleaser build --snapshot --clean --single-target

# Find the built binary
BINARY_PATH=$(find dist/ -name "rune" -type f | head -1)

if [[ -n "$BINARY_PATH" ]]; then
    echo ""
    echo "=== Testing Built Binary ==="
    echo "Binary path: $BINARY_PATH"
    chmod +x "$BINARY_PATH"
    
    echo ""
    echo "Testing with debug output:"
    RUNE_DEBUG=true "$BINARY_PATH" --version
    
    echo ""
    echo "Testing basic functionality:"
    "$BINARY_PATH" --help > /dev/null && echo "✅ Binary works correctly"
else
    echo "❌ No binary found in dist/ directory"
    ls -la dist/ || echo "dist/ directory not found"
    exit 1
fi

echo ""
echo "=== Test Complete ==="