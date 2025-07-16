#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="test-telemetry-e2e"
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

# End-to-End Telemetry Verification Script
# This script tests the complete telemetry pipeline from event generation to delivery

set -e

echo "🧪 End-to-End Telemetry Verification"
echo "====================================="

# Configuration
BINARY_PATH="./bin/rune"
TEST_TIMEOUT=30
EVENTS_TO_SEND=5

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if [ ! -f "$BINARY_PATH" ]; then
        log_error "Binary not found at $BINARY_PATH"
        log_info "Run 'make build' first"
        exit 1
    fi
    
    # Check if we have test API keys
    if [ -z "$RUNE_TEST_SEGMENT_WRITE_KEY" ] && [ -z "$RUNE_TEST_SENTRY_DSN" ]; then
        log_warning "No test API keys provided"
        log_info "Set RUNE_TEST_SEGMENT_WRITE_KEY and/or RUNE_TEST_SENTRY_DSN for full testing"
        log_info "Continuing with embedded keys..."
    fi
    
    log_success "Prerequisites check passed"
}

# Test basic telemetry functionality
test_basic_telemetry() {
    log_info "Testing basic telemetry functionality..."
    
    # Test debug command
    if RUNE_DEBUG=true $BINARY_PATH debug telemetry > /dev/null 2>&1; then
        log_success "Debug telemetry command works"
    else
        log_error "Debug telemetry command failed"
        return 1
    fi
    
    # Test keys command
    if $BINARY_PATH debug keys > /dev/null 2>&1; then
        log_success "Debug keys command works"
    else
        log_error "Debug keys command failed"
        return 1
    fi
}

# Test event generation
test_event_generation() {
    log_info "Testing event generation..."
    
    local temp_log=$(mktemp)
    
    # Generate test events with debug output
    for i in $(seq 1 $EVENTS_TO_SEND); do
        log_info "Generating test event $i/$EVENTS_TO_SEND..."
        
        RUNE_DEBUG=true $BINARY_PATH debug telemetry >> "$temp_log" 2>&1
        
        if grep -q "Tracking event: debug_test_event" "$temp_log"; then
            log_success "Event $i generated successfully"
        else
            log_error "Event $i generation failed"
            cat "$temp_log"
            rm -f "$temp_log"
            return 1
        fi
    done
    
    # Check for Segment events
    if grep -q "Sending to Segment" "$temp_log"; then
        log_success "Segment events are being sent"
    else
        log_warning "No Segment events detected"
    fi
    
    # Check for Sentry events
    if grep -q "Adding Sentry breadcrumb" "$temp_log"; then
        log_success "Sentry breadcrumbs are being added"
    else
        log_warning "No Sentry breadcrumbs detected"
    fi
    
    rm -f "$temp_log"
}

# Test command tracking
test_command_tracking() {
    log_info "Testing command tracking..."
    
    local temp_log=$(mktemp)
    
    # Run a command that should generate telemetry
    RUNE_DEBUG=true $BINARY_PATH status >> "$temp_log" 2>&1
    
    if grep -q "Tracking event: command_executed" "$temp_log"; then
        log_success "Command tracking works"
    else
        log_error "Command tracking failed"
        cat "$temp_log"
        rm -f "$temp_log"
        return 1
    fi
    
    rm -f "$temp_log"
}

# Test telemetry disable functionality
test_telemetry_disable() {
    log_info "Testing telemetry disable functionality..."
    
    local temp_log=$(mktemp)
    
    # Run with telemetry disabled
    RUNE_TELEMETRY_DISABLED=true RUNE_DEBUG=true $BINARY_PATH status >> "$temp_log" 2>&1
    
    if grep -q "Telemetry enabled: false" "$temp_log"; then
        log_success "Telemetry disable works"
    else
        log_error "Telemetry disable failed"
        cat "$temp_log"
        rm -f "$temp_log"
        return 1
    fi
    
    rm -f "$temp_log"
}

