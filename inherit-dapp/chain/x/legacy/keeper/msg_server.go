package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/inherit-dapp/chain/x/legacy/types"
)

// HandleMsgCreateLegacyPlan handles the create legacy plan message
func (k Keeper) HandleMsgCreateLegacyPlan(ctx sdk.Context, msg types.MsgCreateLegacyPlan) (*sdk.Result, error) {
	beneficiaries := make([]types.Beneficiary, len(msg.Beneficiaries))
	for i, b := range msg.Beneficiaries {
		beneficiaries[i] = types.Beneficiary{
			Address: b.Address,
			Share:   b.Share,
		}
	}

	plan := types.LegacyPlan{
		Creator:           msg.Creator,
		Beneficiaries:     beneficiaries,
		HeartbeatInterval: msg.HeartbeatInterval,
		LastHeartbeat:     ctx.BlockTime(),
		TriggerTime:       ctx.BlockTime().Add(time.Duration(msg.HeartbeatInterval) * time.Second),
		CreatedAt:         ctx.BlockTime(),
	}

	planID, err := k.CreateLegacyPlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
		sdk.NewEvent(
			"create_legacy_plan",
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", planID)),
			sdk.NewAttribute("creator", msg.Creator.String()),
			sdk.NewAttribute("heartbeat_interval", fmt.Sprintf("%d", msg.HeartbeatInterval)),
			sdk.NewAttribute("beneficiary_count", fmt.Sprintf("%d", len(msg.Beneficiaries))),
		),
	})

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", planID)),
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

// HandleMsgAddAsset handles the add asset message
func (k Keeper) HandleMsgAddAsset(ctx sdk.Context, msg types.MsgAddAsset) (*sdk.Result, error) {
	plan, found := k.GetPlan(ctx, msg.PlanID)
	if !found {
		return nil, types.ErrPlanNotFound
	}
	if !plan.Creator.Equals(msg.Owner) {
		return nil, types.ErrUnauthorized
	}

	asset := types.Asset{
		PlanID:        msg.PlanID,
		Owner:         msg.Owner,
		AssetType:     msg.AssetType,
		Denom:         msg.Denom,
		Amount:        msg.Amount,
		IPFSCid:       msg.IPFSCid,
		Metadata:      msg.Metadata,
		EncryptedData: msg.EncryptedData,
	}

	assetID, err := k.AddAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
		sdk.NewEvent(
			"add_asset",
			sdk.NewAttribute("asset_id", fmt.Sprintf("%d", assetID)),
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", msg.PlanID)),
			sdk.NewAttribute("asset_type", msg.AssetType),
			sdk.NewAttribute("denom", msg.Denom),
			sdk.NewAttribute("amount", msg.Amount.String()),
		),
	})

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", assetID)),
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

// HandleMsgHeartbeat handles the heartbeat message
func (k Keeper) HandleMsgHeartbeat(ctx sdk.Context, msg types.MsgHeartbeat) (*sdk.Result, error) {
	plan, found := k.GetPlan(ctx, msg.PlanID)
	if !found {
		return nil, types.ErrPlanNotFound
	}
	if !plan.Creator.Equals(msg.Creator) {
		return nil, types.ErrUnauthorized
	}

	if err := k.UpdateHeartbeat(ctx, msg.PlanID); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
		sdk.NewEvent(
			"heartbeat",
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", msg.PlanID)),
			sdk.NewAttribute("creator", msg.Creator.String()),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

// HandleMsgClaimInheritance handles the claim inheritance message
func (k Keeper) HandleMsgClaimInheritance(ctx sdk.Context, msg types.MsgClaimInheritance) (*sdk.Result, error) {
	if err := k.ClaimInheritance(ctx, msg.PlanID, msg.Beneficiary); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Beneficiary.String()),
		),
		sdk.NewEvent(
			"claim_inheritance",
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", msg.PlanID)),
			sdk.NewAttribute("beneficiary", msg.Beneficiary.String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

// HandleMsgUpdateBeneficiaries handles the update beneficiaries message
func (k Keeper) HandleMsgUpdateBeneficiaries(ctx sdk.Context, msg types.MsgUpdateBeneficiaries) (*sdk.Result, error) {
	plan, found := k.GetPlan(ctx, msg.PlanID)
	if !found {
		return nil, types.ErrPlanNotFound
	}
	if !plan.Creator.Equals(msg.Creator) {
		return nil, types.ErrUnauthorized
	}
	if plan.Status != types.PlanStatusActive {
		return nil, types.ErrPlanNotActive
	}

	beneficiaries := make([]types.Beneficiary, len(msg.Beneficiaries))
	for i, b := range msg.Beneficiaries {
		beneficiaries[i] = types.Beneficiary{
			Address: b.Address,
			Share:   b.Share,
		}
	}
	plan.Beneficiaries = beneficiaries
	k.SetPlan(ctx, plan)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
		sdk.NewEvent(
			"update_beneficiaries",
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", msg.PlanID)),
			sdk.NewAttribute("beneficiary_count", fmt.Sprintf("%d", len(msg.Beneficiaries))),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
