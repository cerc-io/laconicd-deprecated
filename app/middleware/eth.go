package middleware

import (
	context "context"
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/middleware"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethermint "github.com/tharsis/ethermint/types"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	"github.com/tharsis/ethermint/x/evm/statedb"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// EthSetupContextMiddleware is adapted from SetUpContextMiddleware from cosmos-sdk, it ignores gas consumption
// by setting the gas meter to infinite
type EthSetupContextMiddleware struct {
	next      tx.Handler
	evmKeeper EVMKeeper
}

// gasContext returns a new context with a gas meter set from a given context.
func gasContext(ctx sdk.Context, tx sdk.Tx, isSimulate bool) (sdk.Context, error) {
	// all transactions must implement GasTx
	gasTx, ok := tx.(middleware.GasTx)
	if !ok {
		// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
		// during runTx.
		newCtx := setGasMeter(ctx, 0, isSimulate)
		return newCtx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be GasTx")
	}

	return setGasMeter(ctx, gasTx.GetGas(), isSimulate), nil
}

// setGasMeter returns a new context with a gas meter set from a given context.
func setGasMeter(ctx sdk.Context, gasLimit uint64, simulate bool) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
}

// CheckTx implements tx.Handler
func (esc EthSetupContextMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	sdkCtx, err := gasContext(sdk.UnwrapSDKContext(ctx), req.Tx, false)
	if err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	// Reset transient gas used to prepare the execution of current cosmos tx.
	// Transient gas-used is necessary to sum the gas-used of cosmos tx, when it contains multiple eth msgs.
	esc.evmKeeper.ResetTransientGasUsed(sdkCtx)
	return esc.next.CheckTx(sdkCtx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (esc EthSetupContextMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	sdkCtx, err := gasContext(sdk.UnwrapSDKContext(ctx), req.Tx, false)
	if err != nil {
		return tx.Response{}, err
	}

	// Reset transient gas used to prepare the execution of current cosmos tx.
	// Transient gas-used is necessary to sum the gas-used of cosmos tx, when it contains multiple eth msgs.
	esc.evmKeeper.ResetTransientGasUsed(sdkCtx)
	return esc.next.DeliverTx(sdkCtx, req)
}

// SimulateTx implements tx.Handler
func (esc EthSetupContextMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	sdkCtx, err := gasContext(sdk.UnwrapSDKContext(ctx), req.Tx, false)
	if err != nil {
		return tx.Response{}, err
	}

	// Reset transient gas used to prepare the execution of current cosmos tx.
	// Transient gas-used is necessary to sum the gas-used of cosmos tx, when it contains multiple eth msgs.
	esc.evmKeeper.ResetTransientGasUsed(sdkCtx)
	return esc.next.SimulateTx(sdkCtx, req)
}

var _ tx.Handler = EthSetupContextMiddleware{}

func NewEthSetUpContextMiddleware(evmKeeper EVMKeeper) tx.Middleware {
	return func(txh tx.Handler) tx.Handler {
		return EthSetupContextMiddleware{
			next:      txh,
			evmKeeper: evmKeeper,
		}
	}
}

// EthMempoolFeeMiddleware will check if the transaction's effective fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, Middleware returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeMiddleware
type EthMempoolFeeMiddleware struct {
	next      tx.Handler
	evmKeeper EVMKeeper
}

// CheckTx implements tx.Handler
func (mfd EthMempoolFeeMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := mfd.evmKeeper.GetParams(sdkCtx)
	ethCfg := params.ChainConfig.EthereumConfig(mfd.evmKeeper.ChainID())
	baseFee := mfd.evmKeeper.BaseFee(sdkCtx, ethCfg)
	if baseFee == nil {
		for _, msg := range req.Tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
			}

			evmDenom := params.EvmDenom
			feeAmt := ethMsg.GetFee()
			glDec := sdk.NewDec(int64(ethMsg.GetGas()))
			requiredFee := sdkCtx.MinGasPrices().AmountOf(evmDenom).Mul(glDec)
			if sdk.NewDecFromBigInt(feeAmt).LT(requiredFee) {
				return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeAmt, requiredFee)
			}
		}
	}

	return mfd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (mfd EthMempoolFeeMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return mfd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (mfd EthMempoolFeeMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return mfd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = EthMempoolFeeMiddleware{}

func NewEthMempoolFeeMiddleware(ek EVMKeeper) tx.Middleware {
	return func(txh tx.Handler) tx.Handler {
		return EthMempoolFeeMiddleware{
			next:      txh,
			evmKeeper: ek,
		}
	}
}

// EthValidateBasicMiddleware is adapted from ValidateBasicMiddleware from cosmos-sdk, it ignores ErrNoSignatures
type EthValidateBasicMiddleware struct {
	next      tx.Handler
	evmKeeper EVMKeeper
}

// CheckTx implements tx.Handler
func (vbd EthValidateBasicMiddleware) CheckTx(cx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	reqTx := req.Tx

	// no need to validate basic on recheck tx, call next antehandler
	if ctx.IsReCheckTx() {
		return vbd.next.CheckTx(ctx, req, checkReq)
	}

	err := reqTx.ValidateBasic()
	// ErrNoSignatures is fine with eth tx
	if err != nil && !errors.Is(err, sdkerrors.ErrNoSignatures) {
		return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(err, "tx basic validation failed")
	}

	// For eth type cosmos tx, some fields should be verified as zero values,
	// since we will only verify the signature against the hash of the MsgEthereumTx.Data
	if wrapperTx, ok := reqTx.(protoTxProvider); ok {
		protoTx := wrapperTx.GetProtoTx()
		body := protoTx.Body
		if body.Memo != "" || body.TimeoutHeight != uint64(0) || len(body.NonCriticalExtensionOptions) > 0 {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
				"for eth tx body Memo TimeoutHeight NonCriticalExtensionOptions should be empty")
		}

		if len(body.ExtensionOptions) != 1 {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "for eth tx length of ExtensionOptions should be 1")
		}

		txFee := sdk.Coins{}
		txGasLimit := uint64(0)

		for _, msg := range protoTx.GetMsgs() {
			msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
			}
			txGasLimit += msgEthTx.GetGas()

			txData, err := evmtypes.UnpackTxData(msgEthTx.Data)
			if err != nil {
				return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(err, "failed to unpack MsgEthereumTx Data")
			}

			params := vbd.evmKeeper.GetParams(ctx)
			chainID := vbd.evmKeeper.ChainID()
			ethCfg := params.ChainConfig.EthereumConfig(chainID)
			baseFee := vbd.evmKeeper.BaseFee(ctx, ethCfg)
			if baseFee == nil && txData.TxType() == ethtypes.DynamicFeeTxType {
				return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(ethtypes.ErrTxTypeNotSupported, "dynamic fee tx not supported")
			}

			txFee = txFee.Add(sdk.NewCoin(params.EvmDenom, sdk.NewIntFromBigInt(txData.Fee())))
		}

		authInfo := protoTx.AuthInfo
		if len(authInfo.SignerInfos) > 0 {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "for eth tx AuthInfo SignerInfos should be empty")
		}

		if authInfo.Fee.Payer != "" || authInfo.Fee.Granter != "" {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "for eth tx AuthInfo Fee payer and granter should be empty")
		}

		if !authInfo.Fee.Amount.IsEqual(txFee) {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid AuthInfo Fee Amount (%s != %s)", authInfo.Fee.Amount, txFee)
		}

		if authInfo.Fee.GasLimit != txGasLimit {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid AuthInfo Fee GasLimit (%d != %d)", authInfo.Fee.GasLimit, txGasLimit)
		}

		sigs := protoTx.Signatures
		if len(sigs) > 0 {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "for eth tx Signatures should be empty")
		}
	}

	return vbd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (vbd EthValidateBasicMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return vbd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (vbd EthValidateBasicMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return vbd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = EthValidateBasicMiddleware{}

// NewEthValidateBasicMiddleware creates a new EthValidateBasicMiddleware
func NewEthValidateBasicMiddleware(ek EVMKeeper) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return EthValidateBasicMiddleware{
			next:      h,
			evmKeeper: ek,
		}
	}

}

