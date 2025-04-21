#!/bin/bash

# Copyright 2024 Galactica Network
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
##########################################################################

MAIN_PATH_HOME=${1:-"./.galactica"}
MAIN_PATH_CONFIG=$MAIN_PATH_HOME/config
# {identifier}_{EIP155}-{version}
CHAIN_ID=${2:-"galactica_41246"}
KEYRING_BACKEND=${3:-"test"}
BASE_DENOM=${4:-"agnet"}
DISPLAY_DENOM=${5:-"gnet"}
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# localKey address 0x7cb61d4117ae31a12e393a1cfa3bac666481d02e
PREDEFINED_KEY_NAME="localkey"
PREDEFINED_KEY_MNEMONIC="gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"


function gala() {
    ./build/galacticad --home "$MAIN_PATH_HOME" "$@"
}

function add_key() {
    local key_name=$1

    gala keys add "$key_name" --algo eth_secp256k1 --keyring-backend $KEYRING_BACKEND --keyring-dir $MAIN_PATH_HOME/
}

function add_key_predefined() {
  echo $PREDEFINED_KEY_MNEMONIC | gala keys add \
    $PREDEFINED_KEY_NAME \
    --recover \
    --keyring-backend $KEYRING_BACKEND \
    --algo "eth_secp256k1" \
    --keyring-dir $MAIN_PATH_HOME/
}

function configure_app() {
    # Configure app settings
    sed -i '/\[api\]/,+3 s/enable = false/enable = true/' $MAIN_PATH_CONFIG/app.toml
    sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' $MAIN_PATH_CONFIG/app.toml
    # sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' $MAIN_PATH_CONFIG/app.toml
    sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' $MAIN_PATH_CONFIG/app.toml
    # sed -i '/\[grpc-web\]/,+7 s/address = "localhost:9091"/address = "0.0.0.0:9091"/' $MAIN_PATH_CONFIG/app.toml
    sed -i 's/pruning = "default"/pruning = "nothing"/g'  $MAIN_PATH_CONFIG/app.toml
    sed -i "s|minimum-gas-prices = \".*\"|minimum-gas-prices = \"10${BASE_DENOM}\"|g" "$MAIN_PATH_CONFIG/app.toml"
    sed -i '/\[telemetry\]/,+8 s/enabled = false/enabled = true/' $MAIN_PATH_CONFIG/app.toml
    sed -i '/\[telemetry\]/,+20 s/prometheus-retention-time = 0/prometheus-retention-time = 60/' $MAIN_PATH_CONFIG/app.toml
    sed -i '/global-labels = \[/a\
  \["chain_id", "'$CHAIN_ID'"\],
' $MAIN_PATH_CONFIG/app.toml

    sed -i 's/timeout_propose = ".*"/timeout_propose = "3s"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/timeout_propose_delta = ".*"/timeout_propose_delta = "500ms"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/timeout_prevote = ".*"/timeout_prevote = "1s"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/timeout_prevote_delta = ".*"/timeout_prevote_delta = "500ms"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/timeout_precommit = ".*"/timeout_precommit = "1s"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/timeout_precommit_delta = ".*"/timeout_precommit_delta = "500ms"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/timeout_commit = ".*"/timeout_commit = "5s"/g'  $MAIN_PATH_CONFIG/config.toml


    # configure config settings
    sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26658"/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/max_num_inbound_peers = 40/max_num_inbound_peers = 120/g'  $MAIN_PATH_CONFIG/config.toml
    sed -i 's/max_num_outbound_peers = 10/max_num_outbound_peers = 60/g'  $MAIN_PATH_CONFIG/config.toml
}

