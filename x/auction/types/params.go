package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace.
const (
	DefaultParamspace = ModuleName
)

var (
	DefaultCommitsDuration = 5 * time.Minute
	DefaultRevealsDuration = 5 * time.Minute
	DefaultCommitFee       = sdk.Coin{Amount: sdk.NewInt(10), Denom: sdk.DefaultBondDenom}
	DefaultRevealFee       = sdk.Coin{Amount: sdk.NewInt(10), Denom: sdk.DefaultBondDenom}
	DefaultMinimumBid      = sdk.Coin{Amount: sdk.NewInt(1000), Denom: sdk.DefaultBondDenom}

	ParamStoreKeyCommitsDuration = []byte("CommitsDuration")
	ParamStoreKeyRevealsDuration = []byte("RevealsDuration")
	ParamStoreKeyCommitFee       = []byte("CommitFee")
	ParamStoreKeyRevealFee       = []byte("RevealFee")
	ParamStoreKeyMinimumBid      = []byte("MinimumBid")
)

var _ types.ParamSet = &Params{}

func NewParams(commitsDuration time.Duration, revealsDuration time.Duration, commitFee sdk.Coin, revealFee sdk.Coin, minimumBid sdk.Coin) Params {
	return Params{
		CommitsDuration: commitsDuration,
		RevealsDuration: revealsDuration,
		CommitFee:       commitFee,
		RevealFee:       revealFee,
		MinimumBid:      minimumBid,
	}
}

// ParamKeyTable - ParamTable for bond module.
func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		types.NewParamSetPair(ParamStoreKeyCommitsDuration, &p.CommitsDuration, validateCommitsDuration),
		types.NewParamSetPair(ParamStoreKeyRevealsDuration, &p.RevealsDuration, validateRevealsDuration),
		types.NewParamSetPair(ParamStoreKeyCommitFee, &p.CommitFee, validateCommitFee),
		types.NewParamSetPair(ParamStoreKeyRevealFee, &p.RevealFee, validateRevealFee),
		types.NewParamSetPair(ParamStoreKeyMinimumBid, &p.MinimumBid, validateMinimumBid),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		CommitsDuration: DefaultCommitsDuration,
		RevealsDuration: DefaultRevealsDuration,
		CommitFee:       DefaultCommitFee,
		RevealFee:       DefaultRevealFee,
		MinimumBid:      DefaultMinimumBid,
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("CommitsDuration: %s\n", p.CommitsDuration.String()))
	sb.WriteString(fmt.Sprintf("RevealsDuration: %s\n", p.RevealsDuration.String()))
	sb.WriteString(fmt.Sprintf("CommitFee: %s\n", p.CommitFee.String()))
	sb.WriteString(fmt.Sprintf("RevealFee: %s\n", p.RevealFee.String()))
	sb.WriteString(fmt.Sprintf("MinimumBid: %s\n", p.MinimumBid.String()))
	return sb.String()
}

func validateCommitsDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return errors.New("commits duration cannot be negative")
	}

	return nil
}

func validateRevealsDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return errors.New("commits duration cannot be negative")
	}

	return nil
}

func validateCommitFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("commit fee must be positive")
	}

	return nil
}

func validateRevealFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("reveal fee must be positive")
	}

	return nil
}

func validateMinimumBid(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("minimum bid must be positive")
	}

	return nil
}

// Validate a set of params.
func (p Params) Validate() error {
	if err := validateCommitsDuration(p.CommitsDuration); err != nil {
		return err
	}

	if err := validateRevealsDuration(p.RevealsDuration); err != nil {
		return err
	}

	if err := validateCommitFee(p.CommitFee); err != nil {
		return err
	}

	if err := validateRevealFee(p.RevealFee); err != nil {
		return err
	}

	if err := validateMinimumBid(p.MinimumBid); err != nil {
		return err
	}

	return nil
}
