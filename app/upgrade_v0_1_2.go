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

package app

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/Galactica-corp/galactica/app/upgrades/v0_1_2"
)

// applyUpgrade_v0_1_2 checks and applies the upgrade plan if necessary.
func (app *App) applyUpgrade_v0_1_2() {
	latestBlock := app.LastBlockHeight()
	logger := app.Logger().With("upgrade", v0_1_2.UpgradeName)

	ctx, err := app.CreateQueryContext(latestBlock, false)
	if err != nil {
		logger.Error("Failed to create query context with block", "error", err, "block", latestBlock)
		return
	}

	plan, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil || plan.Height < v0_1_2.UpgradeBlockHeight {
		logger.Info("Applying upgrade plan", "info", plan.Info)
		app.UpgradeKeeper.SetUpgradeHandler(v0_1_2.UpgradeName, app.upgradeHandler_v0_1_2())
		app.UpgradeKeeper.ApplyUpgrade(ctx, v0_1_2.Plan)
	}
}

// upgradeHandler_v0_1_2 returns a handler function for processing the upgrade.
func (app *App) upgradeHandler_v0_1_2() func(
	ctx sdk.Context,
	_ upgradetypes.Plan,
	fromVM module.VersionMap,
) (module.VersionMap, error) {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", plan.Name)

		if plan.Name != v0_1_2.UpgradeName {
			logger.Error("Invalid upgrade plan", "expected", v0_1_2.UpgradeName, "got", plan.Name)
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		}

		validators := app.StakingKeeper.GetAllValidators(ctx)

		for _, validator := range validators {
			if err := app.updateValidatorPowerIndex(ctx, validator); err != nil {
				logger.Error("Failed to update validator power index", "error", err, "validator", validator.OperatorAddress)
				return nil, err
			}
			logger.Info("Validator power index updated", "validator", validator.OperatorAddress)
		}

		logger.Info("All validators updated successfully.")

		if err := app.UpgradeKeeper.DumpUpgradeInfoToDisk(plan.Height, plan); err != nil {
			logger.Error("Failed to dump upgrade info to disk", "error", err)
			return nil, err
		}

		return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
	}
}

// updateValidatorPowerIndex updates the power index for a single validator.
func (app *App) updateValidatorPowerIndex(ctx sdk.Context, validator stakingtypes.Validator) error {
	store := ctx.KVStore(app.GetKey(stakingtypes.StoreKey))
	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.ValidatorsByPowerIndexKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		valAddr := stakingtypes.ParseValidatorPowerRankKey(iterator.Key())
		if bytes.Equal(valAddr, validator.GetOperator()) {
			store.Delete(iterator.Key())
			break // Assuming unique power index key per validator.
		}
	}

	app.StakingKeeper.SetValidatorByPowerIndex(ctx, validator)
	if _, err := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx); err != nil {
		return err
	}

	return nil
}
