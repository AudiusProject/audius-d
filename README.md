# audius-d

Run your own node.

## Installation
### Latest release

Download from the [releases page](https://github.com/AudiusProject/audius-d/releases), **OR** run the following:

```bash
ARCH=x86 # Linux x86_64
ARCH=arm # MacOS
curl -L "https://github.com/AudiusProject/audius-d/releases/latest/download/audius-ctl-$ARCH" -o ~/.local/bin/audius-ctl && chmod +x ~/.local/bin/audius-ctl
```

### From build

```bash
make

make install  # installs to ~/.local/bin/

## OR ##

sudo make install  # installs to /usr/local/bin/
```

### Uninstall

```bash
sudo make uninstall
```

#### Build macos version with statusbar feature (Experimental)

```bash
make audius-ctl-arm-mac
```

## Quickstart

### Run a local devnet

Run Audius nodes and chains in a sandbox on your local machine.

Devnet uses a local nginx container on 80/443 to act as a layer 7 load balancer. Hence we need to add the hosts so we may intelligently route on localhost.
```
sudo sh -c 'echo "127.0.0.1       creator-1.devnet.audius-d discovery-1.devnet.audius-d identity.devnet.audius-d eth-ganache.devnet.audius-d acdc-ganache.devnet.audius-d solana-test-validator.devnet.audius-d" >> /etc/hosts'
```

Instruct audius-ctl what services to create and how to configure them. More on this concept below.
```
audius-ctl config create-context devnet -f configs/templates/devnet.yaml
```

Install the devnet certificate to avoid https warnings when connecting to local nodes
```
# MacOS
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain deployments/tls/devnet-cert.pem

# Ubuntu
sudo cp deployments/tls/devnet-cert.pem /usr/local/share/ca-certificates/devnet.audius-d.crt
sudo update-ca-certificates
```

Stand up audius nodes
```
audius-ctl up
```

Test context to verify it is all working.
```
audius-ctl status
...
https://creator-1.audius-d   [ /health_check .data.healthy                    ] true
https://discovery-1.audius-d [ /health_check .data.discovery_provider_healthy ] true
https://identity.audius-d    [ /health_check .healthy                         ] true
```

## Run

### Run installed binary

```bash
audius-ctl help
```

### Run built binary (without installation)

From git project directory:

```bash
# automatically builds and runs the correct binary for your system
./audius-ctl help

## OR ##

# Manually select binary after running `make`
bin/audius-ctl-x86 help  # linux
bin/audius-ctl-arm help  # mac
```

## Usage

### Create a content node 

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

### Switch between contexts

Contexts are modeled after `kubectl`. See:

```bash
audius-ctl config --help
```

Switch contexts

```bash
audius-ctl config use-context my-existing-context
```

Create new contexts

```bash
audius-ctl config create-context my-new-sandbox-context -f configs/templates/devnet.toml
```

Use contexts to experiment with different setups without clobbering changes.

### Migrate from audius-docker-compose to audius-d

Already running audius via [audius-docker-compose](https://github.com/AudiusProject/audius-docker-compose)?

```bash
audius-ctl config migrate-context my-new-migrated-context path/to/audius-docker-compose
```

## Releases

1. Increment version in audius-d/.version.json
1. Commit (and ideally push, review, land) changes
1. Ensure you are authenticated with the github cli (`gh auth status || gh auth login`)
1. Run `make release-audius-ctl`
1. Check the [releases page](https://github.com/AudiusProject/audius-d/releases)
