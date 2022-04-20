package middleware

import (
	context "context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// RejectMessagesDecorator prevents invalid msg types from being executed
type RejectMessagesDecorator struct {
	next tx.Handler
}

func NewRejectMessagesDecorator() tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return RejectMessagesDecorator{
			next: h,
		}
	}
}

// Middleware rejects messages that requires ethereum-specific authentication.
// For example `MsgEthereumTx` requires fee to be deducted in the antehandler in
// order to perform the refund.

// CheckTx implements tx.Handler
func (rmd RejectMessagesDecorator) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	reqTx := req.Tx
	for _, msg := range reqTx.GetMsgs() {
		if _, ok := msg.(*evmtypes.MsgEthereumTx); ok {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(
				sdkerrors.ErrInvalidType,
				"MsgEthereumTx needs to be contained within a tx with 'ExtensionOptionsEthereumTx' option",
			)
		}
	}
	return rmd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (rmd RejectMessagesDecorator) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return rmd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (rmd RejectMessagesDecorator) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return rmd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = RejectMessagesDecorator{}
