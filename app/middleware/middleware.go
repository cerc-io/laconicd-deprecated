package middleware

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/middleware"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

const (
	secp256k1VerifyCost uint64 = 21000
)

type MD struct {
	ethMiddleware    tx.Handler
	cosmosMiddleware tx.Handler
	cosmoseip792     tx.Handler
}

var _ tx.Handler = MD{}

func NewMiddleware(indexEventsStr []string, options HandlerOptions) (tx.Handler, error) {
	ethMiddleware, err := newEthAuthMiddleware(options)
	if err != nil {
		return nil, err
	}
	cosmoseip792, err := newCosmosAnteHandlerEip712(options)
	if err != nil {
		return nil, err
	}
	cosmosMiddleware, err := newCosmosAuthMiddleware(options)
	if err != nil {
		return nil, err
	}
	return MD{
		ethMiddleware:    ethMiddleware,
		cosmosMiddleware: cosmosMiddleware,
		cosmoseip792:     cosmoseip792,
	}, nil
}

// CheckTx implements tx.Handler
func (md MD) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	var anteHandler tx.Handler
	reqTx := req.Tx
	txWithExtensions, ok := reqTx.(authante.HasExtensionOptionsTx)
	if ok {
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) > 0 {
			switch typeURL := opts[0].GetTypeUrl(); typeURL {
			case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
				// handle as *evmtypes.MsgEthereumTx
				anteHandler = md.ethMiddleware
			case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
				// handle as normal Cosmos SDK tx, except signature is checked for EIP712 representation
				anteHandler = md.cosmoseip792
			default:
				return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(
					sdkerrors.ErrUnknownExtensionOptions,
					"rejecting tx with unsupported extension option: %s", typeURL,
				)
			}

			return anteHandler.CheckTx(ctx, req, checkReq)
		}
	}

	// // handle as totally normal Cosmos SDK tx
	// _, ok = reqTx.(sdk.Tx)
	// if !ok {
	// 	return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", reqTx)
	// }

	anteHandler = md.cosmosMiddleware

	return anteHandler.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (md MD) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	var anteHandler tx.Handler
	reqTx := req.Tx
	txWithExtensions, ok := reqTx.(authante.HasExtensionOptionsTx)
	if ok {
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) > 0 {
			switch typeURL := opts[0].GetTypeUrl(); typeURL {
			case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
				// handle as *evmtypes.MsgEthereumTx
				anteHandler = md.ethMiddleware
			case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
				// handle as normal Cosmos SDK tx, except signature is checked for EIP712 representation
				anteHandler = md.cosmoseip792
			default:
				return tx.Response{}, sdkerrors.Wrapf(
					sdkerrors.ErrUnknownExtensionOptions,
					"rejecting tx with unsupported extension option: %s", typeURL,
				)
			}

			return anteHandler.DeliverTx(ctx, req)
		}
	}

	anteHandler = md.cosmosMiddleware
	return anteHandler.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (md MD) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	var anteHandler tx.Handler
	reqTx := req.Tx
	txWithExtensions, ok := reqTx.(authante.HasExtensionOptionsTx)
	if ok {
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) > 0 {
			switch typeURL := opts[0].GetTypeUrl(); typeURL {
			case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
				// handle as *evmtypes.MsgEthereumTx
				anteHandler = md.ethMiddleware
			case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
				// handle as normal Cosmos SDK tx, except signature is checked for EIP712 representation
				anteHandler = md.cosmoseip792
			default:
				return tx.Response{}, sdkerrors.Wrapf(
					sdkerrors.ErrUnknownExtensionOptions,
					"rejecting tx with unsupported extension option: %s", typeURL,
				)
			}

			return anteHandler.SimulateTx(ctx, req)
		}
	}

	anteHandler = md.cosmosMiddleware
	return anteHandler.SimulateTx(ctx, req)
}

var _ authante.SignatureVerificationGasConsumer = DefaultSigVerificationGasConsumer

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

	return authante.DefaultSigVerificationGasConsumer(meter, sig, params)
}
