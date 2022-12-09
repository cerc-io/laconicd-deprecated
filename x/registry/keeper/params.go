package keeper

import (
	"github.com/cerc-io/laconicd/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams - Get all parameters as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}