// EthSigVerificationMiddleware validates an ethereum signatures
type EthSigVerificationMiddleware struct {
	next      tx.Handler
	evmKeeper EVMKeeper
}

func ethSigVerificationMiddleware(esvd EthSigVerificationMiddleware, cx context.Context, req tx.Request) (tx.Response, error) {
	chainID := esvd.evmKeeper.ChainID()
	ctx := sdk.UnwrapSDKContext(cx)
	reqTx := req.Tx

	params := esvd.evmKeeper.GetParams(ctx)

	ethCfg := params.ChainConfig.EthereumConfig(chainID)
	blockNum := big.NewInt(ctx.BlockHeight())
	signer := ethtypes.MakeSigner(ethCfg, blockNum)

	for _, msg := range reqTx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		sender, err := signer.Sender(msgEthTx.AsTransaction())
		if err != nil {
			return tx.Response{}, sdkerrors.Wrapf(
				sdkerrors.ErrorInvalidSigner,
				"couldn't retrieve sender address ('%s') from the ethereum transaction: %s",
				msgEthTx.From,
				err.Error(),
			)
		}

		// set up the sender to the transaction field if not already
		msgEthTx.From = sender.Hex()
	}

	return tx.Response{}, nil
}

// CheckTx implements tx.Handler
func (esvd EthSigVerificationMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	if _, err := ethSigVerificationMiddleware(esvd, ctx, req); err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	return esvd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (esvd EthSigVerificationMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := ethSigVerificationMiddleware(esvd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return esvd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (esvd EthSigVerificationMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := ethSigVerificationMiddleware(esvd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return esvd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = EthSigVerificationMiddleware{}

// NewEthSigVerificationMiddleware creates a new EthSigVerificationMiddleware
func NewEthSigVerificationMiddleware(ek EVMKeeper) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return EthSigVerificationMiddleware{
			next:      h,
			evmKeeper: ek,
		}
	}
}

// EthAccountVerificationMiddleware validates an account balance checks
type EthAccountVerificationMiddleware struct {
	next       tx.Handler
	ak         evmtypes.AccountKeeper
	bankKeeper evmtypes.BankKeeper
	evmKeeper  EVMKeeper
}

// CheckTx implements tx.Handler
func (avd EthAccountVerificationMiddleware) CheckTx(cx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	reqTx := req.Tx
	ctx := sdk.UnwrapSDKContext(cx)
	if !ctx.IsCheckTx() {
		return avd.next.CheckTx(cx, req, checkReq)
	}

	for i, msg := range reqTx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		txData, err := evmtypes.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(err, "failed to unpack tx data any for tx %d", i)
		}

		// sender address should be in the tx cache from the previous AnteHandle call
		from := msgEthTx.GetFrom()
		if from.Empty() {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from address cannot be empty")
		}

		// check whether the sender address is EOA
		fromAddr := common.BytesToAddress(from)
		acct := avd.evmKeeper.GetAccount(ctx, fromAddr)

		if acct == nil {
			acc := avd.ak.NewAccountWithAddress(ctx, from)
			avd.ak.SetAccount(ctx, acc)
			acct = statedb.NewEmptyAccount()
		} else if acct.IsContract() {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType,
				"the sender is not EOA: address %s, codeHash <%s>", fromAddr, acct.CodeHash)
		}

		if err := evmkeeper.CheckSenderBalance(sdk.NewIntFromBigInt(acct.Balance), txData); err != nil {
			return tx.Response{}, tx.ResponseCheckTx{}, sdkerrors.Wrap(err, "failed to check sender balance")
		}
	}
	return avd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (avd EthAccountVerificationMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return avd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (avd EthAccountVerificationMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	return avd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = EthAccountVerificationMiddleware{}

// NewEthAccountVerificationMiddleware creates a new EthAccountVerificationMiddleware
func NewEthAccountVerificationMiddleware(ak evmtypes.AccountKeeper, bankKeeper evmtypes.BankKeeper, ek EVMKeeper) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return EthAccountVerificationMiddleware{
			next:       h,
			ak:         ak,
			bankKeeper: bankKeeper,
			evmKeeper:  ek,
		}
	}
}

// EthGasConsumeMiddleware validates enough intrinsic gas for the transaction and
// gas consumption.
type EthGasConsumeMiddleware struct {
	next         tx.Handler
	evmKeeper    EVMKeeper
	maxGasWanted uint64
}

func ethGasMiddleware(egcd EthGasConsumeMiddleware, cx context.Context, req tx.Request) (context.Context, tx.Response, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	reqTx := req.Tx

	params := egcd.evmKeeper.GetParams(ctx)

	ethCfg := params.ChainConfig.EthereumConfig(egcd.evmKeeper.ChainID())

	blockHeight := big.NewInt(ctx.BlockHeight())
	homestead := ethCfg.IsHomestead(blockHeight)
	istanbul := ethCfg.IsIstanbul(blockHeight)
	london := ethCfg.IsLondon(blockHeight)
	evmDenom := params.EvmDenom
	gasWanted := uint64(0)
	var events sdk.Events

	for _, msg := range reqTx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		txData, err := evmtypes.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return ctx, tx.Response{}, sdkerrors.Wrap(err, "failed to unpack tx data")
		}

		if ctx.IsCheckTx() {
			// We can't trust the tx gas limit, because we'll refund the unused gas.
			if txData.GetGas() > egcd.maxGasWanted {
				gasWanted += egcd.maxGasWanted
			} else {
				gasWanted += txData.GetGas()
			}
		} else {
			gasWanted += txData.GetGas()
		}

		fees, err := egcd.evmKeeper.DeductTxCostsFromUserBalance(
			ctx,
			*msgEthTx,
			txData,
			evmDenom,
			homestead,
			istanbul,
			london,
		)
		if err != nil {
			return ctx, tx.Response{}, sdkerrors.Wrapf(err, "failed to deduct transaction costs from user balance")
		}

		events = append(events, sdk.NewEvent(sdk.EventTypeTx, sdk.NewAttribute(sdk.AttributeKeyFee, fees.String())))
	}

	// TODO: change to typed events
	ctx.EventManager().EmitEvents(events)

	// TODO: deprecate after https://github.com/cosmos/cosmos-sdk/issues/9514  is fixed on SDK
	blockGasLimit := ethermint.BlockGasLimit(ctx)

	// NOTE: safety check
	if blockGasLimit > 0 {
		// generate a copy of the gas pool (i.e block gas meter) to see if we've run out of gas for this block
		// if current gas consumed is greater than the limit, this funcion panics and the error is recovered on the Baseapp
		gasPool := sdk.NewGasMeter(blockGasLimit)
		gasPool.ConsumeGas(ctx.GasMeter().GasConsumedToLimit(), "gas pool check")
	}

	// Set ctx.GasMeter with a limit of GasWanted (gasLimit)
	gasConsumed := ctx.GasMeter().GasConsumed()
	ctx = ctx.WithGasMeter(ethermint.NewInfiniteGasMeterWithLimit(gasWanted))
	ctx.GasMeter().ConsumeGas(gasConsumed, "copy gas consumed")

	return ctx, tx.Response{}, nil
}

