package keeper

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	auctiontypes "github.com/cerc-io/laconicd/x/auction/types"
	"github.com/cerc-io/laconicd/x/nameservice/helpers"
	"github.com/cerc-io/laconicd/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func getAuthorityPubKey(pubKey cryptotypes.PubKey) string {
	if pubKey != nil {
		return helpers.BytesToBase64(pubKey.Bytes())
	}
	return ""
}

// GetNameAuthorityIndexKey Generates name -> NameAuthority index key.
func GetNameAuthorityIndexKey(name string) []byte {
	return append(PrefixNameAuthorityRecordIndex, []byte(name)...)
}

// GetNameRecordIndexKey Generates CRN -> NameRecord index key.
func GetNameRecordIndexKey(crn string) []byte {
	return append(PrefixCRNToNameRecordIndex, []byte(crn)...)
}

func GetCIDToNamesIndexKey(id string) []byte {
	return append(PrefixCIDToNamesIndex, []byte(id)...)
}

func SetNameAuthority(ctx sdk.Context, store sdk.KVStore, codec codec.BinaryCodec, name string, authority *types.NameAuthority) {
	store.Set(GetNameAuthorityIndexKey(name), codec.MustMarshal(authority))
	updateBlockChangeSetForNameAuthority(ctx, codec, store, name)
}

// SetNameAuthority creates the NameAuthority record.
func (k Keeper) SetNameAuthority(ctx sdk.Context, name string, authority *types.NameAuthority) {
	SetNameAuthority(ctx, ctx.KVStore(k.storeKey), k.cdc, name, authority)
}

func removeAuctionToAuthorityMapping(store sdk.KVStore, auctionID string) {
	store.Delete(GetAuctionToAuthorityIndexKey(auctionID))
}

func (k Keeper) RemoveAuctionToAuthorityMapping(ctx sdk.Context, auctionID string) {
	removeAuctionToAuthorityMapping(ctx.KVStore(k.storeKey), auctionID)
}

// GetNameAuthority - gets a name authority from the store.
func GetNameAuthority(store sdk.KVStore, codec codec.BinaryCodec, name string) types.NameAuthority {
	authorityKey := GetNameAuthorityIndexKey(name)
	if !store.Has(authorityKey) {
		return types.NameAuthority{}
	}

	bz := store.Get(authorityKey)
	var obj types.NameAuthority
	codec.MustUnmarshal(bz, &obj)

	return obj
}

// GetNameAuthority - gets a name authority from the store.
func (k Keeper) GetNameAuthority(ctx sdk.Context, name string) types.NameAuthority {
	return GetNameAuthority(ctx.KVStore(k.storeKey), k.cdc, name)
}

// HasNameAuthority - checks if a name authority entry exists.
func HasNameAuthority(store sdk.KVStore, name string) bool {
	return store.Has(GetNameAuthorityIndexKey(name))
}

// HasNameAuthority - checks if a name/authority exists.
func (k Keeper) HasNameAuthority(ctx sdk.Context, name string) bool {
	return HasNameAuthority(ctx.KVStore(k.storeKey), name)
}

func getBondIDToAuthoritiesIndexKey(bondID string, name string) []byte {
	return append(append(PrefixBondIDToAuthoritiesIndex, []byte(bondID)...), []byte(name)...)
}

// AddBondToAuthorityIndexEntry adds the Bond ID -> [Authority] index entry.
func (k Keeper) AddBondToAuthorityIndexEntry(ctx sdk.Context, bondID string, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getBondIDToAuthoritiesIndexKey(bondID, name), []byte{})
}

// RemoveBondToAuthorityIndexEntry removes the Bond ID -> [Authority] index entry.
func (k Keeper) RemoveBondToAuthorityIndexEntry(ctx sdk.Context, bondID string, name string) {
	RemoveBondToAuthorityIndexEntry(ctx.KVStore(k.storeKey), bondID, name)
}

func RemoveBondToAuthorityIndexEntry(store sdk.KVStore, bondID string, name string) {
	store.Delete(getBondIDToAuthoritiesIndexKey(bondID, name))
}

