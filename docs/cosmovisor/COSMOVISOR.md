# Cosmovisor for galactica

Cosmovisor is necessary for automatic node updates. Cosmovisor [docs](https://docs.cosmos.network/v0.50/build/tooling/cosmovisor)

## Installing cosmovizor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0
```

```bash
cosmovisor version
```

Set the required environment variables:

* export DAEMON_NAME=galacticad (the name of the galactica binary)
* export DAEMON_HOME=$HOME/.galactica (home directory)
* export DAEMON_ALLOW_DOWNLOAD_BINARIES=true (enable auto-downloading of new binaries)
* export UNSAFE_SKIP_BACKUP=true (upgrades directly without performing a backup)

## Initialization cosmosvisor:

This command creates the folder structure required for using cosmovisor.

```bash 
cosmovisor init <path to binary file>
```

## Start cosmovisor 

cosmovisor run \<galacticad start command\>