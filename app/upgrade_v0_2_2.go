package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/Galactica-corp/galactica/app/upgrades/v0_2_2"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func (app *App) applyUpgrade_v0_2_2() {
	planName := v0_2_2.UpgradeName

	logger := app.Logger().With("upgrade", planName)

	app.UpgradeKeeper.SetUpgradeHandler(planName, app.upgradeHandler_v0_2_2())
	
	ctx := app.NewContext(true).WithBlockHeader(tmproto.Header{Height: app.LastBlockHeight()}).
		WithMultiStore(app.CommitMultiStore())

	doneHeight, err := app.UpgradeKeeper.GetDoneHeight(ctx, planName)
	if err != nil {
		logger.Error("Error with GetDoneHeight", "error", err)
		return
	}

	if doneHeight != 0 {
		logger.Info("upgrade v0.2.2 done")
		return
	}

	logger.Info("Schedule upgrade plan", "name", planName)

	if err := app.UpgradeKeeper.ScheduleUpgrade(ctx, v0_2_2.Plan); err != nil {
		panic(err)
	}
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
