package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Auction status values.
const (
	// Auction is in commit phase.
	AuctionStatusCommitPhase = "commit"

	// Auction is in reveal phase.
	AuctionStatusRevealPhase = "reveal"

	// Auction has ended (no reveals allowed).
	AuctionStatusExpired = "expired"

	// Auction has completed (winner selected).
	AuctionStatusCompleted = "completed"
)

// Bid status values.
const (
	BidStatusCommitted = "commit"
	BidStatusRevealed  = "reveal"
)

// AuctionID simplifies generation of auction IDs.
type AuctionID struct {
	Address  sdk.Address
	AccNum   uint64
	Sequence uint64
}

// Generate creates the auction ID.
func (auctionID AuctionID) Generate() string {
	hasher := sha256.New()
	str := fmt.Sprintf("%s:%d:%d", auctionID.Address.String(), auctionID.AccNum, auctionID.Sequence)
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (auction Auction) GetCreateTime() string {
	return string(sdk.FormatTimeBytes(auction.CreateTime))
}

func (auction Auction) GetCommitsEndTime() string {
	return string(sdk.FormatTimeBytes(auction.CommitsEndTime))
}

func (auction Auction) GetRevealsEndTime() string {
	return string(sdk.FormatTimeBytes(auction.RevealsEndTime))
}

func (bid Bid) GetCommitTime() string {
	return string(sdk.FormatTimeBytes(bid.CommitTime))
}

func (bid Bid) GetRevealTime() string {
	return string(sdk.FormatTimeBytes(bid.RevealTime))
}
