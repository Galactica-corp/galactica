package v0_2_2

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
)

const (
	UpgradeName        = "0.2.2"
	UpgradeBlockHeight = 10
)

// Plan defines the upgrade plan for addressing the staking PowerReduction issue.
var Plan = upgradetypes.Plan{
	Name:   UpgradeName,
	Height: UpgradeBlockHeight,
	Info:   "migration " + UpgradeName,
}
