package bond_test

import (
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	app2 "github.com/tharsis/ethermint/app"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
	"testing"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := app2.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(abcitypes.RequestInitChain{
		AppStateBytes: []byte("{}"),
		ChainId:       "test-chain-id",
	})

	acc := app.AccountKeeper.GetModuleAccount(ctx, bondtypes.ModuleName)
	require.NotNil(t, acc)
}
