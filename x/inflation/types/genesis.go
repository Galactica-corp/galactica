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
	"cosmossdk.io/math"
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
		ValidatorsShare: math.LegacyMustNewDecFromStr("1.0"),
		OtherShares:     []*InflationShare{},
	}
}

func DefaultPeriodMintProvisions() []sdk.DecCoin {
	// Cummulative inflation over 36 years is 300,000,000.00 tokens (30% of total supply)
	amountPerPeriod := []string{
		"548402880.90", "448224670.43", "366346280.06", "299424830.09", "244728100.04",
		"200022960.72", "163484230.97", "133620130.87", "109211390.25", "89261450.69",
		"72955820.91", "59628790.37", "48736240.33", "39833460.41", "32556970.92",
		"26609700.92", "21748840.28", "17775920.39", "14528740.86", "11874740.34",
		"9705550.24", "7932610.34", "6483540.18", "5299170.59", "4331160.14",
		"3539970.66", "2893310.97", "2364780.93", "1932800.69", "1579730.59",
		"1291160.14", "1055300.15", "862520.68", "704960.67", "576180.86",
		"470930.47",
	}

	tokensPerPeriod := make([]sdk.DecCoin, len(amountPerPeriod))
	for i, amount := range amountPerPeriod {
		tokensPerPeriod[i] = sdk.NormalizeDecCoin(
			sdk.NewDecCoinFromDec(defaultDenom, math.LegacyMustNewDecFromStr(amount)),
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
