package types

import (
	"cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSetRecord{}
	_ sdk.Msg = &MsgRenewRecord{}
	_ sdk.Msg = &MsgAssociateBond{}
	_ sdk.Msg = &MsgDissociateBond{}
	_ sdk.Msg = &MsgDissociateRecords{}
	_ sdk.Msg = &MsgReAssociateRecords{}

	_ cdctypes.UnpackInterfacesMessage = &MsgSetRecord{}
)

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload Payload, bondID string, signer sdk.AccAddress) MsgSetRecord {
	return MsgSetRecord{
		Payload: payload,
		BondId:  bondID,
		Signer:  signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgSetRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetRecord) Type() string { return "set-record" }

func (msg MsgSetRecord) ValidateBasic() error {
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	owners := msg.Payload.Record.Owners
	for _, owner := range owners {
		if owner == "" {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "record owner not set")
		}
	}

	if len(msg.BondId) == 0 {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "bond ID is required")
	}
	return nil
}

func (msg MsgSetRecord) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgSetRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgSetRecord) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	var attr Attributes
	return unpacker.UnpackAny(msg.Payload.Record.Attributes, &attr)
}

// NewMsgRenewRecord is the constructor function for MsgRenewRecord.
func NewMsgRenewRecord(recordID string, signer sdk.AccAddress) MsgRenewRecord {
	return MsgRenewRecord{
		RecordId: recordID,
		Signer:   signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgRenewRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRenewRecord) Type() string { return "renew-record" }

// ValidateBasic Implements Msg.
func (msg MsgRenewRecord) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "record id is required")
	}

	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgRenewRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgRenewRecord) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgAssociateBond is the constructor function for MsgAssociateBond.
func NewMsgAssociateBond(recordID, bondID string, signer sdk.AccAddress) MsgAssociateBond {
	return MsgAssociateBond{
		BondId:   bondID,
		RecordId: recordID,
		Signer:   signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgAssociateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAssociateBond) Type() string { return "associate-bond" }

// ValidateBasic Implements Msg.
func (msg MsgAssociateBond) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "record id is required")
	}
	if len(msg.BondId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "bond id is required")
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgAssociateBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgAssociateBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgDissociateBond is the constructor function for MsgDissociateBond.
func NewMsgDissociateBond(recordID string, signer sdk.AccAddress) MsgDissociateBond {
	return MsgDissociateBond{
		RecordId: recordID,
		Signer:   signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgDissociateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDissociateBond) Type() string { return "dissociate-bond" }

// ValidateBasic Implements Msg.
func (msg MsgDissociateBond) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "record id is required")
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgDissociateBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDissociateBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgDissociateRecords is the constructor function for MsgDissociateRecords.
func NewMsgDissociateRecords(bondID string, signer sdk.AccAddress) MsgDissociateRecords {
	return MsgDissociateRecords{
		BondId: bondID,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgDissociateRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDissociateRecords) Type() string { return "dissociate-records" }

// ValidateBasic Implements Msg.
func (msg MsgDissociateRecords) ValidateBasic() error {
	if len(msg.BondId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "bond id is required")
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgDissociateRecords) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDissociateRecords) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgReAssociateRecords is the constructor function for MsgReAssociateRecords.
func NewMsgReAssociateRecords(oldBondID, newBondID string, signer sdk.AccAddress) MsgReAssociateRecords {
	return MsgReAssociateRecords{
		OldBondId: oldBondID,
		NewBondId: newBondID,
		Signer:    signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgReAssociateRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgReAssociateRecords) Type() string { return "reassociate-records" }

// ValidateBasic Implements Msg.
func (msg MsgReAssociateRecords) ValidateBasic() error {
	if len(msg.OldBondId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "old-bond-id is required")
	}
	if len(msg.NewBondId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "new-bond-id is required")
	}
	if len(msg.Signer) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for Msg
func (msg MsgReAssociateRecords) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgReAssociateRecords) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}
