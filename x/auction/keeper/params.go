package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cerc-io/laconicd/x/auction/types"
)

// GetParams - Get all parameteras as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}
