# Galactica

**Galactica** is a Layer 1 protocol that leverages zero-knowledge cryptography to achieve Sybil resistance,
compliant privacy and infuse robust reputation primitives into DeFi and DAOs.

## Project Build

To compile and prepare the **Galactica** project for deployment, ensure you have the necessary tools and follow the outlined steps.

### Prerequisites

- **Go**: Galactica requires Go version 1.21. Verify your installation by running `go version` in your terminal.

### Building the Project

1. **Clone the Repository**: Get the Galactica source code by cloning the repository.

   ```sh
   git clone https://github.com/Galactica-corp/galactica.git
   cd galactica
   ```

2. **Build the Binary**: Compile the source code into an executable binary with the `make` command.

   ```sh
   make build
   ```

   The binary will be located in the `build/` directory.


3. **Install the Binary**: Install the Galactica binary into a directory of your choice. By default, the binary is installed to `/usr/local/bin`, but you can specify a custom directory by setting the `BINDIR` environment variable before running the `make install` command. This step may require administrative privileges.

   ```sh
   make install
   ```

   For a custom directory, such as `/opt/galactica/bin`, run:
   
   ```sh
   BINDIR=/opt/galactica/bin make install
   ```

Ensure the installation path is included in your system's PATH environment variable to run Galactica from any terminal.

## galactica-testnet-v1

After successfully building and installing the Galactica binary, the next step is to connect to the Galactica testnet v1. This involves generating a key, setting up your node's home directory, and ultimately submitting a genesis transaction (gentx) to become a validator. Follow these steps to join the network:

1. **Generate Your Validator Key**: Within the same repository, you'll find scripts and guidance on generating your validator key. This key is essential for securing your node and participating in the network.

2. **Set Up Node's Home Directory**: Navigate to the [Galactica Networks repository](https://github.com/Galactica-corp/networks) and follow the instructions to install the necessary home directory for your node. This step is crucial for configuring your node with the correct network parameters.

3. **Submit Your Address for Genesis Airdrop**: After generating your key, submit your address to the Galactica network team. Your address will be included in the genesis airdrop, enabling you to stake tokens and participate in the network from the initial block.

4. **Generate and Submit Your gentx**: With your address included in the genesis airdrop, generate your gentx using the provided scripts. Submit this gentx as a pull request to the [Galactica-corp/networks](https://github.com/Galactica-corp/networks) repository. This step registers you as a validator on the network.

By following these streamlined steps, you can transition smoothly from building the Galactica binary to actively participating in the Galactica testnet v1 as a validator. This approach ensures that new validators have a clear and direct path to joining the network, contributing to its security and decentralization from the very first block.

## Usage

### Run local network

Local network allows you to run your blockchain locally. It contains 3 validators nodes.

To build docker for local network run:

```sh
make localnet-build
```

To initialize local network validators configuration files run:

```sh
make localnet-init
```

Files will be saved to `./.galactica` folder.

To start local network run:

```sh
make localnet-start
```

You can find network parameters for validators in `localnet/init-configs.sh` file.

### Interact with the Blockchain

To interact with the blockchain, you can use [CLI](https://docs.cosmos.network/v0.50/learn/advanced/cli)
and [gRPC or REST endpoints](https://docs.cosmos.network/v0.50/learn/advanced/grpc_rest). Additionally, the node
provides an [Ethereum RPC endpoint](https://ethereum.org/en/developers/docs/apis/json-rpc).

## Joining Network

Our [wiki](https://github.com/Galactica-corp/galactica/wiki) gives you information on how to run Galactica validators to
take part in the decentralized infrastructure and blockchain consensus.

## Contributing

Contributions to the blockchain project are welcome! Feel free to open issues for bug fixes, feature requests, or submit
pull requests to contribute directly.

## License

This project is licensed under the [Apache License 2.0](LICENSE), except for the modified third-party modules `x/epochs`
and `x/inflation`described below.

Certain modules within this project are adapted versions sourced from
the [EVMOS](https://github.com/evmos/evmos/tree/v12.1.6) library. Specifically, the modules `x/epochs`
and `x/inflation`, obtained from version `v12.1.6` as of July 4, 2023, have been customized to meet the requirements of
this project. These modifications entail the removal of unused fields and methods, as well as the addition of new ones
tailored to our project's needs. These modified modules are licensed under
the [GNU Lesser General Public License v3.0](COPYING.LESSER). Please refer to the corresponding license file for more
details.

### Replace LGPL-3.0 Libraries Implementation

To comply with the LGPL-3.0 license for the `go-ethereum` library, you may need to have a possibility to replace it with
another interface-compatible library. This can be achieved using the `replace` directive in the `go.mod` file.

Follow these steps to replace the `go-ethereum` library:

1. **Identify Replacement Library**: Find an alternative library that provides similar functionality and is compatible
   with the interfaces used in the project.

2. **Update `go.mod` File**: Add a `replace` directive in the `go.mod` file to specify the replacement library. For
   example:

   ```
   replace github.com/ethereum/go-ethereum => github.com/your/ethereum-library v1.0.0
   ```

   Replace `github.com/your/ethereum-library` with the actual import path of the replacement library, and `v1.0.0` with
   the appropriate version tag.

3. **Verify Compatibility**: Ensure that the replacement library is interface-compatible with the existing codebase to
   prevent any breaking changes.

By following these steps, you can replace the go-ethereum library with another interface-compatible library while
ensuring compliance with the LGPL-3.0 license.
