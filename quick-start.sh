#!/bin/bash

# SnapAsk Quick Start Script
# This script helps you get started with SnapAsk development

set -e

echo "🚀 SnapAsk Quick Start Setup"
echo "=============================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
    fi
}

# Check prerequisites
echo -e "${BLUE}Checking prerequisites...${NC}"
echo ""

# Check Go
if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_status 0 "Go installed: $GO_VERSION"
else
    print_status 1 "Go not found. Please install Go 1.22+ from https://golang.org/dl/"
    exit 1
fi

# Check Node.js
if command_exists node; then
    NODE_VERSION=$(node --version)
    print_status 0 "Node.js installed: $NODE_VERSION"
else
    print_status 1 "Node.js not found. Please install Node.js 18+ from https://nodejs.org/"
    exit 1
fi

# Check npm
if command_exists npm; then
    NPM_VERSION=$(npm --version)
    print_status 0 "npm installed: $NPM_VERSION"
else
    print_status 1 "npm not found"
    exit 1
fi

# Check Wails
if command_exists wails; then
    WAILS_VERSION=$(wails version | head -n 1)
    print_status 0 "Wails installed: $WAILS_VERSION"
else
    echo -e "${YELLOW}⚠${NC} Wails CLI not found. Installing..."
    if go install github.com/wailsapp/wails/v2/cmd/wails@latest; then
        GO_BIN_PATH=$(go env GOPATH)/bin

        # Add Go bin to PATH for current session if needed
        if [[ ":$PATH:" != *":$GO_BIN_PATH:"* ]]; then
            export PATH="$GO_BIN_PATH:$PATH"
        fi

        # Verify installation and report status
        WAILS_VERSION=$(wails version | head -n 1)
        print_status 0 "Wails installed: $WAILS_VERSION"
    else
        print_status 1 "Failed to install Wails CLI automatically."
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}All prerequisites are installed. You're ready to develop with Korner!${NC}"