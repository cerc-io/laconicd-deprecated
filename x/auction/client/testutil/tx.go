package testutil

import (
	"fmt"

	"github.com/cerc-io/laconicd/x/auction/client/cli"
	"github.com/cerc-io/laconicd/x/auction/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

const (
	sampleCommitTime     = "90s"
	sampleRevealTime     = "5s"
	placeholderAuctionID = "placeholder_auction_id"
)

func (suite *IntegrationTestSuite) TestTxCreateAuction() {
	sr := suite.Require()
	testCases := []struct {
		msg       string
		args      []string
		expectErr bool
	}{
		{
			"create auction with missing arguments",
			[]string{sampleCommitTime, sampleRevealTime},
			true,
		},
		{
			"create auction with proper arguments",
			[]string{
				sampleCommitTime, sampleRevealTime,
				fmt.Sprintf("10%s", suite.cfg.BondDenom),
				fmt.Sprintf("10%s", suite.cfg.BondDenom),
				fmt.Sprintf("100%s", suite.cfg.BondDenom),
			},
			false,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			resp, err := suite.executeTx(cli.GetCmdCreateAuction(), test.args, ownerAccount)
			if test.expectErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				sr.Zero(resp.Code)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestTxCommitBid() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg           string
		args          []string
		createAuction bool
	}{
		{
			"commit bid with missing args",
			[]string{fmt.Sprintf("200%s", suite.cfg.BondDenom)},
			false,
		},
		{
			"commit bid with valid args",
			[]string{
				placeholderAuctionID,
				fmt.Sprintf("200%s", suite.cfg.BondDenom),
			},
			true,
		},
	}

	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				auctionArgs := []string{
					sampleCommitTime, sampleRevealTime,
					fmt.Sprintf("10%s", suite.cfg.BondDenom),
					fmt.Sprintf("10%s", suite.cfg.BondDenom),
					fmt.Sprintf("100%s", suite.cfg.BondDenom),
				}
				_, err := suite.executeTx(cli.GetCmdCreateAuction(), auctionArgs, ownerAccount)
				sr.NoError(err)

				out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetCmdList(),
					[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)})
				sr.NoError(err)
				var queryResponse types.AuctionsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)
				sr.NotNil(queryResponse.GetAuctions())
				test.args[0] = queryResponse.GetAuctions().Auctions[0].Id
			}

			resp, err := suite.executeTx(cli.GetCmdCommitBid(), test.args, bidderAccount)
			if test.createAuction {
				sr.NoError(err)
				sr.Zero(resp.Code)
			} else {
				sr.Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) executeTx(cmd *cobra.Command, args []string, caller string) (sdk.TxResponse, error) {
	val := suite.network.Validators[0]
	additionalArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, caller),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", suite.cfg.BondDenom)),
	}
	args = append(args, additionalArgs...)

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var resp sdk.TxResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	err = suite.network.WaitForNextBlock()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}
