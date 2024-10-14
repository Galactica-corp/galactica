# Cosmovisor for galactica

Cosmovisor is necessary for automatic node updates. Cosmovisor [docs](https://docs.cosmos.network/v0.50/build/tooling/cosmovisor)

## Installing cosmovizor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0
cosmovisor version
```

Set the required environment variables:

1. export DAEMON_NAME=galacticad
2. export DAEMON_HOME=$HOME/.galactica
3. export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
4. export UNSAFE_SKIP_BACKUP=true

## Initialization cosmosvisor:

This command creates the folder structure required for using cosmovisor.

```bash 
cosmovisor init <path to binary file>
```

## Start cosmovisor 

cosmovisor run \<galacticad start command\>