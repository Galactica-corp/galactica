# Upgrade to v0.1.2

## Guide for Upgrading from v0.1.1 to v0.1.2

1. Halt the currently running `galacticad` process.
2. Fetch the latest updates from the repository and switch to the `v0.1.2` tag.
    ```bash
    git fetch --all --tags
    git checkout v0.1.2
    ```
3. Build the updated source code.
    ```bash
    make install
    ```
4. Execute the upgrade script located at ./scripts/upgrade_v0_1_2.sh
    ```bash
    ./scripts/upgrade_v0_1_2.sh
    ```
   Ensure the following environment variable is defined:
    - `GALACTICA_HOME`: Specifies the home directory for the `galacticad` node.

   The script executes the following operations:
    - Updates the node storage to version v0.1.2.
    - Reverts the state of the latest block.
    - Creates a backup of the current priv_validator_state.json and replaces it with a new version with default settings.
5. Restart the `galacticad` node.
