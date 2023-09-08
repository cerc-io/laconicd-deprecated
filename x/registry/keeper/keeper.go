package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	errorsmod "cosmossdk.io/errors"
	auctionkeeper "github.com/cerc-io/laconicd/x/auction/keeper"
	bondkeeper "github.com/cerc-io/laconicd/x/bond/keeper"
	"github.com/cerc-io/laconicd/x/registry/helpers"
	"github.com/cerc-io/laconicd/x/registry/types"
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

	//  PrefixAttributesIndex is the prefix for the registry Record.Attribute -> []Record.ID index
	PrefixAttributesIndex = []byte{0x07}

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

// NewKeeper creates new instances of the registry Keeper
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
	decodeRecordNames(store, &record)
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
			var record types.Record
			k.cdc.MustUnmarshal(bz, &record)
			decodeRecordNames(store, &record)
			records = append(records, record)
		}
	}

	return records
}

// RecordsFromAttributes gets a list of records whose attributes match all provided values
func (k Keeper) RecordsFromAttributes(ctx sdk.Context, attributes []*types.QueryListRecordsRequest_KeyValueInput, all bool) ([]types.Record, error) {
	resultRecordIds := []string{}
	for i, attr := range attributes {
		val := GetAttributeValue(attr.Value)
		attributeIndex := GetAttributesIndexKey(attr.Key, val)
		recordIds, err := k.GetAttributeMapping(ctx, attributeIndex)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			resultRecordIds = recordIds
		} else {
			resultRecordIds = getIntersection(recordIds, resultRecordIds)
		}
	}

	records := []types.Record{}
	for _, id := range resultRecordIds {
		record := k.GetRecord(ctx, id)
		if record.Deleted {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		decodeRecordNames(store, &record)
		if !all && len(record.Names) == 0 {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func GetAttributeValue(input *types.QueryListRecordsRequest_ValueInput) interface{} {
	if input.Type == "int" {
		return input.GetInt()
	}
	if input.Type == "float" {
		return input.GetFloat()
	}
	if input.Type == "string" {
		return input.GetString_()
	}
	if input.Type == "boolean" {
		return input.GetBoolean()
	}
	if input.Type == "reference" {
		return input.GetReference().GetId()
	}
	return nil
}

func getIntersection(a []string, b []string) []string {
	result := []string{}
	if len(a) < len(b) {
		for _, str := range a {
			if contains(b, str) {
				result = append(result, str)
			}
		}
	} else {
		for _, str := range b {
			if contains(a, str) {
				result = append(result, str)
			}
		}
	}
	return result
}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
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
func (k Keeper) ProcessSetRecord(ctx sdk.Context, msg types.MsgSetRecord) (*types.RecordEncodable, error) {
	payload := msg.Payload.ToReadablePayload()
	record := types.RecordEncodable{Attributes: payload.Record, BondID: msg.BondId}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	cid, err := record.GetCID()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid record JSON")
	}

	record.ID = cid

	if exists := k.HasRecord(ctx, record.ID); exists {
		// Immutable record already exists. No-op.
		return &record, nil
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := legacy.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Invalid public key.")
		}

		sigOK := pubKey.VerifySignature(resourceSignBytes, helpers.BytesFromBase64(sig.Sig))
		if !sigOK {
			fmt.Println("Signature mismatch: ", sig.PubKey)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Invalid signature.")
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

func (k Keeper) processRecord(ctx sdk.Context, record *types.RecordEncodable, isRenewal bool) error {
	params := k.GetParams(ctx)
	rent := params.RecordRent

	err := k.bondKeeper.TransferCoinsToModuleAccount(
		ctx, record.BondID, types.RecordRentModuleAccountName, sdk.NewCoins(rent),
	)
	if err != nil {
		return err
	}

	record.CreateTime = ctx.BlockHeader().Time.Format(time.RFC3339)
	record.ExpiryTime = ctx.BlockHeader().Time.Add(params.RecordRentDuration).Format(time.RFC3339)
	record.Deleted = false

	recordObj, err := record.ToRecordObj()
	if err != nil {
		return err
	}
	k.PutRecord(ctx, recordObj)

	// TODO process type here
	// recordType, ok := record.Attributes["type"].(string)

	if err := k.processAttributes(ctx, record.Attributes, record.ID, ""); err != nil {
		return err
	}

	expiryTimeKey := GetAttributesIndexKey(ExpiryTimeAttributeName, record.ExpiryTime)
	if err := k.SetAttributeMapping(ctx, expiryTimeKey, record.ID); err != nil {
		return err
	}

	k.InsertRecordExpiryQueue(ctx, recordObj)

	// Renewal doesn't change the name and bond indexes.
	if !isRenewal {
		k.AddBondToRecordIndexEntry(ctx, record.BondID, record.ID)
	}

	return nil
}

// PutRecord - saves a record to the store and updates ID -> Record index.
func (k Keeper) PutRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRecordIndexKey(record.Id), k.cdc.MustMarshal(&record))
	k.updateBlockChangeSetForRecord(ctx, record.Id)
}

func (k Keeper) processAttributes(ctx sdk.Context, attrs map[string]any, id string, prefix string) error {
	for key, value := range attrs {
		if subRecord, ok := value.(map[string]any); ok {
			k.processAttributes(ctx, subRecord, id, key)
		} else {
			indexKey := GetAttributesIndexKey(prefix+key, value)
			if err := k.SetAttributeMapping(ctx, indexKey, id); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetAttributesIndexKey(key string, value interface{}) []byte {
	keyString := fmt.Sprintf("%s%s", key, value)
	return append(PrefixAttributesIndex, []byte(keyString)...)
}

func (k Keeper) SetAttributeMapping(ctx sdk.Context, key []byte, recordID string) error {
	store := ctx.KVStore(k.storeKey)
	var recordIds []string
	if store.Has(key) {
		err := json.Unmarshal(store.Get(key), &recordIds)
		if err != nil {
			return fmt.Errorf("cannot unmarshal byte array, error, %w", err)
		}
	} else {
		recordIds = []string{}
	}
	recordIds = append(recordIds, recordID)
	bz, err := json.Marshal(recordIds)
	if err != nil {
		return fmt.Errorf("cannot marshal string array, error, %w", err)
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetAttributeMapping(ctx sdk.Context, key []byte) ([]string, error) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(key) {
		return nil, fmt.Errorf("store doesn't have key")
	}

	var recordIds []string
	if err := json.Unmarshal(store.Get(key), &recordIds); err != nil {
		return nil, fmt.Errorf("cannot unmarshal byte array, error, %w", err)
	}

	return recordIds, nil
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

// GetModuleBalances gets the registry module account(s) balances.
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

func decodeRecordNames(store sdk.KVStore, record *types.Record) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(record.Id)

	if store.Has(reverseNameIndexKey) {
		names, err := helpers.BytesArrToStringArr(store.Get(reverseNameIndexKey))
		if err != nil {
			panic(err)
		}

		record.Names = names
	}
}
