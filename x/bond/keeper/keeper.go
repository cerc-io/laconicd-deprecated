package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/cerc-io/laconicd/x/bond/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// prefixIDToBondIndex is the prefix for ID -> Bond index in the KVStore.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var prefixIDToBondIndex = []byte{0x00}

// prefixOwnerToBondsIndex is the prefix for the Owner -> [Bond] index in the KVStore.
var prefixOwnerToBondsIndex = []byte{0x01}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	// Track bond usage in other cosmos-sdk modules (more like a usage tracker).
	usageKeepers []types.BondUsageKeeper

	storeKey storetypes.StoreKey

	cdc codec.BinaryCodec

	paramSubspace paramtypes.Subspace
}

// NewKeeper creates new instances of the bond Keeper
func NewKeeper(cdc codec.BinaryCodec, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, usageKeepers []types.BondUsageKeeper, storeKey storetypes.StoreKey, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		usageKeepers:  usageKeepers,
		paramSubspace: ps,
	}
}

// Generates Bond ID -> Bond index key.
func getBondIndexKey(id string) []byte {
	return append(prefixIDToBondIndex, []byte(id)...)
}

// Generates Owner -> Bonds index key.
func getOwnerToBondsIndexKey(owner string, bondID string) []byte {
	return append(append(prefixOwnerToBondsIndex, []byte(owner)...), []byte(bondID)...)
}

// BondID simplifies generation of bond IDs.
type BondID struct {
	Address  sdk.Address
	AccNum   uint64
	Sequence uint64
}

// Generate creates the bond ID.
func (bondID BondID) Generate() string {
	hasher := sha256.New()
	str := fmt.Sprintf("%s:%d:%d", bondID.Address.String(), bondID.AccNum, bondID.Sequence)
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// MatchBonds - get all matching bonds.
func (k Keeper) MatchBonds(ctx sdk.Context, matchFn func(*types.Bond) bool) []*types.Bond {
	var bonds []*types.Bond

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToBondIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshal(bz, &obj)
			if matchFn(&obj) {
				bonds = append(bonds, &obj)
			}
		}
	}

	return bonds
}

// CreateBond creates a new bond.
func (k Keeper) CreateBond(ctx sdk.Context, ownerAddress sdk.AccAddress, coins sdk.Coins) (*types.Bond, error) {
	// Check if account has funds.
	for _, coin := range coins {
		balance := k.bankKeeper.HasBalance(ctx, ownerAddress, coin)
		if !balance {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "failed to create bond; Insufficient funds")
		}
	}

	// Generate bond ID.
	account := k.accountKeeper.GetAccount(ctx, ownerAddress)
	bondID := BondID{
		Address:  ownerAddress,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	maxBondAmount, err := k.getMaxBondAmount(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid max bond amount.")
	}

	bond := types.Bond{Id: bondID, Owner: ownerAddress.String(), Balance: coins}
	if bond.Balance.IsAnyGT(maxBondAmount) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Max bond amount exceeded.")
	}

	// Move funds into the bond account module.
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, bond.Balance)
	if err != nil {
		return nil, err
	}

	// Save bond in store.
	k.SaveBond(ctx, &bond)

	return &bond, nil
}

// SaveBond - saves a bond to the store.
func (k Keeper) SaveBond(ctx sdk.Context, bond *types.Bond) {
	store := ctx.KVStore(k.storeKey)

	// Bond ID -> Bond index.
	store.Set(getBondIndexKey(bond.Id), k.cdc.MustMarshal(bond))

	// Owner -> [Bond] index.
	store.Set(getOwnerToBondsIndexKey(bond.Owner, bond.Id), []byte{})
}

// HasBond - checks if a bond by the given ID exists.
func (k Keeper) HasBond(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(getBondIndexKey(id))
}

// GetBond - gets a record from the store.
func (k Keeper) GetBond(ctx sdk.Context, id string) types.Bond {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getBondIndexKey(id))
	var obj types.Bond
	k.cdc.MustUnmarshal(bz, &obj)

	return obj
}

// DeleteBond - deletes the bond.
func (k Keeper) DeleteBond(ctx sdk.Context, bond types.Bond) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getBondIndexKey(bond.Id))
	store.Delete(getOwnerToBondsIndexKey(bond.Owner, bond.Id))
}

