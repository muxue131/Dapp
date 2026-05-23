package module

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/inherit-dapp/chain/x/legacy/keeper"
	"github.com/inherit-dapp/chain/x/legacy/types"
)

// NewHandler returns an sdk.Handler for the legacy module
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateLegacyPlan:
			return msgServer.CreateLegacyPlan(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgAddAsset:
			return msgServer.AddAsset(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgHeartbeat:
			return msgServer.Heartbeat(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgClaimInheritance:
			return msgServer.ClaimInheritance(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgUpdateBeneficiaries:
			return msgServer.UpdateBeneficiaries(sdk.WrapSDKContext(ctx), msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// MsgServerImpl implements the legacy module message server
type MsgServerImpl struct {
	keeper.Keeper
}

func NewMsgServerImpl(keeper keeper.Keeper) *MsgServerImpl {
	return &MsgServerImpl{Keeper: keeper}
}

func (m MsgServerImpl) CreateLegacyPlan(goCtx context.Context, msg *types.MsgCreateLegacyPlan) (*sdk.Result, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.HandleMsgCreateLegacyPlan(ctx, *msg)
}

func (m MsgServerImpl) AddAsset(goCtx context.Context, msg *types.MsgAddAsset) (*sdk.Result, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.HandleMsgAddAsset(ctx, *msg)
}

func (m MsgServerImpl) Heartbeat(goCtx context.Context, msg *types.MsgHeartbeat) (*sdk.Result, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.HandleMsgHeartbeat(ctx, *msg)
}

func (m MsgServerImpl) ClaimInheritance(goCtx context.Context, msg *types.MsgClaimInheritance) (*sdk.Result, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.HandleMsgClaimInheritance(ctx, *msg)
}

func (m MsgServerImpl) UpdateBeneficiaries(goCtx context.Context, msg *types.MsgUpdateBeneficiaries) (*sdk.Result, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.HandleMsgUpdateBeneficiaries(ctx, *msg)
}
