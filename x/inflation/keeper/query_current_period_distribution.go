// Copyright 2025 Galactica Network
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Galactica-corp/galactica/x/inflation/types"
)

// CurrentEpochDistribution returns distribution of minted inflation tokens for the current epoch.
func (k Keeper) CurrentEpochDistribution(
	goCtx context.Context,
	_ *types.QueryCurrentEpochDistributionRequest,
) (*types.QueryCurrentEpochDistributionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	periodMintProvisions, err := k.GetPeriodMintProvisions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "get period mint provisions")
	}

	period := k.GetPeriod(ctx)
	epochsPerPeriod := k.GetEpochsPerPeriod(ctx)

	params := k.GetParams(ctx)
	amount := types.CalculateEpochMintProvision(periodMintProvisions, period, epochsPerPeriod)

	mintedCoin := sdk.NewCoin(params.MintDenom, amount.TruncateInt())

	distribution, found := k.GetInflationDistribution(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "inflation distribution not found")
	}

	var currentDistribution types.CurrentEpochInflationDistribution
	currentDistribution.ValidatorsShare = k.GetProportions(ctx, mintedCoin, distribution.ValidatorsShare)

	for _, share := range distribution.OtherShares {
		currentDistribution.OtherShares[share.Name] = types.CurrentEpochInflationShare{
			Address: share.Address,
			Name:    share.Name,
			Share:   k.GetProportions(ctx, mintedCoin, share.Share),
		}
	}

	epochID := k.GetEpochIdentifier(ctx)

	epochInfo, ok := k.epochsKeeper.GetEpochInfo(ctx, epochID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "epoch %s not found", epochID)
	}

	return &types.QueryCurrentEpochDistributionResponse{
		Distribution:  currentDistribution,
		EpochDuration: epochInfo.Duration,
	}, nil
}
