# Infra

Infrastructure automation is provided by [Pulumi](https://www.pulumi.com/) and accessible via the `infra` subcommand.

**NOTE**
The infra command is experimental and is subject to rapid change.
Currently only AWS and Cloudflare are supported.

```bash
$ audius-ctl infra --help

Manage audius-d instances

Usage:
  audius-ctl infra [command] [flags]
  audius-ctl infra [command]

Available Commands:
  cancel      Cancel the current in progress update
  destroy     Destroy the current context
  output      Get a specific output from the stack
  update      Update (deploy) the current context
```

### Configuration

Run `audius-ctl config edit`. Add the infra section to your context config.

```yaml
network:
  deployOn: mainnet
  infra:                             <------
    cloudflareAPIKey: my_cf_apikey   <------
    cloudflareZoneId: my_cf_zoneid   <------
    cloudflareTld: example.com       <------
    awsAccessKeyID: my_id            <------
    awsSecretAccessKey: my_secret    <------
    awsRegion: us-east-2             <------
nodes:
  audius-d-creator.example.com:
    type: creator
    version: prerelease
    privateKey: 21118f9a6de181061a2abd549511105adb4877cf9026f271092e6813b7cf58ab
    wallet: 0x0D38e653eC28bdea5A2296fD5940aaB2D0B8875c
    rewardsWallet: 0xb3c66e682Bf9a85F6800c769AC5A05c18C3F331d
```

### Deploy

To provision infrastructure for your context, run
```bash
audius-ctl infra update [-y]
```

To wire this up with audius-ctl (on your local machine), run the below script to add to your ssh config (defaults to `~/.ssh/config` if not specified).
```bash
./scripts/infra_ssh_config.sh ~/.ssh/config_personal

# also make the key readonly by the owner
chmod 0400 $(audius-ctl infra output instancePrivateKeyFilePath)
```

Now, startup the service
```bash
audius-ctl up
```

Visit your node
```bash
open https://$(audius-ctl infra output cloudflareRecordHostname)
```


### Cleanup

To destroy all associated infra with your context, run
```bash
audius-ctl infra destroy [-y]
```

If your update/destroy got stuck and into a bad state. You can try to recover with
```bash
audius-ctl infra cancel
```

### Backup

All infrastruture related state (including ssh keys to server instances) is stored by default in the `~/.audius` directory

```bash
ls -al ~/.audius/.pulumi/
```
