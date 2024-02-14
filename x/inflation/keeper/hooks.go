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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	epochstypes "github.com/Galactica-corp/galactica/x/epochs/types"
	"github.com/Galactica-corp/galactica/x/inflation/types"
)

// TODO: switch BeforeEpochStart and AfterEpochEnd
func (k Keeper) AfterEpochEnd(ctx sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd mints and allocates coins at the end of each epoch end
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	params := k.GetParams(ctx)

	expEpochID := k.GetEpochIdentifier(ctx)
	if epochIdentifier != expEpochID {
		return
	}

	period := k.GetPeriod(ctx)
	epochsPerPeriod := k.GetEpochsPerPeriod(ctx)
	periodMintProvisions, err := k.GetPeriodMintProvisions(ctx)
	if err != nil {
		k.Logger(ctx).Error(
			"SKIPPING INFLATION: error getting period mint provisions",
			"error", err.Error(),
		)
		return
	}

	epochMintProvision := types.CalculateEpochMintProvision(
		periodMintProvisions,
		period,
		epochsPerPeriod,
	)

	if !epochMintProvision.IsPositive() {
		k.Logger(ctx).Error(
			"SKIPPING INFLATION: negative epoch mint provision",
			"value", epochMintProvision.String(),
		)
		return
	}

	mintedCoin := sdk.Coin{
		Denom:  params.MintDenom,
		Amount: epochMintProvision.TruncateInt(),
	}
	// skip as no coins need to be minted
	if mintedCoin.Amount.IsNil() || !mintedCoin.Amount.IsPositive() {
		return
	}
	err = k.MintCoins(ctx, mintedCoin)
	if err != nil {
		k.Logger(ctx).Error(
			"SKIPPING INFLATION: error minting coins",
			"error", err.Error(),
			"coin", mintedCoin.String(),
			"denom", mintedCoin.Denom,
		)
		return
	}

	// Allocate staking rewards into fee collector account
	distribution, found := k.GetInflationDistribution(ctx)
	if !found {
		k.Logger(ctx).Error("SKIPPING INFLATION: inflation distribution not found")
		return
	}

	k.Logger(ctx).With(
		"ValidatorsShare", distribution.ValidatorsShare.String(),
		"OtherSharesLen", len(distribution.OtherShares),
	).Info("INFLATION MODULE: distribution")

	if len(distribution.OtherShares) > 0 {
		k.Logger(ctx).With(
			"OtherShares[0].Name", distribution.OtherShares[0].Name,
			"OtherShares[0].Address", distribution.OtherShares[0].Address,
			"OtherShares[0].Share", distribution.OtherShares[0].Share.String(),
		).Info("INFLATION MODULE: OtherShares")
	}

	staking := sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.ValidatorsShare)}
	if !staking.IsZero() {
		if err := k.bankKeeper.SendCoinsFromModuleToModule(
			ctx,
			types.ModuleName,
			authtypes.FeeCollectorName,
			staking,
		); err != nil {
			k.Logger(ctx).Error(
				"SKIPPING INFLATION: error sending coins to validators from module to module",
				"error", err.Error(),
			)
			return
		} else {
			k.Logger(ctx).Info(
				"INFLATION MODULE: sent coins to validators from module to module",
				"coins", staking.String(),
			)
		}
	}

	// allocate inflation to other addresses
	for _, share := range distribution.OtherShares {
		other := sdk.Coins{k.GetProportions(ctx, mintedCoin, share.Share)}
		otherAddress, err := sdk.AccAddressFromBech32(share.Address)
		if err != nil {
			k.Logger(ctx).Error(
				"SKIPPING INFLATION: error getting address from bech32",
				"error", err.Error(),
			)
			return
		}

		if !other.IsZero() {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				types.ModuleName,
				otherAddress,
				other,
			); err != nil {
				k.Logger(ctx).Error(
					"SKIPPING INFLATION: error sending coins to other from module",
					"error", err.Error(),
				)
			} else {
				k.Logger(ctx).Info(
					"INFLATION MODULE: sent coins to other from module",
					"coins", other.String(),
					"otherAddress", otherAddress.String(),
				)
			}

			// if err := k.distrKeeper.FundCommunityPool(
			//	ctx,
			//	other,
			//	otherAddress,
			// ); err != nil {
			//	k.Logger(ctx).Error(
			//		"SKIPPING INFLATION: error sending coins to other from module to module",
			//		"error", err.Error(),
			//	)
			//	return
			// }

		} else {
			k.Logger(ctx).Info(
				"INFLATION MODULE: other is zero",
			)
		}
	}

	// If period is passed, update the period. A period is
	// passed if the current epoch number surpasses the epochsPerPeriod for the
	// current period.
	//
	// Examples:
	// Given, epochNumber = 1, period = 0, epochPerPeriod = 365
	//   => 1 - 365 * 0 - 0 < 365 --- nothing to do here
	// Given, epochNumber = 741, period = 1, epochPerPeriod = 365
	//   => 741 - 1 * 365 - 10 > 365 --- a period has passed! we set a new period
	if epochNumber-epochsPerPeriod*int64(period) > epochsPerPeriod {
		period++
		k.SetPeriod(ctx, period)
	}

	// TODO: telemetry

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochNumber)),
			sdk.NewAttribute(types.AttributeKeyEpochProvisions, epochMintProvision.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)

	k.Logger(ctx).Info("******* AFTER EPOCH END *******")
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
