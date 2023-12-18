# This script should only be run from the Makefile

set -eo pipefail

ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    BINARY_NAME="audius-ctl-arm"
else
    BINARY_NAME="audius-ctl-x86"
fi

if ! [ -f "bin/$BINARY_NAME" ]; then
    echo "No build artifact '$BINARY_NAME' in bin/"
    echo "Please run 'make' first"
    exit 1
fi

if [ -w /usr/local/bin ]; then
    TARGET_DIR="/usr/local/bin"
elif echo "$PATH" | grep -q "$HOME/.local/bin"; then
    TARGET_DIR="$HOME/.local/bin"
elif echo "$PATH" | grep -q "$HOME/bin"; then
    TARGET_DIR="$HOME/bin"
else
    echo 'Insufficient permissions and/or no suitable directory found in $PATH'
    echo 'Try `sudo make install` or add $HOME/.local/bin to your $PATH'
    exit 1
fi

cp "bin/$BINARY_NAME" "$TARGET_DIR/audius-ctl"

echo "$BINARY_NAME has been installed to $TARGET_DIR/audius-ctl"
echo "You can run it using: $ audius-ctl"