// CheckTx implements tx.Handler
func (egcd EthGasConsumeMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	newCtx, _, err := ethGasMiddleware(egcd, ctx, req)
	if err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	return egcd.next.CheckTx(newCtx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (egcd EthGasConsumeMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	newCtx, _, err := ethGasMiddleware(egcd, ctx, req)
	if err != nil {
		return tx.Response{}, err
	}

	return egcd.next.DeliverTx(newCtx, req)
}

// SimulateTx implements tx.Handler
func (egcd EthGasConsumeMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	newCtx, _, err := ethGasMiddleware(egcd, ctx, req)
	if err != nil {
		return tx.Response{}, err
	}

	return egcd.next.SimulateTx(newCtx, req)
}

var _ tx.Handler = EthGasConsumeMiddleware{}

// NewEthGasConsumeMiddleware creates a new EthGasConsumeMiddleware
func NewEthGasConsumeMiddleware(
	evmKeeper EVMKeeper,
	maxGasWanted uint64,
) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return EthGasConsumeMiddleware{
			h,
			evmKeeper,
			maxGasWanted,
		}
	}
}

// CanTransferMiddleware checks if the sender is allowed to transfer funds according to the EVM block
// context rules.
type CanTransferMiddleware struct {
	next      tx.Handler
	evmKeeper EVMKeeper
}

