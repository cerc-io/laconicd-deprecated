package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cerc-io/laconicd/x/auction/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	params "github.com/cosmos/cosmos-sdk/x/params/types"

	wnsUtils "github.com/cerc-io/laconicd/utils"
)

// CompletedAuctionDeleteTimeout => Completed auctions are deleted after this timeout (after reveals end time).
const CompletedAuctionDeleteTimeout = time.Hour * 24

// PrefixIDToAuctionIndex is the prefix for Id -> Auction index in the KVStore.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var PrefixIDToAuctionIndex = []byte{0x00}

// prefixOwnerToAuctionsIndex is the prefix for the Owner -> [Auction] index in the KVStore.
var prefixOwnerToAuctionsIndex = []byte{0x01}

// PrefixAuctionBidsIndex is the prefix for the (auction, bidder) -> Bid index in the KVStore.
var PrefixAuctionBidsIndex = []byte{0x02}

// PrefixBidderToAuctionsIndex is the prefix for the Bidder -> [Auction] index in the KVStore.
var PrefixBidderToAuctionsIndex = []byte{0x03}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	// Track auction usage in other cosmos-sdk modules (more like a usage tracker).
	usageKeepers []types.AuctionUsageKeeper

	storeKey storetypes.StoreKey // Unexposed key to access store from sdk.Context

	cdc codec.BinaryCodec // The wire codec for binary encoding/decoding.

	paramSubspace params.Subspace
}

// AuctionClientKeeper is the subset of functionality exposed to other modules.
type AuctionClientKeeper interface {
	HasAuction(ctx sdk.Context, id string) bool
	GetAuction(ctx sdk.Context, id string) types.Auction
	MatchAuctions(ctx sdk.Context, matchFn func(*types.Auction) bool) []*types.Auction
}

// NewKeeper creates new instances of the auction Keeper
func NewKeeper(accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, ps params.Subspace) Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		paramSubspace: ps,
	}
}

func (k *Keeper) SetUsageKeepers(usageKeepers []types.AuctionUsageKeeper) {
	k.usageKeepers = usageKeepers
}

func (k Keeper) GetUsageKeepers() []types.AuctionUsageKeeper {
	return k.usageKeepers
}

// Generates Auction Id -> Auction index key.
func GetAuctionIndexKey(id string) []byte {
	return append(PrefixIDToAuctionIndex, []byte(id)...)
}

// Generates Owner -> Auctions index key.
func GetOwnerToAuctionsIndexKey(owner string, auctionID string) []byte {
	return append(append(prefixOwnerToAuctionsIndex, []byte(owner)...), []byte(auctionID)...)
}

func GetBidderToAuctionsIndexKey(bidder string, auctionID string) []byte {
	return append(append(PrefixBidderToAuctionsIndex, []byte(bidder)...), []byte(auctionID)...)
}

func GetBidIndexKey(auctionID string, bidder string) []byte {
	return append(GetAuctionBidsIndexPrefix(auctionID), []byte(bidder)...)
}

func GetAuctionBidsIndexPrefix(auctionID string) []byte {
	return append(append(PrefixAuctionBidsIndex, []byte(auctionID)...))
}

// SaveAuction - saves a auction to the store.
func (k Keeper) SaveAuction(ctx sdk.Context, auction *types.Auction) {
	store := ctx.KVStore(k.storeKey)

	// Auction Id -> Auction index.
	store.Set(GetAuctionIndexKey(auction.Id), k.cdc.MustMarshal(auction))

	// Owner -> [Auction] index.
	store.Set(GetOwnerToAuctionsIndexKey(auction.OwnerAddress, auction.Id), []byte{})

	// Notify interested parties.
	for _, keeper := range k.usageKeepers {
		keeper.OnAuction(ctx, auction.Id)
	}
}

func (k Keeper) SaveBid(ctx sdk.Context, bid *types.Bid) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetBidIndexKey(bid.AuctionId, bid.BidderAddress), k.cdc.MustMarshal(bid))
	store.Set(GetBidderToAuctionsIndexKey(bid.BidderAddress, bid.AuctionId), []byte{})

	// Notify interested parties.
	for _, keeper := range k.usageKeepers {
		keeper.OnAuctionBid(ctx, bid.AuctionId, bid.BidderAddress)
	}
}

