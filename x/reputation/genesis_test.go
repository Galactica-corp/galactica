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

package reputation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/Galactica-corp/galactica/testutil/keeper"
	"github.com/Galactica-corp/galactica/testutil/nullify"
	"github.com/Galactica-corp/galactica/x/reputation"
	"github.com/Galactica-corp/galactica/x/reputation/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ReputationKeeper(t)
	reputation.InitGenesis(ctx, k, genesisState)
	got := reputation.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
