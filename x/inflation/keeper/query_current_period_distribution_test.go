// Copyright 2025 Galactica Network
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

package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Galactica-corp/galactica/testutil/keeper"
	"github.com/Galactica-corp/galactica/x/inflation"
	"github.com/Galactica-corp/galactica/x/inflation/types"
)

func TestKeeper_CurrentEpochDistribution(t *testing.T) {
	// BUG: Write some actual tests after the fix. Currently, tests are broken:
	//		cosmossdk.io/store/iavl.(*Store).Get(0xc000744a00?, {0x10c217d00?, 0x10ab2f507?, 0x10b321c28?})
	//		  ~/go/pkg/mod/github.com/crypto-org-chain/cosmos-sdk/store@v0.0.0-20240415105151-0108877a3201/iavl/store.go:191 +0x36
	//		cosmossdk.io/store/gaskv.(*GStore[...]).Get(0x1e, {0x10c217d00?, 0x20, 0x20})
	//		  ~/go/pkg/mod/github.com/crypto-org-chain/cosmos-sdk/store@v0.0.0-20240415105151-0108877a3201/gaskv/store.go:65 +0x68
	//		github.com/Galactica-corp/galactica/x/inflation/keeper.Keeper.GetPeriodMintProvisions({{0x10b3274c0, 0xc00074f3b0}, {0x10b303ef0, 0xc00074f290}, {0x10b303f18, 0xc00074f2a0}, {0xc0001b93b0, 0x2d}, {0x10b325a40, 0xc00059d698}, ...}, ...)
	//		  ~/galactica/x/inflation/keeper/periods.go:45 +0x7d

	keeper, _, _, _, ctx := keepertest.InflationKeeper(t)

	inflation.InitGenesis(ctx, keeper, *types.DefaultGenesis())

	_, err := keeper.CurrentEpochDistribution(ctx, nil)
	require.NoError(t, err)
}
