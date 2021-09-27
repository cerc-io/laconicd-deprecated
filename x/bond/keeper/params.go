package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/bond/types"
)

// GetMaxBondAmount max bond amount
func (k Keeper) GetMaxBondAmount(ctx sdk.Context) (res sdk.Coin) {
	k.paramSubspace.Get(ctx, types.ParamStoreKeyMaxBondAmount, &res)
	return
}

// GetParams - Get all parameter as as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	getMaxBondAmount := k.GetMaxBondAmount(ctx)
	return types.Params{MaxBondAmount: getMaxBondAmount}
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}
