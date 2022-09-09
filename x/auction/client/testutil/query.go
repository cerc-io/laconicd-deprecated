package testutil

import (
	"fmt"

	"github.com/cerc-io/laconicd/x/auction/client/cli"
	"github.com/cerc-io/laconicd/x/auction/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

var queryJSONFlag = []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

func (suite *IntegrationTestSuite) TestGetCmdQueryParams() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	suite.Run(fmt.Sprintf("Case %s", "fetch query params"), func() {
		out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdQueryParams(), queryJSONFlag)
		sr.NoError(err)
		var params types.QueryParamsResponse
		err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &params)
		sr.NoError(err)
		sr.Equal(types.DefaultParams(), *params.Params)
	})
}

func (suite *IntegrationTestSuite) TestGetCmdBalance() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	testCases := []struct {
		msg                 string
		createAuctionAndBid bool
	}{
		{
			"fetch module balance without creating auction and bid",
			false,
		},
		{
			"fetch module balance with valid auction and bid",
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndBid {
				suite.createAuctionAndBid(false, true)
			}

			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdBalance(), queryJSONFlag)
			sr.NoError(err)
			var balance types.BalanceResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &balance)
			sr.NoError(err)
			if test.createAuctionAndBid {
				sr.NotZero(len(balance.Balance))
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetCmdList() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	testCases := []struct {
		msg           string
		createAuction bool
	}{
		{
			"list auctions when no auctions exist",
			false,
		},
		{
			"list auctions after creating an auction",
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdList(), queryJSONFlag)
			sr.NoError(err)
			var auctions types.AuctionsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &auctions)
			sr.NoError(err)
			if test.createAuction {
				sr.NotZero(len(auctions.Auctions.Auctions))
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetCmdGetBid() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	testCases := []struct {
		msg                 string
		args                []string
		createAuctionAndBid bool
	}{
		{
			"get bid without creating auction",
			[]string{},
			false,
		},
		{
			"get bid after creating auction and bid",
			[]string{},
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndBid {
				auctionID := suite.createAuctionAndBid(false, true)
				test.args = append(test.args, auctionID)
				getBidsArgs := []string{auctionID, queryJSONFlag[0]}
				out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdGetBids(), getBidsArgs)
				sr.NoError(err)
				var bids types.BidsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &bids)
				sr.NoError(err)
				test.args = append(test.args, bids.GetBids()[0].BidderAddress)
			}

			test.args = append(test.args, queryJSONFlag...)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdGetBid(), test.args)
			if test.createAuctionAndBid {
				sr.NoError(err)
				var bid types.BidResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &bid)
				sr.NoError(err)
				sr.NotNil(bid.GetBid())
			} else {
				sr.Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetCmdGetBids() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	testCases := []struct {
		msg                 string
		args                []string
		createAuctionAndBid bool
	}{
		{
			"get bids without creating auction",
			[]string{},
			false,
		},
		{
			"get bids after creating auction and bid",
			[]string{},
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndBid {
				auctionID := suite.createAuctionAndBid(false, true)
				test.args = append(test.args, auctionID)
			}

			test.args = append(test.args, queryJSONFlag...)
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdGetBids(), test.args)
			if test.createAuctionAndBid {
				sr.NoError(err)
				var bids types.BidsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &bids)
				sr.NoError(err)
				sr.NotZero(len(bids.Bids))
			} else {
				sr.Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetCmdGetAuction() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	testCases := []struct {
		msg           string
		auctionID     string
		createAuction bool
	}{
		{
			"get auction with empty auction ID",
			"",
			false,
		},
		{
			"get auction with valid auction ID",
			"",
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				test.auctionID = suite.defaultAuctionID
			}

			args := []string{test.auctionID, queryJSONFlag[0]}
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdGetAuction(), args)
			if test.createAuction {
				sr.NoError(err)
				var auction types.AuctionResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &auction)
				sr.NoError(err)
				sr.NotNil(auction.GetAuction())
				sr.Equal(test.auctionID, auction.GetAuction().Id)
			} else {
				sr.Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetCmdAuctionsByBidder() {
	val := suite.network.Validators[0]
	sr := suite.Require()

	testCases := []struct {
		msg                 string
		createAuctionAndBid bool
		bidderAddress       string
	}{
		{
			"get auctions by bidder without creating auctions",
			false,
			"",
		},
		{
			"get auctions by bidder for valid bidder address",
			true,
			"",
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndBid {
				auctionID := suite.createAuctionAndBid(false, true)
				args := []string{auctionID, queryJSONFlag[0]}
				out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdGetBids(), args)
				sr.NoError(err)
				var bids types.BidsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &bids)
				sr.NoError(err)
				test.bidderAddress = bids.Bids[0].BidderAddress
			}

			getByBidderArgs := []string{test.bidderAddress, queryJSONFlag[0]}
			_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdAuctionsByBidder(), getByBidderArgs)
			if test.createAuctionAndBid {
				sr.NoError(err)
			} else {
				sr.Error(err)
			}
		})
	}
}

func (suite IntegrationTestSuite) createAuctionAndBid(createAuction, createBid bool) string {
	val := suite.network.Validators[0]
	sr := suite.Require()
	auctionID := ""

	if createAuction {
		auctionArgs := []string{
			sampleCommitTime, sampleRevealTime,
			fmt.Sprintf("10%s", suite.cfg.BondDenom),
			fmt.Sprintf("10%s", suite.cfg.BondDenom),
			fmt.Sprintf("100%s", suite.cfg.BondDenom),
		}

		resp, err := suite.executeTx(cli.GetCmdCreateAuction(), auctionArgs, ownerAccount)
		sr.NoError(err)
		sr.Zero(resp.Code)
		out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdList(), queryJSONFlag)
		sr.NoError(err)
		var queryResponse types.AuctionsResponse
		err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
		sr.NoError(err)
		auctionID = queryResponse.Auctions.Auctions[0].Id
	} else {
		auctionID = suite.defaultAuctionID
	}

	if createBid {
		bidArgs := []string{auctionID, fmt.Sprintf("200%s", suite.cfg.BondDenom)}
		resp, err := suite.executeTx(cli.GetCmdCommitBid(), bidArgs, bidderAccount)
		sr.NoError(err)
		sr.Zero(resp.Code)
	}

	return auctionID
}
