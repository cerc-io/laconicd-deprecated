package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cerc-io/laconicd/x/auction/types"
)

// RegisterInvariants registers all auction module invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
}

// ModuleAccountInvariant checks that the 'auction' module account balance is non-negative.
func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
		if k.bankKeeper.GetAllBalances(ctx, moduleAddress).IsAnyNegative() {
			return sdk.FormatInvariant(
					types.ModuleName,
					"module-account",
					fmt.Sprintf("Module account '%s' has negative balance.", types.ModuleName)),
				true
		}

		return "", false
	}
}

// AllInvariants runs all invariants of the auctions module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return ModuleAccountInvariant(k)(ctx)
	}
}
