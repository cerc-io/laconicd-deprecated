package middleware

import (
	context "context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// RejectMessagesMiddleware prevents invalid msg types from being executed
type RejectMessagesMiddleware struct {
	next tx.Handler
}

func NewRejectMessagesMiddleware(txh tx.Handler) tx.Handler {
	return RejectMessagesMiddleware{
		next: txh,
	}
}

// Middleware rejects messages that requires ethereum-specific authentication.
// For example `MsgEthereumTx` requires fee to be deducted in the antehandler in
// order to perform the refund.

// CheckTx implements tx.Handler
func (rmd RejectMessagesMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	if _, err := reject(req); err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	return rmd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (rmd RejectMessagesMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := reject(req); err != nil {
		return tx.Response{}, err
	}

	return rmd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (rmd RejectMessagesMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := reject(req); err != nil {
		return tx.Response{}, err
	}

	return rmd.next.SimulateTx(ctx, req)
}

func reject(req tx.Request) (tx.Response, error) {
	reqTx := req.Tx
	for _, msg := range reqTx.GetMsgs() {
		if _, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			return tx.Response{}, sdkerrors.Wrapf(
				sdkerrors.ErrInvalidType,
				"MsgEthereumTx needs to be contained within a tx with 'ExtensionOptionsEthereumTx' option",
			)
		}
	}

	return tx.Response{}, nil
}

var _ tx.Handler = RejectMessagesMiddleware{}
