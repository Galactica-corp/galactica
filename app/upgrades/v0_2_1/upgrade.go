package v0_2_1

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
)

const (
	UpgradeName        = "0.2.1"
	// UpgradeBlockHeight = 4_183_890
	UpgradeBlockHeight = 100
)

// Plan defines the upgrade plan for addressing the staking PowerReduction issue.
var Plan = upgradetypes.Plan{
	Name:   UpgradeName,
	Height: UpgradeBlockHeight,
	Info:   "gov module migration",
}
