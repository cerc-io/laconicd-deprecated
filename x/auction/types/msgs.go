package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateAuction{}
	_ sdk.Msg = &MsgCommitBid{}
	_ sdk.Msg = &MsgRevealBid{}
)

// NewMsgCreateAuction is the constructor function for MsgCreateAuction.
func NewMsgCreateAuction(params Params, signer sdk.AccAddress) MsgCreateAuction {
	return MsgCreateAuction{
		CommitsDuration: params.CommitsDuration,
		RevealsDuration: params.RevealsDuration,
		CommitFee:       params.CommitFee,
		RevealFee:       params.RevealFee,
		MinimumBid:      params.MinimumBid,
		Signer:          signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgCreateAuction) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateAuction) Type() string { return "create" }

// ValidateBasic Implements Msg.
func (msg MsgCreateAuction) ValidateBasic() error {
	if msg.Signer == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	if msg.CommitsDuration <= 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "commit phase duration invalid.")
	}

	if msg.RevealsDuration <= 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "reveal phase duration invalid.")
	}

	if !msg.MinimumBid.IsPositive() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "minimum bid should be greater than zero.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateAuction) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgCommitBid is the constructor function for MsgCommitBid.
func NewMsgCommitBid(auctionID string, commitHash string, signer sdk.AccAddress) MsgCommitBid {
	return MsgCommitBid{
		AuctionId:  auctionID,
		CommitHash: commitHash,
		Signer:     signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgCommitBid) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCommitBid) Type() string { return "commit" }

// ValidateBasic Implements Msg.
func (msg MsgCommitBid) ValidateBasic() error {
	if msg.Signer == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer address.")
	}

	if msg.AuctionId == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid auction ID.")
	}

	if msg.CommitHash == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid commit hash.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCommitBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgCommitBid) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgRevealBid is the constructor function for MsgRevealBid.
func NewMsgRevealBid(auctionID string, reveal string, signer sdk.AccAddress) MsgRevealBid {
	return MsgRevealBid{
		AuctionId: auctionID,
		Reveal:    reveal,
		Signer:    signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgRevealBid) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRevealBid) Type() string { return "reveal" }

// ValidateBasic Implements Msg.
func (msg MsgRevealBid) ValidateBasic() error {
	if msg.Signer == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer address.")
	}

	if msg.AuctionId == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid auction ID.")
	}

	if msg.Reveal == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid reveal data.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRevealBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRevealBid) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}
