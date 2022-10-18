package bond

import (
	"github.com/cerc-io/laconicd/x/bond/keeper"
	"github.com/cerc-io/laconicd/x/bond/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper, data types.GenesisState,
) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)

	for _, bond := range data.Bonds {
		k.SaveBond(ctx, bond)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	bonds := keeper.ListBonds(ctx)

	return types.GenesisState{Params: params, Bonds: bonds}
}

// ValidateGenesis - validating the genesis data
func ValidateGenesis(data types.GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
