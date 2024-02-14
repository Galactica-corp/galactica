// Copyright 2024 Galactica Network
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ethkeys

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strings"

	crypto2 "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/armor"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/bcrypt"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/xsalsa20symmetric"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/spf13/cobra"
)

const (
	blockTypePrivKey = "TENDERMINT PRIVATE KEY"
	defaultAlgo      = "eth_secp256k1"
	headerType       = "type"
)

// UnsafeExportEthereumKeyCommand exports a key with the given name as a private key in hex format.
func UnsafeExportEthereumKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-export-eth-key [name]",
		Short: "**UNSAFE** Export an Ethereum private key",
		Long:  `**UNSAFE** Export an Ethereum private key unencrypted to use in dev tooling`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// clientCtx := client.GetClientContextFromCmd(cmd).WithKeyringOptions(hd.EthSecp256k1Option())
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			decryptPassword := ""
			conf := true

			inBuf := bufio.NewReader(cmd.InOrStdin())
			switch clientCtx.Keyring.Backend() {
			case keyring.BackendFile:
				decryptPassword, err = input.GetPassword(
					"**WARNING this is an unsafe way to export your unencrypted private key**\nEnter key password:",
					inBuf)
			case keyring.BackendOS:
				conf, err = input.GetConfirmation(
					"**WARNING** this is an unsafe way to export your unencrypted private key, are you sure?",
					inBuf, cmd.ErrOrStderr())
			}
			if err != nil || !conf {
				return err
			}

			// Exports private key from keybase using password
			arm, err := clientCtx.Keyring.ExportPrivKeyArmor(args[0], decryptPassword)
			if err != nil {
				return err
			}

			privKey, algo, err := UnarmorDecryptPrivKey(arm, decryptPassword)
			if err != nil {
				return err
			}

			if algo != ethsecp256k1.KeyType {
				return fmt.Errorf("invalid key algorithm, got %s, expected %s", algo, ethsecp256k1.KeyType)
			}

			// Converts key to Ethermint secp256k1 implementation
			ethPrivKey, ok := privKey.(*ethsecp256k1.PrivKey)
			if !ok {
				return fmt.Errorf("invalid private key type %T, expected %T", privKey, &ethsecp256k1.PrivKey{})
			}

			key, err := ethPrivKey.ToECDSA()
			if err != nil {
				return err
			}

			// Formats key for output
			privB := ethcrypto.FromECDSA(key)
			keyS := strings.ToUpper(hexutil.Encode(privB)[2:])

			fmt.Println(keyS)

			return nil
		},
	}
}

// UnarmorDecryptPrivKey returns the privkey byte slice, a string of the algo type, and an error
func UnarmorDecryptPrivKey(
	armorStr string,
	passphrase string,
) (privKey cryptotypes.PrivKey, algo string, err error) {
	blockType, header, encBytes, err := armor.DecodeArmor(armorStr)
	if err != nil {
		return privKey, "", err
	}

	if blockType != blockTypePrivKey {
		return privKey, "", fmt.Errorf("unrecognized armor type: %v", blockType)
	}

	if header["kdf"] != "bcrypt" {
		return privKey, "", fmt.Errorf("unrecognized KDF type: %v", header["kdf"])
	}

	if header["salt"] == "" {
		return privKey, "", fmt.Errorf("missing salt bytes")
	}

	saltBytes, err := hex.DecodeString(header["salt"])
	if err != nil {
		return privKey, "", fmt.Errorf("error decoding salt: %v", err.Error())
	}

	privKey, err = decryptPrivKey(saltBytes, encBytes, passphrase)

	if header[headerType] == "" {
		header[headerType] = defaultAlgo
	}

	return privKey, header[headerType], err
}

func decryptPrivKey(
	saltBytes []byte,
	encBytes []byte,
	passphrase string,
) (privKey cryptotypes.PrivKey, err error) {
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), crypto.BcryptSecurityParameter)
	if err != nil {
		return privKey, sdkerrors.Wrap(err, "error generating bcrypt key from passphrase")
	}

	key = crypto2.Sha256(key) // Get 32 bytes

	privKeyBytes, err := xsalsa20symmetric.DecryptSymmetric(encBytes, key)
	if err != nil && err.Error() == "Ciphertext decryption failed" {
		return privKey, sdkerrors.ErrWrongPassword
	} else if err != nil {
		return privKey, err
	}

	ethsecpKey := &ethsecp256k1.PrivKey{}
	if err := ethsecpKey.UnmarshalAmino(privKeyBytes[1:]); err != nil {
		return privKey, err
	}

	return ethsecpKey, nil
}
