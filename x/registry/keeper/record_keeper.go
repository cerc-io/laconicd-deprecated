package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/errors"
	auctionkeeper "github.com/cerc-io/laconicd/x/auction/keeper"
	auctiontypes "github.com/cerc-io/laconicd/x/auction/types"
	bondtypes "github.com/cerc-io/laconicd/x/bond/types"
	"github.com/cerc-io/laconicd/x/registry/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	auctionKeeper auctionkeeper.Keeper
	storeKey      storetypes.StoreKey // Unexposed key to access store from sdk.Context
	cdc           codec.BinaryCodec   // The wire codec for binary encoding/decoding.
}

func (k RecordKeeper) UsesAuction(ctx sdk.Context, auctionID string) bool {
	return k.GetAuctionToAuthorityMapping(ctx, auctionID) != ""
}

func (k RecordKeeper) OnAuction(ctx sdk.Context, auctionID string) {
	updateBlockChangeSetForAuction(ctx, k, auctionID)
}

func (k RecordKeeper) OnAuctionBid(ctx sdk.Context, auctionID string, bidderAddress string) {
	updateBlockChangeSetForAuctionBid(ctx, k, auctionID, bidderAddress)
}

func (k RecordKeeper) OnAuctionWinnerSelected(ctx sdk.Context, auctionID string) {
	// Update authority status based on auction status/winner.
	name := k.GetAuctionToAuthorityMapping(ctx, auctionID)
	if name == "" {
		// We don't know about this auction, ignore.
		ctx.Logger().Info(fmt.Sprintf("Ignoring auction notification, name mapping not found: %s", auctionID))
		return
	}

	store := ctx.KVStore(k.storeKey)
	if !HasNameAuthority(store, name) {
		// We don't know about this authority, ignore.
		ctx.Logger().Info(fmt.Sprintf("Ignoring auction notification, authority not found: %s", auctionID))
		return
	}

	authority := GetNameAuthority(store, k.cdc, name)
	auctionObj := k.auctionKeeper.GetAuction(ctx, auctionID)

	if auctionObj.Status == auctiontypes.AuctionStatusCompleted {
		store := ctx.KVStore(k.storeKey)

		if auctionObj.WinnerAddress != "" {
			// Mark authority owner and change status to active.
			authority.OwnerAddress = auctionObj.WinnerAddress
			authority.Status = types.AuthorityActive

			// Reset bond ID if required, as owner has changed.
			if authority.BondId != "" {
				RemoveBondToAuthorityIndexEntry(store, authority.BondId, name)
				authority.BondId = ""
			}

			// Update height for updated/changed authority (owner).
			// Can be used to check if names are older than the authority itself (stale names).
			authority.Height = uint64(ctx.BlockHeight())

			ctx.Logger().Info(fmt.Sprintf("Winner selected, marking authority as active: %s", name))
		} else {
			// Mark as expired.
			authority.Status = types.AuthorityExpired

			ctx.Logger().Info(fmt.Sprintf("No winner, marking authority as expired: %s", name))
		}

		authority.AuctionId = ""
		SetNameAuthority(ctx, store, k.cdc, name, &authority)

		// Forget about this auction now, we no longer need it.
		removeAuctionToAuthorityMapping(store, auctionID)
	} else {
		ctx.Logger().Info(fmt.Sprintf("Ignoring auction notification, status: %s", auctionObj.Status))
	}
}

// Record keeper implements the bond usage keeper interface.
var (
	_ bondtypes.BondUsageKeeper       = (*RecordKeeper)(nil)
	_ auctiontypes.AuctionUsageKeeper = (*RecordKeeper)(nil)
)

// ModuleName returns the module name.
func (k RecordKeeper) ModuleName() string {
	return types.ModuleName
}

func (k RecordKeeper) GetAuctionToAuthorityMapping(ctx sdk.Context, auctionID string) string {
	store := ctx.KVStore(k.storeKey)

	auctionToAuthorityIndexKey := GetAuctionToAuthorityIndexKey(auctionID)
	if store.Has(auctionToAuthorityIndexKey) {
		bz := store.Get(auctionToAuthorityIndexKey)
		return string(bz)
	}
	return ""
}

// UsesBond returns true if the bond has associated records.
func (k RecordKeeper) UsesBond(ctx sdk.Context, bondID string) bool {
	bondIDPrefix := append(PrefixBondIDToRecordsIndex, []byte(bondID)...) //nolint: all
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	return itr.Valid()
}

// RemoveBondToRecordIndexEntry removes the Bond ID -> [Record] index entry.
func (k Keeper) RemoveBondToRecordIndexEntry(ctx sdk.Context, bondID string, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getBondIDToRecordsIndexKey(bondID, id))
}

// NewRecordKeeper creates new instances of the registry RecordKeeper
func NewRecordKeeper(auctionKeeper auctionkeeper.Keeper, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) RecordKeeper {
	return RecordKeeper{
		auctionKeeper: auctionKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
	}
}

