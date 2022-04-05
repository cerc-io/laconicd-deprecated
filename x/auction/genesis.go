package auction

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tharsis/ethermint/x/auction/keeper"
	"github.com/tharsis/ethermint/x/auction/types"
)

// func NewGenesisState(params types.Params, auctions []types.Auction) types.GenesisState {
// 	return types.GenesisState{Params: params, Auctions: &types.Auctions{Auctions: auctions}}
// }

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	for _, auction := range data.Auctions {
		keeper.SaveAuction(ctx, auction)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	auctions := keeper.ListAuctions(ctx)

	var genesisAuctions []*types.Auction
	for _, auction := range auctions {
		genesisAuctions = append(genesisAuctions, &auction)
	}
	return types.GenesisState{Params: params, Auctions: genesisAuctions}
}

func ValidateGenesis(data types.GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
