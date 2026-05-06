#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Testing Presidium Offline Build${NC}"
echo -e "${YELLOW}========================================${NC}"

# Step 1: Prepare themes
echo -e "\n${YELLOW}Step 1: Preparing themes bundle...${NC}"
make prepare-themes
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Failed to prepare themes${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Themes prepared successfully${NC}"

# Step 2: Build Docker image with multi-stage build
echo -e "\n${YELLOW}Step 2: Building Docker test image (this may take a while)...${NC}"
# Create a temporary multi-stage Dockerfile
cat > Dockerfile.offline-test.tmp <<'EOF'
# Stage 1: Build the binary
FROM golang:1.25-alpine AS builder

# Install build dependencies for extended Hugo (LibSass)
RUN apk add --no-cache git gcc g++ musl-dev

WORKDIR /build

# Copy go module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and themes
COPY . .

# Build the binary with extended support
# Use cache mounts to avoid filling disk
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -tags extended -o presidium .

# Stage 2: Runtime image
FROM alpine:3.19

# Install runtime dependencies including C++ standard library for LibSass
RUN apk add --no-cache ca-certificates git libstdc++

# Create working directory
WORKDIR /test

# Copy the binary from builder stage
COPY --from=builder /build/presidium /usr/local/bin/presidium

# Copy docs directory
COPY docs /test/docs

# Make binary executable
RUN chmod +x /usr/local/bin/presidium

# Set working directory to docs
WORKDIR /test/docs

# Run the build command
CMD ["presidium", "hugo"]
EOF

# Enable Docker BuildKit for cache mount support
export DOCKER_BUILDKIT=1

docker build -f Dockerfile.offline-test.tmp -t presidium-offline-test .
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Failed to build Docker image${NC}"
    rm -f Dockerfile.offline-test.tmp
    exit 1
fi
rm -f Dockerfile.offline-test.tmp
echo -e "${GREEN}✓ Docker image built successfully${NC}"

# Step 3: Run container with no network access
echo -e "\n${YELLOW}Step 3: Running build in isolated container (no network)...${NC}"
docker run --rm \
    --network none \
    --name presidium-offline-test-run \
    presidium-offline-test

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Build failed in container${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Build succeeded without network access${NC}"

# Step 4: Verify output by running again and extracting files
echo -e "\n${YELLOW}Step 4: Verifying output structure...${NC}"
docker run --rm \
    --network none \
    --name presidium-offline-verify \
    presidium-offline-test \
    sh -c "presidium hugo --quiet && ls -la public/ | head -20"

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Failed to verify output${NC}"
    exit 1
fi

# Step 5: Check for expected theme files
echo -e "\n${YELLOW}Step 5: Checking for theme-specific files...${NC}"
EXPECTED_FILES=(
    "public/index.html"
    "public/presidium.js"
    "public/links.js"
    "public/assets"
    "public/images"
)

for file in "${EXPECTED_FILES[@]}"; do
    docker run --rm --network none presidium-offline-test sh -c "test -e $file"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Found: $file${NC}"
    else
        echo -e "${RED}✗ Missing: $file${NC}"
        exit 1
    fi
done

# Step 6: Verify no network calls were attempted
echo -e "\n${YELLOW}Step 6: Running with network monitoring...${NC}"
echo -e "${YELLOW}(This confirms no network access was attempted)${NC}"

# Run one more time and capture any network-related errors
OUTPUT=$(docker run --rm --network none presidium-offline-test 2>&1)
if echo "$OUTPUT" | grep -qi "network\|connection\|resolve\|lookup\|dial"; then
    echo -e "${RED}✗ Network access was attempted!${NC}"
    echo "$OUTPUT"
    exit 1
fi
echo -e "${GREEN}✓ No network access attempted${NC}"

# Success!
echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}✓ All offline tests passed!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "\nPresidium successfully built a Hugo site without any network access."
echo -e "This confirms that themes are properly embedded in the binary."

# Cleanup
echo -e "\n${YELLOW}Cleaning up Docker image and test artifacts...${NC}"
docker rmi presidium-offline-test >/dev/null 2>&1 || true
rm -f Dockerfile.offline-test.tmp
echo -e "${GREEN}✓ Cleanup complete${NC}"

exit 0
