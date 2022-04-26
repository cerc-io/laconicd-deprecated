package app

import (
	"testing"

	memdb "github.com/cosmos/cosmos-sdk/db/memdb"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/tharsis/ethermint/encoding"
)

func TestEthermintAppExport(t *testing.T) {
	encCfg := encoding.MakeConfig(ModuleBasics)
	db := memdb.NewDB()
	logger, _ := log.NewDefaultLogger("plain", "info", false)
	app := NewTestAppWithCustomOptions(t, false, SetupOptions{
		Logger:             logger,
		DB:                 db,
		InvCheckPeriod:     0,
		EncConfig:          encCfg,
		HomePath:           DefaultNodeHome,
		SkipUpgradeHeights: map[int64]bool{},
		AppOpts:            EmptyAppOptions{},
	})

	for acc := range allowedReceivingModAcc {
		// check module account is not blocked in bank
		require.False(
			t,
			app.BankKeeper.BlockedAddr(app.AccountKeeper.GetModuleAddress(acc)),
			"ensure that blocked addresses %s are properly set in bank keeper",
		)
	}

	app.Commit()
	logger2, _ := log.NewDefaultLogger("plain", "info", false)

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewEthermintApp(logger2, db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})
	require.NoError(t, app2.Init())
	_, err := app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
