# SSH Configuration

All audius-ctl commands should be run from your *local machine.*

BUT you must have ssh access for audius-ctl to work.

```bash
ssh my.audius.node.example.com     # if this does not work, neither will audius-ctl
```

An example ssh config on your local machine might look like this:

```bash
# ~/.ssh/config

Host my.audius.node.example.com
  HostName 35.135.531.53           # <-- External IP address of the server, not the domain
  User ubuntu                      # <-- The host user whose ~/.ssh/authorized_keys file contains your PUBLIC key 
  Port 22
  IdentityFile ~/.ssh/my_rsa_key   # <-- Your local private key

```

### Tips
* Ensure that port 22 is open on your server's firewall.
* Ensure that your *public* key is added to the `~/.ssh/authorized_keys` file on the host under the home directory of the specified `User`.
* If your domain name is behind a proxy (e.g. Cloudflare), ensure the `HostName` matches the external IP address of your server, not the the proxy address.
* If using macOS and your ssh key has a passphrase, see [instructions for permanently adding it to keychain](https://apple.stackexchange.com/questions/48502/how-can-i-permanently-add-my-ssh-private-key-to-keychain-so-it-is-automatically)
