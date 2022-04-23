package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuctionUsageKeeper keep track of auction usage in other modules.
// Used to, for example, prevent deletion of a auction that's in use.
type AuctionUsageKeeper interface {
	ModuleName() string
	UsesAuction(ctx sdk.Context, auctionID string) bool

	OnAuction(ctx sdk.Context, auctionID string)
	OnAuctionBid(ctx sdk.Context, auctionID string, bidderAddress string)
	OnAuctionWinnerSelected(ctx sdk.Context, auctionID string)
}
