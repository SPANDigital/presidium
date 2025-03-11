#!/bin/bash

# Define the packages you want to download - we can improve this by reading from a file
PACKAGES=(
    "github.com/spandigital/presidium-layouts-base@v0.2.1"
    "github.com/spandigital/presidium-styling-base@v0.2.2"
    "github.com/spandigital/presidium-layouts-base@v0.2.3"
    "github.com/spandigital/presidium-styling-base@v0.2.6"
)

# Get the directory of the script
SCRIPT_DIR=$(dirname "$(realpath "$0")")

# Create a temporary directory
TEMP_DIR=$(mktemp -d)

# Download the specified packages into the temporary directory
for PACKAGE in "${PACKAGES[@]}"; do
    GOPATH=$TEMP_DIR go mod download $PACKAGE
done

# Define the destination directory for the embedded packages relative to the script path
DEST_DIR="$SCRIPT_DIR/hugo-files"

# Create the destination directory if it doesn't exist
mkdir -p $DEST_DIR

# Copy the downloaded packages to the destination directory
cp -r $TEMP_DIR/pkg/mod/* $DEST_DIR

# Change ownership of the destination directory to the current user
#sudo chown -R $(whoami) $DEST_DIR

# Set the correct permissions for the destination directory
chmod -R 755 $DEST_DIR

# Add a dummy file to each empty directory to ensure embedded packages are included in the Hugo build
sudo find $DEST_DIR -type d -exec sh -c 'if [ -z "$(find "$1" -maxdepth 1 -type f)" ]; then touch "$1/empty"; fi' _ {} \;

echo "Specified packages have been downloaded and copied to the $DEST_DIR directory."

echo $SCRIPT_DIR