func (k Keeper) DeleteBid(ctx sdk.Context, bid types.Bid) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetBidIndexKey(bid.AuctionId, bid.BidderAddress))
	store.Delete(GetOwnerToAuctionsIndexKey(bid.BidderAddress, bid.AuctionId))
}

// HasAuction - checks if a auction by the given Id exists.
func (k Keeper) HasAuction(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetAuctionIndexKey(id))
}

func (k Keeper) HasBid(ctx sdk.Context, id string, bidder string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetBidIndexKey(id, bidder))
}

// DeleteAuction - deletes the auction.
func (k Keeper) DeleteAuction(ctx sdk.Context, auction types.Auction) {
	// Delete all bids first.
	bids := k.GetBids(ctx, auction.Id)
	for _, bid := range bids {
		k.DeleteBid(ctx, *bid)
	}

	// Delete the auction itself.
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetAuctionIndexKey(auction.Id))
	store.Delete(GetOwnerToAuctionsIndexKey(auction.OwnerAddress, auction.Id))
}

// GetAuction - gets a record from the store.
func (k Keeper) GetAuction(ctx sdk.Context, id string) *types.Auction {
	store := ctx.KVStore(k.storeKey)
	auctionKey := GetAuctionIndexKey(id)
	if !store.Has(auctionKey) {
		return nil
	}

	bz := store.Get(auctionKey)
	var obj types.Auction
	k.cdc.MustUnmarshal(bz, &obj)

	return &obj
}

// GetBids gets the auction bids.
func (k Keeper) GetBids(ctx sdk.Context, id string) []*types.Bid {
	store := ctx.KVStore(k.storeKey)
	bids := []*types.Bid{}
	itr := sdk.KVStorePrefixIterator(store, GetAuctionBidsIndexPrefix(id))
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Bid
			k.cdc.MustUnmarshal(bz, &obj)
			bids = append(bids, &obj)
		}
	}

	return bids
}

func (k Keeper) GetBid(ctx sdk.Context, id string, bidder string) types.Bid {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetBidIndexKey(id, bidder))
	var obj types.Bid
	k.cdc.MustUnmarshal(bz, &obj)

	return obj
}

// ListAuctions - get all auctions.
func (k Keeper) ListAuctions(ctx sdk.Context) []types.Auction {
	var auctions []types.Auction

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixIDToAuctionIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshal(bz, &obj)
			auctions = append(auctions, obj)
		}
	}

	return auctions
}

// QueryAuctionsByOwner - query auctions by owner.
func (k Keeper) QueryAuctionsByOwner(ctx sdk.Context, ownerAddress string) []types.Auction {
	auctions := []types.Auction{}

	ownerPrefix := append(prefixOwnerToAuctionsIndex, []byte(ownerAddress)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, ownerPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		auctionID := itr.Key()[len(ownerPrefix):]
		bz := store.Get(append(PrefixIDToAuctionIndex, auctionID...))
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshal(bz, &obj)
			auctions = append(auctions, obj)
		}
	}

	return auctions
}

// QueryAuctionsByBidder - query auctions by bidder
func (k Keeper) QueryAuctionsByBidder(ctx sdk.Context, bidderAddress string) []types.Auction {
	auctions := []types.Auction{}

	bidderPrefix := append(PrefixBidderToAuctionsIndex, []byte(bidderAddress)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, []byte(bidderPrefix))
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		auctionID := itr.Key()[len(bidderPrefix):]
		bz := store.Get(append(PrefixIDToAuctionIndex, auctionID...))
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshal(bz, &obj)
			auctions = append(auctions, obj)
		}
	}

	return auctions
}

// MatchAuctions - get all matching auctions.
func (k Keeper) MatchAuctions(ctx sdk.Context, matchFn func(*types.Auction) bool) []*types.Auction {
	var auctions []*types.Auction

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixIDToAuctionIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Auction
			k.cdc.MustUnmarshal(bz, &obj)
			if matchFn(&obj) {
				auctions = append(auctions, &obj)
			}
		}
	}

	return auctions
}

