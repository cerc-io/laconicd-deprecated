package utils

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CtxWithCustomKVGasConfig(ctx *sdk.Context) *sdk.Context {
	updatedCtx := ctx.WithKVGasConfig(storetypes.GasConfig{
		HasCost:          0,
		DeleteCost:       0,
		ReadCostFlat:     0,
		ReadCostPerByte:  0,
		WriteCostFlat:    0,
		WriteCostPerByte: 0,
		IterNextCostFlat: 0,
	})

	return &updatedCtx
}
