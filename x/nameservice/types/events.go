package types

const (
	EventTypeSetRecord            = "set"
	EventTypeDeleteName           = "delete-name"
	EventTypeReserveNameAuthority = "reserve-authority"
	EventTypeAuthorityBond        = "authority-bond"
	EventTypeRenewRecord          = "renew-record"
	EventTypeAssociateBond        = "associate-bond"
	EventTypeDissociateBond       = "dissociate-bond"
	EventTypeDissociateRecords    = "dissociate-record"
	EventTypeReAssociateRecords   = "re-associate-records"

	AttributeKeySigner     = "signer"
	AttributeKeyOwner      = "owner"
	AttributeKeyBondId     = "bond-id"
	AttributeKeyPayload    = "payload"
	AttributeKeyOldBondId  = "old-bond-id"
	AttributeKeyNewBondId  = "new-bond-id"
	AttributeKeyCID        = "cid"
	AttributeKeyName       = "name"
	AttributeKeyWRN        = "wrn"
	AttributeKeyRecordId   = "record-id"
	AttributeValueCategory = ModuleName
)