// CreateAuction creates a new auction.
func (k Keeper) CreateAuction(ctx sdk.Context, msg types.MsgCreateAuction) (*types.Auction, error) {
	// Might be called from another module directly, always validate.
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	// Generate auction Id.
	account := k.accountKeeper.GetAccount(ctx, signerAddress)
	if account == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Account not found.")
	}

	auctionID := types.AuctionID{
		Address:  signerAddress,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	// Compute timestamps.
	now := ctx.BlockTime()
	commitsEndTime := now.Add(time.Duration(msg.CommitsDuration))
	revealsEndTime := now.Add(time.Duration(msg.CommitsDuration + msg.RevealsDuration))

	auction := types.Auction{
		Id:             auctionID,
		Status:         types.AuctionStatusCommitPhase,
		OwnerAddress:   signerAddress.String(),
		CreateTime:     now,
		CommitsEndTime: commitsEndTime,
		RevealsEndTime: revealsEndTime,
		CommitFee:      msg.CommitFee,
		RevealFee:      msg.RevealFee,
		MinimumBid:     msg.MinimumBid,
	}

	// Save auction in store.
	k.SaveAuction(ctx, &auction)

	return &auction, nil
}

func (k Keeper) CommitBid(ctx sdk.Context, msg types.MsgCommitBid) (*types.Bid, error) {
	if !k.HasAuction(ctx, msg.AuctionId) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Auction not found.")
	}

	auction := k.GetAuction(ctx, msg.AuctionId)
	if auction.Status != types.AuctionStatusCommitPhase {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Auction is not in commit phase.")
	}

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	// Take auction fees from account.
	totalFee := auction.CommitFee.Add(auction.RevealFee)
	sdkErr := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signerAddress, types.ModuleName, sdk.NewCoins(totalFee))
	if sdkErr != nil {
		return nil, sdkErr
	}

	// Check if an old bid already exists, if so, return old bids auction fee (update bid scenario).
	bidder := signerAddress.String()
	if k.HasBid(ctx, msg.AuctionId, bidder) {
		oldBid := k.GetBid(ctx, msg.AuctionId, bidder)
		oldTotalFee := oldBid.CommitFee.Add(oldBid.RevealFee)
		sdkErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, signerAddress, sdk.NewCoins(oldTotalFee))
		if sdkErr != nil {
			return nil, sdkErr
		}
	}

	// Save new bid.
	bid := types.Bid{
		AuctionId:     msg.AuctionId,
		BidderAddress: bidder,
		Status:        types.BidStatusCommitted,
		CommitHash:    msg.CommitHash,
		CommitTime:    ctx.BlockTime(),
		CommitFee:     auction.CommitFee,
		RevealFee:     auction.RevealFee,
	}

	k.SaveBid(ctx, &bid)

	return &bid, nil
}

// RevealBid reeals a bid comitted earlier.
func (k Keeper) RevealBid(ctx sdk.Context, msg types.MsgRevealBid) (*types.Auction, error) {
	if !k.HasAuction(ctx, msg.AuctionId) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Auction not found.")
	}

	auction := k.GetAuction(ctx, msg.AuctionId)
	if auction.Status != types.AuctionStatusRevealPhase {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Auction is not in reveal phase.")
	}

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	if !k.HasBid(ctx, msg.AuctionId, signerAddress.String()) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bid not found.")
	}

	bid := k.GetBid(ctx, auction.Id, signerAddress.String())
	if bid.Status != types.BidStatusCommitted {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bid not in committed state.")
	}

	revealBytes, err := hex.DecodeString(msg.Reveal)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal string.")
	}

	cid, err := wnsUtils.CIDFromJSONBytes(revealBytes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal JSON.")
	}

	if bid.CommitHash != cid {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit hash mismatch.")
	}

	var reveal map[string]interface{}
	err = json.Unmarshal(revealBytes, &reveal)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Reveal JSON unmarshal error.")
	}

	chainID, err := wnsUtils.GetAttributeAsString(reveal, "chainId")
	if err != nil || chainID != ctx.ChainID() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal chainID.")
	}

	auctionID, err := wnsUtils.GetAttributeAsString(reveal, "auctionId")
	if err != nil || auctionID != msg.AuctionId {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal auction Id.")
	}

	bidderAddress, err := wnsUtils.GetAttributeAsString(reveal, "bidderAddress")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal bid address.")
	}

	if bidderAddress != signerAddress.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Reveal bid address mismatch.")
	}

	bidAmountStr, err := wnsUtils.GetAttributeAsString(reveal, "bidAmount")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal bid amount.")
	}

	bidAmount, err := sdk.ParseCoinNormalized(bidAmountStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal bid amount.")
	}

	if bidAmount.IsLT(auction.MinimumBid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bid is lower than minimum bid.")
	}

	// Lock bid amount.
	sdkErr := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signerAddress, types.ModuleName, sdk.NewCoins(bidAmount))
	if sdkErr != nil {
		return nil, sdkErr
	}

	// Update bid.
	bid.BidAmount = bidAmount
	bid.RevealTime = ctx.BlockTime()
	bid.Status = types.BidStatusRevealed
	k.SaveBid(ctx, &bid)

	return auction, nil
}

