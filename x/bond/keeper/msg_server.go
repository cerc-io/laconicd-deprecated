package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/bond/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bond MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateBond(c context.Context, msg *types.MsgCreateBond) (*types.MsgCreateBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	_, err = k.Keeper.CreateBond(ctx, signerAddress, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	return &types.MsgCreateBondResponse{}, nil
}

func (k msgServer) RefillBond(c context.Context, msg *types.MsgRefillBond) (*types.MsgRefillBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.RefillBond(ctx, msg.Id, signerAddress, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRefillBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.Id),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	return &types.MsgRefillBondResponse{}, nil
}

func (k msgServer) WithdrawBond(c context.Context, msg *types.MsgWithdrawBond) (*types.MsgWithdrawBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.WithdrawBond(ctx, msg.Id, signerAddress, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.Id),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	return &types.MsgWithdrawBondResponse{}, nil
}

func (k msgServer) CancelBond(c context.Context, msg *types.MsgCancelBond) (*types.MsgCancelBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	_, err = k.Keeper.CancelBond(ctx, msg.Id, signerAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.Id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	return &types.MsgCancelBondResponse{}, nil
}
