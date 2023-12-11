# audius-d

run your own node.

## Installation
### Latest release

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
audius-ctl config create-context my-new-sandbox-context -f config/templates/devnet.toml
```

Use contexts to experiment with different setups without clobbering changes.

### Development Sandbox

```bash
audius-ctl config use-context my-new-sandbox-context  # created in previous step
audius-ctl devnet  # start local eth and solana chains
audius-ctl up      # start creator, discovery, and identity nodes
```

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
