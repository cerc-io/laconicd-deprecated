package keeper_test

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/cerc-io/laconicd/app"
	"github.com/cerc-io/laconicd/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
)

const testCommitHash = "71D8CF34026E32A3A34C2C2D4ADF25ABC8D7943A4619761BE27F196603D91B9D"

var (
	seed = int64(233)
)

func (suite *KeeperTestSuite) TestGrpcGetAllAuctions() {
	client, ctx, k := suite.queryClient, suite.ctx, suite.app.AuctionKeeper

	testCases := []struct {
		msg            string
		req            *types.AuctionsRequest
		createAuctions bool
		auctionCount   int
	}{
		{
			"fetch auctions when no auctions exist",
			&types.AuctionsRequest{},
			false,
			0,
		},

		{
			"fetch auctions with one auction created",
			&types.AuctionsRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctions {
				r := rand.New(rand.NewSource(seed))
				accs := app.RandomAccounts(r, 1)
				account := accs[0].Address
				err := testutil.FundAccount(suite.app.BankKeeper, ctx, account, sdk.NewCoins(
					sdk.Coin{Amount: sdk.NewInt(100), Denom: sdk.DefaultBondDenom},
				))

				_, err = k.CreateAuction(ctx, types.NewMsgCreateAuction(k.GetParams(ctx), account))
				suite.Require().NoError(err)
			}

			resp, _ := client.Auctions(context.Background(), test.req)
			suite.Require().Equal(test.auctionCount, len(resp.GetAuctions().Auctions))
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcQueryParams() {
	testCases := []struct {
		msg string
		req *types.QueryParamsRequest
	}{
		{
			"fetch params",
			&types.QueryParamsRequest{},
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			resp, err := suite.queryClient.QueryParams(context.Background(), test.req)
			suite.Require().Nil(err)
			suite.Require().Equal(*(resp.Params), types.DefaultParams())
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetAuction() {
	testCases := []struct {
		msg           string
		req           *types.AuctionRequest
		createAuction bool
	}{
		{
			"fetch auction with empty auction ID",
			&types.AuctionRequest{},
			false,
		},
		{
			"fetch auction with valid auction ID",
			&types.AuctionRequest{},
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				auction, _, err := suite.createAuctionAndCommitBid(false)
				suite.Require().NoError(err)
				test.req.Id = auction.Id
			}

			resp, err := suite.queryClient.GetAuction(context.Background(), test.req)
			if test.createAuction {
				suite.Require().Nil(err)
				suite.Require().NotNil(resp.GetAuction())
				suite.Require().Equal(test.req.Id, resp.GetAuction().Id)
			} else {
				suite.Require().NotNil(err)
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetBids() {
	testCases := []struct {
		msg           string
		req           *types.BidsRequest
		createAuction bool
		commitBid     bool
		bidCount      int
	}{
		{
			"fetch all bids when no auction exists",
			&types.BidsRequest{},
			false,
			false,
			0,
		},
		{
			"fetch all bids for valid auction but no added bids",
			&types.BidsRequest{},
			true,
			false,
			0,
		},
		{
			"fetch all bids for valid auction and valid bid",
			&types.BidsRequest{},
			true,
			true,
			1,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				auction, _, err := suite.createAuctionAndCommitBid(test.commitBid)
				suite.Require().NoError(err)
				test.req.AuctionId = auction.Id
			}

			resp, err := suite.queryClient.GetBids(context.Background(), test.req)
			if test.createAuction {
				suite.Require().Nil(err)
				suite.Require().Equal(test.bidCount, len(resp.GetBids()))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetBid() {
	testCases := []struct {
		msg                 string
		req                 *types.BidRequest
		createAuctionAndBid bool
	}{
		{
			"fetch bid when bid does not exist",
			&types.BidRequest{},
			false,
		},
		{
			"fetch bid when valid bid exists",
			&types.BidRequest{},
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndBid {
				auction, bid, err := suite.createAuctionAndCommitBid(test.createAuctionAndBid)
				suite.Require().NoError(err)
				test.req.AuctionId = auction.Id
				test.req.Bidder = bid.BidderAddress
			}

			resp, err := suite.queryClient.GetBid(context.Background(), test.req)
			if test.createAuctionAndBid {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp.Bid)
				suite.Require().Equal(test.req.Bidder, resp.Bid.BidderAddress)
			} else {
				suite.Require().NotNil(err)
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetAuctionsByBidder() {
	testCases := []struct {
		msg                       string
		req                       *types.AuctionsByBidderRequest
		createAuctionAndCommitBid bool
		auctionCount              int
	}{
		{
			"get auctions by bidder with invalid bidder address",
			&types.AuctionsByBidderRequest{},
			false,
			0,
		},
		{
			"get auctions by bidder with valid auction and bid",
			&types.AuctionsByBidderRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndCommitBid {
				_, bid, err := suite.createAuctionAndCommitBid(test.createAuctionAndCommitBid)
				suite.Require().NoError(err)
				test.req.BidderAddress = bid.BidderAddress
			}

			resp, err := suite.queryClient.AuctionsByBidder(context.Background(), test.req)
			if test.createAuctionAndCommitBid {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp.Auctions)
				suite.Require().Equal(test.auctionCount, len(resp.Auctions.Auctions))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetAuctionsByOwner() {
	testCases := []struct {
		msg           string
		req           *types.AuctionsByOwnerRequest
		createAuction bool
		auctionCount  int
	}{
		{
			"get auctions by owner with invalid owner address",
			&types.AuctionsByOwnerRequest{},
			false,
			0,
		},
		{
			"get auctions by owner with valid auction",
			&types.AuctionsByOwnerRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				auction, _, err := suite.createAuctionAndCommitBid(false)
				suite.Require().NoError(err)
				test.req.OwnerAddress = auction.OwnerAddress
			}

			resp, err := suite.queryClient.AuctionsByOwner(context.Background(), test.req)
			if test.createAuction {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp.Auctions)
				suite.Require().Equal(test.auctionCount, len(resp.Auctions.Auctions))
			} else {
				suite.Require().NotNil(err)
				suite.Require().Error(err)
			}
		})
	}
}

func (suite KeeperTestSuite) TestGrpcQueryBalance() {
	testCases := []struct {
		msg           string
		req           *types.BalanceRequest
		createAuction bool
		auctionCount  int
	}{
		{
			"get balance with no auctions created",
			&types.BalanceRequest{},
			false,
			0,
		},
		{
			"get balance with single auction created",
			&types.BalanceRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		if test.createAuction {
			_, _, err := suite.createAuctionAndCommitBid(true)
			suite.Require().NoError(err)
		}

		resp, err := suite.queryClient.Balance(context.Background(), test.req)
		suite.Require().NoError(err)
		suite.Require().Equal(test.auctionCount, len(resp.GetBalance()))
	}
}

func (suite *KeeperTestSuite) createAuctionAndCommitBid(commitBid bool) (*types.Auction, *types.Bid, error) {
	ctx, k := suite.ctx, suite.app.AuctionKeeper
	accCount := 1
	if commitBid {
		accCount++
	}

	r := rand.New(rand.NewSource(seed))
	accounts := app.RandomAccounts(r, accCount)
	for _, account := range accounts {
		err := testutil.FundAccount(suite.app.BankKeeper, ctx, account.Address, sdk.NewCoins(
			sdk.Coin{Amount: sdk.NewInt(100), Denom: sdk.DefaultBondDenom},
		))
		if err != nil {
			return nil, nil, err
		}
	}

	auction, err := k.CreateAuction(ctx, types.NewMsgCreateAuction(k.GetParams(ctx), accounts[0].Address))
	if err != nil {
		return nil, nil, err
	}

	if commitBid {
		bid, err := k.CommitBid(ctx, types.NewMsgCommitBid(auction.Id, testCommitHash, accounts[1].Address))
		if err != nil {
			return nil, nil, err
		}

		return auction, bid, nil
	}

	return auction, nil, nil
}
