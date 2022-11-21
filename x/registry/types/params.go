package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values.
var (
	// DefaultRecordRent is the default record rent for 1 time period (see expiry time).
	DefaultRecordRent = sdk.NewInt(1000000)

	// DefaultRecordExpiryTime is the default record expiry time (1 year).
	DefaultRecordExpiryTime = time.Hour * 24 * 365

	DefaultAuthorityRent        = sdk.NewInt(1000000)
	DefaultAuthorityExpiryTime  = time.Hour * 24 * 365
	DefaultAuthorityGracePeriod = time.Hour * 24 * 2

	DefaultAuthorityAuctionEnabled = false
	DefaultCommitsDuration         = time.Hour * 24
	DefaultRevealsDuration         = time.Hour * 24
	DefaultCommitFee               = sdk.NewInt(1000000)
	DefaultRevealFee               = sdk.NewInt(1000000)
	DefaultMinimumBid              = sdk.NewInt(5000000)
)

// Keys for parameter access
var (
	KeyRecordRent         = []byte("RecordRent")
	KeyRecordRentDuration = []byte("RecordRentDuration")

	KeyAuthorityRent         = []byte("AuthorityRent")
	KeyAuthorityRentDuration = []byte("AuthorityRentDuration")
	KeyAuthorityGracePeriod  = []byte("AuthorityGracePeriod")

	KeyAuthorityAuctionEnabled = []byte("AuthorityAuctionEnabled")
	KeyCommitsDuration         = []byte("AuthorityAuctionCommitsDuration")
	KeyRevealsDuration         = []byte("AuthorityAuctionRevealsDuration")
	KeyCommitFee               = []byte("AuthorityAuctionCommitFee")
	KeyRevealFee               = []byte("AuthorityAuctionRevealFee")
	KeyMinimumBid              = []byte("AuthorityAuctionMinimumBid")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable ParamTable for staking module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRecordRent, &p.RecordRent, validateRecordRent),
		paramtypes.NewParamSetPair(KeyRecordRentDuration, &p.RecordRentDuration, validateRecordRentDuration),

		paramtypes.NewParamSetPair(KeyAuthorityRent, &p.AuthorityRent, validateAuthorityRent),
		paramtypes.NewParamSetPair(KeyAuthorityRentDuration, &p.AuthorityRentDuration, validateAuthorityRentDuration),
		paramtypes.NewParamSetPair(KeyAuthorityGracePeriod, &p.AuthorityGracePeriod, validateAuthorityGracePeriod),

		paramtypes.NewParamSetPair(KeyAuthorityAuctionEnabled, &p.AuthorityAuctionEnabled, validateAuthorityAuctionEnabled),
		paramtypes.NewParamSetPair(KeyCommitsDuration, &p.AuthorityAuctionCommitsDuration, validateCommitsDuration),
		paramtypes.NewParamSetPair(KeyRevealsDuration, &p.AuthorityAuctionRevealsDuration, validateRevealsDuration),
		paramtypes.NewParamSetPair(KeyCommitFee, &p.AuthorityAuctionCommitFee, validateCommitFee),
		paramtypes.NewParamSetPair(KeyRevealFee, &p.AuthorityAuctionRevealFee, validateRevealFee),
		paramtypes.NewParamSetPair(KeyMinimumBid, &p.AuthorityAuctionMinimumBid, validateMinimumBid),
	}
}

