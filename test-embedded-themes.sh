#!/bin/bash
# Test script to verify presidium works with embedded themes (no network access)

set -e

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║  Testing Presidium with Embedded Themes (Offline Mode)       ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Check if docs site exists
if [ ! -d "docs" ]; then
    echo "❌ Error: Documentation site not found at docs/"
    echo "   Please ensure the docs directory exists."
    exit 1
fi

echo "✓ Found docs directory"
echo ""

# Clean previous build artifacts
echo "🧹 Cleaning previous build artifacts..."
rm -rf docs/public \
       docs/.hugo_build.lock \
       docs/resources

# Build the Docker image (builds presidium binary and tests it)
echo "🐋 Building Docker image (compiling presidium for Linux)..."
docker build -f Dockerfile.test -t presidium-offline-test .

echo ""
echo "🔒 Running test with network disabled..."
echo ""

# Run the container with network disabled to prove themes work offline
# The --network=none flag ensures absolutely no network access
docker run --rm --network=none presidium-offline-test

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ SUCCESS! Presidium works with embedded themes (no network)  "
echo "═══════════════════════════════════════════════════════════════"
