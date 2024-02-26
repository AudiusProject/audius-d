# audius-d

Run your own node.

## Installation

```bash
curl -sSL https://install.audius.org | sh
```

#### Uninstall

```bash
sudo make uninstall
```

## Run a Node

#### Content Node 

On your local computer

```bash
audius-ctl config edit
```

Write the following:

```yaml
network:
  deployOn: mainnet
nodes:
  my.domain.example.com:
    type: creator
    privateKey: abc123          # <--- YOUR PRIVATE KEY HERE
    wallet: 0xABC123            # <--- YOUR WALLET HERE
    rewardsWallet: 0xABC123     # <--- YOUR WALLET HERE
```

This assumes you own a server at my.domain.example.com which has your ssh key and docker installed.

Stand up the node

```bash
audius-ctl up
```

Tear down the node

```bash
audius-ctl down my.domain.example.com
```

## Migrate from audius-docker-compose

Already running audius via [audius-docker-compose](https://github.com/AudiusProject/audius-docker-compose)? Use the below to create an audius-ctl [context](./docs/development.md#contexts) based on your audius-docker-compose environment configuration.

```bash
audius-ctl config migrate-context my-new-migrated-context path/to/audius-docker-compose
```

## Contributing

- [Development](./docs/development.md)
- [Releases](./docs/releases.md)
