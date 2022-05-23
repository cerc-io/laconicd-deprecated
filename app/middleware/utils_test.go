package middleware_test

import (
	context "context"
	"math"
	"testing"
	"time"

	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types3 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cast"
	"github.com/tharsis/ethermint/app/middleware"
	"github.com/tharsis/ethermint/ethereum/eip712"
	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tharsis/ethermint/app"
	"github.com/tharsis/ethermint/encoding"
	"github.com/tharsis/ethermint/tests"
	"github.com/tharsis/ethermint/x/evm/statedb"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

// customTxHandler is a test middleware that will run a custom function.
type customTxHandler struct {
	fn func(context.Context, tx.Request) (tx.Response, error)
}

var _ tx.Handler = customTxHandler{}

func (h customTxHandler) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return h.fn(ctx, req)
}
func (h customTxHandler) CheckTx(ctx context.Context, req tx.Request, _ tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	res, err := h.fn(ctx, req)
	return res, tx.ResponseCheckTx{}, err
}
func (h customTxHandler) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return h.fn(ctx, req)
}

// noopTxHandler is a test middleware that returns an empty response.
var noopTxHandler = customTxHandler{func(_ context.Context, _ tx.Request) (tx.Response, error) {
	return tx.Response{}, nil
}}

type MiddlewareTestSuite struct {
	suite.Suite

	ctx             sdk.Context
	app             *app.EthermintApp
	clientCtx       client.Context
	anteHandler     tx.Handler
	ethSigner       ethtypes.Signer
	enableFeemarket bool
	enableLondonHF  bool
}

func (suite *MiddlewareTestSuite) StateDB() *statedb.StateDB {
	return statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))
}

func (suite *MiddlewareTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(suite.T(), checkTx, func(app *app.EthermintApp, genesis simapp.GenesisState) simapp.GenesisState {
		if suite.enableFeemarket {
			// setup feemarketGenesis params
			feemarketGenesis := feemarkettypes.DefaultGenesisState()
			feemarketGenesis.Params.EnableHeight = 1
			feemarketGenesis.Params.NoBaseFee = false
			// Verify feeMarket genesis
			err := feemarketGenesis.Validate()
			suite.Require().NoError(err)
			genesis[feemarkettypes.ModuleName] = app.AppCodec().MustMarshalJSON(feemarketGenesis)
		}
		if !suite.enableLondonHF {
			evmGenesis := evmtypes.DefaultGenesisState()
			maxInt := sdk.NewInt(math.MaxInt64)
			evmGenesis.Params.ChainConfig.LondonBlock = &maxInt
			evmGenesis.Params.ChainConfig.ArrowGlacierBlock = &maxInt
			evmGenesis.Params.ChainConfig.MergeForkBlock = &maxInt
			genesis[evmtypes.ModuleName] = app.AppCodec().MustMarshalJSON(evmGenesis)
		}
		return genesis
	})

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 2, ChainID: "ethermint_9000-1", Time: time.Now().UTC()})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(evmtypes.DefaultEVMDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1000000000000000000))
	suite.app.EvmKeeper.WithChainID(suite.ctx)

	infCtx := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	suite.app.AccountKeeper.SetParams(infCtx, authtypes.DefaultParams())

	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	// We're using TestMsg amino encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	maxGasWanted := cast.ToUint64(config.DefaultMaxTxGasWanted)

	options := middleware.HandlerOptions{
		TxDecoder:        suite.clientCtx.TxConfig.TxDecoder(),
		AccountKeeper:    suite.app.AccountKeeper,
		BankKeeper:       suite.app.BankKeeper,
		EvmKeeper:        suite.app.EvmKeeper,
		FeegrantKeeper:   suite.app.FeeGrantKeeper,
		Codec:            suite.app.AppCodec(),
		Debug:            suite.app.Trace(),
		LegacyRouter:     suite.app.LegacyRouter,
		MsgServiceRouter: suite.app.MsgSvcRouter,
		FeeMarketKeeper:  suite.app.FeeMarketKeeper,
		SignModeHandler:  suite.clientCtx.TxConfig.SignModeHandler(),
		SigGasConsumer:   middleware.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:   maxGasWanted,
	}

	suite.Require().NoError(options.Validate())
	suite.anteHandler = middleware.NewTxHandler(options)
	suite.ethSigner = ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &MiddlewareTestSuite{
		enableLondonHF: true,
	})
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *MiddlewareTestSuite) CreateTestTx(
	msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool,
	unsetExtensionOptions ...bool,
) authsigning.Tx {
	return suite.CreateTestTxBuilder(msg, priv, accNum, signCosmosTx).GetTx()
}

