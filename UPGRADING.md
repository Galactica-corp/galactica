# Upgrading Galactica
Updating a node with a new version of Cosmos SDK

1. Stopping Nodes
* ```sudo kill -15 $(sudo lsof -t -i:1317)```
 or forcefully: ```sudo kill -9 $(sudo lsof -t -i:1317)```
* If the node didn't stop, you can output the API port from config/app.toml:
```sed -n '/\[api\]/,/^\[/{/address/p}' .galactica/config/app.toml | sed 's/.*= //' | sed 's/.*://; s/"//g'```
Replace .galactica with your home directory if it's different And replace 1317 with the port you obtained

2. Adding Parameters to app.toml
```sed -i '/^\[evm\]/a block-executor = "sequential"' .galactica/config/app.toml```
Instead of .galactica, it could be another path to the node's home directory

3. Restart the `galacticad` node.
