package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cerc-io/laconicd/x/auction/types"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Auctions queries all auctions
func (q Querier) Auctions(c context.Context, req *types.AuctionsRequest) (*types.AuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	resp := q.Keeper.ListAuctions(ctx)
	return &types.AuctionsResponse{Auctions: &types.Auctions{Auctions: resp}}, nil
}

// GetAuction queries an auction
func (q Querier) GetAuction(c context.Context, req *types.AuctionRequest) (*types.AuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req.Id == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "auction ID is required")
	}

	resp := q.Keeper.GetAuction(ctx, req.Id)
	return &types.AuctionResponse{Auction: resp}, nil
}

// GetBid queries and auction bid
func (q Querier) GetBid(c context.Context, req *types.BidRequest) (*types.BidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req.AuctionId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "auction ID is required")
	}
	if req.Bidder == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bidder address is required")
	}
	resp := q.Keeper.GetBid(ctx, req.AuctionId, req.Bidder)
	return &types.BidResponse{Bid: &resp}, nil
}

// GetBids queries all auction bids
func (q Querier) GetBids(c context.Context, req *types.BidsRequest) (*types.BidsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req.AuctionId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "auction ID is required")
	}
	resp := q.Keeper.GetBids(ctx, req.AuctionId)
	return &types.BidsResponse{Bids: resp}, nil
}

// AuctionsByBidder queries auctions by bidder
func (q Querier) AuctionsByBidder(c context.Context, req *types.AuctionsByBidderRequest) (*types.AuctionsByBidderResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req.BidderAddress == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bidder address is required")
	}
	resp := q.Keeper.QueryAuctionsByBidder(ctx, req.BidderAddress)
	return &types.AuctionsByBidderResponse{Auctions: &types.Auctions{Auctions: resp}}, nil
}

// AuctionsByOwner queries auctions by owner
func (q Querier) AuctionsByOwner(c context.Context, req *types.AuctionsByOwnerRequest) (*types.AuctionsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req.OwnerAddress == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "owner address is required")
	}
	resp := q.Keeper.QueryAuctionsByOwner(ctx, req.OwnerAddress)
	return &types.AuctionsByOwnerResponse{Auctions: &types.Auctions{Auctions: resp}}, nil
}

// QueryParams implements the params query command
func (q Querier) QueryParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	resp := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: &resp}, nil
}

// Balance queries the auction module account balance
func (q Querier) Balance(c context.Context, req *types.BalanceRequest) (*types.BalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	resp := q.Keeper.GetAuctionModuleBalances(ctx)
	return &types.BalanceResponse{Balance: resp}, nil
}
