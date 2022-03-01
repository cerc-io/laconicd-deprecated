package app

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/db/memdb"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tharsis/ethermint/encoding"
)

func TestEthermintAppExport(t *testing.T) {
	db := memdb.NewDB()
	logger, err := log.NewDefaultLogger("plain", "info", false)
	if err != nil {
		t.Fatal(err)
	}
	app := NewEthermintApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})

	genesisState := NewDefaultGenesisState()
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			ChainId:       "ethermint_9000-1",
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	logger, err = log.NewDefaultLogger("plain", "info", false)
	if err != nil {
		t.Fatal(err)
	}
	app2 := NewEthermintApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	_, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