func (k Keeper) updateBlockChangeSetForName(ctx sdk.Context, crn string) {
	changeSet := k.getOrCreateBlockChangeSet(ctx, ctx.BlockHeight())
	changeSet.Names = append(changeSet.Names, crn)
	k.saveBlockChangeSet(ctx, changeSet)
}

func (k Keeper) getAuthority(ctx sdk.Context, crn string) (string, *url.URL, *types.NameAuthority, error) {
	parsedCRN, err := url.Parse(crn)
	if err != nil {
		return "", nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid CRN.")
	}

	name := parsedCRN.Host
	if !k.HasNameAuthority(ctx, name) {
		return name, nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name authority not found.")
	}
	authority := k.GetNameAuthority(ctx, name)
	return name, parsedCRN, &authority, nil
}

func (k Keeper) checkCRNAccess(ctx sdk.Context, signer sdk.AccAddress, crn string) error {
	name, parsedCRN, authority, err := k.getAuthority(ctx, crn)
	if err != nil {
		return err
	}

	formattedCRN := fmt.Sprintf("crn://%s%s", name, parsedCRN.RequestURI())
	if formattedCRN != crn {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid CRN.")
	}

	if authority.OwnerAddress != signer.String() {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Access denied.")
	}

	if authority.Status != types.AuthorityActive {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Authority is not active.")
	}

	if authority.BondId == "" || len(authority.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Authority bond not found.")
	}

	if authority.OwnerPublicKey == "" {
		// Try to set owner public key if account has it available now.
		ownerAccount := k.accountKeeper.GetAccount(ctx, signer)
		pubKey := ownerAccount.GetPubKey()
		if pubKey != nil {
			// Update public key in authority record.
			authority.OwnerPublicKey = getAuthorityPubKey(pubKey)
			k.SetNameAuthority(ctx, name, authority)
		}
	}

	return nil
}

// HasNameRecord - checks if a name record exists.
func (k Keeper) HasNameRecord(ctx sdk.Context, crn string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetNameRecordIndexKey(crn))
}

// GetNameRecord - gets a name record from the store.
func GetNameRecord(store sdk.KVStore, codec codec.BinaryCodec, crn string) *types.NameRecord {
	nameRecordKey := GetNameRecordIndexKey(crn)
	if !store.Has(nameRecordKey) {
		return nil
	}

	bz := store.Get(nameRecordKey)
	var obj types.NameRecord
	codec.MustUnmarshal(bz, &obj)

	return &obj
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, crn string) *types.NameRecord {
	_, _, authority, err := k.getAuthority(ctx, crn)
	if err != nil || authority.Status != types.AuthorityActive {
		// If authority is not active (or any other error), lookup fails.
		return nil
	}

	nameRecord := GetNameRecord(ctx.KVStore(k.storeKey), k.cdc, crn)

	// Name record may not exist.
	if nameRecord == nil {
		return nil
	}

	// Name lookup should fail if the name record is stale.
	// i.e. authority was registered later than the name.
	if authority.Height > nameRecord.Latest.Height {
		return nil
	}

	return nameRecord
}

// RemoveRecordToNameMapping removes a name from the record ID -> []names index.
func RemoveRecordToNameMapping(store sdk.KVStore, id string, crn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	names, _ := helpers.BytesArrToStringArr(store.Get(reverseNameIndexKey))
	nameSet := helpers.SliceToSet(names)
	nameSet.Remove(crn)

	if nameSet.Cardinality() == 0 {
		// Delete as storing empty slice throws error from baseapp.
		store.Delete(reverseNameIndexKey)
	} else {
		data, _ := helpers.StrArrToBytesArr(helpers.SetToSlice(nameSet))
		store.Set(reverseNameIndexKey, data)
	}
}

// AddRecordToNameMapping adds a name to the record ID -> []names index.
func AddRecordToNameMapping(store sdk.KVStore, id string, crn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	if store.Has(reverseNameIndexKey) {
		names, _ = helpers.BytesArrToStringArr(store.Get(reverseNameIndexKey))
	}

	nameSet := helpers.SliceToSet(names)
	nameSet.Add(crn)
	bz, _ := helpers.StrArrToBytesArr(helpers.SetToSlice(nameSet))
	store.Set(reverseNameIndexKey, bz)
}

