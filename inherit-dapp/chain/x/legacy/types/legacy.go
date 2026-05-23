package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "legacy"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
	QuerierRoute = ModuleName

	TypeMsgCreateLegacyPlan = "create_legacy_plan"
	TypeMsgAddAsset         = "add_asset"
	TypeMsgHeartbeat        = "heartbeat"
	TypeMsgClaimInheritance  = "claim_inheritance"
	TypeMsgUpdateBeneficiaries = "update_beneficiaries"

	// Plan status
	PlanStatusActive   = "active"
	PlanStatusTriggered = "triggered"
	PlanStatusClaimed  = "claimed"

	// Asset types
	AssetTypeNative  = "native"
	AssetTypeCW20    = "cw20"
	AssetTypeNFT     = "nft"
	AssetTypeIPFS    = "ipfs"
)

// LegacyPlan represents an inheritance plan
type LegacyPlan struct {
	PlanID         uint64         `json:"plan_id"`
	Creator        sdk.AccAddress `json:"creator"`
	Beneficiaries  []Beneficiary  `json:"beneficiaries"`
	HeartbeatInterval int64       `json:"heartbeat_interval"` // in seconds
	LastHeartbeat  time.Time      `json:"last_heartbeat"`
	TriggerTime    time.Time      `json:"trigger_time"`       // when heartbeat expires
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	EncryptedKey   []byte         `json:"encrypted_key"`      // Shamir-encrypted master key
	IPFSCid        string         `json:"ipfs_cid"`           // CID of encrypted metadata
}

// Beneficiary represents a beneficiary with their share
type Beneficiary struct {
	Address sdk.AccAddress `json:"address"`
	Share   sdk.Dec        `json:"share"`   // percentage share (0-1)
	KeyShare []byte        `json:"key_share"` // Shamir key share
}

// Asset represents an asset in the inheritance plan
type Asset struct {
	AssetID     uint64         `json:"asset_id"`
	PlanID      uint64         `json:"plan_id"`
	Owner       sdk.AccAddress `json:"owner"`
	AssetType   string         `json:"asset_type"`
	Denom       string         `json:"denom"`       // token denom or contract address
	Amount      sdk.Int        `json:"amount"`
	IPFSCid     string         `json:"ipfs_cid"`    // encrypted file CID
	Metadata    string         `json:"metadata"`    // JSON metadata
	EncryptedData []byte       `json:"encrypted_data"`
}

// HeartbeatRecord tracks heartbeat history
type HeartbeatRecord struct {
	PlanID    uint64         `json:"plan_id"`
	Validator sdk.AccAddress `json:"validator"`
	Timestamp time.Time      `json:"timestamp"`
	TxHash    string         `json:"tx_hash"`
}

// GenesisState defines the legacy module's genesis state
type GenesisState struct {
	Plans       []LegacyPlan `json:"plans"`
	Assets      []Asset      `json:"assets"`
	NextPlanID  uint64       `json:"next_plan_id"`
	NextAssetID uint64       `json:"next_asset_id"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Plans:       []LegacyPlan{},
		Assets:      []Asset{},
		NextPlanID:  1,
		NextAssetID: 1,
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	planIDs := make(map[uint64]bool)
	for _, plan := range data.Plans {
		if planIDs[plan.PlanID] {
			return ErrDuplicatePlanID
		}
		planIDs[plan.PlanID] = true
		if plan.Creator.Empty() {
			return ErrInvalidCreator
		}
		if len(plan.Beneficiaries) == 0 {
			return ErrNoBeneficiaries
		}
		totalShare := sdk.ZeroDec()
		for _, b := range plan.Beneficiaries {
			totalShare = totalShare.Add(b.Share)
			if b.Address.Empty() {
				return ErrInvalidBeneficiary
			}
		}
		if !totalShare.Equal(sdk.OneDec()) {
			return ErrInvalidShareTotal
		}
	}
	return nil
}
