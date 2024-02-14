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

package inflation

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Galactica-corp/galactica/x/inflation/keeper"
	"github.com/Galactica-corp/galactica/x/inflation/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set genesis state
	params := genState.Params
	// this line is used by starport scaffolding # genesis/module/init
	err := k.SetParams(ctx, params)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting params"))
	}

	period := genState.Period
	k.SetPeriod(ctx, period)

	epochIdentifier := genState.EpochIdentifier
	k.SetEpochIdentifier(ctx, epochIdentifier)

	epochsPerPeriod := genState.EpochsPerPeriod
	k.SetEpochsPerPeriod(ctx, epochsPerPeriod)

	skippedEpochs := genState.SkippedEpochs
	k.SetSkippedEpochs(ctx, skippedEpochs)

	periodMintProvisions := genState.PeriodMintProvisions
	err = k.SetPeriodMintProvisions(ctx, periodMintProvisions)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting period mint provisions"))
	}

	inflationDistribution := genState.InflationDistribution
	err = k.SetInflationDistribution(ctx, inflationDistribution)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting inflation distribution"))
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	periodMintProvisions, err := k.GetPeriodMintProvisions(ctx)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error getting period mint provisions"))
	}

	return &types.GenesisState{
		Params:          k.GetParams(ctx),
		Period:          k.GetPeriod(ctx),
		EpochIdentifier: k.GetEpochIdentifier(ctx),
		EpochsPerPeriod: k.GetEpochsPerPeriod(ctx),
		SkippedEpochs:   k.GetSkippedEpochs(ctx),

		PeriodMintProvisions: periodMintProvisions,
	}
}
