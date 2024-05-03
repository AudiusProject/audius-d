# audius-d

Run your own node.

## Installation

Install audius-ctl on your local computer

```bash
curl -sSL https://install.audius.org | sh
```

#### Uninstall

```bash
rm -f $(which audius-ctl)
```

## Run a Node

On your local computer

```bash
audius-ctl config edit
```

Write the following

```yaml
network:
  deployOn: mainnet
nodes:
  creator-1.example.com:
    type: creator
    privateKey: abc123          # <--- UNIQUE PRIV KEY USED BY THIS NODE TO SIGN RESPONSES
    wallet: 0xABC123            # <--- UNIQUE WALLET ADDRESS OF ABOVE PRIV KEY
    rewardsWallet: 0xABC123     # <--- ADDRESS OF WALLET HOLDING STAKED TOKENS
  discovery-1.example.com:
    type: discovery
    privateKey: abc123          # <--- UNIQUE PRIV KEY USED BY THIS NODE TO SIGN RESPONSES
    wallet: 0xABC123            # <--- UNIQUE WALLET ADDRESS OF ABOVE PRIV KEY
    rewardsWallet: 0xABC123     # <--- ADDRESS OF WALLET HOLDING STAKED TOKENS
```

You MUST:
* have docker installed on your server(s)
* have simple ssh access to your server(s) (see [SSH Configuration](./docs/ssh.md))

Stand up the node(s)

```bash
audius-ctl up
```

Restart a node

```bash
audius-ctl restart discovery-1.example.com
```

Tear down a node

```bash
audius-ctl down creator-1.example.com
```

## Contributing

- [Development](./docs/development.md)
- [Releases](./docs/releases.md)
