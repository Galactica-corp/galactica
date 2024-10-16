package v0_2_4

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
)

const (
	UpgradeName        = "0.2.4"
	UpgradeBlockHeight = 1
)

// Plan defines the upgrade plan for addressing the staking PowerReduction issue.
var Plan = upgradetypes.Plan{
	Name:   UpgradeName,
	Height: UpgradeBlockHeight,
	Info:   "migration " + UpgradeName,
}
