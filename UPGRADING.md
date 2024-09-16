# Upgrading Galactica
Updating a node with a new version of Cosmos SDK

* Stopping Nodes
```sudo kill -15 $(sudo lsof -t -i:1317)```
 or forcefully: ```sudo kill -9 $(sudo lsof -t -i:1317)```

* Adding Parameters to app.toml
```sed -i '/^\[evm\]/a block-executor = "sequential"' .galactica/config/app.toml```

* Restart the `galacticad` node.