package keeper_test

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/cerc-io/laconicd/app"
	"github.com/cerc-io/laconicd/x/bond/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
)

var seed = int64(233)

func (suite *KeeperTestSuite) TestGrpcQueryBondsList() {
	grpcClient, ctx, k := suite.queryClient, suite.ctx, suite.app.BondKeeper

	testCases := []struct {
		msg         string
		req         *types.QueryGetBondsRequest
		resp        *types.QueryGetBondsResponse
		noOfBonds   int
		createBonds bool
	}{
		{
			"empty request",
			&types.QueryGetBondsRequest{},
			&types.QueryGetBondsResponse{},
			0,
			false,
		},
		{
			"Get Bonds",
			&types.QueryGetBondsRequest{},
			&types.QueryGetBondsResponse{},
			1,
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createBonds {
				r := rand.New(rand.NewSource(seed))
				accs := app.RandomAccounts(r, 1)
				account := accs[0].Address
				err := testutil.FundAccount(suite.app.BankKeeper, ctx, account, sdk.NewCoins(sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(1000),
				}))
				suite.Require().NoError(err)
				_, err = k.CreateBond(ctx, account, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))))
				suite.Require().NoError(err)
			}
			resp, _ := grpcClient.Bonds(context.Background(), test.req)
			suite.Require().Equal(test.noOfBonds, len(resp.GetBonds()))
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcQueryParams() {
	grpcClient := suite.queryClient

	testCases := []struct {
		msg string
		req *types.QueryParamsRequest
	}{
		{
			"Get Params",
			&types.QueryParamsRequest{},
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			resp, _ := grpcClient.Params(context.Background(), test.req)
			suite.Require().Equal(resp.GetParams().MaxBondAmount, types.DefaultParams().MaxBondAmount)
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcQueryBondBondId() {
	grpcClient, ctx, k, suiteRequire := suite.queryClient, suite.ctx, suite.app.BondKeeper, suite.Require()

	testCases := []struct {
		msg         string
		req         *types.QueryGetBondByIDRequest
		createBonds bool
		errResponse bool
		bondId      string
	}{
		{
			"empty request",
			&types.QueryGetBondByIDRequest{},
			false,
			true,
			"",
		},
		{
			"Get Bond By ID",
			&types.QueryGetBondByIDRequest{},
			true,
			false,
			"",
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createBonds {
				r := rand.New(rand.NewSource(seed))
				accs := app.RandomAccounts(r, 1)
				account := accs[0].Address
				err := testutil.FundAccount(suite.app.BankKeeper, ctx, account, sdk.NewCoins(sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(1000),
				}))
				suite.Require().NoError(err)
				bond, err := k.CreateBond(ctx, account, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))))
				suiteRequire.NoError(err)
				test.req.Id = bond.Id
			}
			resp, err := grpcClient.GetBondByID(context.Background(), test.req)
			if !test.errResponse {
				suiteRequire.Nil(err)
				suiteRequire.NotNil(resp.GetBond())
				suiteRequire.Equal(test.req.Id, resp.GetBond().GetID())
			} else {
				suiteRequire.NotNil(err)
				suiteRequire.Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetBondsByOwner() {
	grpcClient, ctx, k, suiteRequire := suite.queryClient, suite.ctx, suite.app.BondKeeper, suite.Require()

	testCases := []struct {
		msg         string
		req         *types.QueryGetBondsByOwnerRequest
		noOfBonds   int
		createBonds bool
		errResponse bool
		bondId      string
	}{
		{
			"empty request",
			&types.QueryGetBondsByOwnerRequest{},
			0,
			false,
			true,
			"",
		},
		{
			"Get Bond By Owner",
			&types.QueryGetBondsByOwnerRequest{},
			1,
			true,
			false,
			"",
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createBonds {
				r := rand.New(rand.NewSource(seed))
				accs := app.RandomAccounts(r, 1)
				account := accs[0].Address
				_ = testutil.FundAccount(suite.app.BankKeeper, ctx, account, sdk.NewCoins(sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(1000),
				}))
				_, err := k.CreateBond(ctx, account, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))))
				suiteRequire.NoError(err)
				test.req.Owner = account.String()
			}
			resp, err := grpcClient.GetBondsByOwner(context.Background(), test.req)
			if !test.errResponse {
				suiteRequire.Nil(err)
				suiteRequire.NotNil(resp.GetBonds())
				suiteRequire.Equal(test.noOfBonds, len(resp.GetBonds()))
			} else {
				suiteRequire.NotNil(err)
				suiteRequire.Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetModuleBalance() {
	grpcClient, ctx, k, suiteRequire := suite.queryClient, suite.ctx, suite.app.BondKeeper, suite.Require()

	testCases := []struct {
		msg         string
		req         *types.QueryGetBondModuleBalanceRequest
		noOfBonds   int
		createBonds bool
		errResponse bool
	}{
		{
			"empty request",
			&types.QueryGetBondModuleBalanceRequest{},
			0,
			true,
			false,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createBonds {
				r := rand.New(rand.NewSource(seed))
				accs := app.RandomAccounts(r, 1)
				account := accs[0].Address
				_ = testutil.FundAccount(suite.app.BankKeeper, ctx, account, sdk.NewCoins(sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(1000),
				}))
				_, err := k.CreateBond(ctx, account, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))))
				suiteRequire.NoError(err)
			}
			resp, err := grpcClient.GetBondsModuleBalance(context.Background(), test.req)
			if !test.errResponse {
				suiteRequire.Nil(err)
				suiteRequire.NotNil(resp.GetBalance())
				suiteRequire.Equal(resp.GetBalance(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))))
			} else {
				suiteRequire.NotNil(err)
				suiteRequire.Error(err)
			}
		})
	}
}
