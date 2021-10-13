package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tharsis/ethermint/app"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
	nameservicekeeper "github.com/tharsis/ethermint/x/nameservice/keeper"
	"github.com/tharsis/ethermint/x/nameservice/types"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite
	app         *app.EthermintApp
	ctx         sdk.Context
	queryClient types.QueryClient
	accounts    []sdk.AccAddress
	bond        bondtypes.Bond
}

func (suite *KeeperTestSuite) SetupTest() {
	testApp := app.Setup(false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	querier := nameservicekeeper.Querier{Keeper: testApp.NameServiceKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.accounts = app.CreateRandomAccounts(1)
	account := suite.accounts[0]
	_ = simapp.FundAccount(testApp.BankKeeper, ctx, account, sdk.NewCoins(sdk.Coin{
		Denom:  sdk.DefaultBondDenom,
		Amount: sdk.NewInt(100000000000),
	}))

	bond, err := testApp.BondKeeper.CreateBond(ctx, account, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000))))
	if err != nil {
		return
	}
	suite.bond = *bond
	suite.app, suite.ctx, suite.queryClient = testApp, ctx, queryClient
}

func TestParams(t *testing.T) {
	testApp := app.Setup(false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	expParams := types.DefaultParams()
	params := testApp.NameServiceKeeper.GetParams(ctx)
	require.True(t, params.Equal(expParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
