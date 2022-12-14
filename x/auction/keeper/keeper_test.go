package keeper_test

import (
	"testing"

	"github.com/cerc-io/laconicd/app"
	auctionkeeper "github.com/cerc-io/laconicd/x/auction/keeper"
	"github.com/cerc-io/laconicd/x/auction/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite
	app         *app.LaconicApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	testApp := app.Setup(false, func(ea *app.LaconicApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	querier := auctionkeeper.Querier{Keeper: testApp.AuctionKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app, suite.ctx, suite.queryClient = testApp, ctx, queryClient
}

func TestParams(t *testing.T) {
	testApp := app.Setup(false, func(ea *app.LaconicApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	expParams := types.DefaultParams()
	params := testApp.AuctionKeeper.GetParams(ctx)
	require.Equal(t, expParams.CommitsDuration, params.CommitsDuration)
	require.Equal(t, expParams.RevealsDuration, params.RevealsDuration)
	require.Equal(t, expParams.CommitFee, params.CommitFee)
	require.Equal(t, expParams.RevealFee, params.RevealFee)
	require.Equal(t, expParams.MinimumBid, params.MinimumBid)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
