package types

// bond module event types

const (
	EventTypeCreateBond   = "crate_bond"
	EventTypeRefillBond   = "refill_bond"
	EventTypeCancelBond   = "cancel_bond"
	EventTypeWithdrawBond = "withdraw_bond"

	AttributeKeySigner     = "signer"
	AttributeKeyAmount     = "amount"
	AttributeKeyBondID     = "bond_id"
	AttributeValueCategory = ModuleName
)