function check_config_app {
    check_config_change "$MAIN_PATH_CONFIG/app.toml" 'enable = true' "API not enabled"
    check_config_change "$MAIN_PATH_CONFIG/app.toml" 'address = "tcp://0.0.0.0:1317"' "API address not updated"
    # check_config_change "$MAIN_PATH_CONFIG/app.toml" 'enabled-unsafe-cors = true' "CORS not enabled"  # Uncomment if needed
    check_config_change "$MAIN_PATH_CONFIG/app.toml" 'address = "0.0.0.0:9090"' "gRPC address not updated"
    # check_config_change "$MAIN_PATH_CONFIG/app.toml" 'address = "0.0.0.0:9091"' "gRPC-web address not updated"
    check_config_change "$MAIN_PATH_CONFIG/app.toml" 'pruning = "nothing"' "Pruning not set to 'nothing'"
    check_config_change "$MAIN_PATH_CONFIG/app.toml" "minimum-gas-prices = \"10$BASE_DENOM\"" "Minimum gas price not set correctly"
    check_config_change "$MAIN_PATH_CONFIG/app.toml" 'enabled = true' "Telemetry not enabled"
    check_config_change "$MAIN_PATH_CONFIG/app.toml" 'prometheus-retention-time = 60' "Prometheus retention time not set"
    check_config_change "$MAIN_PATH_CONFIG/app.toml" "\[\"chain_id\", \"$CHAIN_ID\"\]" "Telemetry chain_id label not added"

    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_propose = "3s"' "timeout_propose not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_propose_delta = "500ms"' "timeout_propose_delta not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_prevote = "1s"' "timeout_prevote not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_prevote_delta = "500ms"' "timeout_prevote_delta not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_precommit = "1s"' "timeout_precommit not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_precommit_delta = "500ms"' "timeout_precommit_delta not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'timeout_commit = "5s"' "timeout_commit not updated"

    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'laddr = "tcp://0.0.0.0:26657"' "RPC address not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'proxy_app = "tcp://127.0.0.1:26658"' "Proxy app address not set"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'cors_allowed_origins = \["\*"\]' "CORS allowed origins not set"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'max_num_inbound_peers = 120' "Max inbound peers not updated"
    check_config_change "$MAIN_PATH_CONFIG/config.toml" 'max_num_outbound_peers = 60' "Max outbound peers not updated"
}

check_config_change() {
    local file=$1
    local pattern=$2
    local message=$3

    grep -q "$pattern" "$file" || echo "❌ $message"
}


function update_genesis_json() {
    local jq_command=$1
    local file_path=${2:-"$MAIN_PATH_CONFIG/genesis.json"}

    jq "$jq_command" "$file_path" > "${file_path}.tmp" && mv "${file_path}.tmp" "$file_path"
}

function configure_genesis() {
    local staking_min_deposit="$1"
    # local total_supply="$2"
    local faucet_address="$2"
    local voting_period="600s"
    local expedited_voting_period="60s"
    local unbonding_time="30s"
    local max_deposit_period="600s"
    local block_max_bytes="22020096"
    local block_max_gas="40000000"
    local time_iota_ms="1000"
    local inflation_validators_share="0.99933"
    local inflation_faucet_share="0.00067"

    local symbol=$(echo $DISPLAY_DENOM | tr '[:lower:]' '[:upper:]')

    update_genesis_json '.consensus_params.block.max_bytes = "'$block_max_bytes'"'
    update_genesis_json '.consensus_params.block.max_gas = "'$block_max_gas'"'
    update_genesis_json '.consensus_params.block.time_iota_ms = "'$time_iota_ms'"'

    update_genesis_json '.app_state.gov.voting_params.voting_period = "'$voting_period'"'
    update_genesis_json '.app_state.staking.params.bond_denom = "'$BASE_DENOM'"'
    update_genesis_json '.app_state.staking.params.unbonding_time = "'$unbonding_time'"'
    update_genesis_json '.app_state.crisis.constant_fee.denom = "'$BASE_DENOM'"'
    update_genesis_json '.app_state.gov.deposit_params.min_deposit[0].denom = "'$BASE_DENOM'"'
    update_genesis_json '.app_state.gov.params.min_deposit[0] = {"denom": "'$BASE_DENOM'", "amount": "'$staking_min_deposit'"}'
    update_genesis_json '.app_state.gov.params.max_deposit_period = "'$max_deposit_period'"'
    update_genesis_json '.app_state.gov.params.voting_period = "'$voting_period'"'
    update_genesis_json '.app_state.gov.params.expedited_voting_period = "'$expedited_voting_period'"'
    update_genesis_json '.app_state.mint.params.mint_denom = "'$BASE_DENOM'"'
    update_genesis_json '.app_state.bank.denom_metadata[0] = {
        "description": "The native staking token of the Galactica Network.",
        "denom_units": [
            {"denom": "'$BASE_DENOM'", "exponent": 0, "aliases": ["attognet"]},
            {"denom": "ugnet", "exponent": 6, "aliases": ["micrognet"]},
            {"denom": "'$DISPLAY_DENOM'", "exponent": 18}
        ],
        "base": "'$BASE_DENOM'",
        "display": "'$DISPLAY_DENOM'",
        "name": "Galactica Network",
        "symbol": "'$symbol'",
        "uri": "",
        "uri_hash": ""
    }'
    update_genesis_json '.app_state.bank.send_enabled[0] = {"denom": "'$BASE_DENOM'", "enabled": true}'
