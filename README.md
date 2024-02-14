# Galactica

**Galactica** is a Layer 1 protocol that leverages zero-knowledge cryptography to achieve Sybil resistance,
compliant privacy and infuse robust reputation primitives into DeFi and DAOs.

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
