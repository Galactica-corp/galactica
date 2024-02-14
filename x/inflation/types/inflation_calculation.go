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
)

// CalculateEpochMintProvision returns mint provision per epoch
func CalculateEpochMintProvision(
	periodMintProvisions sdk.DecCoins,
	period uint64,
	epochsPerPeriod int64,
) sdk.Dec {
	if period > uint64(len(periodMintProvisions)) {
		return sdk.ZeroDec()
	}

	periodProvision := periodMintProvisions[period]
	epochsPerPeriodDec := sdk.NewDec(epochsPerPeriod)
	epochProvision := periodProvision.Amount.Quo(epochsPerPeriodDec)

	// Multiply epochProvision with power reduction (10^18 for gnet) as the calculation
	// is based on `gnet` and the issued tokens need to be given in `agnet`
	// epochProvision = epochProvision.Mul(sdk.NewDecFromInt(sdk.DefaultPowerReduction))
	return epochProvision
}
