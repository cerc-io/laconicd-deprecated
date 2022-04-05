package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BondUsageKeeper keep track of bond usage in other modules.
// Used to, for example, prevent deletion of a bond that's in use.
type BondUsageKeeper interface {
	ModuleName() string
	UsesBond(ctx sdk.Context, bondId string) bool
}
