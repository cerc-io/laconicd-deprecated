package testutil

import (
	"fmt"

	"github.com/cerc-io/laconicd/crypto/hd"
	"github.com/cerc-io/laconicd/testutil/network"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg              network.Config
	network          *network.Network
	defaultAuctionID string
}

var (
	ownerAccount  = "owner"
	bidderAccount = "bidder"
	ownerAddress  string
	bidderAddress string
)

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() { //nolint: all
	s.T().Log("setting up integration test suite")

	var err error

	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// setting up random owner and bidder accounts
	s.createAccountWithBalance(ownerAccount, &ownerAddress)
	s.createAccountWithBalance(bidderAccount, &bidderAddress)

	s.defaultAuctionID = s.createAuctionAndBid(true, false)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) createAccountWithBalance(accountName string, accountAddress *string) {
	val := s.network.Validators[0]
	sr := s.Require()

	info, _, err := val.ClientCtx.Keyring.NewMnemonic(accountName, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.EthSecp256k1) //nolint:lll
	sr.NoError(err)

	newAddr, _ := info.GetAddress()
	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200000))),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	sr.NoError(err)
	*accountAddress = newAddr.String()
}
