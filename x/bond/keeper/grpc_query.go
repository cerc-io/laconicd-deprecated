package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tharsis/ethermint/x/bond/types"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (q Querier) Bonds(c context.Context, _ *types.QueryGetBondsRequest) (*types.QueryGetBondsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	resp := q.Keeper.ListBonds(ctx)
	return &types.QueryGetBondsResponse{Bonds: resp}, nil
}

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: &params}, nil
}

func (q Querier) GetBondById(c context.Context, req *types.QueryGetBondByIdRequest) (*types.QueryGetBondByIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bondId := req.GetId()
	if len(bondId) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bond id required")
	}
	bond := q.Keeper.GetBond(ctx, req.GetId())
	return &types.QueryGetBondByIdResponse{Bond: &bond}, nil
}

func (q Querier) GetBondsByOwner(c context.Context, req *types.QueryGetBondsByOwnerRequest) (*types.QueryGetBondsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	owner := req.GetOwner()
	if len(owner) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "owner id required")
	}
	bonds := q.Keeper.QueryBondsByOwner(ctx, owner)
	return &types.QueryGetBondsByOwnerResponse{Bonds: bonds}, nil
}

func (q Querier) GetBondsModuleBalance(c context.Context, _ *types.QueryGetBondModuleBalanceRequest) (*types.QueryGetBondModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balance := q.Keeper.GetBondModuleBalances(ctx)
	return &types.QueryGetBondModuleBalanceResponse{Balance: balance}, nil
}
