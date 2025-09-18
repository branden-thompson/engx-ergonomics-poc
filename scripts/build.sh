#!/bin/bash

# EngX Ergonomics POC - Build Script
# Cross-platform build script for the engx CLI tool

set -e

# Configuration
BINARY_NAME="engx"
VERSION=${VERSION:-"dev"}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}
DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "🛩️ Building DPX Web Ergonomics POC"
echo "Version: ${VERSION}"
echo "Commit: ${COMMIT}"
echo "Date: ${DATE}"
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf dist/
mkdir -p dist/

# Build for current platform
echo "🔨 Building for current platform..."
go build -ldflags "${LDFLAGS}" -o "dist/${BINARY_NAME}" ./cmd/dpx-web

# Build for multiple platforms (optional)
if [[ "${CROSS_COMPILE}" == "true" ]]; then
    echo "🌍 Cross-compiling for multiple platforms..."

    # macOS
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "dist/${BINARY_NAME}-darwin-amd64" ./cmd/dpx-web
    GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o "dist/${BINARY_NAME}-darwin-arm64" ./cmd/dpx-web

    # Linux
    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "dist/${BINARY_NAME}-linux-amd64" ./cmd/dpx-web
    GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o "dist/${BINARY_NAME}-linux-arm64" ./cmd/dpx-web

    # Windows
    GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "dist/${BINARY_NAME}-windows-amd64.exe" ./cmd/dpx-web
fi

echo "✅ Build complete!"
echo "Binaries available in ./dist/"
ls -la dist/