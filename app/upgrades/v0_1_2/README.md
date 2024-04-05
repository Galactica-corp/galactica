# Upgrade to v0.1.2

## Guide for Upgrading from v0.1.1 to v0.1.2

1. Halt the currently running `galacticad` process.
2. **!!! Backup the data directory.** !!!
3. Fetch the latest updates from the repository and switch to the `v0.1.2` tag.
    ```bash
    git fetch --all --tags
    git checkout v0.1.2
    ```
4. Build the updated source code.
    ```bash
    make install
    ```
5. Execute the upgrade script located at ./scripts/upgrade_v0_1_2.sh
    ```bash
    ./scripts/upgrade_v0_1_2.sh
    ```
   Ensure the following environment variable is defined:
    - `GALACTICA_HOME`: Specifies the home directory for the `galacticad` node.

   The script executes the following operations:
    - Updates the node storage to version v0.1.2.
    - Roll back the state of the most recent block.
    - Backup of the current priv_validator_state.json and replaces it with a new version with default settings.

   Alternatively, these steps can be performed manually:
   ```bash
   galacticad --home $GALACTICA_HOME rollback --hard
   echo '{"height":"0","round":0,"step":0}' > $GALACTICA_HOME/data/priv_validator_state.json
    ```
   Executing the rollback command will also automatically execute the storage upgrade.

6. Restart the `galacticad` node.

