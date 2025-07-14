#!/bin/bash
# NimsForest Package Manager Universal Installation Script
# Usage: curl -fsSL get.nimsforest.com | sh
# This script detects the OS and runs the appropriate installer

set -e

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

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and redirect to appropriate installer
detect_and_install() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    case $os in
        linux)
            log "Detected Linux, running Linux installer..."
            curl -fsSL get.nimsforest.com/linux | sh
            ;;
        darwin)
            log "Detected macOS, running macOS installer..."
            curl -fsSL get.nimsforest.com/macos | sh
            ;;
        mingw*|msys*|cygwin*)
            log "Detected Windows (Git Bash/MSYS), please use PowerShell instead:"
            echo "  irm get.nimsforest.com/windows | iex"
            exit 1
            ;;
        *)
            error "Unsupported operating system: $os"
            ;;
    esac
}

# Main function
main() {
    log "NimsForest Package Manager Universal Installer"
    log "Detecting operating system..."
    
    detect_and_install
}

# Run if executed directly (not sourced)
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi