package middleware

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authmiddleware "github.com/cosmos/cosmos-sdk/x/auth/middleware"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// HandlerOptions extend the SDK's TxHandler options by requiring the IBC
// channel keeper, EVM Keeper and Fee Market Keeper.
type HandlerOptions struct {
	Debug bool

	Codec codec.Codec
	// TxDecoder is used to decode the raw tx bytes into a sdk.Tx.
	TxDecoder sdk.TxDecoder

	// IndexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs Tendermint what to index. If empty, all events will be indexed.
	IndexEvents map[string]struct{}

	LegacyRouter     sdk.Router
	MsgServiceRouter *authmiddleware.MsgServiceRouter

	ExtensionOptionChecker authmiddleware.ExtensionOptionChecker
	TxFeeChecker           authmiddleware.TxFeeChecker

	AccountKeeper   evmtypes.AccountKeeper
	BankKeeper      evmtypes.BankKeeper
	FeeMarketKeeper evmtypes.FeeMarketKeeper
	EvmKeeper       EVMKeeper
	FeegrantKeeper  authmiddleware.FeegrantKeeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
	MaxTxGasWanted  uint64
}

func (options HandlerOptions) Validate() error {
	if options.TxDecoder == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "txDecoder is required for middlewares")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for middlewares")
	}
	if options.AccountKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for middlewares")
	}
	if options.BankKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for middlewares")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for middlewares")
	}
	if options.FeeMarketKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee market keeper is required for middlewares")
	}
	if options.EvmKeeper == nil {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "evm keeper is required for middlewares")
	}
	return nil
}

func newEthAuthMiddleware(options HandlerOptions) tx.Middleware {
	stack := []tx.Middleware{
		NewEthSetUpContextMiddleware(options.EvmKeeper),
		NewEthMempoolFeeMiddleware(options.EvmKeeper),
		NewEthValidateBasicMiddleware(options.EvmKeeper),
		NewEthSigVerificationMiddleware(options.EvmKeeper),
		NewEthAccountVerificationMiddleware(options.AccountKeeper, options.BankKeeper, options.EvmKeeper),
		NewEthGasConsumeMiddleware(options.EvmKeeper, options.MaxTxGasWanted),
		NewCanTransferMiddleware(options.EvmKeeper),
		NewEthIncrementSenderSequenceMiddleware(options.AccountKeeper),
	}
	return func(txh tx.Handler) tx.Handler {
		return authmiddleware.ComposeMiddlewares(txh, stack...)
	}
}

