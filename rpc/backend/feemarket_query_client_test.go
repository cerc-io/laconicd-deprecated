package backend

import (
	"github.com/cerc-io/laconicd/rpc/backend/mocks"
	rpc "github.com/cerc-io/laconicd/rpc/types"
	feemarkettypes "github.com/cerc-io/laconicd/x/feemarket/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ feemarkettypes.QueryClient = &mocks.FeeMarketQueryClient{}

// Params
func RegisterFeeMarketParams(feeMarketClient *mocks.FeeMarketQueryClient, height int64) {
	feeMarketClient.On("Params", rpc.ContextWithHeight(height), &feemarkettypes.QueryParamsRequest{}).
		Return(&feemarkettypes.QueryParamsResponse{Params: feemarkettypes.DefaultParams()}, nil)
}

func RegisterFeeMarketParamsError(feeMarketClient *mocks.FeeMarketQueryClient, height int64) {
	feeMarketClient.On("Params", rpc.ContextWithHeight(height), &feemarkettypes.QueryParamsRequest{}).
		Return(nil, sdkerrors.ErrInvalidRequest)
}
