package middleware

import (
	context "context"

	"github.com/cosmos/cosmos-sdk/types/tx"
)

type MD struct {
	next tx.Handler
}

var _ tx.Handler = MD{}

func NewMiddleware(indexEventsStr []string, options HandlerOptions) (tx.Handler, error) {
	return newEthAuthMiddleware(options)
}

// CheckTx implements tx.Handler
func (md MD) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	return md.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (md MD) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return md.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (md MD) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return md.next.SimulateTx(ctx, req)
}
