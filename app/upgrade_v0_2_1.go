package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/Galactica-corp/galactica/app/upgrades/v0_2_1"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func (app *App) applyUpgrade_v0_2_1() {
	logger := app.Logger().With("upgrade", v0_2_1.UpgradeName)

	blockHeader := tmproto.Header{Height: app.LastBlockHeight()}
	commitStore := app.CommitMultiStore()
	isCheckTx := false

	ctx := app.NewContext(isCheckTx).
		WithBlockHeader(blockHeader).
		WithMultiStore(commitStore)

	if v0_2_1.Plan.Height == app.LastBlockHeight() {
		logger.Info("Schedule upgrade plan", "info", v0_2_1.Plan.Info)

		if err := app.UpgradeKeeper.ScheduleUpgrade(ctx, v0_2_1.Plan); err != nil {
			panic(err)
		}

		planName := v0_2_1.UpgradeName
		app.UpgradeKeeper.SetUpgradeHandler(planName, app.upgradeHandler_v0_2_1())
	}
}

func (app *App) upgradeHandler_v0_2_1() func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
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
