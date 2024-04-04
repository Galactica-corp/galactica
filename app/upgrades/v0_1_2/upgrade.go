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

package v0_1_2

import (
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName        = "0.1.2"
	UpgradeBlockHeight = 16951
)

// Plan defines the upgrade plan for addressing the staking PowerReduction issue.
var Plan = upgradetypes.Plan{
	Name:   UpgradeName,
	Height: UpgradeBlockHeight,
	Info: "Addresses a critical staking PowerReduction issue by mutating validators' " +
		"power for accurate voting power recalibration and network integrity.",
}
