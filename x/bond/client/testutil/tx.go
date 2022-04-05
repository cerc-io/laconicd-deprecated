package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tharsis/ethermint/x/bond/client/cli"
	"github.com/tharsis/ethermint/x/bond/types"
)

func (s *IntegrationTestSuite) TestTxCreateBond() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			"without deposit",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
		},
		{
			"create bond",
			[]string{
				fmt.Sprintf("10%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			cmd := cli.NewCreateBondCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.Nil(err)
				sr.NoError(err)
				sr.Zero(d.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTxRefillBond() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name      string
		args      []string
		getBondId bool
		err       bool
	}{
		{
			"without refill amount and bond id",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			true,
		},
		{
			"refill bond",
			[]string{
				fmt.Sprintf("10%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			cmd := cli.RefillBondCmd()
			if tc.getBondId {
				cmd := cli.GetQueryBondLists()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				})
				sr.NoError(err)
				var queryResponse types.QueryGetBondsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				// extract bond id from bonds list
				bond := queryResponse.GetBonds()[0]
				tc.args = append([]string{bond.GetId()}, tc.args...)
			}
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				// checking the balance of bond
				cmd := cli.GetBondByIdCmd()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
					fmt.Sprintf(tc.args[0]),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				})
				sr.NoError(err)
				var queryResponse types.QueryGetBondByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				sr.True(queryResponse.GetBond().GetBalance().IsEqual(
					sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(20)))))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTxWithdrawAmountFromBond() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name      string
		args      []string
		getBondId bool
		err       bool
	}{
		{
			"without withdraw amount and bond id",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			true,
		},
		{
			"withdraw amount from bond",
			[]string{
				fmt.Sprintf("10%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			cmd := cli.WithdrawBondCmd()
			if tc.getBondId {
				cmd := cli.GetQueryBondLists()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				})
				sr.NoError(err)
				var queryResponse types.QueryGetBondsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				// extract bond id from bonds list
				bond := queryResponse.GetBonds()[0]
				tc.args = append([]string{bond.GetId()}, tc.args...)
			}
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				// checking the balance of bond
				cmd := cli.GetBondByIdCmd()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
					fmt.Sprintf(tc.args[0]),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				})
				sr.NoError(err)
				var queryResponse types.QueryGetBondByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				sr.True(queryResponse.GetBond().GetBalance().IsEqual(
					sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestTxCancelBond() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name      string
		args      []string
		getBondId bool
		err       bool
	}{
		{
			"without bond id",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			true,
		},
		{
			"cancel bond",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			cmd := cli.CancelBondCmd()
			if tc.getBondId {
				cmd := cli.GetQueryBondLists()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				})
				sr.NoError(err)
				var queryResponse types.QueryGetBondsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				// extract bond id from bonds list
				bond := queryResponse.GetBonds()[0]
				tc.args = append([]string{bond.GetId()}, tc.args...)
			}
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				// checking the bond exists or not after cancel
				cmd := cli.GetBondByIdCmd()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
					fmt.Sprintf(tc.args[0]),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				})
				sr.NoError(err)
				var queryResponse types.QueryGetBondByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				sr.Zero(queryResponse.GetBond().GetId())
			}
		})
	}
}
