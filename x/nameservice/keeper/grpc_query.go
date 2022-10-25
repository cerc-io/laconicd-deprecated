package keeper

import (
	"context"

	"github.com/cerc-io/laconicd/x/nameservice/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// BondIDAttributeName denotes the record bond ID.
const BondIDAttributeName = "bondId"

// ExpiryTimeAttributeName denotes the record expiry time.
const ExpiryTimeAttributeName = "expiryTime"

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: &params}, nil
}

func (q Querier) ListRecords(c context.Context, req *types.QueryListRecordsRequest) (*types.QueryListRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	attributes := req.GetAttributes()
	all := req.GetAll()
	var records []types.Record
	if len(attributes) > 0 {
		var err error
		records, err = q.Keeper.RecordsFromAttributes(ctx, attributes, all)
		if err != nil {
			return nil, err
		}
	} else {
		records = q.Keeper.ListRecords(ctx)
	}

	return &types.QueryListRecordsResponse{Records: records}, nil
}

func (q Querier) GetRecord(c context.Context, req *types.QueryRecordByIDRequest) (*types.QueryRecordByIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := req.GetId()
	if !q.Keeper.HasRecord(ctx, id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Record not found.")
	}
	record := q.Keeper.GetRecord(ctx, id)
	return &types.QueryRecordByIDResponse{Record: record}, nil
}

func (q Querier) GetRecordByBondID(c context.Context, req *types.QueryRecordByBondIDRequest) (*types.QueryRecordByBondIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	records := q.recordKeeper.QueryRecordsByBond(ctx, req.GetId())
	return &types.QueryRecordByBondIDResponse{Records: records}, nil
}

func (q Querier) GetNameServiceModuleBalance(c context.Context, _ *types.GetNameServiceModuleBalanceRequest) (*types.GetNameServiceModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balances := q.Keeper.GetModuleBalances(ctx)
	return &types.GetNameServiceModuleBalanceResponse{
		Balances: balances,
	}, nil
}

func (q Querier) ListNameRecords(c context.Context, _ *types.QueryListNameRecordsRequest) (*types.QueryListNameRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nameRecords := q.Keeper.ListNameRecords(ctx)
	return &types.QueryListNameRecordsResponse{Names: nameRecords}, nil
}

func (q Querier) Whois(c context.Context, request *types.QueryWhoisRequest) (*types.QueryWhoisResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nameAuthority := q.Keeper.GetNameAuthority(ctx, request.GetName())
	return &types.QueryWhoisResponse{NameAuthority: nameAuthority}, nil
}

func (q Querier) LookupCrn(c context.Context, req *types.QueryLookupCrn) (*types.QueryLookupCrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	crn := req.GetCrn()
	if !q.Keeper.HasNameRecord(ctx, crn) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "CRN not found.")
	}
	nameRecord := q.Keeper.GetNameRecord(ctx, crn)
	if nameRecord == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "name record not found.")
	}
	return &types.QueryLookupCrnResponse{Name: nameRecord}, nil
}

func (q Querier) ResolveCrn(c context.Context, req *types.QueryResolveCrn) (*types.QueryResolveCrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	crn := req.GetCrn()
	record := q.Keeper.ResolveCRN(ctx, crn)
	if record == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "record not found.")
	}
	return &types.QueryResolveCrnResponse{Record: record}, nil
}

func (q Querier) GetRecordExpiryQueue(c context.Context, _ *types.QueryGetRecordExpiryQueue) (*types.QueryGetRecordExpiryQueueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	records := q.Keeper.GetRecordExpiryQueue(ctx)
	return &types.QueryGetRecordExpiryQueueResponse{Records: records}, nil
}

func (q Querier) GetAuthorityExpiryQueue(c context.Context, _ *types.QueryGetAuthorityExpiryQueue) (*types.QueryGetAuthorityExpiryQueueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	authorities := q.Keeper.GetAuthorityExpiryQueue(ctx)
	return &types.QueryGetAuthorityExpiryQueueResponse{Authorities: authorities}, nil
}