// SetNameRecord - sets a name record.
func SetNameRecord(store sdk.KVStore, codec codec.BinaryCodec, crn string, id string, height int64) {
	nameRecordIndexKey := GetNameRecordIndexKey(crn)

	var nameRecord types.NameRecord
	if store.Has(nameRecordIndexKey) {
		bz := store.Get(nameRecordIndexKey)
		codec.MustUnmarshal(bz, &nameRecord)
		nameRecord.History = append(nameRecord.History, nameRecord.Latest)

		// Update old CID -> []Name index.
		if nameRecord.Latest.Id != "" || len(nameRecord.Latest.Id) != 0 {
			RemoveRecordToNameMapping(store, nameRecord.Latest.Id, crn)
		}
	}

	nameRecord.Latest = &types.NameRecordEntry{
		Id:     id,
		Height: uint64(height),
	}

	store.Set(nameRecordIndexKey, codec.MustMarshal(&nameRecord))

	// Update new CID -> []Name index.
	if id != "" {
		AddRecordToNameMapping(store, id, crn)
	}
}

// SetNameRecord - sets a name record.
func (k Keeper) SetNameRecord(ctx sdk.Context, crn string, id string) {
	SetNameRecord(ctx.KVStore(k.storeKey), k.cdc, crn, id, ctx.BlockHeight())

	// Update changeSet for name.
	k.updateBlockChangeSetForName(ctx, crn)
}

// ProcessSetName creates a CRN -> Record ID mapping.
func (k Keeper) ProcessSetName(ctx sdk.Context, msg types.MsgSetName) error {
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	err = k.checkCRNAccess(ctx, signerAddress, msg.Crn)
	if err != nil {
		return err
	}

	nameRecord := k.GetNameRecord(ctx, msg.Crn)
	if nameRecord != nil && nameRecord.Latest.Id == msg.Cid {
		return nil
	}

	k.SetNameRecord(ctx, msg.Crn, msg.Cid)

	return nil
}

// ListNameRecords - get all name records.
func (k Keeper) ListNameRecords(ctx sdk.Context) []types.NameEntry {
	var nameEntries []types.NameEntry
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixCRNToNameRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			k.cdc.MustUnmarshal(bz, &record)
			nameEntries = append(nameEntries, types.NameEntry{
				Name:  string(itr.Key()[len(PrefixCRNToNameRecordIndex):]),
				Entry: &record,
			})
		}
	}

	return nameEntries
}

// ProcessReserveSubAuthority reserves a sub-authority.
func (k Keeper) ProcessReserveSubAuthority(ctx sdk.Context, name string, msg types.MsgReserveAuthority) error {
	// Get parent authority name.
	names := strings.Split(name, ".")
	parent := strings.Join(names[1:], ".")

	// Check if parent authority exists.
	if !k.HasNameAuthority(ctx, parent) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Parent authority not found.")
	}
	parentAuthority := k.GetNameAuthority(ctx, parent)

	// Sub-authority creator needs to be the owner of the parent authority.
	if parentAuthority.OwnerAddress != msg.Signer {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Access denied.")
	}

	// Sub-authority owner defaults to parent authority owner.
	subAuthorityOwner := msg.Signer
	if len(msg.Owner) != 0 {
		// Override sub-authority owner if provided in message.
		subAuthorityOwner = msg.Owner
	}

	sdkErr := k.createAuthority(ctx, name, subAuthorityOwner, false)
	if sdkErr != nil {
		return sdkErr
	}

	return nil
}

func GetAuctionToAuthorityIndexKey(auctionID string) []byte {
	return append(PrefixAuctionToAuthorityNameIndex, []byte(auctionID)...)
}

func (k Keeper) AddAuctionToAuthorityMapping(ctx sdk.Context, auctionID string, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetAuctionToAuthorityIndexKey(auctionID), []byte(name))
}

