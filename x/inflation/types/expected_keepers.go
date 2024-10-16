// Galactica is a Layer 1 protocol with zero-knowledge and privacy features.
// Copyright (C) 2024 Galactica Network
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types

import (
	context "context"

	"cosmossdk.io/math"
	"cosmossdk.io/x/feegrant"
	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

type DistrKeeper interface {
	// TODO Add methods imported from distr should be defined here
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error // TODO: Удалить, если не нужно
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	HasSupply(ctx context.Context, denom string) bool
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

// StakingKeeper defines the expected interface for the Staking module.
type StakingKeeper interface {
	TotalBondedTokens(context.Context) math.Int
	// BondedRatio the fraction of the staking tokens which are currently bonded
	BondedRatio(ctx context.Context) math.LegacyDec
	StakingTokenSupply(ctx context.Context) math.Int
}

// SlashingKeeper defines the expected interface for the Slashing module.
type SlashingKeeper interface {
	Slash(ctx context.Context, consAddr sdk.ConsAddress, fraction math.LegacyDec, power, distributionHeight int64)
	// Methods imported from account should be defined here
}

// DistributionKeeper defines the expected interface for the Distribution module.
type DistributionKeeper interface {
	GetFeePoolCommunityCoins(context.Context) sdk.DecCoins
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// MintKeeper defines the expected interface for the Mint module.
type MintKeeper interface {
	MintCoins(context.Context, sdk.Coins) error
	// Methods imported from account should be defined here
}

// AuthzKeeper defines the expected interface for the Authz module.
type AuthzKeeper interface {
	GetAuthorizations(context.Context, sdk.AccAddress, sdk.AccAddress) ([]authz.Authorization, error)
	// Methods imported from account should be defined here
}

// FeegrantKeeper defines the expected interface for the FeeGrant module.
type FeegrantKeeper interface {
	GrantAllowance(context.Context, sdk.AccAddress, sdk.AccAddress, feegrant.FeeAllowanceI) error
	// Methods imported from account should be defined here
}

// GroupKeeper defines the expected interface for the Group module.
type GroupKeeper interface {
	GetGroupSequence(context.Context) uint64
	// Methods imported from account should be defined here
}

// NftKeeper defines the expected interface for the NFT module.
type NftKeeper interface {
	Mint(context.Context, nft.NFT, sdk.AccAddress) error
	// Methods imported from account should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
