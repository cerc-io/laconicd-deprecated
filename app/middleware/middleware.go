package middleware

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authmiddleware "github.com/cosmos/cosmos-sdk/x/auth/middleware"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

const (
	secp256k1VerifyCost uint64 = 21000
)

type txRouter struct {
	eth, cosmos, eip712 tx.Handler
}

var _ tx.Handler = txRouter{}

func NewTxHandler(options HandlerOptions) tx.Handler {
	return authmiddleware.ComposeMiddlewares(
		authmiddleware.NewRunMsgsTxHandler(options.MsgServiceRouter, options.LegacyRouter),
		authmiddleware.NewTxDecoderMiddleware(options.TxDecoder),
		NewTxRouterMiddleware(options),
	)
}

func NewTxRouterMiddleware(options HandlerOptions) tx.Middleware {
	ethMiddleware := newEthAuthMiddleware(options)
	cosmoseip712 := newCosmosMiddlewareEip712(options)
	cosmosMiddleware := newCosmosAuthMiddleware(options)
	return func(txh tx.Handler) tx.Handler {
		return txRouter{
			eth:    ethMiddleware(txh),
			cosmos: cosmosMiddleware(txh),
			eip712: cosmoseip712(txh),
		}
	}
}

// CheckTx implements tx.Handler
func (txh txRouter) route(req tx.Request) (tx.Handler, error) {
	txWithExtensions, ok := req.Tx.(authmiddleware.HasExtensionOptionsTx)
	if ok {
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) > 0 {
			var next tx.Handler
			switch typeURL := opts[0].GetTypeUrl(); typeURL {
			case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
				// handle as *evmtypes.MsgEthereumTx
				next = txh.eth
			case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
				// handle as normal Cosmos SDK tx, except signature is checked for EIP712 representation
				next = txh.eip712
			default:
				return nil, sdkerrors.Wrapf(
					sdkerrors.ErrUnknownExtensionOptions,
					"rejecting tx with unsupported extension option: %s", typeURL,
				)
			}
			return next, nil
		}
	}
	// // handle as totally normal Cosmos SDK tx
	// if _, ok = reqTx.(sdk.Tx); !ok {
	// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", reqTx)
	// }
	return txh.cosmos, nil
}

// CheckTx implements tx.Handler
func (txh txRouter) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (res tx.Response, rct tx.ResponseCheckTx, err error) {
	next, err := txh.route(req)
	if err != nil {
		return
	}
	return next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (txh txRouter) DeliverTx(ctx context.Context, req tx.Request) (res tx.Response, err error) {
	next, err := txh.route(req)
	if err != nil {
		return
	}
	return next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (txh txRouter) SimulateTx(ctx context.Context, req tx.Request) (res tx.Response, err error) {
	next, err := txh.route(req)
	if err != nil {
		return
	}
	return next.SimulateTx(ctx, req)
}

var _ authmiddleware.SignatureVerificationGasConsumer = DefaultSigVerificationGasConsumer

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params,
) error {
	// support for ethereum ECDSA secp256k1 keys
	_, ok := sig.PubKey.(*ethsecp256k1.PubKey)
	if ok {
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: eth_secp256k1")
		return nil
	}

	return authmiddleware.DefaultSigVerificationGasConsumer(meter, sig, params)
}
