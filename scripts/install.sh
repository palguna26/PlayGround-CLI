#!/bin/sh
# PlayGround CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
#
# This script installs the PlayGround CLI (pg) binary.
# It detects OS/arch, downloads the correct binary, verifies checksum, and installs.

set -e

# Configuration
REPO="palguna26/PlayGround-CLI"
BINARY_NAME="pg"
INSTALL_DIR=""

# Colors (only if terminal supports it)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

# Print functions
info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}[OK]${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$OS" in
        linux*)  OS="linux" ;;
        darwin*) OS="darwin" ;;
        mingw*|msys*|cygwin*) OS="windows" ;;
        *)       error "Unsupported operating system: $OS" ;;
    esac
    echo "$OS"
}

# Detect architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)  ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l)        ARCH="arm" ;;
        i386|i686)     ARCH="386" ;;
        *)             error "Unsupported architecture: $ARCH" ;;
    esac
    echo "$ARCH"
}

# Get latest release version from GitHub
get_latest_version() {
    RELEASE_URL="https://api.github.com/repos/${REPO}/releases/latest"
    
    RESPONSE=$(curl -fsSL "$RELEASE_URL" 2>/dev/null) || {
        error "Failed to fetch releases from GitHub.\n\nPossible causes:\n  - No releases exist yet (check GitHub Actions)\n  - Network connectivity issue\n  - Rate limited by GitHub API\n\nCheck: https://github.com/${REPO}/releases"
    }
    
    VERSION=$(echo "$RESPONSE" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        error "No releases found.\n\nThe release workflow may still be running.\nCheck: https://github.com/${REPO}/actions"
    fi
    echo "$VERSION"
}

# Determine install directory
get_install_dir() {
    # Try /usr/local/bin first (requires sudo on most systems)
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
        return
    fi
    
    # Fall back to ~/.local/bin
    LOCAL_BIN="$HOME/.local/bin"
    mkdir -p "$LOCAL_BIN" 2>/dev/null || true
    
    if [ -w "$LOCAL_BIN" ]; then
        echo "$LOCAL_BIN"
        return
    fi
    
    error "No writable install directory found. Try running with sudo."
}

# Download and verify binary
download_and_install() {
    OS="$1"
    ARCH="$2"
    VERSION="$3"
    INSTALL_DIR="$4"
    
    # Construct download URL
    ARCHIVE_NAME="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    info "Downloading PlayGround CLI ${VERSION} for ${OS}/${ARCH}..."
    
    # Download archive
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_NAME"; then
        error "Failed to download $DOWNLOAD_URL"
    fi
    
    # Download and verify checksum
    info "Verifying checksum..."
    if curl -fsSL "$CHECKSUM_URL" -o "$TMP_DIR/checksums.txt" 2>/dev/null; then
        EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "$TMP_DIR/checksums.txt" | awk '{print $1}')
        
        if [ -n "$EXPECTED_CHECKSUM" ]; then
            if command -v sha256sum >/dev/null 2>&1; then
                ACTUAL_CHECKSUM=$(sha256sum "$TMP_DIR/$ARCHIVE_NAME" | awk '{print $1}')
            elif command -v shasum >/dev/null 2>&1; then
                ACTUAL_CHECKSUM=$(shasum -a 256 "$TMP_DIR/$ARCHIVE_NAME" | awk '{print $1}')
            else
                warn "No checksum tool found, skipping verification"
                ACTUAL_CHECKSUM="$EXPECTED_CHECKSUM"
            fi
            
            if [ "$EXPECTED_CHECKSUM" != "$ACTUAL_CHECKSUM" ]; then
                error "Checksum verification failed!\nExpected: $EXPECTED_CHECKSUM\nActual: $ACTUAL_CHECKSUM"
            fi
            success "Checksum verified"
        else
            warn "Checksum not found for $ARCHIVE_NAME, skipping verification"
        fi
    else
        warn "Could not download checksums, skipping verification"
    fi
    
    # Extract archive
    info "Extracting..."
    tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"
    
    # Find the binary (might be in root or in a subdirectory)
    if [ -f "$TMP_DIR/$BINARY_NAME" ]; then
        BINARY_PATH="$TMP_DIR/$BINARY_NAME"
    elif [ -f "$TMP_DIR/${BINARY_NAME}_${OS}_${ARCH}/$BINARY_NAME" ]; then
        BINARY_PATH="$TMP_DIR/${BINARY_NAME}_${OS}_${ARCH}/$BINARY_NAME"
    else
        # Search for it
        BINARY_PATH=$(find "$TMP_DIR" -name "$BINARY_NAME" -type f | head -1)
        if [ -z "$BINARY_PATH" ]; then
            error "Binary not found in archive"
        fi
    fi
    
    # Make executable
    chmod +x "$BINARY_PATH"
    
    # Install
    info "Installing to ${INSTALL_DIR}..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    success "Installed $BINARY_NAME to $INSTALL_DIR/$BINARY_NAME"
}

# Check if binary is in PATH
check_path() {
    INSTALL_DIR="$1"
    
    case ":$PATH:" in
        *":$INSTALL_DIR:"*) 
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Main installation flow
main() {
    printf "\n"
    printf "${BLUE}╔════════════════════════════════════════╗${NC}\n"
    printf "${BLUE}║   PlayGround CLI Installer             ║${NC}\n"
    printf "${BLUE}╚════════════════════════════════════════╝${NC}\n"
    printf "\n"
    
    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        error "curl is required but not installed"
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        error "tar is required but not installed"
    fi
    
    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected: $OS/$ARCH"
    
    # Get latest version
    VERSION=$(get_latest_version)
    info "Latest version: $VERSION"
    
    # Get install directory
    INSTALL_DIR=$(get_install_dir)
    info "Install directory: $INSTALL_DIR"
    
    # Download and install
    download_and_install "$OS" "$ARCH" "$VERSION" "$INSTALL_DIR"
    
    # Verify installation
    if [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        INSTALLED_VERSION=$("$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        success "PlayGround CLI installed successfully!"
        printf "\n"
        printf "  Version: ${GREEN}%s${NC}\n" "$INSTALLED_VERSION"
        printf "  Binary:  ${GREEN}%s${NC}\n" "$INSTALL_DIR/$BINARY_NAME"
        printf "\n"
    else
        error "Installation verification failed"
    fi
    
    # Check PATH
    if ! check_path "$INSTALL_DIR"; then
        printf "${YELLOW}NOTE:${NC} $INSTALL_DIR is not in your PATH.\n"
        printf "Add it by running:\n"
        printf "\n"
        printf "  ${BLUE}export PATH=\"\$PATH:$INSTALL_DIR\"${NC}\n"
        printf "\n"
        printf "Or add this line to your ~/.bashrc or ~/.zshrc\n"
        printf "\n"
    fi
    
    # Quick start
    printf "Get started:\n"
    printf "  ${BLUE}pg setup${NC}         # Configure API key\n"
    printf "  ${BLUE}pg agent${NC}         # Start interactive mode\n"
    printf "  ${BLUE}pg --help${NC}        # Show all commands\n"
    printf "\n"
}

# Run main
main
