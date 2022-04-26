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

	// for acc := range maccPerms {
	// 	fmt.Println(acc)
	// 	require.True(
	// 		t,
	// 		app.BankKeeper.BlockedAddr(app.AccountKeeper.GetModuleAddress(acc)),
	// 		fmt.Sprintf("ensure that blocked addresses %s are properly set in bank keeper", acc),
	// 	)
	// }

	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewEthermintApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})
	_, err := app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
