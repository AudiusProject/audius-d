# audius-d

run your own node.

## Installation
### Latest release

Download from the [releases page](https://github.com/AudiusProject/audius-d/releases), **OR** run the following:

```bash
ARCH=x86 # Linux x86_64
ARCH=arm # MacOS
gh release download -R https://github.com/AudiusProject/audius-d --clobber --output ~/.local/bin/audius-ctl --pattern audius-ctl-$ARCH
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
sh -c 'echo "127.0.0.1       creator-1.devnet.audius-d discovery-1.devnet.audius-d identity-1.devnet.audius-d eth-ganache.devnet.audius-d acdc-ganache.devnet.audius-d solana-test-validator.devnet.audius-d" >> /etc/hosts'
```

Instruct audius-ctl what services to create and how to configure them. More on this concept below.
```
audius-ctl config create-context devnet -f configs/templates/devnet.toml
```

Optionally, install the devnet certificate to avoid https warnings when connecting to local nodes
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

Visit local health checks to verify it is all working.
```
# TODO: audius-ctl config test-context 

curl -s https://creator-1.devnet.audius-d/health_check | jq .data.healthy
true

curl -s https://discovery-1.devnet.audius-d/health_check | jq .data.discovery_provider_healthy
true

curl -s https://identity-1.devnet.audius-d/health_check | jq .healthy
true
```

Note: if you choose not to install the devnet.audius-d certificate, use `curl -sk` 

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

> **Note:**
> Use of templates assumes you are in the git project directory

Using a template

```bash
audius-ctl config create-context my-creator-node -f configs/templates/operator.creator.toml
```

Or manually

```bash
audius-ctl config edit
```

Append the following config:

```toml
[CreatorNodes.creator-node]
InternalHttpPort = 80
ExternalHttpPort = 80
InternalHttpsPort = 443
ExternalHttpsPort = 443
Host = "http://localhost"
OperatorPrivateKey = ""      # <--- YOUR PRIVATE KEY HERE
OperatorWallet = ""          # <--- YOUR WALLET HERE
OperatorRewardsWallet = ""   # <--- YOUR WALLET HERE
```

Stand up the node

```bash
audius-ctl up
```

Tear down the node

```bash
audius-ctl down 
```

### Create a discovery node 

Using a brand new context

```bash
audius-ctl config create-context my-discovery-node -f configs/templates/operator.discovery.toml
audius-ctl up
```

OR Using an existing context

```bash
audius-ctl config edit
```

Add the following to the config:

```toml
[DiscoveryNodes.discovery-node]
InternalHttpPort = 5000
ExternalHttpPort = 5000
InternalHttpsPort = 5001
ExternalHttpsPort = 5001
Host = "http://localhost"
OperatorPrivateKey = ""     # <--- YOUR PRIVATE KEY HERE 
OperatorWallet = ""         # <--- YOUR WALLET HERE
OperatorRewardsWallet = ""  # <--- YOUR WALLET HERE
```

Stand up the node(s)

```bash
audius-ctl up
```

Tear down the node(s)

```bash
audius-ctl down 
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

### Run the gui

View transaction info in the browser

```bash
audius-ctl gui
```

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
