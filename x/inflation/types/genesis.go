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

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # genesis/types/import
	epochstypes "github.com/Galactica-corp/galactica/x/epochs/types"
)

// DefaultIndex is the default global index
const (
	DefaultIndex uint64 = 1

	defaultDenom = "gnet"
)

func DefaultInflationDistribution() InflationDistribution {
	return InflationDistribution{
		// If no other shares are specified, validators get 100% of the inflation
		ValidatorsShare: sdk.MustNewDecFromStr("1.0"),
		OtherShares:     []*InflationShare{},
	}
}

func DefaultPeriodMintProvisions() []sdk.DecCoin {
	// Cummulative inflation over 36 years is 300,000,000.00 tokens (30% of total supply)
	amountPerPeriod := []string{
		"54840288.90", "44822467.43", "36634628.06", "29942483.09", "24472810.04",
		"20002296.72", "16348423.97", "13362013.87", "10921139.25", "8926145.69",
		"7295582.91", "5962879.37", "4873624.33", "3983346.41", "3255697.92",
		"2660970.92", "2174884.28", "1777592.39", "1452874.86", "1187474.34",
		"970555.24", "793261.34", "648354.18", "529917.59", "433116.14",
		"353997.66", "289331.97", "236478.93", "193280.69", "157973.59",
		"129116.14", "105530.15", "86252.68", "70496.67", "57618.86",
		"47093.47",
	}

	tokensPerPeriod := make([]sdk.DecCoin, len(amountPerPeriod))
	for i, amount := range amountPerPeriod {
		tokensPerPeriod[i] = sdk.NormalizeDecCoin(
			sdk.NewDecCoinFromDec(defaultDenom, sdk.MustNewDecFromStr(amount)),
		)
	}

	return tokensPerPeriod
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:          DefaultParams(),
		Period:          uint64(0),
		EpochIdentifier: epochstypes.DayEpochID,
		EpochsPerPeriod: 365,
		SkippedEpochs:   0,

		PeriodMintProvisions:  DefaultPeriodMintProvisions(),
		InflationDistribution: DefaultInflationDistribution(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