#    update_genesis_json '.app_state.bank.supply[0] = {"denom": "'$BASE_DENOM'", "amount": "'$total_supply'"}'

    # EVM params
    update_genesis_json '.app_state.evm.params.evm_denom = "'$BASE_DENOM'"'

    # inflation params
    update_genesis_json '.app_state.inflation.params.enable_inflation = true'
    update_genesis_json '.app_state.inflation.params.mint_denom = "'$BASE_DENOM'"'

    update_genesis_json '.app_state.inflation.params.inflation_distribution.validators_share = "'$inflation_validators_share'"'
    update_genesis_json '.app_state.inflation.params.inflation_distribution.other_shares = [
      {"address": "'$faucet_address'","name": "faucet","share": "'$inflation_faucet_share'"}
    ]'

    update_genesis_json '.app_state.inflation.inflation_distribution.validators_share = "'$inflation_validators_share'"'
    update_genesis_json '.app_state.inflation.inflation_distribution.other_shares = [
      {"address": "'$faucet_address'","name": "faucet","share": "'$inflation_faucet_share'"}
    ]'
}

function check_genesis_change() {
    local jq_expr=$1
    local expected_value=$2
    local file_path=${3:-"$MAIN_PATH_CONFIG/genesis.json"}

    local actual_value
    actual_value=$(jq -r "$jq_expr" "$file_path")

    if [[ "$actual_value" != "$expected_value" ]]; then
        echo "❌ Genesis config mismatch for '$jq_expr': expected '$expected_value', got '$actual_value'"
        exit 1
    else
        echo "✅ Verified: $jq_expr = $expected_value"
    fi
}

function verify_genesis() {
    local staking_min_deposit="$1"
    local faucet_address="$2"
    local symbol=$(echo $DISPLAY_DENOM | tr '[:lower:]' '[:upper:]')

    check_genesis_change '.consensus_params.block.max_bytes' "22020096"
    check_genesis_change '.consensus_params.block.max_gas' "40000000"
    check_genesis_change '.consensus_params.block.time_iota_ms' "1000"

    check_genesis_change '.app_state.gov.voting_params.voting_period' "600s"
    check_genesis_change '.app_state.staking.params.bond_denom' "$BASE_DENOM"
    check_genesis_change '.app_state.staking.params.unbonding_time' "30s"
    check_genesis_change '.app_state.crisis.constant_fee.denom' "$BASE_DENOM"
    check_genesis_change '.app_state.gov.deposit_params.min_deposit[0].denom' "$BASE_DENOM"
    check_genesis_change '.app_state.gov.params.min_deposit[0].denom' "$BASE_DENOM"
    check_genesis_change '.app_state.gov.params.min_deposit[0].amount' "$staking_min_deposit"
    check_genesis_change '.app_state.gov.params.max_deposit_period' "600s"
    check_genesis_change '.app_state.gov.params.voting_period' "600s"
    check_genesis_change '.app_state.gov.params.expedited_voting_period' "60s"
    check_genesis_change '.app_state.mint.params.mint_denom' "$BASE_DENOM"

    check_genesis_change '.app_state.bank.denom_metadata[0].base' "$BASE_DENOM"
    check_genesis_change '.app_state.bank.denom_metadata[0].display' "$DISPLAY_DENOM"
    check_genesis_change '.app_state.bank.denom_metadata[0].symbol' "$symbol"
    check_genesis_change '.app_state.bank.send_enabled[0].enabled' "true"
    check_genesis_change '.app_state.bank.send_enabled[0].denom' "$BASE_DENOM"

    check_genesis_change '.app_state.evm.params.evm_denom' "$BASE_DENOM"

    check_genesis_change '.app_state.inflation.params.enable_inflation' "true"
    check_genesis_change '.app_state.inflation.params.mint_denom' "$BASE_DENOM"
    check_genesis_change '.app_state.inflation.params.inflation_distribution.validators_share' "0.99933"
    check_genesis_change '.app_state.inflation.inflation_distribution.validators_share' "0.99933"

    check_genesis_change '.app_state.inflation.params.inflation_distribution.other_shares[0].address' "$faucet_address"
    check_genesis_change '.app_state.inflation.params.inflation_distribution.other_shares[0].share' "0.00067"
    check_genesis_change '.app_state.inflation.inflation_distribution.other_shares[0].address' "$faucet_address"
    check_genesis_change '.app_state.inflation.inflation_distribution.other_shares[0].share' "0.00067"
}

