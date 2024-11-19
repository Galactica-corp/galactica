package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	planName_v0_2_7 = "0.2.7"
)

func (app *App) applyUpgrade_v0_2_7() {
	app.UpgradeKeeper.SetUpgradeHandler(planName_v0_2_7, app.upgradeHandler_v0_2_7())
}

func (app *App) upgradeHandler_v0_2_7() func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := sdk.UnwrapSDKContext(ctx).Logger()

		logger.Info("Starting module migrations...")

		vm, err := app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		if err != nil {
			return vm, err
		}

		logger.Info("Upgrade " + plan.Name + " complete")

		return vm, err
	}
}