func (k Keeper) createAuthority(ctx sdk.Context, name string, owner string, isRoot bool) error {
	moduleParams := k.GetParams(ctx)

	if k.HasNameAuthority(ctx, name) {
		authority := k.GetNameAuthority(ctx, name)
		if authority.Status != types.AuthorityExpired {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name already reserved.")
		}
	}

	ownerAddress, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid owner address.")
	}
	ownerAccount := k.accountKeeper.GetAccount(ctx, ownerAddress)
	if ownerAccount == nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "Account not found.")
	}

	authority := types.NameAuthority{
		OwnerPublicKey: getAuthorityPubKey(ownerAccount.GetPubKey()),
		OwnerAddress:   owner,
		Height:         uint64(ctx.BlockHeight()),
		Status:         types.AuthorityActive,
		AuctionId:      "",
		BondId:         "",
		ExpiryTime:     ctx.BlockTime().Add(moduleParams.AuthorityGracePeriod),
	}

	if isRoot && moduleParams.AuthorityAuctionEnabled {
		// If auctions are enabled, clear out owner fields. They will be set after a winner is picked.
		authority.OwnerAddress = ""
		authority.OwnerPublicKey = ""
		// Reset bond ID if required.
		if authority.BondId != "" || len(authority.BondId) != 0 {
			k.RemoveBondToAuthorityIndexEntry(ctx, authority.BondId, name)
			authority.BondId = ""
		}

		params := auctiontypes.Params{
			CommitsDuration: moduleParams.AuthorityAuctionCommitsDuration,
			RevealsDuration: moduleParams.AuthorityAuctionRevealsDuration,
			CommitFee:       moduleParams.AuthorityAuctionCommitFee,
			RevealFee:       moduleParams.AuthorityAuctionRevealFee,
			MinimumBid:      moduleParams.AuthorityAuctionMinimumBid,
		}

		// Create an auction.
		msg := auctiontypes.NewMsgCreateAuction(params, ownerAddress)

		auction, sdkErr := k.auctionKeeper.CreateAuction(ctx, msg)
		if sdkErr != nil {
			return sdkErr
		}

		// Create auction ID -> authority name index.
		k.AddAuctionToAuthorityMapping(ctx, auction.Id, name)

		authority.Status = types.AuthorityUnderAuction
		authority.AuctionId = auction.Id
		authority.ExpiryTime = auction.RevealsEndTime.Add(moduleParams.AuthorityGracePeriod)
	}

	k.SetNameAuthority(ctx, name, &authority)
	k.InsertAuthorityExpiryQueue(ctx, name, authority.ExpiryTime)

	return nil
}

// ProcessReserveAuthority reserves a name authority.
func (k Keeper) ProcessReserveAuthority(ctx sdk.Context, msg types.MsgReserveAuthority) error {
	crn := fmt.Sprintf("crn://%s", msg.GetName())
	parsedCrn, err := url.Parse(crn)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid name")
	}
	name := parsedCrn.Host
	if fmt.Sprintf("crn://%s", name) != crn {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid name")
	}
	if strings.Contains(name, ".") {
		return k.ProcessReserveSubAuthority(ctx, name, msg)
	}
	err = k.createAuthority(ctx, name, msg.GetSigner(), true)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) ProcessSetAuthorityBond(ctx sdk.Context, msg types.MsgSetAuthorityBond) error {
	name := msg.GetName()
	signer := msg.GetSigner()
	if !k.HasNameAuthority(ctx, name) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name authority not found.")
	}
	authority := k.GetNameAuthority(ctx, name)
	if authority.OwnerAddress != signer {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Access denied")
	}

	if !k.bondKeeper.HasBond(ctx, msg.BondId) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}
	//
	bond := k.bondKeeper.GetBond(ctx, msg.BondId)
	if bond.Owner != signer {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	// No-op if bond hasn't changed.
	if authority.BondId == msg.BondId {
		return nil
	}

	// Remove old bond ID mapping, if any.
	if authority.BondId != "" {
		k.RemoveBondToAuthorityIndexEntry(ctx, authority.BondId, name)
	}

	// Update bond ID for authority.
	authority.BondId = bond.Id
	k.SetNameAuthority(ctx, name, &authority)
	// Add new bond ID mapping.
	k.AddBondToAuthorityIndexEntry(ctx, authority.BondId, name)
	return nil
}

// ProcessDeleteName removes a CRN -> Record ID mapping.
func (k Keeper) ProcessDeleteName(ctx sdk.Context, msg types.MsgDeleteNameAuthority) error {
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	err = k.checkCRNAccess(ctx, signerAddress, msg.Crn)
	if err != nil {
		return err
	}

	if !k.HasNameRecord(ctx, msg.Crn) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name not found.")
	}

	// Set CID to empty string.
	k.SetNameRecord(ctx, msg.Crn, "")

	return nil
}

