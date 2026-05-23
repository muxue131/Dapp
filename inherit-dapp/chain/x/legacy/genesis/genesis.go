package genesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/inherit-dapp/chain/x/legacy/keeper"
	"github.com/inherit-dapp/chain/x/legacy/types"
)

// InitGenesis initializes the legacy module's state from a genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetNextPlanID(ctx, data.NextPlanID)
	k.SetNextAssetID(ctx, data.NextAssetID)

	for _, plan := range data.Plans {
		k.SetPlan(ctx, plan)
	}
	for _, asset := range data.Assets {
		k.SetAsset(ctx, asset)
	}
}

// ExportGenesis exports the legacy module's state to a genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	plans := k.GetAllPlans(ctx)
	assets := k.GetAllAssets(ctx)
	nextPlanID := k.GetNextPlanID(ctx)
	nextAssetID := k.GetNextAssetID(ctx)

	return types.GenesisState{
		Plans:       plans,
		Assets:      assets,
		NextPlanID:  nextPlanID,
		NextAssetID: nextAssetID,
	}
}
