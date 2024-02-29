#!/bin/bash

# Use the first command-line argument as the SSH config file path; default to ~/.ssh/config if not provided
SSH_CONFIG_PATH="${1:-$HOME/.ssh/config}"

# Retrieve the private key file path, hostname, and IP address
echo "Getting stack outputs for current audius-ctl context"
HOSTNAME=$(audius-ctl infra output cloudflareRecordHostname)
IDENTITY_FILE=$(audius-ctl infra output instancePrivateKeyFilePath)
IP=$(audius-ctl infra output instancePublicIp)

# Ensure HOSTNAME is not empty
if [ -z "$HOSTNAME" ]; then
  echo "HOSTNAME is empty, aborting."
  exit 1
fi

# Check if the entry already exists in the SSH config
if grep -q "^Host $HOSTNAME$" "$SSH_CONFIG_PATH"; then
  echo "Updating existing entry for $HOSTNAME"

  # Use awk to update the existing entry without affecting other entries
  awk -v hostname="$HOSTNAME" -v idfile="$IDENTITY_FILE" -v ip="$IP" '
    BEGIN { printit = 0; }
    /^Host / { printit = 0; }
    $0 ~ "^Host " hostname "$" { printit = 1; }
    printit { if (/^  IdentityFile /) $0 = "  IdentityFile " idfile; if (/^  HostName /) $0 = "  HostName " ip; }
    { print; }
    END {
      if (!printit) {
        print "\nHost " hostname;
        print "  IdentityFile " idfile;
        print "  HostName " ip;
        print "  User ubuntu";
        print "  IdentitiesOnly yes";
      }
    }
  ' "$SSH_CONFIG_PATH" > "${SSH_CONFIG_PATH}.tmp" && mv "${SSH_CONFIG_PATH}.tmp" "$SSH_CONFIG_PATH"
else
  echo "Adding new entry for $HOSTNAME"
  # Append new entry to the SSH config
  cat >> "$SSH_CONFIG_PATH" <<EOF

Host $HOSTNAME
  IdentityFile $IDENTITY_FILE
  HostName $IP
  User ubuntu
  IdentitiesOnly yes

EOF
fi

echo "SSH config updated successfully."
echo "Updated file: $SSH_CONFIG_PATH"