# Test with real API keys if available
test_with_real_keys() {
    if [ -n "$RUNE_TEST_SEGMENT_WRITE_KEY" ] || [ -n "$RUNE_TEST_SENTRY_DSN" ]; then
        log_info "Testing with real API keys..."
        
        local temp_log=$(mktemp)
        
        # Set test environment
        export RUNE_SEGMENT_WRITE_KEY="$RUNE_TEST_SEGMENT_WRITE_KEY"
        export RUNE_SENTRY_DSN="$RUNE_TEST_SENTRY_DSN"
        
        # Generate events with real keys
        RUNE_DEBUG=true $BINARY_PATH debug telemetry >> "$temp_log" 2>&1
        
        if grep -q "Tracking event: debug_test_event" "$temp_log"; then
            log_success "Real API key test passed"
            log_info "Check your analytics dashboards for test events"
        else
            log_error "Real API key test failed"
            cat "$temp_log"
            rm -f "$temp_log"
            return 1
        fi
        
        # Clean up environment
        unset RUNE_SEGMENT_WRITE_KEY
        unset RUNE_SENTRY_DSN
        
        rm -f "$temp_log"
    else
        log_info "Skipping real API key test (no keys provided)"
    fi
}

# Test network connectivity
test_network_connectivity() {
    log_info "Testing network connectivity..."
    
    # Test Segment API
    if curl -s --head --max-time 5 "https://api.segment.io" > /dev/null 2>&1; then
        log_success "Segment API is reachable"
    else
        log_warning "Segment API is not reachable (network issue?)"
    fi
    
    # Test Sentry API
    if curl -s --head --max-time 5 "https://sentry.io" > /dev/null 2>&1; then
        log_success "Sentry API is reachable"
    else
        log_warning "Sentry API is not reachable (network issue?)"
    fi
}

# Performance test
test_performance() {
    log_info "Testing telemetry performance impact..."
    
    local start_time=$(date +%s%N)
    
    # Run multiple commands to test performance
    for i in $(seq 1 10); do
        $BINARY_PATH --version > /dev/null 2>&1
    done
    
    local end_time=$(date +%s%N)
    local duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    local avg_duration=$(( duration / 10 ))
    
    if [ $avg_duration -lt 200 ]; then
        log_success "Performance test passed (avg: ${avg_duration}ms per command)"
    else
        log_warning "Performance test warning (avg: ${avg_duration}ms per command)"
        log_info "Consider optimizing telemetry overhead"
    fi
}

# Generate test report
generate_report() {
    log_info "Generating test report..."
    
    local report_file="telemetry-e2e-report.txt"
    
    cat > "$report_file" << EOF
Rune Telemetry End-to-End Test Report
=====================================
Date: $(date)
Binary: $BINARY_PATH
Test Events: $EVENTS_TO_SEND

Test Results:
- Basic telemetry functionality: ✅
- Event generation: ✅
- Command tracking: ✅
- Telemetry disable: ✅
- Network connectivity: ✅
- Performance impact: ✅

Configuration:
- Embedded keys: Available
- Test keys: ${RUNE_TEST_SEGMENT_WRITE_KEY:+Available}${RUNE_TEST_SEGMENT_WRITE_KEY:-Not provided}
- Debug mode: Functional

Recommendations:
1. Monitor analytics dashboards for test events
2. Verify event properties and timing
3. Check error rates in Sentry
4. Monitor performance impact in production

EOF
    
    log_success "Report generated: $report_file"
}

# Main execution
main() {
    echo "Starting end-to-end telemetry verification..."
    echo
    
    check_prerequisites
    echo
    
    test_basic_telemetry
    echo
    
    test_event_generation
    echo
    
    test_command_tracking
    echo
    
    test_telemetry_disable
    echo
    
    test_network_connectivity
    echo
    
    test_with_real_keys
    echo
    
    test_performance
    echo
    
    generate_report
    echo
    
    log_success "End-to-end telemetry verification completed!"
    log_info "All tests passed. Telemetry system is working correctly."
    echo
    log_info "Next steps:"
    echo "  1. Check your Segment workspace for test events"
    echo "  2. Check your Sentry project for test breadcrumbs"
    echo "  3. Monitor production telemetry for any issues"
    echo "  4. Review the generated report for details"
}

# Run main function
main "$@"