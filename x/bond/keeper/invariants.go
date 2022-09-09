package keeper

import (
	"fmt"

	"github.com/cerc-io/laconicd/x/bond/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers all bond invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
}

// ModuleAccountInvariant checks that the 'bond' module account balance is non-negative.
func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
		balances := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
		for _, balance := range balances {
			if balance.IsNegative() {
				return sdk.FormatInvariant(
						types.ModuleName,
						"module-account",
						fmt.Sprintf("Module account '%s' has negative balance.", types.ModuleName)),
					true
			}
		}
		return "", false
	}
}

// AllInvariants runs all invariants of the bonds module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return ModuleAccountInvariant(k)(ctx)
	}
}