func (k Keeper) GetAuthorityExpiryQueue(ctx sdk.Context) []*types.ExpiryQueueRecord {
	var authorities []*types.ExpiryQueueRecord

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixExpiryTimeToAuthoritiesIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		var record []string
		record, err := helpers.BytesArrToStringArr(itr.Value())
		if err != nil {
			return authorities
		}
		authorities = append(authorities, &types.ExpiryQueueRecord{
			Id:    string(itr.Key()[len(PrefixExpiryTimeToAuthoritiesIndex):]),
			Value: record,
		})
	}

	return authorities
}

// ResolveCRN resolves a CRN to a record.
func (k Keeper) ResolveCRN(ctx sdk.Context, crn string) *types.Record {
	_, _, authority, err := k.getAuthority(ctx, crn)
	if err != nil || authority.Status != types.AuthorityActive {
		// If authority is not active (or any other error), resolution fails.
		return nil
	}

	// Name should not resolve if it's stale.
	// i.e. authority was registered later than the name.
	record, nameRecord := ResolveCRN(ctx.KVStore(k.storeKey), crn, k, ctx)
	if authority.Height > nameRecord.Latest.Height {
		return nil
	}

	return record
}

// ResolveCRN resolves a CRN to a record.
func ResolveCRN(store sdk.KVStore, crn string, k Keeper, c sdk.Context) (*types.Record, *types.NameRecord) {
	nameKey := GetNameRecordIndexKey(crn)

	if store.Has(nameKey) {
		bz := store.Get(nameKey)
		var obj types.NameRecord
		k.cdc.MustUnmarshal(bz, &obj)

		recordExists := k.HasRecord(c, obj.Latest.Id)
		if !recordExists || obj.Latest.Id == "" {
			return nil, &obj
		}

		record := k.GetRecord(c, obj.Latest.Id)
		return &record, &obj
	}

	return nil, nil
}

func getAuthorityExpiryQueueTimeKey(timestamp time.Time) []byte {
	timeBytes := sdk.FormatTimeBytes(timestamp)
	return append(PrefixExpiryTimeToAuthoritiesIndex, timeBytes...)
}

func (k Keeper) InsertAuthorityExpiryQueue(ctx sdk.Context, name string, expiryTime time.Time) {
	timeSlice := k.GetAuthorityExpiryQueueTimeSlice(ctx, expiryTime)
	timeSlice = append(timeSlice, name)
	k.SetAuthorityExpiryQueueTimeSlice(ctx, expiryTime, timeSlice)
}

func (k Keeper) GetAuthorityExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getAuthorityExpiryQueueTimeKey(timestamp))
	if bz == nil {
		return []string{}
	}

	names, err := helpers.BytesArrToStringArr(bz)
	if err != nil {
		return []string{}
	}

	return names
}

func (k Keeper) SetAuthorityExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time, names []string) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := helpers.StrArrToBytesArr(names)
	store.Set(getAuthorityExpiryQueueTimeKey(timestamp), bz)
}

// ProcessAuthorityExpiryQueue tries to renew expiring authorities (by collecting rent) else marks them as expired.
func (k Keeper) ProcessAuthorityExpiryQueue(ctx sdk.Context) {
	names := k.GetAllExpiredAuthorities(ctx, ctx.BlockHeader().Time)
	for _, name := range names {
		authority := k.GetNameAuthority(ctx, name)

		// If authority doesn't have an associated bond or if bond no longer exists, mark it expired.
		if authority.BondId == "" || !k.bondKeeper.HasBond(ctx, authority.BondId) {
			authority.Status = types.AuthorityExpired
			k.SetNameAuthority(ctx, name, &authority)
			k.DeleteAuthorityExpiryQueue(ctx, name, authority)

			ctx.Logger().Info(fmt.Sprintf("Marking authority expired as no bond present: %s", name))

			return
		}

		// Try to renew the authority by taking rent.
		k.TryTakeAuthorityRent(ctx, name, authority)
	}
}

