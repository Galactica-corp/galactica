# Updating a Node with a New Version of Cosmos SDK

The update is performed at a specific block height, up to this height, it is not necessary to stop and restart the node. It will stop itself at the block height specified in the proposal.

## Stopping Nodes

To stop the node, run:

```bash
sudo kill -15 $(sudo lsof -t -i:1317)
```

Or, forcefully:

```bash
sudo kill -9 $(sudo lsof -t -i:1317)
```

If the node does not stop, you can obtain the API port from the `config/app.toml` file using the following command:

```bash
sed -n '/\[api\]/,/^\[/{/address/p}' .galactica/config/app.toml | sed 's/.*= //' | sed 's/.*://; s/"//g'
```

*Replace `.galactica` with your home directory if it's different, and replace `1317` with the port you obtained.*

## Adding the `block-executor` Parameter to `app.toml`

To add the `block-executor` parameter, use the following command:

```bash
sed -i '/^\[evm\]/a block-executor = "sequential"' .galactica/config/app.toml
```

*Replace `.galactica` with your home directory if it's different.*

## Restarting the Node

Restart the `galacticad` node.
