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

package cmd

import (
	"math/big"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethermint "github.com/evmos/ethermint/types"

	"github.com/Galactica-corp/galactica/app"
)

const (
	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "gnet"

	// AddrLen is the allowed length (in bytes) for an address.
	//
	// NOTE: In the SDK, the default value is 255.
	AddrLen = 20

	// AttoGnet defines the default coin denomination used in Galactica in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in Ethermint.
	AttoGnet string = "agnet"

	// BaseDenomUnit defines the base denomination unit for Photons.
	// 1 photon = 1x10^{BaseDenomUnit} agnet
	BaseDenomUnit = 18

	MicroGnet      string = "ugnet"
	MicroDenomUnit        = 6

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)

func initSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := app.AccountAddressPrefix + "pub"
	validatorAddressPrefix := app.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := app.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := app.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := app.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(VerifyAddressFormat)

	RegisterDenoms()

	sdk.DefaultPowerReduction = math.NewIntFromBigInt(
		new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit), nil),
	)

	config.Seal()
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                      // Shared
	config.SetFullFundraiserPath(ethermint.BIP44HDPath) // nolint: staticcheck
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, math.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(AttoGnet, math.LegacyNewDecWithPrec(1, BaseDenomUnit)); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(MicroGnet, math.LegacyNewDecWithPrec(1, MicroDenomUnit)); err != nil {
		panic(err)
	}
}

// VerifyAddressFormat verifies whether the address is compatible with Ethereum
func VerifyAddressFormat(bz []byte) error {
	if len(bz) == 0 {
		return errors.Wrap(sdkerrors.ErrUnknownAddress, "invalid address; cannot be empty")
	}
	if len(bz) != AddrLen {
		return errors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"invalid address length; got: %d, expect: %d", len(bz), AddrLen,
		)
	}

	return nil
}
