package keeper_test

import (
	"math/rand"
	"testing"

	"github.com/cerc-io/laconicd/app"
	bondtypes "github.com/cerc-io/laconicd/x/bond/types"
	nameservicekeeper "github.com/cerc-io/laconicd/x/nameservice/keeper"
	"github.com/cerc-io/laconicd/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var seed = int64(233)

type KeeperTestSuite struct {
	suite.Suite
	app         *app.EthermintApp
	ctx         sdk.Context
	queryClient types.QueryClient
	accounts    []sdk.AccAddress
	bond        bondtypes.Bond
}

func (suite *KeeperTestSuite) SetupTest() {
	testApp := app.Setup(false, func(ea *app.EthermintApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	querier := nameservicekeeper.Querier{Keeper: testApp.NameServiceKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	r := rand.New(rand.NewSource(seed))
	accs := app.RandomAccounts(r, 1)
	suite.accounts = []sdk.AccAddress{accs[0].Address}
	account := suite.accounts[0]
	_ = testutil.FundAccount(testApp.BankKeeper, ctx, account, sdk.NewCoins(sdk.Coin{
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
	testApp := app.Setup(false, func(ea *app.EthermintApp, genesis simapp.GenesisState) simapp.GenesisState {
		return genesis
	})
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	expParams := types.DefaultParams()
	params := testApp.NameServiceKeeper.GetParams(ctx)
	require.True(t, params.Equal(expParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
