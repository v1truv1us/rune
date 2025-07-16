#!/bin/bash

# Auto-start tmux if not already running inside one
SESSION_NAME="security-release-check"
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

# Pre-Release Security Verification Script
# This script performs comprehensive security checks before releasing Rune

set -e

echo "🔒 Pre-Release Security Checklist"
echo "=================================="

# Configuration
BINARY_PATH="./bin/rune"
REPORT_FILE="security-release-report.txt"

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

# Initialize report
init_report() {
    cat > "$REPORT_FILE" << EOF
Rune Pre-Release Security Report
===============================
Date: $(date)
Version: $(git describe --tags --always 2>/dev/null || echo "unknown")
Commit: $(git rev-parse HEAD 2>/dev/null || echo "unknown")

Security Checks:
EOF
}

# Update report
update_report() {
    echo "$1" >> "$REPORT_FILE"
}

# 1. Dependency Security Audit
check_dependencies() {
    log_info "1. Checking dependencies for vulnerabilities..."
    
    local temp_log=$(mktemp)
    local status=0
    
    # Check for known vulnerabilities
    if go install golang.org/x/vuln/cmd/govulncheck@latest > /dev/null 2>&1; then
        if govulncheck ./... > "$temp_log" 2>&1; then
            log_success "No known vulnerabilities found"
            update_report "✅ Dependency vulnerabilities: PASS"
        else
            log_error "Vulnerabilities found:"
            cat "$temp_log"
            update_report "❌ Dependency vulnerabilities: FAIL"
            status=1
        fi
    else
        log_warning "govulncheck not available"
        update_report "⚠️  Dependency vulnerabilities: SKIPPED"
    fi
    
    # Check with Nancy if available
    if command -v nancy >/dev/null 2>&1; then
        if go list -json -deps ./... | nancy sleuth --loud > "$temp_log" 2>&1; then
            log_success "Nancy dependency scan passed"
            update_report "✅ Nancy dependency scan: PASS"
        else
            log_warning "Nancy found potential issues:"
            cat "$temp_log"
            update_report "⚠️  Nancy dependency scan: WARNINGS"
        fi
    else
        log_info "Nancy not available, installing..."
        if go install github.com/sonatypecommunity/nancy@latest > /dev/null 2>&1; then
            if go list -json -deps ./... | nancy sleuth --loud > "$temp_log" 2>&1; then
                log_success "Nancy dependency scan passed"
                update_report "✅ Nancy dependency scan: PASS"
            else
                log_warning "Nancy found potential issues"
                update_report "⚠️  Nancy dependency scan: WARNINGS"
            fi
        else
            log_warning "Could not install Nancy"
            update_report "⚠️  Nancy dependency scan: SKIPPED"
        fi
    fi
    
    rm -f "$temp_log"
    return $status
}

# 2. Static Security Analysis
check_static_analysis() {
    log_info "2. Running static security analysis..."
    
    local temp_log=$(mktemp)
    local status=0
    
    # Install and run gosec
    if go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest > /dev/null 2>&1; then
        if gosec ./... > "$temp_log" 2>&1; then
            log_success "Static analysis passed"
            update_report "✅ Static analysis (gosec): PASS"
        else
            log_warning "Static analysis found potential issues:"
            cat "$temp_log"
            update_report "⚠️  Static analysis (gosec): WARNINGS"
        fi
    else
        log_error "Could not install gosec"
        update_report "❌ Static analysis (gosec): FAIL"
        status=1
    fi
    
    rm -f "$temp_log"
    return $status
}

# 3. Secret Scanning
check_secrets() {
    log_info "3. Scanning for secrets..."
    
    local temp_log=$(mktemp)
    local status=0
    
    # Check for common secret patterns in source code
    if grep -r -E "(password|secret|key|token)\s*[:=]\s*['\"][^'\"]{8,}['\"]" . \
        --exclude-dir=.git \
        --exclude-dir=bin \
        --exclude="*.log" \
        --exclude="$REPORT_FILE" > "$temp_log" 2>/dev/null; then
        
        log_warning "Potential secrets found in source code:"
        cat "$temp_log"
        update_report "⚠️  Source code secrets: WARNINGS"
    else
        log_success "No obvious secrets found in source code"
        update_report "✅ Source code secrets: PASS"
    fi
    
    # Check environment variables
    if env | grep -E "(RUNE_.*KEY|RUNE_.*DSN|RUNE_.*SECRET)" > "$temp_log" 2>/dev/null; then
        log_info "Found telemetry environment variables (expected):"
        cat "$temp_log"
        update_report "ℹ️  Environment variables: EXPECTED"
    else
        log_info "No telemetry environment variables found"
        update_report "ℹ️  Environment variables: NONE"
    fi
    
    rm -f "$temp_log"
    return $status
}

