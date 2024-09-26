# audius-d

# 🚨 Project Moved

> **Notice:** The project has moved to the [Audius Protocol monorepo ](https://github.com/AudiusProject/audius-protocol) as part of ongoing work to consolidate audius tooling.
> Visit [the Audius docs](https://docs.audius.org/node-operator/setup/installation) for up-to-date info on installation and usage.

![Moved](https://img.shields.io/badge/status-moved-red)


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
  content-1.example.com:
    type: content
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
audius-ctl down content-1.example.com
```

## Contributing

- [Development](./docs/development.md)
- [Releases](./docs/releases.md)
