package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bond MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) SetRecord(c context.Context, msg *types.MsgSetRecord) (*types.MsgSetRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.ProcessSetRecord(ctx, types.MsgSetRecord{
		BondId:  msg.GetBondId(),
		Signer:  msg.GetSigner(),
		Payload: msg.GetPayload(),
	})
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetRecord,
			sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.GetBondId()),
			sdk.NewAttribute(types.AttributeKeyPayload, msg.Payload.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	return &types.MsgSetRecordResponse{}, nil
}

func (m msgServer) SetName(c context.Context, msg *types.MsgSetName) (*types.MsgSetNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessSetName(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetRecord,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyCRN, msg.Crn),
			sdk.NewAttribute(types.AttributeKeyCID, msg.Cid),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgSetNameResponse{}, nil
}

func (m msgServer) ReserveName(c context.Context, msg *types.MsgReserveAuthority) (*types.MsgReserveAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessReserveAuthority(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReserveNameAuthority,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgReserveAuthorityResponse{}, nil
}

func (m msgServer) SetAuthorityBond(c context.Context, msg *types.MsgSetAuthorityBond) (*types.MsgSetAuthorityBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessSetAuthorityBond(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAuthorityBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.BondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgSetAuthorityBondResponse{}, nil
}

func (m msgServer) DeleteName(c context.Context, msg *types.MsgDeleteNameAuthority) (*types.MsgDeleteNameAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessDeleteName(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeleteName,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyCRN, msg.Crn),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgDeleteNameAuthorityResponse{}, nil
}

func (m msgServer) RenewRecord(c context.Context, msg *types.MsgRenewRecord) (*types.MsgRenewRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	err = m.Keeper.ProcessRenewRecord(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRenewRecord,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyRecordId, msg.RecordId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgRenewRecordResponse{}, nil
}

func (m msgServer) AssociateBond(c context.Context, msg *types.MsgAssociateBond) (*types.MsgAssociateBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.ProcessAssociateBond(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAssociateBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyRecordId, msg.RecordId),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.BondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgAssociateBondResponse{}, nil
}

func (m msgServer) DissociateBond(c context.Context, msg *types.MsgDissociateBond) (*types.MsgDissociateBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessDissociateBond(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDissociateBond,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyRecordId, msg.RecordId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgDissociateBondResponse{}, nil
}

func (m msgServer) DissociateRecords(c context.Context, msg *types.MsgDissociateRecords) (*types.MsgDissociateRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessDissociateRecords(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDissociateRecords,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyBondId, msg.BondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgDissociateRecordsResponse{}, nil
}

func (m msgServer) ReAssociateRecords(c context.Context, msg *types.MsgReAssociateRecords) (*types.MsgReAssociateRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessReAssociateRecords(ctx, *msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReAssociateRecords,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyOldBondId, msg.OldBondId),
			sdk.NewAttribute(types.AttributeKeyNewBondId, msg.NewBondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgReAssociateRecordsResponse{}, nil
}