func canTransfer(ctd CanTransferMiddleware, cx context.Context, req tx.Request) (tx.Response, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	reqTx := req.Tx
	params := ctd.evmKeeper.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(ctd.evmKeeper.ChainID())
	signer := ethtypes.MakeSigner(ethCfg, big.NewInt(ctx.BlockHeight()))

	for _, msg := range reqTx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		baseFee := ctd.evmKeeper.BaseFee(ctx, ethCfg)

		coreMsg, err := msgEthTx.AsMessage(signer, baseFee)
		if err != nil {
			return tx.Response{}, sdkerrors.Wrapf(
				err,
				"failed to create an ethereum core.Message from signer %T", signer,
			)
		}

		// NOTE: pass in an empty coinbase address and nil tracer as we don't need them for the check below
		cfg := &evmtypes.EVMConfig{
			ChainConfig: ethCfg,
			Params:      params,
			CoinBase:    common.Address{},
			BaseFee:     baseFee,
		}
		stateDB := statedb.New(ctx, ctd.evmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes())))
		evm := ctd.evmKeeper.NewEVM(ctx, coreMsg, cfg, evmtypes.NewNoOpTracer(), stateDB)

		// check that caller has enough balance to cover asset transfer for **topmost** call
		// NOTE: here the gas consumed is from the context with the infinite gas meter
		if coreMsg.Value().Sign() > 0 && !evm.Context.CanTransfer(stateDB, coreMsg.From(), coreMsg.Value()) {
			return tx.Response{}, sdkerrors.Wrapf(
				sdkerrors.ErrInsufficientFunds,
				"failed to transfer %s from address %s using the EVM block context transfer function",
				coreMsg.Value(),
				coreMsg.From(),
			)
		}

		if evmtypes.IsLondon(ethCfg, ctx.BlockHeight()) {
			if baseFee == nil {
				return tx.Response{}, sdkerrors.Wrap(
					evmtypes.ErrInvalidBaseFee,
					"base fee is supported but evm block context value is nil",
				)
			}
			if coreMsg.GasFeeCap().Cmp(baseFee) < 0 {
				return tx.Response{}, sdkerrors.Wrapf(
					sdkerrors.ErrInsufficientFee,
					"max fee per gas less than block base fee (%s < %s)",
					coreMsg.GasFeeCap(), baseFee,
				)
			}
		}
	}

	return tx.Response{}, nil
}

