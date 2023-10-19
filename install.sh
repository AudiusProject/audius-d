#!/bin/sh

BASE_URL="https://raw.githubusercontent.com/AudiusProject/dot-slash-audius/main/bin"

ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    BINARY_NAME="audius-arm"
else
    BINARY_NAME="audius-x86"
fi

curl -LO "$BASE_URL/$BINARY_NAME"

chmod +x "$BINARY_NAME"

if echo "$PATH" | grep -q "$HOME/.local/bin"; then
    TARGET_DIR="$HOME/.local/bin"
elif echo "$PATH" | grep -q "$HOME/bin"; then
    TARGET_DIR="$HOME/bin"
else
    TARGET_DIR="/usr/local/bin"
fi

if [ "$TARGET_DIR" = "/usr/local/bin" ]; then
    sudo mv "$BINARY_NAME" "$TARGET_DIR/audius"
else
    mv "$BINARY_NAME" "$TARGET_DIR/audius"
fi

echo "$BINARY_NAME has been installed to $TARGET_DIR.\nYou can run it using: $ audius"