# 4. Binary Analysis
check_binary() {
    log_info "4. Analyzing binary for security issues..."
    
    local status=0
    
    if [ ! -f "$BINARY_PATH" ]; then
        log_error "Binary not found at $BINARY_PATH"
        log_info "Run 'make build' first"
        update_report "❌ Binary analysis: FAIL (binary not found)"
        return 1
    fi
    
    # Check binary for embedded secrets (excluding expected telemetry keys)
    local temp_log=$(mktemp)
    if strings "$BINARY_PATH" | grep -E "(password|secret|key|token)" | \
        grep -v -E "(segmentWriteKey|sentryDSN|RUNE_|keySize|keyHash|keyValue|keyShares)" > "$temp_log" 2>/dev/null; then
        
        log_warning "Potential secrets found in binary:"
        head -20 "$temp_log"  # Show first 20 matches
        if [ $(wc -l < "$temp_log") -gt 20 ]; then
            echo "... and $(( $(wc -l < "$temp_log") - 20 )) more"
        fi
        update_report "⚠️  Binary secrets: WARNINGS (review required)"
    else
        log_success "No unexpected secrets found in binary"
        update_report "✅ Binary secrets: PASS"
    fi
    
    # Check binary size (should be reasonable)
    local binary_size=$(stat -f%z "$BINARY_PATH" 2>/dev/null || stat -c%s "$BINARY_PATH" 2>/dev/null || echo "0")
    local size_mb=$((binary_size / 1024 / 1024))
    
    if [ $size_mb -gt 100 ]; then
        log_warning "Binary size is large: ${size_mb}MB"
        update_report "⚠️  Binary size: ${size_mb}MB (large)"
    else
        log_success "Binary size is reasonable: ${size_mb}MB"
        update_report "✅ Binary size: ${size_mb}MB"
    fi
    
    rm -f "$temp_log"
    return $status
}

# 5. Test Security
check_test_security() {
    log_info "5. Running security-focused tests..."
    
    local temp_log=$(mktemp)
    local status=0
    
    # Run tests with security tags
    if go test -v ./... -tags=security > "$temp_log" 2>&1; then
        log_success "Security tests passed"
        update_report "✅ Security tests: PASS"
    else
        log_warning "Some security tests failed:"
        cat "$temp_log"
        update_report "⚠️  Security tests: WARNINGS"
    fi
    
    # Run telemetry tests
    if go test -v ./internal/telemetry/ > "$temp_log" 2>&1; then
        log_success "Telemetry tests passed"
        update_report "✅ Telemetry tests: PASS"
    else
        log_error "Telemetry tests failed:"
        cat "$temp_log"
        update_report "❌ Telemetry tests: FAIL"
        status=1
    fi
    
    rm -f "$temp_log"
    return $status
}

# 6. Telemetry Security
check_telemetry_security() {
    log_info "6. Verifying telemetry security..."
    
    local temp_log=$(mktemp)
    local status=0
    
    # Test telemetry debug functionality
    if RUNE_DEBUG=true "$BINARY_PATH" debug keys > "$temp_log" 2>&1; then
        # Check that keys are properly masked
        if grep -q "\*\*\*\*" "$temp_log"; then
            log_success "API keys are properly masked in debug output"
            update_report "✅ Telemetry key masking: PASS"
        else
            log_error "API keys may not be properly masked"
            update_report "❌ Telemetry key masking: FAIL"
            status=1
        fi
        
        # Check that telemetry can be disabled
        if RUNE_TELEMETRY_DISABLED=true RUNE_DEBUG=true "$BINARY_PATH" status > "$temp_log" 2>&1; then
            if grep -q "Telemetry enabled: false" "$temp_log"; then
                log_success "Telemetry can be disabled"
                update_report "✅ Telemetry disable: PASS"
            else
                log_error "Telemetry disable may not work"
                update_report "❌ Telemetry disable: FAIL"
                status=1
            fi
        fi
    else
        log_error "Telemetry debug commands failed"
        update_report "❌ Telemetry debug: FAIL"
        status=1
    fi
    
    rm -f "$temp_log"
    return $status
}

