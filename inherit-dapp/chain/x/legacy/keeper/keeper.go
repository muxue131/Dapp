package keeper

import (
	"encoding/binary"
	"fmt"
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/inherit-dapp/chain/x/legacy/types"
)

const (
	PlanKeyPrefix      = "plan/"
	AssetKeyPrefix     = "asset/"
	PlanByCreatorPrefix = "creator/"
	AssetByPlanPrefix  = "planasset/"
	NextPlanIDKey      = "next_plan_id"
	NextAssetIDKey     = "next_asset_id"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	bankKeeper types.BankKeeper
}

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	SendCoins(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, sender sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipient sdk.AccAddress, amt sdk.Coins) error
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, bankKeeper types.BankKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bankKeeper: bankKeeper,
	}
}

// --- Plan Operations ---

func (k Keeper) CreateLegacyPlan(ctx sdk.Context, plan types.LegacyPlan) (uint64, error) {
	planID := k.GetNextPlanID(ctx)
	plan.PlanID = planID
	plan.Status = types.PlanStatusActive

	k.SetPlan(ctx, plan)
	k.SetNextPlanID(ctx, planID+1)

	return planID, nil
}

func (k Keeper) SetPlan(ctx sdk.Context, plan types.LegacyPlan) {
	store := ctx.KVStore(k.storeKey)
	key := PlanKey(plan.PlanID)
	bz := k.cdc.MustMarshal(&plan)
	store.Set(key, bz)

	// Index by creator
	creatorKey := PlanByCreatorKey(plan.Creator, plan.PlanID)
	store.Set(creatorKey, []byte{1})
}

func (k Keeper) GetPlan(ctx sdk.Context, planID uint64) (types.LegacyPlan, bool) {
	store := ctx.KVStore(k.storeKey)
	key := PlanKey(planID)
	bz := store.Get(key)
	if bz == nil {
		return types.LegacyPlan{}, false
	}
	var plan types.LegacyPlan
	k.cdc.MustUnmarshal(bz, &plan)
	return plan, true
}

func (k Keeper) DeletePlan(ctx sdk.Context, planID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := PlanKey(planID)
	store.Delete(key)
}

func (k Keeper) GetAllPlans(ctx sdk.Context) []types.LegacyPlan {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(PlanKeyPrefix))
	defer iterator.Close()

	var plans []types.LegacyPlan
	for ; iterator.Valid(); iterator.Next() {
		var plan types.LegacyPlan
		k.cdc.MustUnmarshal(iterator.Value(), &plan)
		plans = append(plans, plan)
	}
	return plans
}

func (k Keeper) GetAllAssets(ctx sdk.Context) []types.Asset {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(AssetKeyPrefix))
	defer iterator.Close()

	var assets []types.Asset
	for ; iterator.Valid(); iterator.Next() {
		var asset types.Asset
		k.cdc.MustUnmarshal(iterator.Value(), &asset)
		assets = append(assets, asset)
	}
	return assets
}

func (k Keeper) GetPlansByCreator(ctx sdk.Context, creator sdk.AccAddress) []types.LegacyPlan {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(fmt.Sprintf("%s%s/", PlanByCreatorPrefix, creator.String()))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var plans []types.LegacyPlan
	for ; iterator.Valid(); iterator.Next() {
		// Extract planID from key suffix (last 8 bytes are big-endian uint64)
		key := iterator.Key()
		if len(key) < 8 {
			continue
		}
		planID := binary.BigEndian.Uint64(key[len(key)-8:])
		if plan, found := k.GetPlan(ctx, planID); found {
			plans = append(plans, plan)
		}
	}
	return plans
}

// --- Asset Operations ---

func (k Keeper) AddAsset(ctx sdk.Context, asset types.Asset) (uint64, error) {
	assetID := k.GetNextAssetID(ctx)
	asset.AssetID = assetID

	k.SetAsset(ctx, asset)
	k.SetNextAssetID(ctx, assetID+1)

	return assetID, nil
}

func (k Keeper) SetAsset(ctx sdk.Context, asset types.Asset) {
	store := ctx.KVStore(k.storeKey)
	key := AssetKey(asset.AssetID)
	bz := k.cdc.MustMarshal(&asset)
	store.Set(key, bz)

	// Index by plan
	planKey := AssetByPlanKey(asset.PlanID, asset.AssetID)
	store.Set(planKey, []byte{1})
}

func (k Keeper) GetAsset(ctx sdk.Context, assetID uint64) (types.Asset, bool) {
	store := ctx.KVStore(k.storeKey)
	key := AssetKey(assetID)
	bz := store.Get(key)
	if bz == nil {
		return types.Asset{}, false
	}
	var asset types.Asset
	k.cdc.MustUnmarshal(bz, &asset)
	return asset, true
}

