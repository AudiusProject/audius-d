
# Development

Run Audius nodes and chains in a sandbox on your local machine.

### Contexts

Use contexts to experiment with different setups without clobbering changes.
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
audius-ctl config create-context devnet -f ./configs/templates/devnet.yaml
```

#### Build osx version with statusbar feature (Experimental)

```bash
make audius-ctl-arm-osx
```

### Devnet

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
