#!/bin/bash
set -euo pipefail

# Configuration
REPO="wavezync/pulse-bridge"
BINARY_NAME="pulse-bridge"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

get_arch() {
    case $(uname -m) in
        x86_64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        i386|i686) echo "386" ;;
        armv7l) echo "arm" ;;
        *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
    esac
}

get_os() {
    case $(uname -s | tr '[:upper:]' '[:lower:]') in
        linux) echo "linux" ;;
        darwin) echo "darwin" ;;
        freebsd) echo "freebsd" ;;
        openbsd) echo "openbsd" ;;
        netbsd) echo "netbsd" ;;
        dragonfly) echo "dragonfly" ;;
        solaris) echo "solaris" ;;
        mingw*|msys*|cygwin*) echo "windows" ;;
        *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
    esac
}

main() {
    local os arch binary_name download_url temp_dir
    
    os=$(get_os)
    arch=$(get_arch)
    binary_name="$BINARY_NAME"
    
    echo "🔍 Getting latest release info..."
    
    # Get latest release
    local latest_release
    latest_release=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
                    grep '"tag_name":' | \
                    sed -E 's/.*"([^"]+)".*/\1/')
    
    if [[ -z "$latest_release" ]]; then
        echo "❌ Could not find latest release" >&2
        exit 1
    fi
    
    # Construct filename based on your naming convention: pulse-bridge-{version}-{os}-{arch}
    local filename="${BINARY_NAME}-${latest_release}-${os}-${arch}"
    
    if [[ "$os" == "windows" ]]; then
        filename="${filename}.exe"
        binary_name="${binary_name}.exe"
    fi
    
    download_url="https://github.com/$REPO/releases/download/$latest_release/$filename"
    
    echo "📦 Installing $BINARY_NAME $latest_release for $os/$arch..."
    echo "🔗 Downloading: $filename"
    
    # Create temporary directory
    temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT
    
    # Download
    if ! curl -L --fail --progress-bar "$download_url" -o "$temp_dir/$binary_name"; then
        echo "❌ Failed to download binary from: $download_url"
        echo "📋 Available platforms:"
        echo "   Linux: amd64, arm64, 386, mips64, mips64le, ppc64, ppc64le"
        echo "   macOS: amd64, arm64"
        echo "   Windows: amd64, arm64, 386"
        echo "   BSD: FreeBSD, OpenBSD, NetBSD, DragonFly"
        echo "   Solaris: amd64"
        exit 1
    fi
    
    chmod +x "$temp_dir/$binary_name"
    
    # Install
    echo "📍 Installing to $INSTALL_DIR..."
    if [[ ! -d "$INSTALL_DIR" ]]; then
        mkdir -p "$INSTALL_DIR"
    fi
    
    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$temp_dir/$binary_name" "$INSTALL_DIR/"
    else
        echo "🔐 Requesting sudo access to install to $INSTALL_DIR..."
        sudo mv "$temp_dir/$binary_name" "$INSTALL_DIR/"
    fi
    
    echo ""
    echo "🎉 $BINARY_NAME installed successfully!"
    echo "📍 Location: $INSTALL_DIR/$BINARY_NAME"
    echo "🚀 Run '$BINARY_NAME --help' to get started"
    echo ""
    echo "📖 Documentation: https://github.com/$REPO"
}

main "$@"