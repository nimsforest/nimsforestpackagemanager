#!/bin/bash
# NimsForest Package Manager Installation Script for macOS
# Usage: curl -fsSL get.nimsforest.com/macos | sh

set -e

# Configuration
BINARY_NAME="nimsforestpm"
GITHUB_REPO="nimsforest/nimsforestpackagemanager"
INSTALL_DIR="${NIMSFOREST_INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Detect architecture
detect_arch() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            error "Unsupported architecture: $arch"
            ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "Failed to get latest version from GitHub"
    fi
    echo "$version"
}

# Download and install binary
install_binary() {
    local version="$1"
    local arch="$2"
    
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${BINARY_NAME}_darwin_${arch}"
    local temp_file="/tmp/${BINARY_NAME}"
    
    log "Downloading ${BINARY_NAME} ${version} for macOS/${arch}..."
    
    if ! curl -fsSL "$download_url" -o "$temp_file"; then
        error "Failed to download binary from $download_url"
    fi
    
    # Make executable
    chmod +x "$temp_file"
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log "Creating install directory: $INSTALL_DIR"
        if ! mkdir -p "$INSTALL_DIR" 2>/dev/null; then
            if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
                warn "Cannot write to $INSTALL_DIR, trying with sudo..."
                sudo mkdir -p "$INSTALL_DIR"
                sudo mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
            else
                error "Cannot create install directory: $INSTALL_DIR"
            fi
        else
            mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
        fi
    else
        # Move binary to install directory
        if ! mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME" 2>/dev/null; then
            if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
                warn "Cannot write to $INSTALL_DIR, trying with sudo..."
                sudo mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
            else
                error "Cannot write to install directory: $INSTALL_DIR"
            fi
        fi
    fi
    
    success "Installed ${BINARY_NAME} to $INSTALL_DIR/$BINARY_NAME"
}

# Check for Gatekeeper and suggest workaround
check_gatekeeper() {
    if [ "$(uname -s)" = "Darwin" ]; then
        log "Note: If macOS Gatekeeper blocks execution, run:"
        log "  xattr -d com.apple.quarantine $INSTALL_DIR/$BINARY_NAME"
        log "  Or go to System Preferences > Security & Privacy and allow the app"
    fi
}

# Verify installation
verify_installation() {
    log "Verifying installation..."
    
    if ! command -v "$BINARY_NAME" >/dev/null 2>&1; then
        warn "$BINARY_NAME not found in PATH. You may need to add $INSTALL_DIR to your PATH."
        echo "Add this to your shell profile (.bashrc, .zshrc, etc.):"
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
        echo "Or run with full path: $INSTALL_DIR/$BINARY_NAME"
    else
        log "Running system check..."
        "$BINARY_NAME" hello || {
            warn "System check failed. This might be due to Gatekeeper restrictions."
            check_gatekeeper
        }
    fi
}

# Main installation flow
main() {
    log "Installing NimsForest Package Manager for macOS..."
    
    # Check prerequisites
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed"
    fi
    
    # Detect system
    local arch version
    arch=$(detect_arch)
    version=$(get_latest_version)
    
    log "Detected system: macOS/${arch}"
    log "Latest version: $version"
    
    # Install
    install_binary "$version" "$arch"
    
    # Check for Gatekeeper
    check_gatekeeper
    
    # Verify
    verify_installation
    
    success "Installation complete!"
    echo ""
    echo "Get started with:"
    echo "  nimsforestpm hello"
    echo "  nimsforestpm create-organization-workspace my-org"
    echo ""
    echo "For help: nimsforestpm --help"
}

# Run main function
main "$@"