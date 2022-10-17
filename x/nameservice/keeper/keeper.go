package keeper

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	auctionkeeper "github.com/cerc-io/laconicd/x/auction/keeper"
	bondkeeper "github.com/cerc-io/laconicd/x/bond/keeper"
	"github.com/cerc-io/laconicd/x/nameservice/helpers"
	"github.com/cerc-io/laconicd/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (

	// PrefixCIDToRecordIndex is the prefix for CID -> Record index.
	// Note: This is the primary index in the system.
	// Note: Golang doesn't support const arrays.
	PrefixCIDToRecordIndex = []byte{0x00}

	// PrefixNameAuthorityRecordIndex is the prefix for the name -> NameAuthority index.
	PrefixNameAuthorityRecordIndex = []byte{0x01}

	// PrefixCRNToNameRecordIndex is the prefix for the CRN -> NamingRecord index.
	PrefixCRNToNameRecordIndex = []byte{0x02}

	// PrefixBondIDToRecordsIndex is the prefix for the Bond ID -> [Record] index.
	PrefixBondIDToRecordsIndex = []byte{0x03}

	// PrefixBlockChangesetIndex is the prefix for the block changeset index.
	PrefixBlockChangesetIndex = []byte{0x04}

	// PrefixAuctionToAuthorityNameIndex is the prefix for the auction ID -> authority name index.
	PrefixAuctionToAuthorityNameIndex = []byte{0x05}

	// PrefixBondIDToAuthoritiesIndex is the prefix for the Bond ID -> [Authority] index.
	PrefixBondIDToAuthoritiesIndex = []byte{0x06}

	// PrefixExpiryTimeToRecordsIndex is the prefix for the Expiry Time -> [Record] index.
	PrefixExpiryTimeToRecordsIndex = []byte{0x10}

	// PrefixExpiryTimeToAuthoritiesIndex is the prefix for the Expiry Time -> [Authority] index.
	PrefixExpiryTimeToAuthoritiesIndex = []byte{0x11}

	// PrefixCIDToNamesIndex the the reverse index for naming, i.e. maps CID -> []Names.
	// TODO(ashwin): Move out of WNS once we have an indexing service.
	PrefixCIDToNamesIndex = []byte{0xe0}
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	recordKeeper  RecordKeeper
	bondKeeper    bondkeeper.Keeper
	auctionKeeper auctionkeeper.Keeper

	storeKey storetypes.StoreKey // Unexposed key to access store from sdk.Context

	cdc codec.BinaryCodec // The wire codec for binary encoding/decoding.

	paramSubspace paramtypes.Subspace
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(cdc codec.BinaryCodec, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, recordKeeper RecordKeeper,
	bondKeeper bondkeeper.Keeper, auctionKeeper auctionkeeper.Keeper, storeKey storetypes.StoreKey, ps paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		recordKeeper:  recordKeeper,
		bondKeeper:    bondKeeper,
		auctionKeeper: auctionKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		paramSubspace: ps,
	}
}

// GetRecordIndexKey Generates Bond ID -> Bond index key.
func GetRecordIndexKey(id string) []byte {
	return append(PrefixCIDToRecordIndex, []byte(id)...)
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetRecordIndexKey(id))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(ctx sdk.Context, id string) (record types.Record) {
	store := ctx.KVStore(k.storeKey)
	result := store.Get(GetRecordIndexKey(id))
	k.cdc.MustUnmarshal(result, &record)
	return record
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixCIDToRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Record
			k.cdc.MustUnmarshal(bz, &obj)
			records = append(records, recordObjToRecord(store, k.cdc, obj))
		}
	}

	return records
}

// MatchRecords - get all matching records.
func (k Keeper) MatchRecords(ctx sdk.Context, matchFn func(*types.RecordType) bool) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixCIDToRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Record
			k.cdc.MustUnmarshal(bz, &obj)
			obj = recordObjToRecord(store, k.cdc, obj)
			record := obj.ToRecordType()
			if matchFn(&record) {
				records = append(records, obj)
			}
		}
	}

	return records
}

func (k Keeper) GetRecordExpiryQueue(ctx sdk.Context) []*types.ExpiryQueueRecord {
	var records []*types.ExpiryQueueRecord

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixExpiryTimeToRecordsIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		record, err := helpers.BytesArrToStringArr(itr.Value())
		if err != nil {
			return records
		}
		records = append(records, &types.ExpiryQueueRecord{
			Id:    string(itr.Key()[len(PrefixExpiryTimeToRecordsIndex):]),
			Value: record,
		})
	}

	return records
}

