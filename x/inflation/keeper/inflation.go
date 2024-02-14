// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Galactica-corp/galactica/x/inflation/types"
)

// AllocateInflation allocates coins from the inflation to external
// modules according to allocation proportions:
//   - staking rewards -> sdk `auth` module fee collector
//   - usage incentives -> `x/incentives` module
//   - community pool -> `sdk `distr` module community pool
func (k Keeper) AllocateInflation(
	ctx sdk.Context,
	mintedCoin sdk.Coin,
	params types.Params,
) (
	validatorsCoins sdk.Coins,
	otherCoins []sdk.Coins,
	// staking, incentives, communityPool sdk.Coins,
	err error,
) {
	distribution := params.InflationDistribution

	// Allocate validators rewards
	validatorsCoins = sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.ValidatorsShare)}

	// if err := k.bankKeeper.SendCoinsFromModuleToModule(
	// 	ctx,
	// 	types.ModuleName,
	// 	k.feeCollectorName,
	// 	staking,
	// ); err != nil {
	// 	return nil, nil, nil, err
	// }

	// // Allocate usage incentives to incentives module account
	// incentives = sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.UsageIncentives)}

	// if err = k.bankKeeper.SendCoinsFromModuleToModule(
	// 	ctx,
	// 	types.ModuleName,
	// 	incentivestypes.ModuleName,
	// 	incentives,
	// ); err != nil {
	// 	return nil, nil, nil, err
	// }

	// // Allocate community pool amount (remaining module balance) to community
	// // pool address
	// moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	// inflationBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddr)

	// err = k.distrKeeper.FundCommunityPool(
	// 	ctx,
	// 	inflationBalance,
	// 	moduleAddr,
	// )
	// if err != nil {
	// 	return nil, nil, nil, err
	// }

	// return staking, incentives, communityPool, nil

	return validatorsCoins, otherCoins, nil
}

// GetProportion calculates the proportion of coins that is to be
// allocated during inflation for a given distribution.
func (k Keeper) GetProportions(
	_ sdk.Context,
	coin sdk.Coin,
	distribution sdk.Dec,
) sdk.Coin {
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: sdk.NewDecFromInt(coin.Amount).Mul(distribution).TruncateInt(),
	}
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, coin sdk.Coin) error {
	coins := sdk.Coins{coin}
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
}

func (k Keeper) SetInflationDistribution(ctx sdk.Context, dist types.InflationDistribution) error {
	store := ctx.KVStore(k.storeKey)
	b, err := k.cdc.Marshal(&dist)
	if err != nil {
		return err
	}
	store.Set(types.InflationDistributionKey, b)
	return nil
}

func (k Keeper) GetInflationDistribution(ctx sdk.Context) (dist types.InflationDistribution, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.InflationDistributionKey)
	if b == nil {
		return dist, false
	}
	err := k.cdc.Unmarshal(b, &dist)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal inflation distribution: %s", err))
	}
	return dist, true
}