// QueryRecordsByBond - get all records for the given bond.
func (k RecordKeeper) QueryRecordsByBond(ctx sdk.Context, bondID string) []types.Record {
	var records []types.Record

	bondIDPrefix := append(PrefixBondIDToRecordsIndex, []byte(bondID)...) //nolint: all
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		cid := itr.Key()[len(bondIDPrefix):]
		bz := store.Get(append(PrefixCIDToRecordIndex, cid...))
		if bz != nil {
			var obj types.Record
			k.cdc.MustUnmarshal(bz, &obj)
			records = append(records, recordObjToRecord(store, obj))
		}
	}

	return records
}

// ProcessRenewRecord renews a record.
func (k Keeper) ProcessRenewRecord(ctx sdk.Context, msg types.MsgRenewRecord) error {
	if !k.HasRecord(ctx, msg.RecordId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "record not found")
	}

	// Check if renewal is required (i.e. expired record marked as deleted).
	record := k.GetRecord(ctx, msg.RecordId)
	expiryTime, err := time.Parse(time.RFC3339, record.ExpiryTime)
	if err != nil {
		panic(err)
	}

	if !record.Deleted || expiryTime.After(ctx.BlockTime()) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "renewal not required")
	}

	recordType := record.ToRecordType()
	err = k.processRecord(ctx, &recordType, true)
	if err != nil {
		return err
	}

	return nil
}

// ProcessAssociateBond associates a record with a bond.
func (k Keeper) ProcessAssociateBond(ctx sdk.Context, msg types.MsgAssociateBond) error {
	if !k.HasRecord(ctx, msg.RecordId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "record not found")
	}

	if !k.bondKeeper.HasBond(ctx, msg.BondId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "bond not found")
	}

	// Check if already associated with a bond.
	record := k.GetRecord(ctx, msg.RecordId)
	if record.BondId != "" || len(record.BondId) != 0 {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "bond already exists")
	}

	// Only the bond owner can associate a record with the bond.
	bond := k.bondKeeper.GetBond(ctx, msg.BondId)
	if msg.Signer != bond.Owner {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "bond owner mismatch")
	}

	record.BondId = msg.BondId
	k.PutRecord(ctx, record)
	k.AddBondToRecordIndexEntry(ctx, msg.BondId, msg.RecordId)

	// Required so that renewal is triggered (with new bond ID) for expired records.
	if record.Deleted {
		k.InsertRecordExpiryQueue(ctx, record)
	}

	return nil
}

// ProcessDissociateBond dissociates a record from its bond.
func (k Keeper) ProcessDissociateBond(ctx sdk.Context, msg types.MsgDissociateBond) error {
	if !k.HasRecord(ctx, msg.RecordId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "record not found")
	}

	// Check if associated with a bond.
	record := k.GetRecord(ctx, msg.RecordId)
	bondID := record.BondId
	if bondID == "" || len(bondID) == 0 {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "bond not found")
	}

	// Only the bond owner can dissociate a record from the bond.
	bond := k.bondKeeper.GetBond(ctx, bondID)
	if msg.Signer != bond.Owner {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "bond owner mismatch")
	}

	// Clear bond ID.
	record.BondId = ""
	k.PutRecord(ctx, record)
	k.RemoveBondToRecordIndexEntry(ctx, bondID, record.Id)

	return nil
}

// ProcessDissociateRecords dissociates all records associated with a given bond.
func (k Keeper) ProcessDissociateRecords(ctx sdk.Context, msg types.MsgDissociateRecords) error {
	if !k.bondKeeper.HasBond(ctx, msg.BondId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "bond not found")
	}

	// Only the bond owner can dissociate all records from the bond.
	bond := k.bondKeeper.GetBond(ctx, msg.BondId)
	if msg.Signer != bond.Owner {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "bond owner mismatch")
	}

	// Dissociate all records from the bond.
	records := k.recordKeeper.QueryRecordsByBond(ctx, msg.BondId)
	for _, record := range records {
		// Clear bond ID.
		record.BondId = ""
		k.PutRecord(ctx, record)
		k.RemoveBondToRecordIndexEntry(ctx, msg.BondId, record.Id)
	}

	return nil
}

// ProcessReAssociateRecords switches records from and old to new bond.
func (k Keeper) ProcessReAssociateRecords(ctx sdk.Context, msg types.MsgReAssociateRecords) error {
	if !k.bondKeeper.HasBond(ctx, msg.OldBondId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "old bond not found")
	}

	if !k.bondKeeper.HasBond(ctx, msg.NewBondId) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "new bond not found")
	}

	// Only the bond owner can re-associate all records.
	oldBond := k.bondKeeper.GetBond(ctx, msg.OldBondId)
	if msg.Signer != oldBond.Owner {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "old bond owner mismatch")
	}

	newBond := k.bondKeeper.GetBond(ctx, msg.NewBondId)
	if msg.Signer != newBond.Owner {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "new bond owner mismatch")
	}

	// Re-associate all records.
	records := k.recordKeeper.QueryRecordsByBond(ctx, msg.OldBondId)
	for _, record := range records {
		// Switch bond ID.
		record.BondId = msg.NewBondId
		k.PutRecord(ctx, record)

		k.RemoveBondToRecordIndexEntry(ctx, msg.OldBondId, record.Id)
		k.AddBondToRecordIndexEntry(ctx, msg.NewBondId, record.Id)

		// Required so that renewal is triggered (with new bond ID) for expired records.
		if record.Deleted {
			k.InsertRecordExpiryQueue(ctx, record)
		}
	}

	return nil
}
