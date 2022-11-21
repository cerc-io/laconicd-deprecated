package types

func NewGenesisState(params Params, records []Record, authorities []AuthorityEntry, names []NameEntry) GenesisState {
	return GenesisState{
		Params:      params,
		Records:     records,
		Authorities: authorities,
		Names:       names,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and default params and
// chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		Records:     []Record{},
		Authorities: []AuthorityEntry{},
		Names:       []NameEntry{},
	}
}

func ValidateGenesis(data GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
