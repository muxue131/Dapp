package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(ModuleCdc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateLegacyPlan{}, "legacy/MsgCreateLegacyPlan", nil)
	cdc.RegisterConcrete(&MsgAddAsset{}, "legacy/MsgAddAsset", nil)
	cdc.RegisterConcrete(&MsgHeartbeat{}, "legacy/MsgHeartbeat", nil)
	cdc.RegisterConcrete(&MsgClaimInheritance{}, "legacy/MsgClaimInheritance", nil)
	cdc.RegisterConcrete(&MsgUpdateBeneficiaries{}, "legacy/MsgUpdateBeneficiaries", nil)
}

// MarshalJSON implements custom JSON marshaling for LegacyPlan
func (p LegacyPlan) MarshalJSON() ([]byte, error) {
	type Alias LegacyPlan
	return json.Marshal(&struct {
		Alias
		Creator string `json:"creator"`
	}{
		Alias:   Alias(p),
		Creator: p.Creator.String(),
	})
}