// ProcessSetRecord creates a record.
func (k Keeper) ProcessSetRecord(ctx sdk.Context, msg types.MsgSetRecord) (*types.RecordType, error) {
	payload := msg.Payload.ToReadablePayload()
	record := types.RecordType{Attributes: payload.Record, BondId: msg.BondId}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	cid, err := record.GetCID()
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid record JSON")
	}

	record.Id = cid

	if exists := k.HasRecord(ctx, record.Id); exists {
		// Immutable record already exists. No-op.
		return &record, nil
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := legacy.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Invalid public key.")
		}

		sigOK := pubKey.VerifySignature(resourceSignBytes, helpers.BytesFromBase64(sig.Sig))
		if !sigOK {
			fmt.Println("Signature mismatch: ", sig.PubKey)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Invalid signature.")
		}
		record.Owners = append(record.Owners, pubKey.Address().String())
	}

	// Sort owners list.
	sort.Strings(record.Owners)
	sdkErr := k.processRecord(ctx, &record, false)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return &record, nil
}

func (k Keeper) processRecord(ctx sdk.Context, record *types.RecordType, isRenewal bool) error {
	params := k.GetParams(ctx)
	rent := params.RecordRent

	err := k.bondKeeper.TransferCoinsToModuleAccount(ctx, record.BondId, types.RecordRentModuleAccountName, sdk.NewCoins(rent))
	if err != nil {
		return err
	}

	record.CreateTime = ctx.BlockHeader().Time.Format(time.RFC3339)
	record.ExpiryTime = ctx.BlockHeader().Time.Add(params.RecordRentDuration).Format(time.RFC3339)
	record.Deleted = false

	k.PutRecord(ctx, record.ToRecordObj())
	k.InsertRecordExpiryQueue(ctx, record.ToRecordObj())

	// Renewal doesn't change the name and bond indexes.
	if !isRenewal {
		k.AddBondToRecordIndexEntry(ctx, record.BondId, record.Id)
	}

	return nil
}

// PutRecord - saves a record to the store and updates ID -> Record index.
func (k Keeper) PutRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRecordIndexKey(record.Id), k.cdc.MustMarshal(&record))
	k.updateBlockChangeSetForRecord(ctx, record.Id)
}

// AddBondToRecordIndexEntry adds the Bond ID -> [Record] index entry.
func (k Keeper) AddBondToRecordIndexEntry(ctx sdk.Context, bondID string, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getBondIDToRecordsIndexKey(bondID, id), []byte{})
}

// Generates Bond ID -> Records index key.
func getBondIDToRecordsIndexKey(bondID string, id string) []byte {
	return append(append(PrefixBondIDToRecordsIndex, []byte(bondID)...), []byte(id)...)
}

// getRecordExpiryQueueTimeKey gets the prefix for the record expiry queue.
func getRecordExpiryQueueTimeKey(timestamp time.Time) []byte {
	timeBytes := sdk.FormatTimeBytes(timestamp)
	return append(PrefixExpiryTimeToRecordsIndex, timeBytes...)
}

// SetRecordExpiryQueueTimeSlice sets a specific record expiry queue timeslice.
func (k Keeper) SetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time, cids []string) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := helpers.StrArrToBytesArr(cids)
	store.Set(getRecordExpiryQueueTimeKey(timestamp), bz)
}

// DeleteRecordExpiryQueueTimeSlice deletes a specific record expiry queue timeslice.
func (k Keeper) DeleteRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getRecordExpiryQueueTimeKey(timestamp))
}

// GetRecordExpiryQueueTimeSlice gets a specific record queue timeslice.
// A timeslice is a slice of CIDs corresponding to records that expire at a certain time.
func (k Keeper) GetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getRecordExpiryQueueTimeKey(timestamp))
	if bz == nil {
		return []string{}
	}
	cids, err := helpers.BytesArrToStringArr(bz)
	if err != nil {
		return []string{}
	}
	return cids
}

// InsertRecordExpiryQueue inserts a record CID to the appropriate timeslice in the record expiry queue.
func (k Keeper) InsertRecordExpiryQueue(ctx sdk.Context, val types.Record) {
	expiryTime, err := time.Parse(time.RFC3339, val.ExpiryTime)
	if err != nil {
		panic(err)
	}

	timeSlice := k.GetRecordExpiryQueueTimeSlice(ctx, expiryTime)
	timeSlice = append(timeSlice, val.Id)
	k.SetRecordExpiryQueueTimeSlice(ctx, expiryTime, timeSlice)
}

