#!/bin/bash
# Test script to verify presidium works with embedded themes (no network access)

set -e

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║  Testing Presidium with Embedded Themes (Offline Mode)       ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Check if test site exists
if [ ! -d ".tmp/presidium-test-validation" ]; then
    echo "❌ Error: Test site not found at .tmp/presidium-test-validation"
    echo "   Please clone the test site first or provide a test site directory."
    exit 1
fi

echo "✓ Found test site directory"
echo ""

# Clean previous build artifacts
echo "🧹 Cleaning previous build artifacts..."
rm -rf .tmp/presidium-test-validation/public \
       .tmp/presidium-test-validation/.hugo_build.lock \
       .tmp/presidium-test-validation/resources

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
