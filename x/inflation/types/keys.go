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

const (
	// ModuleName defines the module name
	ModuleName = "inflation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_inflation"
)

var (
	ParamsKey                = []byte("params_inflation")
	PeriodKey                = []byte("period_inflation")
	EpochIdentifierKey       = []byte("epoch_identifier_inflation")
	EpochsPerPeriodKey       = []byte("epochs_per_period_inflation")
	SkippedEpochsKey         = []byte("skipped_epochs_inflation")
	PeriodMintProvisionsKey  = []byte("period_mint_provisions_inflation")
	InflationDistributionKey = []byte("inflation_distribution_inflation")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