// CreateTestTxBuilder is a helper function to create a tx builder given multiple inputs.
func (suite *MiddlewareTestSuite) CreateTestTxBuilder(
	msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool,
	unsetExtensionOptions ...bool,
) client.TxBuilder {
	var option *codectypes.Any
	var err error
	if len(unsetExtensionOptions) == 0 {
		option, err = codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
		suite.Require().NoError(err)
	}

	txBuilder := suite.clientCtx.TxConfig.NewTxBuilder()
	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	suite.Require().True(ok)

	if len(unsetExtensionOptions) == 0 {
		builder.SetExtensionOptions(option)
	}

	err = msg.Sign(suite.ethSigner, tests.NewSigner(priv))
	suite.Require().NoError(err)

	err = builder.SetMsgs(msg)
	suite.Require().NoError(err)

	txData, err := evmtypes.UnpackTxData(msg.Data)
	suite.Require().NoError(err)

	fees := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdk.NewIntFromBigInt(txData.Fee())))
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())

	if signCosmosTx {
		// First round: we gather all the signer infos. We use the "set empty
		// signature" hack to do that.
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: txData.GetNonce(),
		}

		sigsV2 := []signing.SignatureV2{sigV2}

		err = txBuilder.SetSignatures(sigsV2...)
		suite.Require().NoError(err)

		// Second round: all signer infos are set, so each signer can sign.

		signerData := authsigning.SignerData{
			ChainID:       suite.ctx.ChainID(),
			AccountNumber: accNum,
			Sequence:      txData.GetNonce(),
		}
		sigV2, err = clientTx.SignWithPrivKey(
			suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData,
			txBuilder, priv, suite.clientCtx.TxConfig, txData.GetNonce(),
		)
		suite.Require().NoError(err)

		sigsV2 = []signing.SignatureV2{sigV2}

		err = txBuilder.SetSignatures(sigsV2...)
		suite.Require().NoError(err)
	}

	return txBuilder
}

func (suite *MiddlewareTestSuite) CreateTestEIP712TxBuilderMsgSend(from sdk.AccAddress, priv cryptotypes.PrivKey, chainId string, gas uint64, gasAmount sdk.Coins) client.TxBuilder {
	// Build MsgSend
	recipient := sdk.AccAddress(common.Address{}.Bytes())
	msgSend := types2.NewMsgSend(from, recipient, sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdk.NewInt(1))))
	return suite.CreateTestEIP712CosmosTxBuilder(from, priv, chainId, gas, gasAmount, msgSend)
}

func (suite *MiddlewareTestSuite) CreateTestEIP712TxBuilderMsgDelegate(from sdk.AccAddress, priv cryptotypes.PrivKey, chainId string, gas uint64, gasAmount sdk.Coins) client.TxBuilder {
	// Build MsgSend
	valEthAddr := tests.GenerateAddress()
	valAddr := sdk.ValAddress(valEthAddr.Bytes())
	msgSend := types3.NewMsgDelegate(from, valAddr, sdk.NewCoin(evmtypes.DefaultEVMDenom, sdk.NewInt(200000)))
	return suite.CreateTestEIP712CosmosTxBuilder(from, priv, chainId, gas, gasAmount, msgSend)
}

func (suite *MiddlewareTestSuite) CreateTestEIP712CosmosTxBuilder(
	from sdk.AccAddress, priv cryptotypes.PrivKey, chainId string, gas uint64, gasAmount sdk.Coins, msg sdk.Msg,
) client.TxBuilder {
	var err error

	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, from)
	suite.Require().NoError(err)

	pc, err := types.ParseChainID(chainId)
	suite.Require().NoError(err)
	ethChainId := pc.Uint64()

	// GenerateTypedData TypedData
	var ethermintCodec codec.ProtoCodecMarshaler
	fee := legacytx.NewStdFee(gas, gasAmount)
	accNumber := suite.app.AccountKeeper.GetAccount(suite.ctx, from).GetAccountNumber()

	data := legacytx.StdSignBytes(chainId, accNumber, nonce, 0, fee, []sdk.Msg{msg}, "", nil)
	typedData, err := eip712.WrapTxToTypedData(ethermintCodec, ethChainId, msg, data, &eip712.FeeDelegationOptions{
		FeePayer: from,
	})
	suite.Require().NoError(err)

	sigHash, err := eip712.ComputeTypedDataHash(typedData)
	suite.Require().NoError(err)

	// Sign typedData
	keyringSigner := tests.NewSigner(priv)
	signature, pubKey, err := keyringSigner.SignByAddress(from, sigHash)
	suite.Require().NoError(err)
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper

	// Add ExtensionOptionsWeb3Tx extension
	var option *codectypes.Any
	option, err = codectypes.NewAnyWithValue(&types.ExtensionOptionsWeb3Tx{
		FeePayer:         from.String(),
		TypedDataChainID: ethChainId,
		FeePayerSig:      signature,
	})
	suite.Require().NoError(err)

	suite.clientCtx.TxConfig.SignModeHandler()
	txBuilder := suite.clientCtx.TxConfig.NewTxBuilder()
	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	suite.Require().True(ok)

	builder.SetExtensionOptions(option)
	builder.SetFeeAmount(gasAmount)
	builder.SetGasLimit(gas)

	sigsV2 := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode: signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		},
		Sequence: nonce,
	}

	err = builder.SetSignatures(sigsV2)
	suite.Require().NoError(err)

	err = builder.SetMsgs(msg)
	suite.Require().NoError(err)

	return builder
}

var _ sdk.Tx = &invalidTx{}

type invalidTx struct{}

func (invalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (invalidTx) ValidateBasic() error { return nil }
