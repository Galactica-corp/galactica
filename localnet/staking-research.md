# Staking research

## Genesis generation

First of all we need to generate genesis file for our network.
I will generate it for 3 validators and 1 treasury address.

```bash
export MAIN_PATH_HOME=./.galactica
export MAIN_PATH_CONFIG=$MAIN_PATH_HOME/config
export CHAIN_ID=galactica_9000-1
export KEYRING_BACKEND=test
alias gala="./build/galacticad --home $MAIN_PATH_HOME"
```

### Init test keyring backend and chain id:
```bash
gala config keyring-backend $KEYRING_BACKEND
gala config chain-id $CHAIN_ID
```

### Add keys for validators and treasury:
```bash
gala keys add validator01 --algo secp256k1 --keyring-backend $KEYRING_BACKEND --keyring-dir $MAIN_PATH_HOME
gala keys add validator02 --algo secp256k1 --keyring-backend $KEYRING_BACKEND --keyring-dir $MAIN_PATH_HOME
gala keys add validator03 --algo secp256k1 --keyring-backend $KEYRING_BACKEND --keyring-dir $MAIN_PATH_HOME
gala keys add treasury --algo secp256k1 --keyring-backend $KEYRING_BACKEND --keyring-dir $MAIN_PATH_HOME
```

## Init chain

```bash
gala init localtestnet --chain-id $CHAIN_ID --default-denom gnet
rm $MAIN_PATH_CONFIG/node_key.json
rm $MAIN_PATH_CONFIG/priv_validator_key.json
```

## Change network parameters

### Change config.toml and app.toml

```bash
sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g'  $MAIN_PATH_CONFIG/config.toml
sed -i '' 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26658"/g'  $MAIN_PATH_CONFIG/config.toml
sed -i '' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/g'  $MAIN_PATH_CONFIG/config.toml

sed -i '' '/\[api\]/,+3 s/enable = false/enable = true/' $MAIN_PATH_CONFIG/app.toml
sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' $MAIN_PATH_CONFIG/app.toml
sed -i '' 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' $MAIN_PATH_CONFIG/app.toml
sed -i '' 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' $MAIN_PATH_CONFIG/app.toml
sed -i '' '/\[grpc-web\]/,+3 s/address = "localhost:9091"/address = "0.0.0.0:9091"/' $MAIN_PATH_CONFIG/app.toml
sed -i '' 's/pruning = "default"/pruning = "nothing"/g'  $MAIN_PATH_CONFIG/app.toml
sed -i '' 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "10ugnet"/g'  $MAIN_PATH_CONFIG/app.toml
```

### change genesis.json

```bash
jq '.app_state.gov.voting_params.voting_period = "600s"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["staking"]["params"]["bond_denom"]="ugnet"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["staking"]["params"]["unbonding_time"]="30s"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["crisis"]["constant_fee"]["denom"]="ugnet"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="ugnet"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["gov"]["params"]["min_deposit"][0]={"denom":"ugnet","amount":"1000000000"}' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["gov"]["params"]["max_deposit_period"]="600s"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["gov"]["params"]["voting_period"]="600s"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["mint"]["params"]["mint_denom"]="ugnet"' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["bank"]["denom_metadata"][0]={"description": "The native staking token of the Galactica Network.","denom_units": [{"denom": "ugnet","exponent": 0,"aliases": ["micrognet"]},{"denom": "gnet","exponent": 6}],"base": "ugnet","display": "gnet","name": "Galactica Network","symbol": "GNET","uri": "","uri_hash": ""}' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["bank"]["send_enabled"][0]={"denom":"ugnet","enabled":true}' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
jq '.app_state["bank"]["supply"][0]={"denom":"ugnet","amount":"1000000000000"}' $MAIN_PATH_CONFIG/genesis.json > $MAIN_PATH_CONFIG/tmp_genesis.json && mv $MAIN_PATH_CONFIG/tmp_genesis.json $MAIN_PATH_CONFIG/genesis.json
```

## Allocate genesis accounts tokens

Allocates tokens to genesis accounts. For each validator account, 10000 GNET is allocated. For the treasury account, 970000 GNET is allocated.

```bash
gala add-genesis-account "$(gala keys show validator01 -a --keyring-dir $MAIN_PATH_HOME)" 10000000000ugnet
gala add-genesis-account "$(gala keys show validator02 -a --keyring-dir $MAIN_PATH_HOME)" 10000000000ugnet
gala add-genesis-account "$(gala keys show validator03 -a --keyring-dir $MAIN_PATH_HOME)" 10000000000ugnet
gala add-genesis-account "$(gala keys show treasury -a --keyring-dir $MAIN_PATH_HOME)" 970000000000ugnet
```

