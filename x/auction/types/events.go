package types

const (
	EventTypeCreateAuction = "create-auction"
	EventTypeCommitBid     = "commit-bid"
	EventTypeRevealBid     = "reveal-bid"

	AttributeKeyCommitsDuration = "commits-duration"
	AttributeKeyRevealsDuration = "reveals-duration"
	AttributeKeyCommitFee       = "commit-fee"
	AttributeKeyRevealFee       = "reveal-fee"
	AttributeKeyMinimumBid      = "minimum-bid"
	AttributeKeySigner          = "signer"
	AttributeKeyAuctionID       = "auction-id"
	AttributeKeyCommitHash      = "commit-hash"
	AttributeKeyReveal          = "reveal"

	AttributeValueCategory = ModuleName
)
