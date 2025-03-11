#!/bin/sh

# Get the directory of the script
SCRIPT_DIR=$(dirname "$(realpath "$0")")

# Target directory for the Hugo binary
TARGET_DIR="$SCRIPT_DIR/hugo-files"

# Build the latest Hugo with the extended version
CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest

mkdir -p "$TARGET_DIR"

# Copy the built Hugo binary to the script directory
cp "$(go env GOPATH)/bin/hugo" "$TARGET_DIR/hugo"

# Test the Hugo binary was built correctly
"$TARGET_DIR/hugo" version
