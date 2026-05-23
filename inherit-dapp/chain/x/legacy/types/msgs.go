package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDuplicatePlanID    = sdkerrors.Register(ModuleName, 2, "duplicate plan ID")
	ErrInvalidCreator     = sdkerrors.Register(ModuleName, 3, "invalid creator address")
	ErrNoBeneficiaries    = sdkerrors.Register(ModuleName, 4, "no beneficiaries specified")
	ErrInvalidBeneficiary = sdkerrors.Register(ModuleName, 5, "invalid beneficiary address")
	ErrInvalidShareTotal  = sdkerrors.Register(ModuleName, 6, "beneficiary shares must sum to 1")
	ErrPlanNotFound       = sdkerrors.Register(ModuleName, 7, "legacy plan not found")
	ErrPlanNotActive      = sdkerrors.Register(ModuleName, 8, "legacy plan is not active")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 9, "unauthorized")
	ErrAssetNotFound      = sdkerrors.Register(ModuleName, 10, "asset not found")
	ErrHeartbeatNotDue    = sdkerrors.Register(ModuleName, 11, "heartbeat not yet due for check")
	ErrInheritanceNotReady = sdkerrors.Register(ModuleName, 12, "inheritance not ready for claiming")
	ErrAlreadyClaimed     = sdkerrors.Register(ModuleName, 13, "inheritance already claimed")
	ErrInvalidAmount      = sdkerrors.Register(ModuleName, 14, "invalid asset amount")
)

// MsgCreateLegacyPlan defines the message to create a legacy plan
type MsgCreateLegacyPlan struct {
	Creator           sdk.AccAddress `json:"creator"`
	Beneficiaries     []BeneficiaryInput `json:"beneficiaries"`
	HeartbeatInterval int64          `json:"heartbeat_interval"`
}

type BeneficiaryInput struct {
	Address sdk.AccAddress `json:"address"`
	Share   sdk.Dec        `json:"share"`
}

func (msg MsgCreateLegacyPlan) Route() string { return RouterKey }
func (msg MsgCreateLegacyPlan) Type() string  { return TypeMsgCreateLegacyPlan }

func (msg MsgCreateLegacyPlan) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator cannot be empty")
	}
	if len(msg.Beneficiaries) == 0 {
		return sdkerrors.Wrap(ErrNoBeneficiaries, "must have at least one beneficiary")
	}
	if msg.HeartbeatInterval <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "heartbeat interval must be positive")
	}

	totalShare := sdk.ZeroDec()
	for _, b := range msg.Beneficiaries {
		if b.Address.Empty() {
			return sdkerrors.Wrap(ErrInvalidBeneficiary, "beneficiary address cannot be empty")
		}
		if b.Share.IsNegative() || b.Share.GT(sdk.OneDec()) {
			return sdkerrors.Wrap(ErrInvalidShareTotal, "share must be between 0 and 1")
		}
		totalShare = totalShare.Add(b.Share)
	}
	if !totalShare.Equal(sdk.OneDec()) {
		return sdkerrors.Wrap(ErrInvalidShareTotal, "shares must sum to 1")
	}
	return nil
}

func (msg MsgCreateLegacyPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateLegacyPlan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgAddAsset defines the message to add an asset to a legacy plan
type MsgAddAsset struct {
	Owner       sdk.AccAddress `json:"owner"`
	PlanID      uint64         `json:"plan_id"`
	AssetType   string         `json:"asset_type"`
	Denom       string         `json:"denom"`
	Amount      sdk.Int        `json:"amount"`
	IPFSCid     string         `json:"ipfs_cid"`
	Metadata    string         `json:"metadata"`
	EncryptedData []byte       `json:"encrypted_data"`
}

func (msg MsgAddAsset) Route() string { return RouterKey }
func (msg MsgAddAsset) Type() string  { return TypeMsgAddAsset }

func (msg MsgAddAsset) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}
	if msg.PlanID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "plan ID cannot be zero")
	}
	if msg.AssetType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset type cannot be empty")
	}
	if msg.Amount.IsNil() || !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAmount, "amount must be positive")
	}
	return nil
}

func (msg MsgAddAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAddAsset) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgHeartbeat defines the heartbeat message
type MsgHeartbeat struct {
	Creator sdk.AccAddress `json:"creator"`
	PlanID  uint64         `json:"plan_id"`
}

func (msg MsgHeartbeat) Route() string { return RouterKey }
func (msg MsgHeartbeat) Type() string  { return TypeMsgHeartbeat }

func (msg MsgHeartbeat) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator cannot be empty")
	}
	if msg.PlanID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "plan ID cannot be zero")
	}
	return nil
}

func (msg MsgHeartbeat) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgHeartbeat) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgClaimInheritance defines the message to claim inheritance
type MsgClaimInheritance struct {
	Beneficiary sdk.AccAddress `json:"beneficiary"`
	PlanID      uint64         `json:"plan_id"`
}

func (msg MsgClaimInheritance) Route() string { return RouterKey }
func (msg MsgClaimInheritance) Type() string  { return TypeMsgClaimInheritance }

func (msg MsgClaimInheritance) ValidateBasic() error {
	if msg.Beneficiary.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "beneficiary cannot be empty")
	}
	if msg.PlanID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "plan ID cannot be zero")
	}
	return nil
}

func (msg MsgClaimInheritance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgClaimInheritance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Beneficiary}
}

// MsgUpdateBeneficiaries defines the message to update beneficiaries
type MsgUpdateBeneficiaries struct {
	Creator       sdk.AccAddress     `json:"creator"`
	PlanID        uint64             `json:"plan_id"`
	Beneficiaries []BeneficiaryInput `json:"beneficiaries"`
}

func (msg MsgUpdateBeneficiaries) Route() string { return RouterKey }
func (msg MsgUpdateBeneficiaries) Type() string  { return TypeMsgUpdateBeneficiaries }

func (msg MsgUpdateBeneficiaries) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator cannot be empty")
	}
	if msg.PlanID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "plan ID cannot be zero")
	}
	if len(msg.Beneficiaries) == 0 {
		return sdkerrors.Wrap(ErrNoBeneficiaries, "must have at least one beneficiary")
	}

	totalShare := sdk.ZeroDec()
	for _, b := range msg.Beneficiaries {
		if b.Address.Empty() {
			return sdkerrors.Wrap(ErrInvalidBeneficiary, "beneficiary address cannot be empty")
		}
		totalShare = totalShare.Add(b.Share)
	}
	if !totalShare.Equal(sdk.OneDec()) {
		return sdkerrors.Wrap(ErrInvalidShareTotal, "shares must sum to 1")
	}
	return nil
}

func (msg MsgUpdateBeneficiaries) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateBeneficiaries) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