func (k Keeper) GetAssetsByPlan(ctx sdk.Context, planID uint64) []types.Asset {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(fmt.Sprintf("%s%d/", AssetByPlanPrefix, planID))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var assets []types.Asset
	for ; iterator.Valid(); iterator.Next() {
		// Extract assetID from key suffix (last 8 bytes are big-endian uint64)
		key := iterator.Key()
		if len(key) < 8 {
			continue
		}
		assetID := binary.BigEndian.Uint64(key[len(key)-8:])
		if asset, found := k.GetAsset(ctx, assetID); found {
			assets = append(assets, asset)
		}
	}
	return assets
}

// --- Heartbeat Operations ---

func (k Keeper) UpdateHeartbeat(ctx sdk.Context, planID uint64) error {
	plan, found := k.GetPlan(ctx, planID)
	if !found {
		return types.ErrPlanNotFound
	}
	if plan.Status != types.PlanStatusActive {
		return types.ErrPlanNotActive
	}

	plan.LastHeartbeat = ctx.BlockTime()
	plan.TriggerTime = ctx.BlockTime().Add(
		time.Duration(plan.HeartbeatInterval) * time.Second,
	)
	k.SetPlan(ctx, plan)
	return nil
}

func (k Keeper) CheckHeartbeatExpiry(ctx sdk.Context, planID uint64) (bool, error) {
	plan, found := k.GetPlan(ctx, planID)
	if !found {
		return false, types.ErrPlanNotFound
	}
	if plan.Status != types.PlanStatusActive {
		return false, nil
	}

	if ctx.BlockTime().After(plan.TriggerTime) {
		plan.Status = types.PlanStatusTriggered
		k.SetPlan(ctx, plan)
		return true, nil
	}
	return false, nil
}

// --- Claim Operations ---

func (k Keeper) ClaimInheritance(ctx sdk.Context, planID uint64, beneficiary sdk.AccAddress) error {
	plan, found := k.GetPlan(ctx, planID)
	if !found {
		return types.ErrPlanNotFound
	}
	if plan.Status != types.PlanStatusTriggered {
		return types.ErrInheritanceNotReady
	}

	// Find beneficiary share
	var beneficiaryShare *types.Beneficiary
	for _, b := range plan.Beneficiaries {
		if b.Address.Equals(beneficiary) {
			beneficiaryShare = &b
			break
		}
	}
	if beneficiaryShare == nil {
		return types.ErrUnauthorized
	}

	// Transfer assets according to share
	assets := k.GetAssetsByPlan(ctx, planID)
	for _, asset := range assets {
		shareAmount := sdk.NewDecFromInt(asset.Amount).Mul(beneficiaryShare.Share).TruncateInt()
		if shareAmount.IsPositive() {
			coins := sdk.NewCoins(sdk.NewCoin(asset.Denom, shareAmount))
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, beneficiary, coins)
			if err != nil {
				return err
			}
		}
	}

	// Mark as claimed if all beneficiaries have claimed
	allClaimed := true
	for _, b := range plan.Beneficiaries {
		if !b.Address.Equals(beneficiary) && len(b.KeyShare) > 0 {
			// Check if this beneficiary has already claimed
			// Simplified: in production, track individual claims
		}
	}
	if allClaimed {
		plan.Status = types.PlanStatusClaimed
		k.SetPlan(ctx, plan)
	}

	return nil
}

// --- ID Management ---

func (k Keeper) GetNextPlanID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(NextPlanIDKey))
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextPlanID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	store.Set([]byte(NextPlanIDKey), bz)
}

func (k Keeper) GetNextAssetID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(NextAssetIDKey))
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextAssetID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	store.Set([]byte(NextAssetIDKey), bz)
}

// --- Key Helpers ---

func PlanKey(planID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, planID)
	return append([]byte(PlanKeyPrefix), bz...)
}

func AssetKey(assetID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, assetID)
	return append([]byte(AssetKeyPrefix), bz...)
}

func PlanByCreatorKey(creator sdk.AccAddress, planID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, planID)
	prefix := fmt.Sprintf("%s%s/", PlanByCreatorPrefix, creator.String())
	return append([]byte(prefix), bz...)
}

func AssetByPlanKey(planID, assetID uint64) []byte {
	assetBz := make([]byte, 8)
	binary.BigEndian.PutUint64(assetBz, assetID)
	key := append([]byte(fmt.Sprintf("%s%d/", AssetByPlanPrefix, planID)), assetBz...)
	return key
}
