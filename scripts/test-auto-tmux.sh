#!/bin/bash

# Auto-tmux Functionality Test Suite
# This script tests the auto-tmux functionality added to all shell scripts

# Skip auto-tmux for this test script itself
export SKIP_AUTO_TMUX=1

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TEST_SESSION_PREFIX="test-auto-tmux"
TEST_SCRIPT_NAME="test-script.sh"
TEMP_DIR=""
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    ((TESTS_FAILED++))
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Cleanup function
cleanup() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
    
    # Kill any test tmux sessions
    tmux list-sessions 2>/dev/null | grep "^${TEST_SESSION_PREFIX}" | cut -d: -f1 | while read session; do
        tmux kill-session -t "$session" 2>/dev/null || true
    done
}

# Set up cleanup trap
trap cleanup EXIT

# Initialize test environment
setup_test_env() {
    log_info "Setting up test environment..."
    
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    log_success "Test environment created at: $TEMP_DIR"
}

# Create a test script with auto-tmux functionality
create_test_script() {
    local script_name="$1"
    local session_name="$2"
    local script_content="$3"
    
    cat > "$script_name" << EOF
#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="$session_name"
if [ -z "\$TMUX" ]; then
  # Check if tmux is available
  if command -v tmux >/dev/null 2>&1; then
    # If session exists, attach; otherwise create new session
    if tmux has-session -t "\$SESSION_NAME" 2>/dev/null; then
      echo "Attaching to existing tmux session: \$SESSION_NAME"
      exec tmux attach-session -t "\$SESSION_NAME"
    else
      echo "Starting new tmux session: \$SESSION_NAME"
      exec tmux new-session -s "\$SESSION_NAME" "\$0" "\$@"
    fi
  else
    # Fallback: check if running in an interactive terminal
    if [ ! -t 0 ]; then
      echo "This script must be run in an interactive terminal (tmux not available)."
      exit 1
    fi
  fi
fi

$script_content
EOF
    
    chmod +x "$script_name"
}

# Test 1: Check if tmux is available
test_tmux_availability() {
    log_info "Test 1: Checking tmux availability..."
    
    if command -v tmux >/dev/null 2>&1; then
        log_success "tmux is available for testing"
        return 0
    else
        log_error "tmux is not available - some tests will be skipped"
        return 1
    fi
}

# Test 2: Test session creation
test_session_creation() {
    log_info "Test 2: Testing tmux session creation..."
    
    if ! command -v tmux >/dev/null 2>&1; then
        log_warning "Skipping session creation test (tmux not available)"
        return 0
    fi
    
    local session_name="${TEST_SESSION_PREFIX}-creation"
    local script_name="test-creation.sh"
    
    # Create test script that exits immediately
    create_test_script "$script_name" "$session_name" "echo 'Test script executed'; exit 0"
    
    # Run the script in background and capture output
    local output
    output=$(timeout 10s bash -c "./$script_name" 2>&1 || true)
    
    if echo "$output" | grep -q "Starting new tmux session: $session_name"; then
        log_success "Session creation message detected"
    else
        log_error "Session creation message not found. Output: $output"
    fi
    
    # Clean up session
    tmux kill-session -t "$session_name" 2>/dev/null || true
}

# Test 3: Test session attachment
test_session_attachment() {
    log_info "Test 3: Testing tmux session attachment..."
    
    if ! command -v tmux >/dev/null 2>&1; then
        log_warning "Skipping session attachment test (tmux not available)"
        return 0
    fi
    
    local session_name="${TEST_SESSION_PREFIX}-attachment"
    local script_name="test-attachment.sh"
    
    # Create test script that waits
    create_test_script "$script_name" "$session_name" "echo 'Script running'; sleep 2; echo 'Script done'"
    
    # Create the session manually first
    tmux new-session -d -s "$session_name" "sleep 5"
    
    # Run the script - it should attach to existing session
    local output
    output=$(timeout 10s bash -c "./$script_name" 2>&1 || true)
    
    if echo "$output" | grep -q "Attaching to existing tmux session: $session_name"; then
        log_success "Session attachment message detected"
    else
        log_error "Session attachment message not found. Output: $output"
    fi
    
    # Clean up session
    tmux kill-session -t "$session_name" 2>/dev/null || true
}

# Test 4: Test fallback when tmux is not available
test_tmux_fallback() {
    log_info "Test 4: Testing fallback when tmux is not available..."
    
    local session_name="${TEST_SESSION_PREFIX}-fallback"
    local script_name="test-fallback.sh"
    
    # Create test script
    create_test_script "$script_name" "$session_name" "echo 'Fallback test executed'"
    
    # Mock tmux as unavailable by temporarily renaming it
    local tmux_path
    tmux_path=$(which tmux 2>/dev/null || echo "")
    
    if [ -n "$tmux_path" ]; then
        # Create a temporary PATH without tmux
        local temp_path
        temp_path=$(echo "$PATH" | tr ':' '\n' | grep -v "$(dirname "$tmux_path")" | tr '\n' ':' | sed 's/:$//')
        
        # Run script with modified PATH (simulating tmux not available)
        local output
        output=$(PATH="$temp_path" timeout 5s bash -c "./$script_name" 2>&1 || true)
        
        if echo "$output" | grep -q "Fallback test executed"; then
            log_success "Fallback execution works when tmux unavailable"
        else
            log_error "Fallback execution failed. Output: $output"
        fi
    else
        log_warning "tmux not found, testing fallback with TTY check"
        
        # Test TTY fallback
        local output
        output=$(echo "" | timeout 5s bash -c "./$script_name" 2>&1 || true)
        
        if echo "$output" | grep -q "This script must be run in an interactive terminal"; then
            log_success "TTY fallback error message works"
        else
            log_error "TTY fallback error message not found. Output: $output"
        fi
    fi
}

