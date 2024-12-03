# Cosmovisor for galactica

Cosmovisor is necessary for automatic node updates. Cosmovisor [docs](https://docs.cosmos.network/v0.50/build/tooling/cosmovisor)

## Installing cosmovizor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0
```

```bash
cosmovisor version
```

Set the required environment variables
In case that you you systemd to run galalacticad service make changes in unit file(example: Environment=DAEMON_NAME=galacticad)

* export DAEMON_NAME=galacticad (the name of the galactica binary)
* export DAEMON_HOME=$HOME/.galactica (home directory)
* export DAEMON_ALLOW_DOWNLOAD_BINARIES=true (enable auto-downloading of new binaries)
* export UNSAFE_SKIP_BACKUP=true (upgrades directly without performing a backup)

## Initialization cosmosvisor

This command creates the folder structure required for using cosmovisor.

```bash
cosmovisor init <path to binary file>
```

## Start cosmovisor

* If you use systemd you need to replace the name of the binary file with cosmovisor and add run command 

```bash
cosmovisor run \<galacticad start command\>
```
