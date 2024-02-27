#!/bin/sh

# This script is intended to be invoked via
# curl -sSL https://install.audius.org | sh

set -e

# Determine architecture
ARCH=$(uname -m)
BINARY_NAME="audius-ctl-${ARCH}"
BINARY_URL="https://github.com/AudiusProject/audius-d/releases/latest/download/${BINARY_NAME}"

# Try to determine the target directory
if [ -w /usr/local/bin ]; then
    TARGET_DIR="/usr/local/bin"
elif echo "$PATH" | grep -q "$HOME/.local/bin"; then
    TARGET_DIR="$HOME/.local/bin"
elif echo "$PATH" | grep -q "$HOME/bin"; then
    TARGET_DIR="$HOME/bin"
else
    echo 'Insufficient permissions and/or no suitable directory found in $PATH.'
    echo 'Please manually add $HOME/.local/bin or $HOME/bin to your $PATH, then rerun this script.'
    exit 1
fi

# Create target directory if it doesn't exist
if [ ! -d "$TARGET_DIR" ]; then
    echo "Creating directory $TARGET_DIR"
    mkdir -p "$TARGET_DIR"
fi

# Download the binary
echo "Downloading ${BINARY_NAME} to ${TARGET_DIR}"
curl -sSL "$BINARY_URL" -o "${TARGET_DIR}/audius-ctl"
chmod +x "${TARGET_DIR}/audius-ctl"

echo "${BINARY_NAME} has been installed to ${TARGET_DIR}/audius-ctl"
echo "You can run it using: audius-ctl"

# Inform user about PATH addition if necessary
if ! echo ":$PATH:" | grep -q ":$TARGET_DIR:" ; then
    echo "To use audius-ctl from any location, add ${TARGET_DIR} to your PATH."
    echo "For bash users, add this line to your ~/.bash_profile or ~/.bashrc:"
    echo "export PATH=\"\$PATH:${TARGET_DIR}\""
    echo "For zsh users, add the line to your ~/.zshrc instead."
    echo "After adding the line, restart your terminal or run 'source <file>' on the modified file to update your current session."
fi
