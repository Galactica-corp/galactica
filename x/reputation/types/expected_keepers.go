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

package types

import (
	"cosmossdk.io/math"
	"cosmossdk.io/x/feegrant"
	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(sdk.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// StakingKeeper defines the expected interface for the Staking module.
type StakingKeeper interface {
	TotalBondedTokens(sdk.Context) math.Int
	// Methods imported from account should be defined here
}

// SlashingKeeper defines the expected interface for the Slashing module.
type SlashingKeeper interface {
	Slash(ctx sdk.Context, consAddr sdk.ConsAddress, fraction sdk.Dec, power, distributionHeight int64)
	// Methods imported from account should be defined here
}

// DistributionKeeper defines the expected interface for the Distribution module.
type DistributionKeeper interface {
	GetFeePoolCommunityCoins(sdk.Context) sdk.DecCoins
	// Methods imported from account should be defined here
}

// MintKeeper defines the expected interface for the Mint module.
type MintKeeper interface {
	MintCoins(sdk.Context, sdk.Coins) error
	// Methods imported from account should be defined here
}

// AuthzKeeper defines the expected interface for the Authz module.
type AuthzKeeper interface {
	GetAuthorizations(sdk.Context, sdk.AccAddress, sdk.AccAddress) ([]authz.Authorization, error)
	// Methods imported from account should be defined here
}

// FeegrantKeeper defines the expected interface for the FeeGrant module.
type FeegrantKeeper interface {
	GrantAllowance(sdk.Context, sdk.AccAddress, sdk.AccAddress, feegrant.FeeAllowanceI) error
	// Methods imported from account should be defined here
}

// GroupKeeper defines the expected interface for the Group module.
type GroupKeeper interface {
	GetGroupSequence(sdk.Context) uint64
	// Methods imported from account should be defined here
}

// NftKeeper defines the expected interface for the NFT module.
type NftKeeper interface {
	Mint(sdk.Context, nft.NFT, sdk.AccAddress) error
	// Methods imported from account should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(sdk.Context, []byte, interface{})
	Set(sdk.Context, []byte, interface{})
}