## Sign genesis transaction

Each validator stakes 6000 GNET.

Check the genesis staking description: https://hub.cosmos.network/main/resources/genesis.html#staking

### Validator 01


```bash
export VALIDATOR_MONIKER=validator01
export VALIDATOR_IP=192.168.10.2
export VALIDATOR_P2P_PORT=26656
export VALIDATOR_STAKING_AMOUNT=6000000000ugnet
export VALIDATOR_HOME=$MAIN_PATH_HOME/validators/$VALIDATOR_MONIKER

gala init \
  $VALIDATOR_MONIKER \
  --chain-id $CHAIN_ID \
  --default-denom gnet \
  --home $VALIDATOR_HOME
  
cp $MAIN_PATH_CONFIG/genesis.json $VALIDATOR_HOME/config/genesis.json

gala gentx $VALIDATOR_MONIKER $VALIDATOR_STAKING_AMOUNT \
  --ip $VALIDATOR_IP \
  --p2p-port $VALIDATOR_P2P_PORT \
  --home $VALIDATOR_HOME \
  --keyring-dir $MAIN_PATH_HOME  \
  --keyring-backend $KEYRING_BACKEND

mkdir -p $MAIN_PATH_CONFIG/gentx/
cp $VALIDATOR_HOME/config/gentx/* $MAIN_PATH_CONFIG/gentx/
```

### Validator 02

```bash
export VALIDATOR_MONIKER=validator02
export VALIDATOR_IP=192.168.10.3
export VALIDATOR_P2P_PORT=26656
export VALIDATOR_STAKING_AMOUNT=6000000000ugnet
export VALIDATOR_HOME=$MAIN_PATH_HOME/validators/$VALIDATOR_MONIKER

gala init \
  $VALIDATOR_MONIKER \
  --chain-id $CHAIN_ID \
  --default-denom gnet \
  --home $VALIDATOR_HOME
  
cp $MAIN_PATH_CONFIG/genesis.json $VALIDATOR_HOME/config/genesis.json

gala gentx $VALIDATOR_MONIKER $VALIDATOR_STAKING_AMOUNT \
  --ip $VALIDATOR_IP \
  --p2p-port $VALIDATOR_P2P_PORT \
  --home $VALIDATOR_HOME \
  --keyring-dir $MAIN_PATH_HOME  \
  --keyring-backend $KEYRING_BACKEND

mkdir -p $MAIN_PATH_CONFIG/gentx/
cp $VALIDATOR_HOME/config/gentx/* $MAIN_PATH_CONFIG/gentx/
```

### Validator 03

```bash
export VALIDATOR_MONIKER=validator03
export VALIDATOR_IP=192.168.10.4
export VALIDATOR_P2P_PORT=26656
export VALIDATOR_STAKING_AMOUNT=6000000000ugnet
export VALIDATOR_HOME=$MAIN_PATH_HOME/validators/$VALIDATOR_MONIKER

gala init \
  $VALIDATOR_MONIKER \
  --chain-id $CHAIN_ID \
  --default-denom gnet \
  --home $VALIDATOR_HOME
  
cp $MAIN_PATH_CONFIG/genesis.json $VALIDATOR_HOME/config/genesis.json

gala gentx $VALIDATOR_MONIKER $VALIDATOR_STAKING_AMOUNT \
  --ip $VALIDATOR_IP \
  --p2p-port $VALIDATOR_P2P_PORT \
  --home $VALIDATOR_HOME \
  --keyring-dir $MAIN_PATH_HOME  \
  --keyring-backend $KEYRING_BACKEND

mkdir -p $MAIN_PATH_CONFIG/gentx/
cp $VALIDATOR_HOME/config/gentx/* $MAIN_PATH_CONFIG/gentx/
```

```bash
gala collect-gentxs
rm $MAIN_PATH_CONFIG/node_key.json
rm $MAIN_PATH_CONFIG/priv_validator_key.json
```

```bash
gala validate-genesis
```

## Populate validators config folder with actual config files

### Validator 01

