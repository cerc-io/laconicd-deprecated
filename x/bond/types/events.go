package types

// bond module event types

const (
	EventTypeCreateBond   = "crate_bond"
	EventTypeRefillBond   = "refill_bond"
	EventTypeCancelBond   = "cancel_bond"
	EventTypeWithdrawBond = "withdraw_bond"

	AttributeKeySigner     = "signer"
	AttributeKeyAmount     = "amount"
	AttributeKeyBondId     = "bond_id"
	AttributeValueCategory = ModuleName
)