// DeleteAuthorityExpiryQueueTimeSlice deletes a specific authority expiry queue timeslice.
func (k Keeper) DeleteAuthorityExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getAuthorityExpiryQueueTimeKey(timestamp))
}

// DeleteAuthorityExpiryQueue deletes an authority name from the authority expiry queue.
func (k Keeper) DeleteAuthorityExpiryQueue(ctx sdk.Context, name string, authority types.NameAuthority) {
	timeSlice := k.GetAuthorityExpiryQueueTimeSlice(ctx, authority.ExpiryTime)
	newTimeSlice := []string{}

	for _, existingName := range timeSlice {
		if !bytes.Equal([]byte(existingName), []byte(name)) {
			newTimeSlice = append(newTimeSlice, existingName)
		}
	}

	if len(newTimeSlice) == 0 {
		k.DeleteAuthorityExpiryQueueTimeSlice(ctx, authority.ExpiryTime)
	} else {
		k.SetAuthorityExpiryQueueTimeSlice(ctx, authority.ExpiryTime, newTimeSlice)
	}
}

// GetAllExpiredAuthorities returns a concatenated list of all the timeslices before currTime.
func (k Keeper) GetAllExpiredAuthorities(ctx sdk.Context, currTime time.Time) (expiredAuthorityNames []string) {
	// Gets an iterator for all timeslices from time 0 until the current block header time.
	itr := k.AuthorityExpiryQueueIterator(ctx, ctx.BlockHeader().Time)
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		timeslice := []string{}
		timeslice, err := helpers.BytesArrToStringArr(itr.Value())
		if err != nil {
			panic(err)
		}

		expiredAuthorityNames = append(expiredAuthorityNames, timeslice...)
	}

	return expiredAuthorityNames
}

// AuthorityExpiryQueueIterator returns all the authority expiry queue timeslices from time 0 until endTime.
func (k Keeper) AuthorityExpiryQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	rangeEndBytes := sdk.InclusiveEndBytes(getAuthorityExpiryQueueTimeKey(endTime))
	return store.Iterator(PrefixExpiryTimeToAuthoritiesIndex, rangeEndBytes)
}

// TryTakeAuthorityRent tries to take rent from the authority bond.
func (k Keeper) TryTakeAuthorityRent(ctx sdk.Context, name string, authority types.NameAuthority) {
	ctx.Logger().Info(fmt.Sprintf("Trying to take rent for authority: %s", name))

	params := k.GetParams(ctx)
	rent := params.AuthorityRent
	sdkErr := k.bondKeeper.TransferCoinsToModuleAccount(ctx, authority.BondId, types.AuthorityRentModuleAccountName, sdk.NewCoins(rent))

	if sdkErr != nil {
		// Insufficient funds, mark authority as expired.
		authority.Status = types.AuthorityExpired
		k.SetNameAuthority(ctx, name, &authority)
		k.DeleteAuthorityExpiryQueue(ctx, name, authority)

		ctx.Logger().Info(fmt.Sprintf("Insufficient funds in owner account to pay authority rent, marking as expired: %s", name))

		return
	}

	// Delete old expiry queue entry, create new one.
	k.DeleteAuthorityExpiryQueue(ctx, name, authority)
	authority.ExpiryTime = ctx.BlockTime().Add(params.AuthorityRentDuration)
	k.InsertAuthorityExpiryQueue(ctx, name, authority.ExpiryTime)

	// Save authority.
	authority.Status = types.AuthorityActive
	k.SetNameAuthority(ctx, name, &authority)
	k.AddBondToAuthorityIndexEntry(ctx, authority.BondId, name)

	ctx.Logger().Info(fmt.Sprintf("Authority rent paid successfully: %s", name))
}

// ListNameAuthorityRecords - get all name authority records.
func (k Keeper) ListNameAuthorityRecords(ctx sdk.Context) map[string]types.NameAuthority {
	nameAuthorityRecords := make(map[string]types.NameAuthority)
	store := ctx.KVStore(k.storeKey)

	itr := sdk.KVStorePrefixIterator(store, PrefixNameAuthorityRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameAuthority
			k.cdc.MustUnmarshal(bz, &record)
			nameAuthorityRecords[string(itr.Key()[len(PrefixNameAuthorityRecordIndex):])] = record
		}
	}

	return nameAuthorityRecords
}