```bash
export VALIDATOR_MONIKER=validator01
export VALIDATOR_HOME=$MAIN_PATH_HOME/validators/$VALIDATOR_MONIKER

cp $MAIN_PATH_CONFIG/app.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/client.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/config.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/genesis.json $VALIDATOR_HOME/config/

sed -i '' 's/moniker = "localtestnet"/moniker = "'$VALIDATOR_MONIKER'"/g'  $VALIDATOR_HOME/config/config.toml
```

### Validator 02

```bash
export VALIDATOR_MONIKER=validator02
export VALIDATOR_HOME=$MAIN_PATH_HOME/validators/$VALIDATOR_MONIKER

cp $MAIN_PATH_CONFIG/app.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/client.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/config.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/genesis.json $VALIDATOR_HOME/config/

sed -i '' 's/moniker = "localtestnet"/moniker = "'$VALIDATOR_MONIKER'"/g'  $VALIDATOR_HOME/config/config.toml
```

### Validator 03

```bash
export VALIDATOR_MONIKER=validator03
export VALIDATOR_HOME=$MAIN_PATH_HOME/validators/$VALIDATOR_MONIKER

cp $MAIN_PATH_CONFIG/app.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/client.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/config.toml $VALIDATOR_HOME/config/
cp $MAIN_PATH_CONFIG/genesis.json $VALIDATOR_HOME/config/

sed -i '' 's/moniker = "localtestnet"/moniker = "'$VALIDATOR_MONIKER'"/g'  $VALIDATOR_HOME/config/config.toml
```

# Run the nodes

```bash
docker compose up -d
```

# Test staking delegation

Let's add delegator account and send some tokens to it from treasury:

```bash
echo "trial coin distance lens rubber attract trumpet wrestle move you marine crazy slender oak depth own witness sudden reunion habit scale electric afford town" | gala keys add delegator01 --recover --algo secp256k1 --keyring-backend $KEYRING_BACKEND --keyring-dir $MAIN_PATH_HOME
```
```bash
TREASURY=$(gala keys show treasury -a --keyring-backend test)
DELEGATOR=$(gala keys show delegator01 -a --keyring-backend test)
echo "Treasury: $TREASURY"
echo "Delegator: $DELEGATOR"
```

```bash
echo "Send 70.5 GNET tokens to delegator from treasury..."
gala tx bank send \
  $TREASURY \
  $DELEGATOR \
  70500000ugnet \
  --fees 2000000ugnet \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --yes
```

```bash
echo "Check delegator balance..."
gala query bank balances $DELEGATOR
```

Let's delegate some tokens to validator01:

```bash
VALIDATOR_01_ADDR=$(gala keys show validator01 --bech val -a --keyring-backend test)
VALIDATOR_02_ADDR=$(gala keys show validator02 --bech val -a --keyring-backend test)
VALIDATOR_03_ADDR=$(gala keys show validator03 --bech val -a --keyring-backend test)

echo "Check validator01 ${VALIDATOR_01_ADDR} delegations..."
gala query staking delegations-to $VALIDATOR_01_ADDR
```

```bash
echo "Delegate 13 GNET tokens to validator01 ($VALIDATOR_01_ADDR)..."

gala tx staking delegate $VALIDATOR_01_ADDR 13000000ugnet \
  --from delegator01 \
  --fees 2000000ugnet \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --yes
```

```bash
echo "Check validator01 ${VALIDATOR_01_ADDR} delegations..."
gala query staking delegations-to $VALIDATOR_01_ADDR
```

Unbond some tokens from validator01 and delegate to validator02:

```bash
echo "Unbond 9 GNET tokens from validator01 ($VALIDATOR_01_ADDR)..."

gala tx staking unbond $VALIDATOR_01_ADDR 9000000ugnet \
  --from delegator01 \
  --fees 2000000ugnet \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --yes
```


```bash
echo "Delegate 2.5 GNET tokens to validator02 ($VALIDATOR_02_ADDR)..."

gala tx staking delegate $VALIDATOR_02_ADDR 2500000ugnet \
  --from delegator01 \
  --fees 2000000ugnet \
  --chain-id $CHAIN_ID \
  --keyring-backend test \
  --yes
```

```bash
echo "Check validator02 ${VALIDATOR_02_ADDR} delegations..."
gala query staking delegations-to $VALIDATOR_02_ADDR
```

Query staking history:

```bash
echo "Query staking history..."
gala query staking historical-info $(gala query block | jq ".block.header.height" -r)
```

```bash
echo "Query staking validators"
gala query staking validators
```

# Keplr wallet

Keplr suggest chain params: https://docs.keplr.app/api/suggest-chain.html
