package testutil

import (
	"fmt"

	"github.com/cerc-io/laconicd/testutil/network"
	"github.com/cerc-io/laconicd/x/bond/client/cli"
	"github.com/cerc-io/laconicd/x/bond/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg            network.Config
	network        *network.Network
	accountName    string
	accountAddress string
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.accountName = "accountName"
	// setting up random account
	s.createAccountWithBalance(s.accountName)
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
			var response types.QueryParamsResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
			s.Require().NoError(err)
			s.Require().Equal(response.Params.MaxBondAmount, types.DefaultParams().MaxBondAmount)
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

	newAddr, _ := info.GetAddress()
	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))),
		fmt.Sprintf("--%s", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	sr.NoError(err)
	s.accountAddress = newAddr.String()
}

func (s *IntegrationTestSuite) createBond() {
	val := s.network.Validators[0]
	sr := s.Require()
	createBondCmd := cli.NewCreateBondCmd()
	args := []string{
		fmt.Sprintf("10%s", s.cfg.BondDenom),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.accountName),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
	}
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, createBondCmd, args)
	sr.NoError(err)
	var d sdk.TxResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
	sr.NoError(err)
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
		createBond bool
		preRun     func()
	}{
		{
			"create and get bond lists",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			func() {
				s.createBond()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if tc.createBond {
				tc.preRun()
			}

			cmd := cli.GetQueryBondLists()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var queryResponse types.QueryGetBondsResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
			sr.NoError(err)
			sr.NotZero(len(queryResponse.GetBonds()))
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondById() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name   string
		args   []string
		err    bool
		preRun func() string
	}{
		{
			"invalid bond id",
			[]string{
				"not_found_bond_id",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
			func() string {
				return ""
			},
		},
		{
			"create and get bond by id",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			func() string {
				// creating the bond
				s.createBond()

				// getting the bonds list and returning the bond-id
				clientCtx := val.ClientCtx
				cmd := cli.GetQueryBondLists()
				args := []string{
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var queryResponse types.QueryGetBondsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				// extract bond id from bonds list
				bond := queryResponse.GetBonds()[0]
				return bond.GetId()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if !tc.err {
				bondId := tc.preRun()
				tc.args = append([]string{bondId}, tc.args...)
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
				sr.Equal(s.accountAddress, queryResponse.GetBond().GetOwner())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondListsByOwner() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name   string
		args   []string
		err    bool
		preRun func()
	}{
		{
			"invalid owner address",
			[]string{
				"not_found_bond_id",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			true,
			func() {

			},
		},
		{
			"get bond lists by owner address",
			[]string{
				fmt.Sprint(s.accountAddress),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
			func() {
				s.createBond()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if !tc.err {
				tc.preRun()
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
				sr.Equal(s.accountAddress, queryResponse.GetBonds()[0].GetOwner())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondModuleBalance() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name   string
		args   []string
		err    bool
		preRun func()
	}{
		{
			"get bond module balance",
			[]string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			false,
			func() {
				s.createBond()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if !tc.err {
				tc.preRun()
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
