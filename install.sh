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
SKIP_VERIFICATION=false
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-homebrew)
            SKIP_HOMEBREW=true
            shift
            ;;
        --skip-verification)
            SKIP_VERIFICATION=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Rune CLI Installation Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --skip-homebrew       Skip Homebrew suggestion even if available"
            echo "  --skip-verification   Skip signature and checksum verification (NOT RECOMMENDED)"
            echo "  --verbose             Enable verbose output for verification process"
            echo "  -h, --help            Show this help message"
            echo ""
            echo "Security:"
            echo "  By default, this script verifies:"
            echo "  • Cosign signatures using keyless verification"
            echo "  • SHA256 checksums of downloaded binaries"
            echo "  • GitHub OIDC certificate validation"
            echo ""
            echo "Examples:"
            echo "  curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh"
            echo "  curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --skip-homebrew"
            echo "  curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --verbose"
            echo ""
            echo "Manual Verification:"
            echo "  If you prefer to verify manually, download the files and run:"
            echo "  • cosign verify-blob --certificate checksums.txt.pem --signature checksums.txt.sig checksums.txt"
            echo "  • sha256sum -c checksums.txt"
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

# Security utility functions
log_verbose() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${YELLOW}[VERBOSE]${NC} $1"
    fi
}

log_security() {
    echo -e "${GREEN}[SECURITY]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are available
check_security_dependencies() {
    local missing_tools=()
    
    if [ "$SKIP_VERIFICATION" = false ]; then
        # Check for cosign
        if ! command -v cosign >/dev/null 2>&1; then
            missing_tools+=("cosign")
        fi
        
        # Check for sha256sum or shasum
        if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
            missing_tools+=("sha256sum/shasum")
        fi
        
        if [ ${#missing_tools[@]} -gt 0 ]; then
            log_error "Missing required security tools: ${missing_tools[*]}"
            echo ""
            echo "To install cosign:"
            echo "  # On macOS with Homebrew:"
            echo "  brew install cosign"
            echo ""
            echo "  # On Linux:"
            echo "  # See: https://docs.sigstore.dev/system_config/installation/"
            echo ""
            echo "Alternatively, you can skip verification (NOT RECOMMENDED):"
            echo "  curl -fsSL https://raw.githubusercontent.com/ferg-cod3s/rune/main/install.sh | sh -s -- --skip-verification"
            echo ""
            exit 1
        fi
    fi
}

# Verify cosign signature
verify_signature() {
    local file="$1"
    local signature="$2"
    local certificate="$3"
    
    log_security "Verifying cosign signature for $file..."
    log_verbose "Signature file: $signature"
    log_verbose "Certificate file: $certificate"
    
    # Verify using cosign with keyless verification
    if ! cosign verify-blob \
        --certificate "$certificate" \
        --signature "$signature" \
        --certificate-identity-regexp "https://github.com/$REPO/.*" \
        --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
        "$file" >/dev/null 2>&1; then
        log_error "Signature verification failed for $file"
        log_error "This could indicate:"
        log_error "• The file has been tampered with"
        log_error "• The signature is invalid or corrupted"
        log_error "• A supply chain attack attempt"
        echo ""
        echo "For your security, installation has been aborted."
        echo "Please report this issue at: https://github.com/$REPO/issues"
        echo ""
        return 1
    fi
    
    log_security "✓ Signature verification passed"
    return 0
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local checksum_file="$2"
    
    log_security "Verifying checksum for $file..."
    
    # Extract expected checksum for our file
    local filename=$(basename "$file")
    local expected_checksum=""
    
    if command -v sha256sum >/dev/null 2>&1; then
        expected_checksum=$(grep "$filename" "$checksum_file" | cut -d' ' -f1)
        if [ -n "$expected_checksum" ]; then
            local actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)
        fi
    elif command -v shasum >/dev/null 2>&1; then
        expected_checksum=$(grep "$filename" "$checksum_file" | cut -d' ' -f1)
        if [ -n "$expected_checksum" ]; then
            local actual_checksum=$(shasum -a 256 "$file" | cut -d' ' -f1)
        fi
    else
        log_error "No checksum tool available"
        return 1
    fi
    
    if [ -z "$expected_checksum" ]; then
        log_error "Checksum not found for $filename in $checksum_file"
        return 1
    fi
    
    log_verbose "Expected: $expected_checksum"
    log_verbose "Actual:   $actual_checksum"
    
    if [ "$expected_checksum" != "$actual_checksum" ]; then
        log_error "Checksum verification failed for $file"
        log_error "Expected: $expected_checksum"
        log_error "Actual:   $actual_checksum"
        log_error "This indicates the file may have been corrupted or tampered with."
        echo ""
        echo "For your security, installation has been aborted."
        echo "Please report this issue at: https://github.com/$REPO/issues"
        echo ""
        return 1
    fi
    
    log_security "✓ Checksum verification passed"
    return 0
}

# Security warning and user consent
show_security_warning() {
    if [ "$SKIP_VERIFICATION" = true ]; then
        echo ""
        log_warning "SECURITY WARNING: Verification has been disabled!"
        log_warning "This script will download and execute binaries without verification."
        log_warning "This is NOT RECOMMENDED and could be dangerous."
        echo ""
        echo "Are you sure you want to continue without verification? [y/N]"
        read -r response
        case "$response" in
            [yY][eE][sS]|[yY]) 
                log_warning "Proceeding without verification as requested..."
                return 0
                ;;
            *)
                echo "Installation cancelled for security reasons."
                echo "Remove --skip-verification to enable secure installation."
                exit 1
                ;;
        esac
    else
        log_security "Security verification enabled"
        log_security "Will verify: signatures, checksums, and certificates"
    fi
}

