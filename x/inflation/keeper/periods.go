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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"

	"github.com/Galactica-corp/galactica/x/inflation/types"
)

// GetPeriod gets current period
func (k Keeper) GetPeriod(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PeriodKey)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetPeriod stores the current period
func (k Keeper) SetPeriod(ctx sdk.Context, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PeriodKey, sdk.Uint64ToBigEndian(period))
}

func (k Keeper) GetPeriodMintProvisions(ctx sdk.Context) (sdk.DecCoins, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PeriodMintProvisionsKey)
	if len(bz) == 0 {
		return sdk.NewDecCoins(), nil
	}

	var genesisState types.GenesisState
	err := proto.Unmarshal(bz, &genesisState)
	if err != nil {
		return nil, err
	}

	return genesisState.PeriodMintProvisions, nil
}

func (k Keeper) SetPeriodMintProvisions(ctx sdk.Context, provisions sdk.DecCoins) error {
	store := ctx.KVStore(k.storeKey)

	genesisState := types.GenesisState{
		PeriodMintProvisions: provisions,
	}

	bz, err := proto.Marshal(&genesisState)
	if err != nil {
		return err
	}

	store.Set(types.PeriodMintProvisionsKey, bz)
	return nil
}
