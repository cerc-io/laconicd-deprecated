package auction

import (
	"github.com/cerc-io/laconicd/x/auction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker is called every block, returns updated validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	k.EndBlockerProcessAuctions(ctx)
	return []abci.ValidatorUpdate{}
}