echo -e "${GREEN}Installing Rune CLI...${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"

# Check dependencies and show security warning
check_security_dependencies
show_security_warning

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

# Construct download URLs
ARCHIVE_NAME="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$ARCHIVE_NAME"
CHECKSUMS_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/checksums.txt"
CHECKSUMS_SIG_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/checksums.txt.sig"
CHECKSUMS_CERT_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/checksums.txt.pem"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download verification files first
if [ "$SKIP_VERIFICATION" = false ]; then
    log_verbose "Downloading verification files..."
    
    echo -e "${YELLOW}Downloading checksums and signatures...${NC}"
    
    # Download checksums file
    if ! curl -sL "$CHECKSUMS_URL" -o checksums.txt; then
        log_error "Failed to download checksums.txt"
        exit 1
    fi
    
    # Download signature file
    if ! curl -sL "$CHECKSUMS_SIG_URL" -o checksums.txt.sig; then
        log_error "Failed to download checksums.txt.sig"
        exit 1
    fi
    
    # Download certificate file
    if ! curl -sL "$CHECKSUMS_CERT_URL" -o checksums.txt.pem; then
        log_error "Failed to download checksums.txt.pem"
        exit 1
    fi
    
    # Verify the checksums file signature first
    if ! verify_signature "checksums.txt" "checksums.txt.sig" "checksums.txt.pem"; then
        exit 1
    fi
fi

# Download and extract binary
echo -e "${YELLOW}Downloading $DOWNLOAD_URL...${NC}"
if ! curl -sL "$DOWNLOAD_URL" -o "$ARCHIVE_NAME"; then
    log_error "Failed to download $ARCHIVE_NAME"
    exit 1
fi

# Verify the archive checksum
if [ "$SKIP_VERIFICATION" = false ]; then
    if ! verify_checksum "$ARCHIVE_NAME" "checksums.txt"; then
        exit 1
    fi
fi

# Extract the archive
log_verbose "Extracting $ARCHIVE_NAME..."
if ! tar xzf "$ARCHIVE_NAME"; then
    log_error "Failed to extract $ARCHIVE_NAME"
    exit 1
fi

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
    
    # Show security summary
    if [ "$SKIP_VERIFICATION" = false ]; then
        echo ""
        log_security "Security verification completed successfully:"
        log_security "✓ Cosign signature verified using keyless verification"
        log_security "✓ SHA256 checksum verified"
        log_security "✓ GitHub OIDC certificate validated"
        echo ""
        echo "Your installation is secure and authentic."
    else
        echo ""
        log_warning "Installation completed WITHOUT security verification"
        log_warning "Future updates should use verification for security"
    fi
    
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