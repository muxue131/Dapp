package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/inherit-dapp/chain/x/legacy/types"
)

// QueryPlan returns a plan by ID
func (k Keeper) QueryPlan(ctx sdk.Context, planID uint64) (types.LegacyPlan, error) {
	plan, found := k.GetPlan(ctx, planID)
	if !found {
		return types.LegacyPlan{}, types.ErrPlanNotFound
	}
	return plan, nil
}

// QueryPlansByCreator returns all plans for a creator
func (k Keeper) QueryPlansByCreator(ctx sdk.Context, creator sdk.AccAddress) []types.LegacyPlan {
	return k.GetPlansByCreator(ctx, creator)
}

// QueryAsset returns an asset by ID
func (k Keeper) QueryAsset(ctx sdk.Context, assetID uint64) (types.Asset, error) {
	asset, found := k.GetAsset(ctx, assetID)
	if !found {
		return types.Asset{}, types.ErrAssetNotFound
	}
	return asset, nil
}

// QueryAssetsByPlan returns all assets for a plan
func (k Keeper) QueryAssetsByPlan(ctx sdk.Context, planID uint64) []types.Asset {
	return k.GetAssetsByPlan(ctx, planID)
}

// QueryPlanWithAssets returns a plan with its assets
func (k Keeper) QueryPlanWithAssets(ctx sdk.Context, planID uint64) (types.PlanResponse, error) {
	plan, found := k.GetPlan(ctx, planID)
	if !found {
		return types.PlanResponse{}, types.ErrPlanNotFound
	}
	assets := k.GetAssetsByPlan(ctx, planID)
	return types.PlanResponse{
		Plan:   plan,
		Assets: assets,
	}, nil
}

// QueryAllPlans returns all plans
func (k Keeper) QueryAllPlans(ctx sdk.Context) []types.LegacyPlan {
	return k.GetAllPlans(ctx)
}
