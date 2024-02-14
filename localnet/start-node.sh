#!/bin/sh

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

LOGLEVEL="info"
HOMEDIR="/root/.galactica"
CHAIN_ID="galactica_9000-1"

galacticad start \
  --log_level $LOGLEVEL \
  --home "$HOMEDIR" \
  --chain-id "$CHAIN_ID" \
  --rpc.unsafe \
  --json-rpc.address 0.0.0.0:8545 \
  --json-rpc.ws-address 0.0.0.0:8546 \
  --json-rpc.api eth,txpool,personal,net,debug,web3 \
  --api.enable

