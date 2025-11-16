#!/bin/bash

# NotesMD Installation Script
# This script installs NotesMD (nmd) to your system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}       NotesMD Installation Script${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "OS: ${GREEN}$OS${NC}"
echo -e "Architecture: ${GREEN}$ARCH${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    echo -e "${YELLOW}Please install Go from https://golang.org/dl/${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓ Go found: $GO_VERSION${NC}"
echo ""

# Installation directory
INSTALL_DIR="$HOME/bin"
if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
fi

echo -e "${YELLOW}Installing to: $INSTALL_DIR${NC}"
echo ""

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Clone and build
echo -e "${BLUE}📦 Downloading NotesMD...${NC}"
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

if ! git clone https://github.com/hrodwolf/notesmd.git; then
    echo -e "${RED}✗ Failed to clone repository${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
fi

cd notesmd

echo -e "${BLUE}🔨 Building NotesMD...${NC}"
if ! go build -o nmd ./cmd/notesmd; then
    echo -e "${RED}✗ Build failed${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo -e "${BLUE}📋 Installing binary to $INSTALL_DIR...${NC}"
mv nmd "$INSTALL_DIR/"

# Cleanup
cd
rm -rf "$TEMP_DIR"

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠ Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add this line to your shell configuration file (~/.bashrc, ~/.zshrc, etc.):"
    echo ""
    echo -e "${GREEN}export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
    echo ""
fi

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ NotesMD installed successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Run ${BLUE}nmd${NC} to get started!"
echo -e "Run ${BLUE}nmd ~/path/to/notes${NC} to open a specific directory"
echo ""
echo -e "For help, run ${BLUE}nmd${NC} and press ${YELLOW}?${NC}"
echo ""