// GetAuctionModuleBalances gets the auction module account(s) balances.
func (k Keeper) GetAuctionModuleBalances(ctx sdk.Context) sdk.Coins {
	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
	return balances
}

func (k Keeper) EndBlockerProcessAuctions(ctx sdk.Context) {
	// Transition auction state (commit, reveal, expired, completed).
	k.processAuctionPhases(ctx)

	// Delete stale auctions.
	k.deleteCompletedAuctions(ctx)
}

func (k Keeper) processAuctionPhases(ctx sdk.Context) {
	auctions := k.MatchAuctions(ctx, func(_ *types.Auction) bool {
		return true
	})

	for _, auction := range auctions {
		// Commit -> Reveal state.
		if auction.Status == types.AuctionStatusCommitPhase && ctx.BlockTime().After(auction.CommitsEndTime) {
			auction.Status = types.AuctionStatusRevealPhase
			k.SaveAuction(ctx, auction)
			ctx.Logger().Info(fmt.Sprintf("Moved auction %s to reveal phase.", auction.Id))
		}

		// Reveal -> Expired state.
		if auction.Status == types.AuctionStatusRevealPhase && ctx.BlockTime().After(auction.RevealsEndTime) {
			auction.Status = types.AuctionStatusExpired
			k.SaveAuction(ctx, auction)
			ctx.Logger().Info(fmt.Sprintf("Moved auction %s to expired state.", auction.Id))
		}

		// If auction has expired, pick a winner from revealed bids.
		if auction.Status == types.AuctionStatusExpired {
			k.pickAuctionWinner(ctx, auction)
		}
	}
}

// Delete completed stale auctions.
func (k Keeper) deleteCompletedAuctions(ctx sdk.Context) {
	auctions := k.MatchAuctions(ctx, func(auction *types.Auction) bool {
		deleteTime := auction.RevealsEndTime.Add(CompletedAuctionDeleteTimeout)
		return auction.Status == types.AuctionStatusCompleted && ctx.BlockTime().After(deleteTime)
	})

	for _, auction := range auctions {
		ctx.Logger().Info(fmt.Sprintf("Deleting completed auction %s after timeout.", auction.Id))
		k.DeleteAuction(ctx, *auction)
	}
}

