package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cerc-io/laconicd/utils"
	"github.com/cerc-io/laconicd/x/auction/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) CreateAuction(c context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := s.Keeper.CreateAuction(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAuction,
			sdk.NewAttribute(types.AttributeKeyCommitsDuration, msg.CommitsDuration.String()),
			sdk.NewAttribute(types.AttributeKeyCommitFee, msg.CommitFee.String()),
			sdk.NewAttribute(types.AttributeKeyRevealFee, msg.RevealFee.String()),
			sdk.NewAttribute(types.AttributeKeyMinimumBid, msg.MinimumBid.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, signerAddress.String()),
		),
	})

	s.logTxGasConsumed(ctx, "CreateAuction")

	return &types.MsgCreateAuctionResponse{Auction: resp}, nil
}

// CommitBid is the command for committing a bid
func (s msgServer) CommitBid(c context.Context, msg *types.MsgCommitBid) (*types.MsgCommitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := s.Keeper.CommitBid(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCommitBid,
			sdk.NewAttribute(types.AttributeKeyAuctionID, msg.AuctionId),
			sdk.NewAttribute(types.AttributeKeyCommitHash, msg.CommitHash),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, signerAddress.String()),
		),
	})

	s.logTxGasConsumed(ctx, "CommitBid")

	return &types.MsgCommitBidResponse{Bid: resp}, nil
}

// RevealBid is the command for revealing a bid
func (s msgServer) RevealBid(c context.Context, msg *types.MsgRevealBid) (*types.MsgRevealBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := s.Keeper.RevealBid(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevealBid,
			sdk.NewAttribute(types.AttributeKeyAuctionID, msg.AuctionId),
			sdk.NewAttribute(types.AttributeKeyReveal, msg.Reveal),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, signerAddress.String()),
		),
	})

	s.logTxGasConsumed(ctx, "RevealBid")

	return &types.MsgRevealBidResponse{Auction: resp}, nil
}

func (s msgServer) logTxGasConsumed(ctx sdk.Context, tx string) {
	gasConsumed := ctx.GasMeter().GasConsumed()
	s.Keeper.Logger(ctx).Info("tx executed", "method", tx, "gas_consumed", fmt.Sprintf("%d", gasConsumed))
}
