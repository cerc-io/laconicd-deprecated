package keeper_test

import (
	"testing"

	"github.com/cerc-io/laconicd/app"
	bondkeeper "github.com/cerc-io/laconicd/x/bond/keeper"
	"github.com/cerc-io/laconicd/x/bond/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite
	app         *app.EthermintApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	testApp := app.Setup(false, func(ea *app.EthermintApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	querier := bondkeeper.Querier{Keeper: testApp.BondKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app, suite.ctx, suite.queryClient = testApp, ctx, queryClient
}

func TestParams(t *testing.T) {
	testApp := app.Setup(false, func(ea *app.EthermintApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	expParams := types.DefaultParams()
	params := testApp.BondKeeper.GetParams(ctx)
	require.Equal(t, expParams.MaxBondAmount, params.MaxBondAmount)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