// NewParams creates a new Params instance
func NewParams(recordRent sdk.Coin, recordRentDuration time.Duration,
	authorityRent sdk.Coin, authorityRentDuration time.Duration, authorityGracePeriod time.Duration,
	authorityAuctionEnabled bool, commitsDuration time.Duration, revealsDuration time.Duration,
	commitFee sdk.Coin, revealFee sdk.Coin, minimumBid sdk.Coin,
) Params {
	return Params{
		RecordRent:         recordRent,
		RecordRentDuration: recordRentDuration,

		AuthorityRent:         authorityRent,
		AuthorityRentDuration: authorityRentDuration,
		AuthorityGracePeriod:  authorityGracePeriod,

		AuthorityAuctionEnabled:         authorityAuctionEnabled,
		AuthorityAuctionCommitsDuration: commitsDuration,
		AuthorityAuctionRevealsDuration: revealsDuration,
		AuthorityAuctionCommitFee:       commitFee,
		AuthorityAuctionRevealFee:       revealFee,
		AuthorityAuctionMinimumBid:      minimumBid,
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
	return NewParams(
		sdk.NewCoin(sdk.DefaultBondDenom, DefaultRecordRent), DefaultRecordExpiryTime,
		sdk.NewCoin(sdk.DefaultBondDenom, DefaultAuthorityRent),
		DefaultAuthorityExpiryTime, DefaultAuthorityGracePeriod, DefaultAuthorityAuctionEnabled, DefaultCommitsDuration,
		DefaultRevealsDuration,
		sdk.NewCoin(sdk.DefaultBondDenom, DefaultCommitFee),
		sdk.NewCoin(sdk.DefaultBondDenom, DefaultRevealFee),
		sdk.NewCoin(sdk.DefaultBondDenom, DefaultMinimumBid),
	)
}

func validateAmount(name string, i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return fmt.Errorf("%s can't be negative", name)
	}

	return nil
}

func validateDuration(name string, i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("%s invalid parameter type: %T", name, i)
	}

	if v <= 0 {
		return fmt.Errorf("%s must be a positive integer", name)
	}

	return nil
}

func validateRecordRent(i interface{}) error {
	return validateAmount("RecordRent", i)
}

func validateRecordRentDuration(i interface{}) error {
	return validateDuration("RecordRentDuration", i)
}

func validateAuthorityRent(i interface{}) error {
	return validateAmount("AuthorityRent", i)
}

func validateAuthorityRentDuration(i interface{}) error {
	return validateDuration("AuthorityRentDuration", i)
}

func validateAuthorityGracePeriod(i interface{}) error {
	return validateDuration("AuthorityGracePeriod", i)
}

func validateAuthorityAuctionEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("%s invalid parameter type: %T", "AuthorityAuctionEnabled", i)
	}

	return nil
}

func validateCommitsDuration(i interface{}) error {
	return validateDuration("AuthorityCommitsDuration", i)
}

func validateRevealsDuration(i interface{}) error {
	return validateDuration("AuthorityRevealsDuration", i)
}

func validateCommitFee(i interface{}) error {
	return validateAmount("AuthorityCommitFee", i)
}

func validateRevealFee(i interface{}) error {
	return validateAmount("AuthorityRevealFee", i)
}

func validateMinimumBid(i interface{}) error {
	return validateAmount("AuthorityMinimumBid", i)
}

// Validate a set of params.
func (p Params) Validate() error {
	if err := validateRecordRent(p.RecordRent); err != nil {
		return err
	}

	if err := validateRecordRentDuration(p.RecordRentDuration); err != nil {
		return err
	}

	if err := validateAuthorityRent(p.AuthorityRent); err != nil {
		return err
	}

	if err := validateAuthorityRentDuration(p.AuthorityRentDuration); err != nil {
		return err
	}

	if err := validateAuthorityGracePeriod(p.AuthorityGracePeriod); err != nil {
		return err
	}

	if err := validateAuthorityAuctionEnabled(p.AuthorityAuctionEnabled); err != nil {
		return err
	}

	if err := validateCommitsDuration(p.AuthorityAuctionCommitsDuration); err != nil {
		return err
	}

	if err := validateRevealsDuration(p.AuthorityAuctionRevealsDuration); err != nil {
		return err
	}

	if err := validateCommitFee(p.AuthorityAuctionCommitFee); err != nil {
		return err
	}

	if err := validateRevealFee(p.AuthorityAuctionRevealFee); err != nil {
		return err
	}

	if err := validateMinimumBid(p.AuthorityAuctionMinimumBid); err != nil {
		return err
	}

	return nil
}