function add_genesis_account() {
    local account_name="$1"
    local amount="$2"

    gala genesis add-genesis-account \
      "$(gala keys show $account_name -a --keyring-dir $MAIN_PATH_HOME)" \
      $amount
}

function initialize_validator() {
    local moniker=$1
    local ip=$2
    local p2p_port=$3
    local staking_amount=$4
    local validator_home=$MAIN_PATH_HOME/validators/$moniker

    gala init \
      $moniker \
      --chain-id $CHAIN_ID \
      --default-denom $BASE_DENOM \
      --home $validator_home

    cp $MAIN_PATH_CONFIG/genesis.json $validator_home/config/genesis.json

    gala genesis gentx \
      $moniker \
      $staking_amount \
      --ip $ip \
      --p2p-port $p2p_port \
      --home $validator_home \
      --keyring-dir $MAIN_PATH_HOME \
      --keyring-backend $KEYRING_BACKEND

    mkdir -p $MAIN_PATH_CONFIG/gentx/
    cp $validator_home/config/gentx/* $MAIN_PATH_CONFIG/gentx/
}

function configure_validator() {
    local moniker=$1
    local ip_address=$2
    local validator_home="$MAIN_PATH_HOME/validators/$moniker"

    echo "Configuring $moniker validator with IP address $ip_address"

    cp $MAIN_PATH_CONFIG/app.toml $validator_home/config/
    cp $MAIN_PATH_CONFIG/client.toml $validator_home/config/
    cp $MAIN_PATH_CONFIG/config.toml $validator_home/config/
    cp $MAIN_PATH_CONFIG/genesis.json $validator_home/config/

    # Filter out the IP address of the current validator from persistent_peers
    local persistent_peers=$(cat $validator_home/config/config.toml | grep persistent_peers | cut -d '=' -f 2 | tr -d '"')
    IFS=',' read -ra parts <<< "$persistent_peers"
    filtered_parts=()
    for part in "${parts[@]}"; do
        if [[ ! "$part" =~ $ip_address ]]; then
            filtered_parts+=("$part")
        fi
    done
    local new_persistent_peers=$(IFS=','; echo "${filtered_parts[*]}")

    echo "Validator $moniker persistent_peers: $new_persistent_peers"
    sed -i "s/\(persistent_peers *= *\"\).*\(\" *\)/\1$new_persistent_peers\2/" $validator_home/config/config.toml
    sed -i 's/moniker = "localtestnet"/moniker = "'$moniker'"/g'  $validator_home/config/config.toml
    grep -rnl 'namespace = "galactica"' ./ | xargs sed -i 's?namespace = "galactica"?namespace = "galactica"?'
    grep -rnl 'prometheus_listen_addr = "validator01:26660"' ./ | xargs sed -i 's?prometheus_listen_addr = "validator01:26660"?prometheus_listen_addr = "'$moniker':26660"?'

    local key=$(gala keys unsafe-export-eth-key --keyring-backend test --keyring-dir ./$MAIN_PATH_HOME $moniker)
    yes '00000000' | gala keys unsafe-import-eth-key --keyring-backend test --keyring-dir ./$MAIN_PATH_HOME/validators/$moniker $moniker $key --chain-id $CHAIN_ID
}

function verify_validator_config() {
    local moniker=$1
    local ip_address=$2
    local validator_home="$MAIN_PATH_HOME/validators/$moniker"
    local config_file="$validator_home/config/config.toml"
    local app_file="$validator_home/config/app.toml"
    local client_file="$validator_home/config/client.toml"
    local genesis_file="$validator_home/config/genesis.json"

    echo "🔍 Verifying validator configuration: $moniker"

    # Check if all required config files exist
    for file in "$config_file" "$app_file" "$client_file" "$genesis_file"; do
        if [[ ! -f "$file" ]]; then
            echo "❌ File not found: $file"
            exit 1
        else
            echo "✅ Found file: $file"
        fi
    done

    # Check if moniker is correctly set
    local actual_moniker
    actual_moniker=$(grep '^moniker' "$config_file" | cut -d '=' -f2 | tr -d ' "')
    if [[ "$actual_moniker" != "$moniker" ]]; then
        echo "❌ Moniker mismatch: expected '$moniker', found '$actual_moniker'"
        exit 1
    else
        echo "✅ Moniker is correctly set: $actual_moniker"
    fi

    # Check that the validator's own IP address is not in persistent_peers
    local persistent_peers
    persistent_peers=$(grep '^persistent_peers' "$config_file" | cut -d '=' -f2- | tr -d ' "')
    if echo "$persistent_peers" | grep -q "$ip_address"; then
        echo "❌ Validator IP address $ip_address is still in persistent_peers"
        exit 1
    else
        echo "✅ Validator IP address $ip_address is not in persistent_peers"
    fi

    # Check if key was imported successfully
    local eth_key_path="$validator_home/keyring-test/$moniker.info"
    if [[ ! -f "$eth_key_path" ]]; then
        echo "❌ Key not found at path: $eth_key_path"
        exit 1
    else
        echo "✅ Key found: $eth_key_path"
    fi

    echo "🎉 Validator configuration verification for $moniker completed successfully!"
}


function configure_faucet() {
    local key=$(gala keys unsafe-export-eth-key --keyring-backend test --keyring-dir ./$MAIN_PATH_HOME faucet)
    yes '00000000' | gala keys unsafe-import-eth-key --keyring-backend test --keyring-dir ./$MAIN_PATH_HOME/faucet faucet $key --chain-id $CHAIN_ID

    echo $key > ./$MAIN_PATH_HOME/faucet/PRIVATE_KEY
}

function main() {
    rm -r $MAIN_PATH_HOME
    
    # init
    gala init localtestnet --chain-id $CHAIN_ID --default-denom $BASE_DENOM

    # configure binary
    echo "Setting up galactica configuration..."
    echo "Home directory: $MAIN_PATH_HOME"
    echo "Config directory: $MAIN_PATH_CONFIG"
    echo "Chain ID: $CHAIN_ID"
    echo "Keyring backend: $KEYRING_BACKEND"

    gala config set client chain-id $CHAIN_ID
	gala config set client keyring-backend $KEYRING_BACKEND

    # Add keys
    add_key "validator01"
    add_key "validator02"
    add_key "validator03"
    add_key "treasury"
    add_key "faucet"
    add_key_predefined

    configure_app
    check_config_app

    local faucet_address_bech32=$(gala keys show faucet -a --keyring-dir $MAIN_PATH_HOME)
    echo "Faucet address bech32: $faucet_address_bech32"

    # $1 = staking min deposit
    # $2 = faucet eth address
    configure_genesis "5000000" $faucet_address_bech32
    verify_genesis "5000000" $faucet_address_bech32

    # Add genesis accounts
    add_genesis_account "validator01" "10000000000000000000000$BASE_DENOM"
    add_genesis_account "validator02" "10000000000000000000000$BASE_DENOM"
    add_genesis_account "validator03" "10000000000000000000000$BASE_DENOM"
    add_genesis_account "faucet" "10000000000000000000000$BASE_DENOM"
    add_genesis_account "localkey" "10000000000000000000000$BASE_DENOM"
    add_genesis_account "treasury" "950000000000000000000000$BASE_DENOM"


    # Initialize validators
    initialize_validator "validator01" "validator01" 26656 "9000000000000000000000$BASE_DENOM"
    initialize_validator "validator02" "validator02" 26656 "9000000000000000000000$BASE_DENOM"
    initialize_validator "validator03" "validator03" 26656 "9000000000000000000000$BASE_DENOM"


    gala genesis collect-gentxs
    rm $MAIN_PATH_CONFIG/node_key.json
    rm $MAIN_PATH_CONFIG/priv_validator_key.json

    gala genesis validate-genesis

    # $1 = moniker
    # $2 = validator ip address
    configure_validator "validator01" "validator01"
    verify_validator_config "validator01" "validator01"

    configure_validator "validator02" "validator02"
    verify_validator_config "validator01" "validator01"

    configure_validator "validator03" "validator03"
    verify_validator_config "validator01" "validator01"

    configure_faucet
    echo "✅ Finish"
}

main
