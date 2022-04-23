package bond_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	app "github.com/tharsis/ethermint/app"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := app.Setup(false, func(ea *app.EthermintApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(abcitypes.RequestInitChain{
		AppStateBytes: []byte("{}"),
		ChainId:       "test-chain-id",
	})

	acc := app.AccountKeeper.GetModuleAccount(ctx, bondtypes.ModuleName)
	require.NotNil(t, acc)
}
