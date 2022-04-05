package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/bond interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(&MsgReserveAuthority{}, "nameservice/ReserveAuthority", nil)
	cdc.RegisterConcrete(&MsgDeleteNameAuthority{}, "nameservice/DeleteAuthority", nil)
	cdc.RegisterConcrete(&MsgSetAuthorityBond{}, "nameservice/SetAuthorityBond", nil)

	cdc.RegisterConcrete(&MsgSetRecord{}, "nameservice/SetRecord", nil)
	cdc.RegisterConcrete(&MsgRenewRecord{}, "nameservice/RenewRecord", nil)
	cdc.RegisterConcrete(&MsgAssociateBond{}, "nameservice/AssociateBond", nil)
	cdc.RegisterConcrete(&MsgDissociateBond{}, "nameservice/DissociateBond", nil)
	cdc.RegisterConcrete(&MsgDissociateRecords{}, "nameservice/DissociateRecords", nil)
	cdc.RegisterConcrete(&MsgReAssociateRecords{}, "nameservice/ReassociateRecords", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetName{},
		&MsgReserveAuthority{},
		&MsgDeleteNameAuthority{},
		&MsgSetAuthorityBond{},

		&MsgSetRecord{},
		&MsgRenewRecord{},
		&MsgAssociateBond{},
		&MsgDissociateBond{},
		&MsgDissociateRecords{},
		&MsgReAssociateRecords{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
