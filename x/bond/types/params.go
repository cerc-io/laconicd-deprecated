package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = &Params{}

// Default parameter values.
var (
	DefaultMaxBondAmountTokens = sdk.NewInt(100000000000)
)

// Parameter keys
var (
	ParamStoreKeyMaxBondAmount = []byte("MaxBondAmount")
)

// ParamKeyTable ParamTable for staking module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(maxBondAmount sdk.Coin) Params {
	return Params{MaxBondAmount: maxBondAmount}
}

// DefaultParams returns default evm parameters
// ExtraEIPs is empty to prevent overriding the latest hard fork instruction set
func DefaultParams() Params {
	return NewParams(sdk.NewCoin(sdk.DefaultBondDenom, DefaultMaxBondAmountTokens))
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyMaxBondAmount, &p.MaxBondAmount, validateMaxBondAmount),
	}
}

func validateMaxBondAmount(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("max bond amount must be positive")
	}

	return nil
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateMaxBondAmount(p.MaxBondAmount); err != nil {
		return err
	}

	return nil
}