# Test 5: Test TTY detection
test_tty_detection() {
    log_info "Test 5: Testing TTY detection..."
    
    local session_name="${TEST_SESSION_PREFIX}-tty"
    local script_name="test-tty.sh"
    
    # Create test script
    create_test_script "$script_name" "$session_name" "echo 'TTY test executed'"
    
    # Test with non-interactive input (should fail)
    local output
    output=$(echo "" | timeout 5s bash -c "PATH=/nonexistent ./$script_name" 2>&1 || true)
    
    if echo "$output" | grep -q "This script must be run in an interactive terminal"; then
        log_success "TTY detection works correctly"
    else
        log_error "TTY detection failed. Output: $output"
    fi
}

# Test 6: Test argument passing
test_argument_passing() {
    log_info "Test 6: Testing argument passing through tmux..."
    
    if ! command -v tmux >/dev/null 2>&1; then
        log_warning "Skipping argument passing test (tmux not available)"
        return 0
    fi
    
    local session_name="${TEST_SESSION_PREFIX}-args"
    local script_name="test-args.sh"
    
    # Create test script that echoes arguments
    create_test_script "$script_name" "$session_name" 'echo "Arguments received: $@"; exit 0'
    
    # Run script with arguments
    local output
    output=$(timeout 10s bash -c "./$script_name arg1 arg2 'arg with spaces'" 2>&1 || true)
    
    if echo "$output" | grep -q "Arguments received: arg1 arg2 arg with spaces"; then
        log_success "Argument passing works correctly"
    else
        log_error "Argument passing failed. Output: $output"
    fi
    
    # Clean up session
    tmux kill-session -t "$session_name" 2>/dev/null || true
}

# Test 7: Test session isolation
test_session_isolation() {
    log_info "Test 7: Testing session isolation..."
    
    if ! command -v tmux >/dev/null 2>&1; then
        log_warning "Skipping session isolation test (tmux not available)"
        return 0
    fi
    
    local session1="${TEST_SESSION_PREFIX}-isolation1"
    local session2="${TEST_SESSION_PREFIX}-isolation2"
    local script1="test-isolation1.sh"
    local script2="test-isolation2.sh"
    
    # Create two different test scripts
    create_test_script "$script1" "$session1" "echo 'Session 1'; sleep 3"
    create_test_script "$script2" "$session2" "echo 'Session 2'; sleep 3"
    
    # Start both scripts in background
    timeout 10s bash -c "./$script1" &
    local pid1=$!
    sleep 1
    timeout 10s bash -c "./$script2" &
    local pid2=$!
    
    sleep 2
    
    # Check that both sessions exist
    local sessions
    sessions=$(tmux list-sessions 2>/dev/null | grep -E "(${session1}|${session2})" | wc -l)
    
    if [ "$sessions" -eq 2 ]; then
        log_success "Session isolation works - both sessions created"
    else
        log_error "Session isolation failed - expected 2 sessions, found $sessions"
    fi
    
    # Clean up
    kill $pid1 $pid2 2>/dev/null || true
    tmux kill-session -t "$session1" 2>/dev/null || true
    tmux kill-session -t "$session2" 2>/dev/null || true
}

# Test 8: Test real script integration
test_real_script_integration() {
    log_info "Test 8: Testing integration with real project scripts..."
    
    # Test one of the actual project scripts
    local test_script="../scripts/test-telemetry-build.sh"
    
    if [ -f "$test_script" ]; then
        # Check that the auto-tmux code is present
        if grep -q "Auto-start tmux if not already running inside one" "$test_script"; then
            log_success "Auto-tmux code found in real script"
        else
            log_error "Auto-tmux code not found in real script"
        fi
        
        # Check session name is set
        if grep -q 'SESSION_NAME=' "$test_script"; then
            log_success "Session name configuration found in real script"
        else
            log_error "Session name configuration not found in real script"
        fi
    else
        log_warning "Real script not found for integration test"
    fi
}

# Main test runner
run_tests() {
    log_info "Starting Auto-tmux Functionality Test Suite"
    echo "========================================"
    
    setup_test_env
    
    # Run all tests
    test_tmux_availability
    test_session_creation
    test_session_attachment
    test_tmux_fallback
    test_tty_detection
    test_argument_passing
    test_session_isolation
    test_real_script_integration
    
    echo ""
    echo "========================================"
    log_info "Test Results Summary"
    echo "========================================"
    
    if [ $TESTS_FAILED -eq 0 ]; then
        log_success "All tests passed! ($TESTS_PASSED/$((TESTS_PASSED + TESTS_FAILED)))"
        echo ""
        log_info "Auto-tmux functionality is working correctly"
        return 0
    else
        log_error "Some tests failed: $TESTS_FAILED failed, $TESTS_PASSED passed"
        echo ""
        log_info "Please review the failed tests above"
        return 1
    fi
}

# Run the tests
run_tests "$@"