package types

// DefaultGenesisState sets default evm genesis state with empty accounts and default params and
// chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Bonds:  []*Bond{},
	}
}

func NewGenesisState(params Params, bonds []*Bond) *GenesisState {
	return &GenesisState{
		Params: params,
		Bonds:  bonds,
	}
}
