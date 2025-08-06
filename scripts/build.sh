#!/bin/bash

# Cross-platform build script for pulse-bridge
# This script builds the Go application for multiple platforms and architectures

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Application name and version
APP_NAME="pulse-bridge"
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Build directory
BUILD_DIR="builds"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Supported platforms and architectures
declare -a PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "linux/ppc64"
    "linux/ppc64le"
    "linux/mips64"
    "linux/mips64le"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
    "windows/386"
    "freebsd/amd64"
    "freebsd/arm64"
    "netbsd/amd64"
    "openbsd/amd64"
    "dragonfly/amd64"
    "solaris/amd64"
)

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to clean build directory
clean_builds() {
    print_status "Cleaning build directory..."
    rm -rf "${BUILD_DIR}"
    mkdir -p "${BUILD_DIR}"
}

# Function to build for a specific platform
build_platform() {
    local platform=$1
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)
    
    # Generate binary name with format: binary-name-version-os-arch[.exe]
    local output_name="${APP_NAME}-${VERSION}-${os}-${arch}"
    
    # Add .exe suffix for Windows builds
    if [ "$os" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${output_name}"
    
    print_status "Building for ${os}/${arch}..."
    
    # Create output directory
    mkdir -p "${BUILD_DIR}"
    
    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0
    
    # Build flags
    local ldflags="-s -w"
    ldflags="${ldflags} -X main.version=${VERSION}"
    ldflags="${ldflags} -X main.commit=${COMMIT}"
    ldflags="${ldflags} -X main.buildTime=${BUILD_TIME}"
    
    # Build the application
    if go build -ldflags "${ldflags}" -o "${output_path}" "${PROJECT_ROOT}/main.go"; then
        # Get file size
        local size=$(du -h "${output_path}" | cut -f1)
        print_success "Built ${os}/${arch} -> ${output_path} (${size})"
    else
        print_error "Failed to build for ${os}/${arch}"
        return 1
    fi
}

# Function to show build summary
show_summary() {
    print_status "Build Summary:"
    echo "=================="
    echo "Application: ${APP_NAME}"
    echo "Version: ${VERSION}"
    echo "Commit: ${COMMIT}"
    echo "Build Time: ${BUILD_TIME}"
    echo ""
    echo "Built binaries:"
    ls -la "${BUILD_DIR}"/ | grep -v '^d' | grep -v '^total' | while read -r line; do
        echo "  $line"
    done
}

# Function to validate Go environment
validate_environment() {
    print_status "Validating environment..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Go version
    local go_version=$(go version | cut -d' ' -f3)
    print_status "Using Go version: ${go_version}"
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the project root or ensure you're in the correct directory."
        exit 1
    fi
    
    print_success "Environment validation passed"
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Cross-platform build script for ${APP_NAME}"
    echo ""
    echo "Options:"
    echo "  -h, --help           Show this help message"
    echo "  -c, --clean          Clean build directory before building"
    echo "  --clean-only         Clean build directory and exit (no building)"
    echo "  -p, --platform OS/ARCH   Build for specific platform (e.g., linux/amd64)"
    echo "  -j, --jobs N         Number of parallel build jobs (default: all platforms in parallel)"
    echo "  -l, --list           List supported platforms"
    echo "  -v, --verbose        Enable verbose output"
    echo ""
    echo "Examples:"
    echo "  $0                   Build for all supported platforms in parallel"
    echo "  $0 -p linux/amd64   Build only for Linux AMD64"
    echo "  $0 -c                Clean and build for all platforms"
    echo "  $0 -j 4              Build with maximum 4 parallel jobs"
    echo "  $0 --clean-only      Clean build directory only"
}

# Function to list supported platforms
list_platforms() {
    echo "Supported platforms:"
    for platform in "${PLATFORMS[@]}"; do
        echo "  $platform"
    done
}

