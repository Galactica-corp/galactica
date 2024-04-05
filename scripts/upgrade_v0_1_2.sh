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

#!/bin/bash

# This script is used to prepare the upgrade from v0.1.1 to v0.1.2

GALACTICAD_VERSION=$(galacticad version)
if [[ $GALACTICAD_VERSION != *"0.1.2"* ]]; then
  echo "Galactica version must be 0.1.2"
  exit 1
fi

# check env var GALACTICA_HOME and if not exists exit
if [ -z "$GALACTICA_HOME" ]; then
  echo "GALACTICA_HOME is not set, using default path"
  GALACTICA_HOME="$HOME/.galactica"
fi

# ask user to confirm:
echo "GALACTICA_HOME: $GALACTICA_HOME"
echo "\nThis script will perform the following actions:\n\
- Backup the existing priv_validator_state.json and replace it with a new one containing default values\n\
- Upgrade the node storage to v0.1.2\n\
- Rollback the latest block state\n"

if [ "$1" != "-y" ]; then
  read -p "Do you want to continue? (y/n): " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "User cancelled the script"
    exit 1
  fi
fi

PRIV_VAL_STATE=$GALACTICA_HOME/data/priv_validator_state.json

UPGRADE_INFO_FILE=$GALACTICA_HOME/data/upgrade-info.json
if [ -f "$UPGRADE_INFO_FILE" ]; then
  echo "upgrade v0.1.2 already applied"
  exit 1
fi

if [ ! -f "$PRIV_VAL_STATE" ]; then
  echo "priv_validator_state.json not found at $PRIV_VAL_STATE"
  exit 1
fi

cp $PRIV_VAL_STATE $PRIV_VAL_STATE.bkp

galacticad --home $GALACTICA_HOME rollback --hard

echo '{"height":"0","round":0,"step":0}' > $PRIV_VAL_STATE

