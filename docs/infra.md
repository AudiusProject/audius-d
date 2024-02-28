# Infra

Infrastructure automation is provided by [Pulumi](https://www.pulumi.com/) and accessible via the `infra` subcommand.

> **NOTE** the infra command is experimental and is subject to rapid change.

```
$ audius-ctl infra --help

Manage audius-d instances

Usage:
  audius-ctl infra [command] [flags]
  audius-ctl infra [command]

Available Commands:
  cancel      Cancel the current in progress update
  destroy     Destroy the current context
  output      Get a specific output from the stack
  preview     Perform a dry-run update for the current context
  ssh         SSH into the EC2 instance and execute commands
  update      Update (deploy) the current context
```
