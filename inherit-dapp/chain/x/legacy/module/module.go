package module

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/inherit-dapp/chain/x/legacy/genesis"
	"github.com/inherit-dapp/chain/x/legacy/keeper"
	"github.com/inherit-dapp/chain/x/legacy/types"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module for the legacy module
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := types.DefaultGenesisState()
	return cdc.MustMarshalJSON(&genesisState)
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ module.TxEncodingConfig, bz json.RawMessage) error {
	var genesisState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesisState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return types.ValidateGenesis(genesisState)
}

// AppModule implements the sdk.AppModule interface
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

func (am AppModule) Name() string { return types.ModuleName }

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	genesis.InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesisState := genesis.ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(&genesisState)
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

func (am AppModule) QuerierRoute() string { return types.QuerierRoute }

func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {}

func (am AppModule) BeginBlock(ctx sdk.Context, _ sdk.RequestBeginBlock) {
	// Check for heartbeat expirations
	plans := am.keeper.GetAllPlans(ctx)
	for _, plan := range plans {
		if plan.Status == types.PlanStatusActive {
			expired, _ := am.keeper.CheckHeartbeatExpiry(ctx, plan.PlanID)
			if expired {
				ctx.EventManager().EmitEvent(sdk.NewEvent(
					"heartbeat_expired",
					sdk.NewAttribute("plan_id", fmt.Sprintf("%d", plan.PlanID)),
					sdk.NewAttribute("creator", plan.Creator.String()),
				))
			}
		}
	}
}

func (am AppModule) EndBlock(ctx sdk.Context, _ sdk.RequestEndBlock) []sdk.ValidatorUpdate {
	return []sdk.ValidatorUpdate{}
}

// ConsensusVersion returns the module consensus version
func (am AppModule) ConsensusVersion() uint64 { return 1 }