func (k Keeper) pickAuctionWinner(ctx sdk.Context, auction *types.Auction) {
	ctx.Logger().Info(fmt.Sprintf("Picking auction %s winner.", auction.Id))

	var highestBid *types.Bid
	var secondHighestBid *types.Bid

	bids := k.GetBids(ctx, auction.Id)
	for _, bid := range bids {
		ctx.Logger().Info(fmt.Sprintf("Processing bid %s %s", bid.BidderAddress, bid.BidAmount.String()))

		// Only consider revealed bids.
		if bid.Status != types.BidStatusRevealed {
			ctx.Logger().Info(fmt.Sprintf("Ignoring unrevealed bid %s %s", bid.BidderAddress, bid.BidAmount.String()))
			continue
		}

		// Init highest bid.
		if highestBid == nil {
			highestBid = bid
			ctx.Logger().Info(fmt.Sprintf("Initializing 1st bid %s %s", bid.BidderAddress, bid.BidAmount.String()))
			continue
		}

		if highestBid.BidAmount.IsLT(bid.BidAmount) {
			ctx.Logger().Info(fmt.Sprintf("New highest bid %s %s", bid.BidderAddress, bid.BidAmount.String()))

			secondHighestBid = highestBid
			highestBid = bid

			ctx.Logger().Info(fmt.Sprintf("Updated 1st bid %s %s", highestBid.BidderAddress, highestBid.BidAmount.String()))
			ctx.Logger().Info(fmt.Sprintf("Updated 2nd bid %s %s", secondHighestBid.BidderAddress, secondHighestBid.BidAmount.String()))

		} else if secondHighestBid == nil || secondHighestBid.BidAmount.IsLT(bid.BidAmount) {
			ctx.Logger().Info(fmt.Sprintf("New 2nd highest bid %s %s", bid.BidderAddress, bid.BidAmount.String()))

			secondHighestBid = bid
			ctx.Logger().Info(fmt.Sprintf("Updated 2nd bid %s %s", secondHighestBid.BidderAddress, secondHighestBid.BidAmount.String()))
		} else {
			ctx.Logger().Info(fmt.Sprintf("Ignoring bid as it doesn't affect 1st/2nd price %s %s", bid.BidderAddress, bid.BidAmount.String()))
		}
	}

	// Highest bid is the winner, but pays second highest bid price.
	auction.Status = types.AuctionStatusCompleted

	if highestBid != nil {
		auction.WinnerAddress = highestBid.BidderAddress
		auction.WinningBid = highestBid.BidAmount

		// Winner pays 2nd price, if a 2nd price exists.
		auction.WinningPrice = highestBid.BidAmount
		if secondHighestBid != nil {
			auction.WinningPrice = secondHighestBid.BidAmount
		}

		ctx.Logger().Info(fmt.Sprintf("Auction %s winner %s.", auction.Id, auction.WinnerAddress))
		ctx.Logger().Info(fmt.Sprintf("Auction %s winner bid %s.", auction.Id, auction.WinningBid.String()))
		ctx.Logger().Info(fmt.Sprintf("Auction %s winner price %s.", auction.Id, auction.WinningPrice.String()))

	} else {
		ctx.Logger().Info(fmt.Sprintf("Auction %s has no valid revealed bids (no winner).", auction.Id))
	}

	k.SaveAuction(ctx, auction)

	for _, bid := range bids {
		bidderAddress, err := sdk.AccAddressFromBech32(bid.BidderAddress)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Invalid bidderAddress address. %v", err))
			panic("Invalid bidder address.")
		}

		if bid.Status == types.BidStatusRevealed {
			// Send reveal fee back to bidders that've revealed the bid.
			sdkErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddress, sdk.NewCoins(bid.RevealFee))
			if sdkErr != nil {
				ctx.Logger().Error(fmt.Sprintf("Auction error returning reveal fee: %v", sdkErr))
				panic(sdkErr)
			}
		}

		// Send back locked bid amount to all bidders.
		sdkErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddress, sdk.NewCoins(bid.BidAmount))
		if sdkErr != nil {
			ctx.Logger().Error(fmt.Sprintf("Auction error returning bid amount: %v", sdkErr))
			panic(sdkErr)
		}
	}

	// Process winner account (if nobody bids, there won't be a winner).
	if auction.WinnerAddress != "" {
		winnerAddress, err := sdk.AccAddressFromBech32(auction.WinnerAddress)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Invalid winner address. %v", err))
			panic("Invalid winner address.")
		}

		// Take 2nd price from winner.
		sdkErr := k.bankKeeper.SendCoinsFromAccountToModule(ctx, winnerAddress, types.ModuleName, sdk.NewCoins(auction.WinningPrice))
		if sdkErr != nil {
			ctx.Logger().Error(fmt.Sprintf("Auction error taking funds from winner: %v", sdkErr))
			panic(sdkErr)
		}

		// Burn anything over the min. bid amount.
		amountToBurn := auction.WinningPrice.Sub(auction.MinimumBid)
		if amountToBurn.IsNegative() {
			ctx.Logger().Error(fmt.Sprintf("Auction coins to burn cannot be negative."))
			panic("Auction coins to burn cannot be negative.")
		}

		// Use auction burn module account instead of actually burning coins to better keep track of supply.
		sdkErr = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.AuctionBurnModuleAccountName, sdk.NewCoins(amountToBurn))
		if sdkErr != nil {
			ctx.Logger().Error(fmt.Sprintf("Auction error burning coins: %v", sdkErr))
			panic(sdkErr)
		}
	}

	// Notify other modules (hook).
	ctx.Logger().Info(fmt.Sprintf("Auction %s notifying %d modules.", auction.Id, len(k.usageKeepers)))
	for _, keeper := range k.usageKeepers {
		ctx.Logger().Info(fmt.Sprintf("Auction %s notifying module %s.", auction.Id, keeper.ModuleName()))
		keeper.OnAuctionWinnerSelected(ctx, auction.Id)
	}
}
