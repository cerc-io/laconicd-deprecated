package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

func GetBlockChangeSetIndexKey(height int64) []byte {
	return append(PrefixBlockChangesetIndex, helpers.Int64ToBytes(height)...)
}

func getOrCreateBlockChangeset(store sdk.KVStore, codec codec.BinaryCodec, height int64) *types.BlockChangeSet {
	bz := store.Get(GetBlockChangeSetIndexKey(height))

	if bz != nil {
		var changeSet types.BlockChangeSet
		err := codec.Unmarshal(bz, &changeSet)
		if err != nil {
			return nil
		}
		return &changeSet
	}

	return &types.BlockChangeSet{
		Height:      height,
		Records:     []string{},
		Names:       []string{},
		Auctions:    []string{},
		AuctionBids: []*types.AuctionBidInfo{},
	}
}

func updateBlockChangeSetForAuction(ctx sdk.Context, k RecordKeeper, id string) {
	changeSet := getOrCreateBlockChangeset(ctx.KVStore(k.storeKey), k.cdc, ctx.BlockHeight())

	found := false
	for _, elem := range changeSet.Auctions {
		if id == elem {
			found = true
			break
		}
	}

	if !found {
		changeSet.Auctions = append(changeSet.Auctions, id)
		saveBlockChangeSet(ctx.KVStore(k.storeKey), k.cdc, changeSet)
	}
}

func saveBlockChangeSet(store sdk.KVStore, codec codec.BinaryCodec, changeset *types.BlockChangeSet) {
	bz := codec.MustMarshal(changeset)
	store.Set(GetBlockChangeSetIndexKey(changeset.Height), bz)
}

func (k Keeper) saveBlockChangeSet(ctx sdk.Context, changeSet *types.BlockChangeSet) {
	saveBlockChangeSet(ctx.KVStore(k.storeKey), k.cdc, changeSet)
}

func (k Keeper) updateBlockChangeSetForRecord(ctx sdk.Context, id string) {
	changeSet := k.getOrCreateBlockChangeSet(ctx, ctx.BlockHeight())
	changeSet.Records = append(changeSet.Records, id)
	k.saveBlockChangeSet(ctx, changeSet)
}

func (k Keeper) getOrCreateBlockChangeSet(ctx sdk.Context, height int64) *types.BlockChangeSet {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetBlockChangeSetIndexKey(height))

	if bz != nil {
		var changeSet types.BlockChangeSet
		err := k.cdc.Unmarshal(bz, &changeSet)
		if err != nil {
			return nil
		}
		return &changeSet
	}

	return &types.BlockChangeSet{
		Height:      height,
		Records:     []string{},
		Names:       []string{},
		Auctions:    []string{},
		AuctionBids: []*types.AuctionBidInfo{},
	}
}

func updateBlockChangeSetForAuctionBid(ctx sdk.Context, k RecordKeeper, id, bidderAddress string) {
	changeSet := getOrCreateBlockChangeset(ctx.KVStore(k.storeKey), k.cdc, ctx.BlockHeight())
	changeSet.AuctionBids = append(changeSet.AuctionBids, &types.AuctionBidInfo{AuctionId: id, BidderAddress: bidderAddress})
	saveBlockChangeSet(ctx.KVStore(k.storeKey), k.cdc, changeSet)
}

func updateBlockChangeSetForNameAuthority(ctx sdk.Context, codec codec.BinaryCodec, store sdk.KVStore, name string) {
	changeSet := getOrCreateBlockChangeset(store, codec, ctx.BlockHeight())
	changeSet.Authorities = append(changeSet.Authorities, name)
	saveBlockChangeSet(store, codec, changeSet)
}
