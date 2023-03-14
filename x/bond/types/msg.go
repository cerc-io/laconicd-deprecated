package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateBond{}
	_ sdk.Msg = &MsgRefillBond{}
	_ sdk.Msg = &MsgWithdrawBond{}
	_ sdk.Msg = &MsgCancelBond{}
)

// NewMsgCreateBond is the constructor function for MsgCreateBond.
func NewMsgCreateBond(coins sdk.Coins, signer sdk.AccAddress) MsgCreateBond {
	return MsgCreateBond{
		Coins:  coins,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgCreateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateBond) Type() string { return "create" }

func (msg MsgCreateBond) ValidateBasic() error {
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgCreateBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgCreateBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// NewMsgRefillBond is the constructor function for MsgRefillBond.
func NewMsgRefillBond(id string, amount sdk.Coin, signer sdk.AccAddress) MsgRefillBond {
	return MsgRefillBond{
		Id:     id,
		Coins:  sdk.NewCoins(amount),
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgRefillBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRefillBond) Type() string { return "refill" }

func (msg MsgRefillBond) ValidateBasic() error {
	if len(msg.Id) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, msg.Id)
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgRefillBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgRefillBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// NewMsgWithdrawBond is the constructor function for NewMsgWithdrawBond.
func NewMsgWithdrawBond(id string, amount sdk.Coin, signer sdk.AccAddress) MsgWithdrawBond {
	return MsgWithdrawBond{
		Id:     id,
		Coins:  sdk.NewCoins(amount),
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgWithdrawBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgWithdrawBond) Type() string { return "withdraw" }

func (msg MsgWithdrawBond) ValidateBasic() error {
	if len(msg.Id) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, msg.Id)
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgWithdrawBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgWithdrawBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// NewMsgCancelBond is the constructor function for CalcelBond.
func NewMsgCancelBond(id string, signer sdk.AccAddress) MsgCancelBond {
	return MsgCancelBond{
		Id:     id,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgCancelBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCancelBond) Type() string { return "cancel" }

func (msg MsgCancelBond) ValidateBasic() error {
	if len(msg.Id) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, msg.Id)
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	return nil
}

func (msg MsgCancelBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgCancelBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}
