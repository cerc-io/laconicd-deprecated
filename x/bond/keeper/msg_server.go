package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cerc-io/laconicd/utils"
	"github.com/cerc-io/laconicd/x/bond/types"
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
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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

	k.logTxGasConsumed(ctx, "CreateBond")

	return &types.MsgCreateBondResponse{}, nil
}

func (k msgServer) RefillBond(c context.Context, msg *types.MsgRefillBond) (*types.MsgRefillBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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
			sdk.NewAttribute(types.AttributeKeyBondID, msg.Id),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	k.logTxGasConsumed(ctx, "RefillBond")

	return &types.MsgRefillBondResponse{}, nil
}

func (k msgServer) WithdrawBond(c context.Context, msg *types.MsgWithdrawBond) (*types.MsgWithdrawBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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
			sdk.NewAttribute(types.AttributeKeyBondID, msg.Id),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	k.logTxGasConsumed(ctx, "WithdrawBond")

	return &types.MsgWithdrawBondResponse{}, nil
}

func (k msgServer) CancelBond(c context.Context, msg *types.MsgCancelBond) (*types.MsgCancelBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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
			sdk.NewAttribute(types.AttributeKeyBondID, msg.Id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	k.logTxGasConsumed(ctx, "CancelBond")

	return &types.MsgCancelBondResponse{}, nil
}

func (k msgServer) logTxGasConsumed(ctx sdk.Context, tx string) {
	gasConsumed := ctx.GasMeter().GasConsumed()
	k.Keeper.Logger(ctx).Info("tx executed", "method", tx, "gas_consumed", fmt.Sprintf("%d", gasConsumed))
}
