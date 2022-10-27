package types

// DefaultGenesisState sets default evm genesis state with empty accounts and default params and
// chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Auctions: []*Auction{},
	}
}

func NewGenesisState(params Params, auctions []*Auction) *GenesisState {
	return &GenesisState{
		Params:   params,
		Auctions: auctions,
	}
}