// CheckTx implements tx.Handler
func (ctd CanTransferMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	if _, err := canTransfer(ctd, ctx, req); err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	return ctd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (ctd CanTransferMiddleware) DeliverTx(cx context.Context, req tx.Request) (tx.Response, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	if _, err := canTransfer(ctd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return ctd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (ctd CanTransferMiddleware) SimulateTx(cx context.Context, req tx.Request) (tx.Response, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	if _, err := canTransfer(ctd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return ctd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = CanTransferMiddleware{}

// NewCanTransferMiddleware creates a new CanTransferMiddleware instance.
func NewCanTransferMiddleware(evmKeeper EVMKeeper) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return CanTransferMiddleware{
			next:      h,
			evmKeeper: evmKeeper,
		}
	}
}

// EthIncrementSenderSequenceMiddleware increments the sequence of the signers.
type EthIncrementSenderSequenceMiddleware struct {
	next tx.Handler
	ak   evmtypes.AccountKeeper
}

func ethIncrementSenderSequence(issd EthIncrementSenderSequenceMiddleware, cx context.Context, req tx.Request) (tx.Response, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	reqTx := req.Tx

	for _, msg := range reqTx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		txData, err := evmtypes.UnpackTxData(msgEthTx.Data)
		if err != nil {
			return tx.Response{}, sdkerrors.Wrap(err, "failed to unpack tx data")
		}

		// increase sequence of sender
		acc := issd.ak.GetAccount(ctx, msgEthTx.GetFrom())
		if acc == nil {
			return tx.Response{}, sdkerrors.Wrapf(
				sdkerrors.ErrUnknownAddress,
				"account %s is nil", common.BytesToAddress(msgEthTx.GetFrom().Bytes()),
			)
		}
		nonce := acc.GetSequence()

		// we merged the nonce verification to nonce increment, so when tx includes multiple messages
		// with same sender, they'll be accepted.
		if txData.GetNonce() != nonce {
			return tx.Response{}, sdkerrors.Wrapf(
				sdkerrors.ErrInvalidSequence,
				"invalid nonce; got %d, expected %d", txData.GetNonce(), nonce,
			)
		}

		if err := acc.SetSequence(nonce + 1); err != nil {
			return tx.Response{}, sdkerrors.Wrapf(err, "failed to set sequence to %d", acc.GetSequence()+1)
		}

		issd.ak.SetAccount(ctx, acc)
	}

	return tx.Response{}, nil
}

// CheckTx implements tx.Handler
func (issd EthIncrementSenderSequenceMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	if _, err := ethIncrementSenderSequence(issd, ctx, req); err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	return issd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (issd EthIncrementSenderSequenceMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := ethIncrementSenderSequence(issd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return issd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (issd EthIncrementSenderSequenceMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := ethIncrementSenderSequence(issd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return issd.next.SimulateTx(ctx, req)
}

var _ tx.Handler = EthIncrementSenderSequenceMiddleware{}

// NewEthIncrementSenderSequenceMiddleware creates a new EthIncrementSenderSequenceMiddleware.
func NewEthIncrementSenderSequenceMiddleware(ak evmtypes.AccountKeeper) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return EthIncrementSenderSequenceMiddleware{
			next: h,
			ak:   ak,
		}
	}
}
