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
	AttributeKeyBondID     = "bond-id"
	AttributeKeyPayload    = "payload"
	AttributeKeyOldBondID  = "old-bond-id"
	AttributeKeyNewBondID  = "new-bond-id"
	AttributeKeyCID        = "cid"
	AttributeKeyName       = "name"
	AttributeKeyCRN        = "crn"
	AttributeKeyRecordID   = "record-id"
	AttributeValueCategory = ModuleName
)
