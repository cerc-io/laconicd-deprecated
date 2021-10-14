package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/testutil/network"
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

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
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
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200000))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	sr.NoError(err)
	*accountAddress = newAddr.String()
}