# 7. Build Security
check_build_security() {
    log_info "7. Verifying build security..."
    
    local status=0
    
    # Check that build is reproducible (basic check)
    local build_info=$(go version -m "$BINARY_PATH" 2>/dev/null || echo "no build info")
    if echo "$build_info" | grep -q "go"; then
        log_success "Build info is available"
        update_report "✅ Build info: AVAILABLE"
    else
        log_warning "Build info not available"
        update_report "⚠️  Build info: NOT AVAILABLE"
    fi
    
    # Check Go version
    local go_version=$(go version | awk '{print $3}')
    log_info "Built with Go version: $go_version"
    update_report "ℹ️  Go version: $go_version"
    
    # Check for debug symbols (should be stripped in release)
    if file "$BINARY_PATH" | grep -q "not stripped"; then
        log_warning "Binary contains debug symbols"
        update_report "⚠️  Debug symbols: PRESENT"
    else
        log_success "Binary is stripped"
        update_report "✅ Debug symbols: STRIPPED"
    fi
    
    return $status
}

# 8. Generate SBOM
generate_sbom() {
    log_info "8. Generating Software Bill of Materials (SBOM)..."
    
    local sbom_file="rune.sbom.json"
    
    # Try to generate SBOM with syft if available
    if command -v syft >/dev/null 2>&1; then
        if syft "$BINARY_PATH" -o spdx-json > "$sbom_file" 2>/dev/null; then
            log_success "SBOM generated: $sbom_file"
            update_report "✅ SBOM generation: SUCCESS"
        else
            log_warning "SBOM generation failed"
            update_report "⚠️  SBOM generation: FAILED"
        fi
    else
        log_info "Syft not available, generating basic dependency list..."
        go list -m all > "dependencies.txt" 2>/dev/null || true
        log_success "Dependency list generated: dependencies.txt"
        update_report "✅ Dependency list: GENERATED"
    fi
}

# 9. Final Security Checklist
final_checklist() {
    log_info "9. Final security checklist..."
    
    update_report ""
    update_report "Final Security Checklist:"
    update_report "========================="
    
    local checklist_items=(
        "All tests pass"
        "No critical vulnerabilities"
        "No exposed secrets"
        "Telemetry is secure and optional"
        "Binary is properly stripped"
        "Dependencies are up to date"
        "SBOM/dependency list generated"
        "Release notes include security info"
    )
    
    for item in "${checklist_items[@]}"; do
        echo "  ☐ $item"
        update_report "☐ $item"
    done
    
    update_report ""
    update_report "Manual Verification Required:"
    update_report "============================"
    update_report "☐ Review all warnings above"
    update_report "☐ Verify telemetry events in dashboards"
    update_report "☐ Test with real API keys"
    update_report "☐ Verify Homebrew formula security"
    update_report "☐ Check GitHub release artifacts"
    update_report "☐ Update security documentation"
}

# Main execution
main() {
    echo "Starting pre-release security verification..."
    echo
    
    init_report
    
    local overall_status=0
    
    check_dependencies || overall_status=1
    echo
    
    check_static_analysis || overall_status=1
    echo
    
    check_secrets || overall_status=1
    echo
    
    check_binary || overall_status=1
    echo
    
    check_test_security || overall_status=1
    echo
    
    check_telemetry_security || overall_status=1
    echo
    
    check_build_security || overall_status=1
    echo
    
    generate_sbom
    echo
    
    final_checklist
    echo
    
    if [ $overall_status -eq 0 ]; then
        log_success "Pre-release security verification completed successfully!"
        update_report ""
        update_report "Overall Status: ✅ PASS"
    else
        log_warning "Pre-release security verification completed with warnings!"
        log_info "Review the issues above before releasing"
        update_report ""
        update_report "Overall Status: ⚠️  WARNINGS"
    fi
    
    log_info "Detailed report saved to: $REPORT_FILE"
    
    return $overall_status
}

# Run main function
main "$@"