// ListBonds - get all bonds.
func (k Keeper) ListBonds(ctx sdk.Context) []*types.Bond {
	var bonds []*types.Bond

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToBondIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshal(bz, &obj)
			bonds = append(bonds, &obj)
		}
	}
	return bonds
}

// QueryBondsByOwner - query bonds by owner.
func (k Keeper) QueryBondsByOwner(ctx sdk.Context, ownerAddress string) []types.Bond {
	var bonds []types.Bond

	ownerPrefix := append(prefixOwnerToBondsIndex, []byte(ownerAddress)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, ownerPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bondID := itr.Key()[len(ownerPrefix):]
		bz := store.Get(append(prefixIDToBondIndex, bondID...))
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshal(bz, &obj)
			bonds = append(bonds, obj)
		}
	}

	return bonds
}

// RefillBond refills an existing bond.
func (k Keeper) RefillBond(ctx sdk.Context, id string, ownerAddress sdk.AccAddress, coins sdk.Coins) (*types.Bond, error) {
	if !k.HasBond(ctx, id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}

	bond := k.GetBond(ctx, id)
	if bond.Owner != ownerAddress.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	// Check if account has funds.
	for _, coin := range coins {
		if !k.bankKeeper.HasBalance(ctx, ownerAddress, coin) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Insufficient funds.")
		}
	}

	maxBondAmount, err := k.getMaxBondAmount(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid max bond amount.")
	}

	updatedBalance := bond.Balance.Add(coins...)
	if updatedBalance.IsAnyGT(maxBondAmount) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Max bond amount exceeded.")
	}

	// Move funds into the bond account module.
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Update bond balance and save.
	bond.Balance = updatedBalance
	k.SaveBond(ctx, &bond)

	return &bond, nil
}

// WithdrawBond withdraws funds from a bond.
func (k Keeper) WithdrawBond(ctx sdk.Context, id string, ownerAddress sdk.AccAddress, coins sdk.Coins) (*types.Bond, error) {
	if !k.HasBond(ctx, id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}

	bond := k.GetBond(ctx, id)
	if bond.Owner != ownerAddress.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	updatedBalance, isNeg := bond.Balance.SafeSub(coins...)
	if isNeg {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Insufficient bond balance.")
	}

	// Move funds from the bond into the account.
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddress, coins)
	if err != nil {
		return nil, err
	}

	// Update bond balance and save.
	bond.Balance = updatedBalance
	k.SaveBond(ctx, &bond)

	return &bond, nil
}

// CancelBond cancels a bond, returning funds to the owner.
func (k Keeper) CancelBond(ctx sdk.Context, id string, ownerAddress sdk.AccAddress) (*types.Bond, error) {
	if !k.HasBond(ctx, id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}

	bond := k.GetBond(ctx, id)
	if bond.Owner != ownerAddress.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	// Check if bond is used in other modules.
	for _, usageKeeper := range k.usageKeepers {
		if usageKeeper.UsesBond(ctx, id) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("Bond in use by the '%s' module.", usageKeeper.ModuleName()))
		}
	}

	// Move funds from the bond into the account.
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddress, bond.Balance)
	if err != nil {
		return nil, err
	}

	k.DeleteBond(ctx, bond)

	return &bond, nil
}

func (k Keeper) getMaxBondAmount(ctx sdk.Context) (sdk.Coins, error) {
	params := k.GetParams(ctx)
	maxBondAmount := params.MaxBondAmount
	return sdk.NewCoins(maxBondAmount), nil
}

// GetBondModuleBalances gets the bond module account(s) balances.
func (k Keeper) GetBondModuleBalances(ctx sdk.Context) sdk.Coins {
	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
	return balances
}

// TransferCoinsToModuleAccount moves funds from the bonds module account to another module account.
func (k Keeper) TransferCoinsToModuleAccount(ctx sdk.Context, id, moduleAccount string, coins sdk.Coins) error {
	if !k.HasBond(ctx, id) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond not found.")
	}

	bondObj := k.GetBond(ctx, id)

	// Deduct rent from bond.
	updatedBalance, isNeg := bondObj.Balance.SafeSub(coins...)

	if isNeg {
		// Check if bond has sufficient funds.
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Insufficient funds.")
	}

	// Move funds from bond module to record rent module.
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, moduleAccount, coins)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Error transferring funds.")
	}

	// Update bond balance.
	bondObj.Balance = updatedBalance
	k.SaveBond(ctx, &bondObj)

	return nil
}
