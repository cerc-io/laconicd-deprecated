package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/cerc-io/laconicd/x/bond/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (q Querier) GetBondByID(c context.Context, req *types.QueryGetBondByIDRequest) (*types.QueryGetBondByIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bondID := req.GetId()
	if len(bondID) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bond id required")
	}
	bond := q.Keeper.GetBond(ctx, req.GetId())
	return &types.QueryGetBondByIDResponse{Bond: &bond}, nil
}

func (q Querier) GetBondsByOwner(c context.Context, req *types.QueryGetBondsByOwnerRequest) (*types.QueryGetBondsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	owner := req.GetOwner()
	if len(owner) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "owner id required")
	}
	bonds := q.Keeper.QueryBondsByOwner(ctx, owner)
	return &types.QueryGetBondsByOwnerResponse{Bonds: bonds}, nil
}

func (q Querier) GetBondsModuleBalance(c context.Context,
	_ *types.QueryGetBondModuleBalanceRequest,
) (*types.QueryGetBondModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balance := q.Keeper.GetBondModuleBalances(ctx)
	return &types.QueryGetBondModuleBalanceResponse{Balance: balance}, nil
}
