#!/bin/bash
set -e
curl -sSL "https://github.com/AudiusProject/audius-d/releases/latest/download/audius-ctl-$(uname -m)" -o ~/.local/bin/audius-ctl
chmod +x ~/.local/bin/audius-ctl
echo "audius-ctl has been installed to ~/.local/bin/audius-ctl"
echo "You can run it using: $ audius-ctl"