# Main function
main() {
    local clean_flag=false
    local clean_only_flag=false
    local specific_platform=""
    local verbose_flag=false
    local max_jobs=0  # 0 means unlimited (all platforms in parallel)
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                clean_flag=true
                shift
                ;;
            --clean-only)
                clean_only_flag=true
                shift
                ;;
            -p|--platform)
                specific_platform="$2"
                shift 2
                ;;
            -j|--jobs)
                max_jobs="$2"
                if ! [[ "$max_jobs" =~ ^[0-9]+$ ]] || [ "$max_jobs" -lt 1 ]; then
                    print_error "Invalid number of jobs: $max_jobs (must be a positive integer)"
                    exit 1
                fi
                shift 2
                ;;
            -l|--list)
                list_platforms
                exit 0
                ;;
            -v|--verbose)
                verbose_flag=true
                set -x
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Change to project root
    cd "${PROJECT_ROOT}"
    
    # Handle clean-only option
    if [ "$clean_only_flag" = true ]; then
        print_status "Cleaning build directory only..."
        clean_builds
        print_success "Build directory cleaned!"
        exit 0
    fi
    
    # Validate environment
    validate_environment
    
    # Clean if requested
    if [ "$clean_flag" = true ]; then
        clean_builds
    else
        mkdir -p "${BUILD_DIR}"
    fi
    
    print_status "Starting cross-platform build for ${APP_NAME}..."
    echo "Version: ${VERSION}"
    echo "Commit: ${COMMIT}"
    echo ""
    
    # Build for specific platform or all platforms
    if [ -n "$specific_platform" ]; then
        # Validate specific platform
        local platform_found=false
        for platform in "${PLATFORMS[@]}"; do
            if [ "$platform" = "$specific_platform" ]; then
                platform_found=true
                break
            fi
        done
        
        if [ "$platform_found" = false ]; then
            print_error "Unsupported platform: $specific_platform"
            print_status "Use -l to list supported platforms"
            exit 1
        fi
        
        build_platform "$specific_platform"
    else
        # Build for all platforms in parallel
        if [ "$max_jobs" -eq 0 ]; then
            print_status "Starting parallel builds for all platforms..."
            local pids=()
            local failed_builds=0
            
            # Start all builds in parallel (unlimited)
            for platform in "${PLATFORMS[@]}"; do
                build_platform "$platform" &
                pids+=($!)
            done
            
            # Wait for all builds to complete and check results
            for i in "${!pids[@]}"; do
                local pid=${pids[$i]}
                local platform=${PLATFORMS[$i]}
                
                if wait "$pid"; then
                    print_status "✓ ${platform} build completed successfully"
                else
                    print_error "✗ ${platform} build failed"
                    ((failed_builds++))
                fi
            done
        else
            print_status "Starting parallel builds with maximum $max_jobs jobs..."
            local pids=()
            local platform_names=()
            local failed_builds=0
            local platform_index=0
            
            # Build with job limit
            while [ $platform_index -lt ${#PLATFORMS[@]} ]; do
                # Start jobs up to the limit
                while [ ${#pids[@]} -lt $max_jobs ] && [ $platform_index -lt ${#PLATFORMS[@]} ]; do
                    local platform=${PLATFORMS[$platform_index]}
                    build_platform "$platform" &
                    pids+=($!)
                    platform_names+=("$platform")
                    ((platform_index++))
                done
                
                # Wait for at least one job to complete
                if [ ${#pids[@]} -gt 0 ]; then
                    local completed_pid=${pids[0]}
                    local completed_platform=${platform_names[0]}
                    
                    if wait "$completed_pid"; then
                        print_status "✓ ${completed_platform} build completed successfully"
                    else
                        print_error "✗ ${completed_platform} build failed"
                        ((failed_builds++))
                    fi
                    
                    # Remove completed job from tracking
                    pids=("${pids[@]:1}")
                    platform_names=("${platform_names[@]:1}")
                fi
            done
            
            # Wait for remaining jobs
            for i in "${!pids[@]}"; do
                local pid=${pids[$i]}
                local platform=${platform_names[$i]}
                
                if wait "$pid"; then
                    print_status "✓ ${platform} build completed successfully"
                else
                    print_error "✗ ${platform} build failed"
                    ((failed_builds++))
                fi
            done
        fi
        
        if [ $failed_builds -gt 0 ]; then
            print_warning "${failed_builds} builds failed"
        fi
    fi
    
    # Show summary
    echo ""
    show_summary
    
    print_success "Cross-platform build completed!"
}

# Run main function with all arguments
main "$@"