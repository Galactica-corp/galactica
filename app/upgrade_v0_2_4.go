package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	planName = "0.2.4"
)

func (app *App) applyUpgrade_v0_2_4() {
	app.UpgradeKeeper.SetUpgradeHandler(planName, app.upgradeHandler_v0_2_4())
}

func (app *App) upgradeHandler_v0_2_4() func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
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

// for andromeda
const (
	planName_0_2_2 = "0.2.2"
)

func (app *App) applyUpgrade_v0_2_2() {
	app.UpgradeKeeper.SetUpgradeHandler(planName_0_2_2, app.upgradeHandler_v0_2_2())
}

func (app *App) upgradeHandler_v0_2_2() func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
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


// solve 0.1.2 update problem on andromeda
const (
	planName_0_2_7 = "0.2.7"
)

func (app *App) applyUpgrade_v0_2_7() {
	app.UpgradeKeeper.SetUpgradeHandler(planName_0_2_7, app.upgradeHandler_v0_2_7())
}

func (app *App) upgradeHandler_v0_2_7() func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdk.UnwrapSDKContext(ctx).Logger().Info("Upgrade " + plan.Name + " complete")
		return fromVM, nil
	}
}