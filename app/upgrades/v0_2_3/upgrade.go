package v0_2_3

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
)

const (
	UpgradeName        = "0.2.3"
	UpgradeBlockHeight = 1
)

// Plan defines the upgrade plan for addressing the staking PowerReduction issue.
var Plan = upgradetypes.Plan{
	Name:   UpgradeName,
	Height: UpgradeBlockHeight,
}