// DeleteRecordExpiryQueue deletes a record CID from the record expiry queue.
func (k Keeper) DeleteRecordExpiryQueue(ctx sdk.Context, record types.Record) {
	expiryTime, err := time.Parse(time.RFC3339, record.ExpiryTime)
	if err != nil {
		panic(err)
	}

	timeSlice := k.GetRecordExpiryQueueTimeSlice(ctx, expiryTime)
	var newTimeSlice []string

	for _, cid := range timeSlice {
		if !bytes.Equal([]byte(cid), []byte(record.Id)) {
			newTimeSlice = append(newTimeSlice, cid)
		}
	}

	if len(newTimeSlice) == 0 {
		k.DeleteRecordExpiryQueueTimeSlice(ctx, expiryTime)
	} else {
		k.SetRecordExpiryQueueTimeSlice(ctx, expiryTime, newTimeSlice)
	}
}

// RecordExpiryQueueIterator returns all the record expiry queue timeslices from time 0 until endTime.
func (k Keeper) RecordExpiryQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	rangeEndBytes := sdk.InclusiveEndBytes(getRecordExpiryQueueTimeKey(endTime))
	return store.Iterator(PrefixExpiryTimeToRecordsIndex, rangeEndBytes)
}

// GetAllExpiredRecords returns a concatenated list of all the timeslices before currTime.
func (k Keeper) GetAllExpiredRecords(ctx sdk.Context, currTime time.Time) (expiredRecordCIDs []string) {
	// Gets an iterator for all timeslices from time 0 until the current block header time.
	itr := k.RecordExpiryQueueIterator(ctx, ctx.BlockHeader().Time)
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		timeslice, err := helpers.BytesArrToStringArr(itr.Value())
		if err != nil {
			panic(err)
		}

		expiredRecordCIDs = append(expiredRecordCIDs, timeslice...)
	}

	return expiredRecordCIDs
}

// ProcessRecordExpiryQueue tries to renew expiring records (by collecting rent) else marks them as deleted.
func (k Keeper) ProcessRecordExpiryQueue(ctx sdk.Context) {
	cids := k.GetAllExpiredRecords(ctx, ctx.BlockHeader().Time)
	for _, cid := range cids {
		record := k.GetRecord(ctx, cid)

		// If record doesn't have an associated bond or if bond no longer exists, mark it deleted.
		if record.BondId == "" || !k.bondKeeper.HasBond(ctx, record.BondId) {
			record.Deleted = true
			k.PutRecord(ctx, record)
			k.DeleteRecordExpiryQueue(ctx, record)

			return
		}

		// Try to renew the record by taking rent.
		k.TryTakeRecordRent(ctx, record)
	}
}

// TryTakeRecordRent tries to take rent from the record bond.
func (k Keeper) TryTakeRecordRent(ctx sdk.Context, record types.Record) {
	params := k.GetParams(ctx)
	rent := params.RecordRent
	sdkErr := k.bondKeeper.TransferCoinsToModuleAccount(ctx, record.BondId, types.RecordRentModuleAccountName, sdk.NewCoins(rent))

	if sdkErr != nil {
		// Insufficient funds, mark record as deleted.
		record.Deleted = true
		k.PutRecord(ctx, record)
		k.DeleteRecordExpiryQueue(ctx, record)

		return
	}

	// Delete old expiry queue entry, create new one.
	k.DeleteRecordExpiryQueue(ctx, record)
	record.ExpiryTime = ctx.BlockHeader().Time.Add(params.RecordRentDuration).Format(time.RFC3339)
	k.InsertRecordExpiryQueue(ctx, record)

	// Save record.
	record.Deleted = false
	k.PutRecord(ctx, record)
	k.AddBondToRecordIndexEntry(ctx, record.BondId, record.Id)
}

// GetModuleBalances gets the nameservice module account(s) balances.
func (k Keeper) GetModuleBalances(ctx sdk.Context) []*types.AccountBalance {
	var balances []*types.AccountBalance
	accountNames := []string{types.RecordRentModuleAccountName, types.AuthorityRentModuleAccountName}

	for _, accountName := range accountNames {
		moduleAddress := k.accountKeeper.GetModuleAddress(accountName)
		moduleAccount := k.accountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			accountBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
			balances = append(balances, &types.AccountBalance{
				AccountName: accountName,
				Balance:     accountBalance,
			})
		}
	}

	return balances
}

func recordObjToRecord(store sdk.KVStore, codec codec.BinaryCodec, record types.Record) types.Record {
	reverseNameIndexKey := GetCIDToNamesIndexKey(record.Id)

	if store.Has(reverseNameIndexKey) {
		names, err := helpers.BytesArrToStringArr(store.Get(reverseNameIndexKey))
		if err != nil {
			panic(err)
		}

		record.Names = names
	}

	return record
}