func newCosmosAuthMiddleware(options HandlerOptions) tx.Middleware {
	stack := []tx.Middleware{
		NewRejectMessagesMiddleware,
		// Set a new GasMeter on sdk.Context.
		//
		// Make sure the Gas middleware is outside of all other middlewares
		// that reads the GasMeter. In our case, the Recovery middleware reads
		// the GasMeter to populate GasInfo.
		authmiddleware.GasTxMiddleware,
		// Recover from panics. Panics outside of this middleware won't be
		// caught, be careful!
		authmiddleware.RecoveryTxMiddleware,
		// Choose which events to index in Tendermint. Make sure no events are
		// emitted outside of this middleware.
		authmiddleware.NewIndexEventsTxMiddleware(options.IndexEvents),
		// Reject all extension options other than the ones needed by the feemarket.
		authmiddleware.NewExtensionOptionsMiddleware(options.ExtensionOptionChecker),
		authmiddleware.ValidateBasicMiddleware,
		authmiddleware.TxTimeoutHeightMiddleware,
		authmiddleware.ValidateMemoMiddleware(options.AccountKeeper),
		authmiddleware.ConsumeTxSizeGasMiddleware(options.AccountKeeper),
		// No gas should be consumed in any middleware above in a "post" handler part. See
		// ComposeMiddlewares godoc for details.
		// `DeductFeeMiddleware` and `IncrementSequenceMiddleware` should be put outside of `WithBranchedStore` middleware,
		// so their storage writes are not discarded when tx fails.
		authmiddleware.DeductFeeMiddleware(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		authmiddleware.SetPubKeyMiddleware(options.AccountKeeper),
		authmiddleware.ValidateSigCountMiddleware(options.AccountKeeper),
		authmiddleware.SigGasConsumeMiddleware(options.AccountKeeper, options.SigGasConsumer),
		authmiddleware.SigVerificationMiddleware(options.AccountKeeper, options.SignModeHandler),
		authmiddleware.IncrementSequenceMiddleware(options.AccountKeeper),
		// Creates a new MultiStore branch, discards downstream writes if the downstream returns error.
		// These kinds of middlewares should be put under this:
		// - Could return error after messages executed succesfully.
		// - Storage writes should be discarded together when tx failed.
		authmiddleware.WithBranchedStore,
		// Consume block gas. All middlewares whose gas consumption after their `next` handler
		// should be accounted for, should go below this middleware.
		authmiddleware.ConsumeBlockGasMiddleware,
		authmiddleware.NewTipMiddleware(options.BankKeeper),
	}
	return func(txh tx.Handler) tx.Handler {
		return authmiddleware.ComposeMiddlewares(txh, stack...)
	}
}

func newCosmosMiddlewareEip712(options HandlerOptions) tx.Middleware {
	stack := []tx.Middleware{
		NewRejectMessagesMiddleware,
		// Set a new GasMeter on sdk.Context.
		//
		// Make sure the Gas middleware is outside of all other middlewares
		// that reads the GasMeter. In our case, the Recovery middleware reads
		// the GasMeter to populate GasInfo.
		authmiddleware.GasTxMiddleware,
		// Recover from panics. Panics outside of this middleware won't be
		// caught, be careful!
		authmiddleware.RecoveryTxMiddleware,
		// Choose which events to index in Tendermint. Make sure no events are
		// emitted outside of this middleware.
		authmiddleware.NewIndexEventsTxMiddleware(options.IndexEvents),
		// Reject all extension options other than the ones needed by the feemarket.
		// authmiddleware.NewExtensionOptionsMiddleware(options.ExtensionOptionChecker),
		authmiddleware.ValidateBasicMiddleware,
		authmiddleware.TxTimeoutHeightMiddleware,
		authmiddleware.ValidateMemoMiddleware(options.AccountKeeper),
		authmiddleware.ConsumeTxSizeGasMiddleware(options.AccountKeeper),
		// No gas should be consumed in any middleware above in a "post" handler part. See
		// ComposeMiddlewares godoc for details.
		// `DeductFeeMiddleware` and `IncrementSequenceMiddleware` should be put outside of `WithBranchedStore` middleware,
		// so their storage writes are not discarded when tx fails.
		authmiddleware.DeductFeeMiddleware(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		authmiddleware.SetPubKeyMiddleware(options.AccountKeeper),
		authmiddleware.ValidateSigCountMiddleware(options.AccountKeeper),
		authmiddleware.SigGasConsumeMiddleware(options.AccountKeeper, options.SigGasConsumer),
		// Note: signature verification uses EIP instead of the cosmos signature validator
		NewEip712SigVerificationMiddleware(options.Codec, options.AccountKeeper, options.SignModeHandler),
		authmiddleware.IncrementSequenceMiddleware(options.AccountKeeper),
		// Creates a new MultiStore branch, discards downstream writes if the downstream returns error.
		// These kinds of middlewares should be put under this:
		// - Could return error after messages executed succesfully.
		// - Storage writes should be discarded together when tx failed.
		authmiddleware.WithBranchedStore,
		// Consume block gas. All middlewares whose gas consumption after their `next` handler
		// should be accounted for, should go below this middleware.
		authmiddleware.ConsumeBlockGasMiddleware,
		authmiddleware.NewTipMiddleware(options.BankKeeper),
	}
	return func(txh tx.Handler) tx.Handler {
		return authmiddleware.ComposeMiddlewares(txh, stack...)
	}
}
