package auction

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cerc-io/laconicd/x/auction/client/cli"
	"github.com/cerc-io/laconicd/x/auction/keeper"
	"github.com/cerc-io/laconicd/x/auction/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module Basics object
type AppModuleBasic struct {
	cdc codec.Codec
}

func (b AppModuleBasic) Name() string {
	return types.ModuleName
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(jsonCodec codec.JSONCodec) json.RawMessage {
	return jsonCodec.MustMarshalJSON(types.DefaultGenesisState())
}

// Validation check of the Genesis
func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	// Once json successfully marshaled, passes along to genesis.go
	return ValidateGenesis(data)
}

// Register rest routes
func (b AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	// No-op.
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

// Get the root query command of this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// Get the root tx command of this module
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         k,
	}
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	keeper.RegisterInvariants(ir, am.keeper)
}

func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) LegacyQuerierHandler(cdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)

	msgServer := keeper.NewMsgServer(am.keeper)
	types.RegisterMsgServer(cfg.MsgServer(), msgServer)
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(&gs)
}
