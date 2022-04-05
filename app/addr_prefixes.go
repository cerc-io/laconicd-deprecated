package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "github.com/tharsis/ethermint/cmd/config"
)

// sdk config
func init() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()

	cmdcfg.RegisterDenoms()
}
