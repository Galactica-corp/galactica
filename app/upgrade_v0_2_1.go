package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/Galactica-corp/galactica/app/upgrades/v0_2_1"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func (app *App) applyUpgrade_v0_2_1() {
	latestBlock := app.LastBlockHeight()
	logger := app.Logger().With("upgrade", v0_2_1.UpgradeName)

	ctx, err := app.CreateQueryContext(latestBlock, false)
	if err != nil {
		logger.Error("Failed to create query context with block", "error", err, "block", latestBlock)
		return
	}

	plan, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		logger.Error("Error reading an upgrade plan err:", err)
	}
	
	if err != nil || plan.Height < v0_2_1.UpgradeBlockHeight {
		logger.Info("Applying upgrade plan", "info", plan.Info)

		app.UpgradeKeeper.SetUpgradeHandler(v0_2_1.UpgradeName, app.upgradeHandler_v0_2_1())
		app.UpgradeKeeper.ApplyUpgrade(ctx, v0_2_1.Plan)

		if err := app.UpgradeKeeper.DumpUpgradeInfoToDisk(plan.Height, plan); err != nil {
			logger.Error("Failed to dump upgrade info to disk", "error", err)
		}
	}
}

func (app *App) upgradeHandler_v0_2_1() func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
	}
}
