package testutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tharsis/ethermint/testutil/network"
	"github.com/tharsis/ethermint/x/bond/client/cli"
	"github.com/tharsis/ethermint/x/bond/types"
	"gopkg.in/yaml.v2"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

var (
	accountName    = "accountName"
	accountAddress string
)

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// setting up random account
	s.createAccountWithBalance(accountName)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGetCmdQueryParams() {
	val := s.network.Validators[0]

	testCases := []struct {
		name       string
		args       []string
		outputType string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"json",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.GetQueryParamsCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			var param types.QueryParamsResponse
			if tc.outputType == "json" {
				err := clientCtx.Codec.UnmarshalJSON(out.Bytes(), &param)
				s.Require().NoError(err)
			} else {
				err := yaml.Unmarshal(out.Bytes(), &param)
				s.Require().NoError(err)
			}
			s.Require().Equal(param.Params.MaxBondAmount, types.DefaultParams().MaxBondAmount)
		})
	}
}

func (s *IntegrationTestSuite) createAccountWithBalance(accountName string) {
	val := s.network.Validators[0]
	sr := s.Require()
	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := s.cfg.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	sr.NoError(err)
	sr.NotNil(consPubKeyBz)

	info, _, err := val.ClientCtx.Keyring.NewMnemonic(accountName, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	sr.NoError(err)

	newAddr := sdk.AccAddress(info.GetPubKey().Address())
	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	sr.NoError(err)
	accountAddress = newAddr.String()
}

func (s *IntegrationTestSuite) createBond() {
	val := s.network.Validators[0]
	sr := s.Require()
	createBondCmd := cli.NewCreateBondCmd()
	args := []string{
		fmt.Sprintf("10%s", s.cfg.BondDenom),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
	}
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, createBondCmd, args)
	sr.NoError(err)
	var d sdk.TxResponse
	val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
	sr.Zero(d.Code)
	err = s.network.WaitForNextBlock()
	sr.NoError(err)
}

func (s *IntegrationTestSuite) TestGetQueryBondLists() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name       string
		args       []string
		outputType string
		createBond bool
	}{
		{
			"create and get bond lists",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"json",
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if tc.createBond {
				s.createBond()
			}
			cmd := cli.GetQueryBondLists()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var queryResponse types.QueryGetBondsResponse
			if tc.outputType == "json" {
				err := clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)
			} else {
				err := yaml.Unmarshal(out.Bytes(), &queryResponse)
				sr.NoError(err)
			}
			sr.NotZero(len(queryResponse.GetBonds()))
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondById() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name       string
		args       []string
		createBond bool
		err        bool
	}{
		{
			"invalid bond id",
			[]string{
				fmt.Sprint("not_found_bond_id"),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
			true,
		},
		{
			"create and get bond by id",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if tc.createBond {
				s.createBond()
				cmd := cli.GetQueryBondLists()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
				sr.NoError(err)
				var queryResponse types.QueryGetBondsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				// extract bond id from bonds list
				bond := queryResponse.GetBonds()[0]
				tc.args = append([]string{bond.GetId()}, tc.args...)

			}
			cmd := cli.GetBondByIdCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var queryResponse types.QueryGetBondByIdResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
			sr.NoError(err)
			if tc.err {
				sr.Zero(len(queryResponse.GetBond().GetId()))
			} else {
				sr.NotZero(len(queryResponse.GetBond().GetId()))
				sr.Equal(accountAddress, queryResponse.GetBond().GetOwner())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondListsByOwner() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name       string
		args       []string
		createBond bool
		err        bool
	}{
		{
			"invalid owner address",
			[]string{
				fmt.Sprint("not_found_bond_id"),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
			true,
		},
		{
			"get bond lists by owner address",
			[]string{
				fmt.Sprint(accountAddress),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if tc.createBond {
				s.createBond()
			}
			cmd := cli.GetBondListByOwnerCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var queryResponse types.QueryGetBondsByOwnerResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
			sr.NoError(err)
			if tc.err {
				sr.Zero(len(queryResponse.GetBonds()))
			} else {
				sr.NotZero(len(queryResponse.GetBonds()))
				sr.Equal(accountAddress, queryResponse.GetBonds()[0].GetOwner())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondModuleBalance() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name       string
		args       []string
		createBond bool
	}{
		{
			"get bond module balance",
			[]string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if tc.createBond {
				s.createBond()
			}
			cmd := cli.GetBondModuleBalanceCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var queryResponse types.QueryGetBondModuleBalanceResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
			sr.NoError(err)
			sr.False(queryResponse.GetBalance().IsZero())
		})
	}
}
