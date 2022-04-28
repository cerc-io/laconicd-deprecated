package middleware

import (
	context "context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/middleware"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/ethereum/eip712"
	ethermint "github.com/tharsis/ethermint/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

var ethermintCodec codec.ProtoCodecMarshaler

func init() {
	registry := codectypes.NewInterfaceRegistry()
	ethermint.RegisterInterfaces(registry)
	ethermintCodec = codec.NewProtoCodec(registry)
}

// Eip712SigVerificationMiddleware Verify all signatures for a tx and return an error if any are invalid. Note,
// the Eip712SigVerificationMiddleware middleware will not get executed on ReCheck.
//
// CONTRACT: Pubkeys are set in context for all signers before this middleware runs
// CONTRACT: Tx must implement SigVerifiableTx interface
type Eip712SigVerificationMiddleware struct {
	appCodec        codec.Codec
	next            tx.Handler
	ak              evmtypes.AccountKeeper
	signModeHandler authsigning.SignModeHandler
}

var _ tx.Handler = Eip712SigVerificationMiddleware{}

// NewEip712SigVerificationMiddleware creates a new Eip712SigVerificationMiddleware
func NewEip712SigVerificationMiddleware(appCodec codec.Codec, ak evmtypes.AccountKeeper, signModeHandler authsigning.SignModeHandler) tx.Middleware {
	return func(h tx.Handler) tx.Handler {
		return Eip712SigVerificationMiddleware{
			appCodec:        appCodec,
			next:            h,
			ak:              ak,
			signModeHandler: signModeHandler,
		}
	}
}

func eipSigVerification(svd Eip712SigVerificationMiddleware, cx context.Context, req tx.Request) (tx.Response, error) {
	ctx := sdk.UnwrapSDKContext(cx)
	reqTx := req.Tx

	sigTx, ok := reqTx.(authsigning.SigVerifiableTx)
	if !ok {
		return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "tx %T doesn't implement authsigning.SigVerifiableTx", reqTx)
	}

	authSignTx, ok := reqTx.(authsigning.Tx)
	if !ok {
		return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "tx %T doesn't implement the authsigning.Tx interface", reqTx)
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return tx.Response{}, err
	}

	signerAddrs := sigTx.GetSigners()

	// EIP712 allows just one signature
	if len(sigs) != 1 {
		return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "invalid number of signers (%d);  EIP712 signatures allows just one signature", len(sigs))
	}

	// check that signer length and signature length are the same
	if len(sigs) != len(signerAddrs) {
		return tx.Response{}, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
	}

	// EIP712 has just one signature, avoid looping here and only read index 0
	i := 0
	sig := sigs[i]

	acc, err := middleware.GetSignerAcc(ctx, svd.ak, signerAddrs[i])
	if err != nil {
		return tx.Response{}, err
	}

	// retrieve pubkey
	pubKey := acc.GetPubKey()
	if pubKey == nil {
		return tx.Response{}, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey on account is not set")
	}

	// Check account sequence number.
	if sig.Sequence != acc.GetSequence() {
		return tx.Response{}, sdkerrors.Wrapf(
			sdkerrors.ErrWrongSequence,
			"account sequence mismatch, expected %d, got %d", acc.GetSequence(), sig.Sequence,
		)
	}

	// retrieve signer data
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()

	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	signerData := authsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      acc.GetSequence(),
	}

	if err := VerifySignature(svd.appCodec, pubKey, signerData, sig.Data, svd.signModeHandler, authSignTx); err != nil {
		errMsg := fmt.Errorf("signature verification failed; please verify account number (%d) and chain-id (%s): %w", accNum, chainID, err)
		return tx.Response{}, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, errMsg.Error())
	}

	return tx.Response{}, nil
}

// CheckTx implements tx.Handler
func (svd Eip712SigVerificationMiddleware) CheckTx(ctx context.Context, req tx.Request, checkReq tx.RequestCheckTx) (tx.Response, tx.ResponseCheckTx, error) {
	if _, err := eipSigVerification(svd, ctx, req); err != nil {
		return tx.Response{}, tx.ResponseCheckTx{}, err
	}

	return svd.next.CheckTx(ctx, req, checkReq)
}

// DeliverTx implements tx.Handler
func (svd Eip712SigVerificationMiddleware) DeliverTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := eipSigVerification(svd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return svd.next.DeliverTx(ctx, req)
}

// SimulateTx implements tx.Handler
func (svd Eip712SigVerificationMiddleware) SimulateTx(ctx context.Context, req tx.Request) (tx.Response, error) {
	if _, err := eipSigVerification(svd, ctx, req); err != nil {
		return tx.Response{}, err
	}

	return svd.next.SimulateTx(ctx, req)
}

// VerifySignature verifies a transaction signature contained in SignatureData abstracting over different signing modes
// and single vs multi-signatures.
func VerifySignature(
	appCodec codec.Codec,
	pubKey cryptotypes.PubKey,
	signerData authsigning.SignerData,
	sigData signing.SignatureData,
	_ authsigning.SignModeHandler,
	tx authsigning.Tx,
) error {
	switch data := sigData.(type) {
	case *signing.SingleSignatureData:
		if data.SignMode != signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON {
			return sdkerrors.Wrapf(sdkerrors.ErrNotSupported, "unexpected SignatureData %T: wrong SignMode", sigData)
		}

		// Note: this prevents the user from sending thrash data in the signature field
		if len(data.Signature) != 0 {
			return sdkerrors.Wrap(sdkerrors.ErrTooManySignatures, "invalid signature value; EIP712 must have the cosmos transaction signature empty")
		}

		// @contract: this code is reached only when Msg has Web3Tx extension (so this custom Ante handler flow),
		// and the signature is SIGN_MODE_LEGACY_AMINO_JSON which is supported for EIP712 for now

		msgs := tx.GetMsgs()
		if len(msgs) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrNoSignatures, "tx doesn't contain any msgs to verify signature")
		}

		txBytes := legacytx.StdSignBytes(
			signerData.ChainID,
			signerData.AccountNumber,
			signerData.Sequence,
			tx.GetTimeoutHeight(),
			legacytx.StdFee{
				Amount: tx.GetFee(),
				Gas:    tx.GetGas(),
			},
			msgs, tx.GetMemo(), tx.GetTip(),
		)

		signerChainID, err := ethermint.ParseChainID(signerData.ChainID)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to parse chainID: %s", signerData.ChainID)
		}

		txWithExtensions, ok := tx.(middleware.HasExtensionOptionsTx)
		if !ok {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownExtensionOptions, "tx doesnt contain any extensions")
		}
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) != 1 {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownExtensionOptions, "tx doesnt contain expected amount of extension options")
		}

		var optIface ethermint.ExtensionOptionsWeb3TxI

		if err := appCodec.UnpackAny(opts[0], &optIface); err != nil {
			return sdkerrors.Wrap(err, "failed to proto-unpack ExtensionOptionsWeb3Tx")
		}

		extOpt, ok := optIface.(*ethermint.ExtensionOptionsWeb3Tx)
		if !ok {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidChainID, "unknown extension option")
		}

		if extOpt.TypedDataChainID != signerChainID.Uint64() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidChainID, "invalid chainID")
		}

		if len(extOpt.FeePayer) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownExtensionOptions, "no feePayer on ExtensionOptionsWeb3Tx")
		}
		feePayer, err := sdk.AccAddressFromBech32(extOpt.FeePayer)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to parse feePayer from ExtensionOptionsWeb3Tx")
		}

		feeDelegation := &eip712.FeeDelegationOptions{
			FeePayer: feePayer,
		}

		typedData, err := eip712.WrapTxToTypedData(appCodec, extOpt.TypedDataChainID, msgs[0], txBytes, feeDelegation)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to pack tx data in EIP712 object")
		}

		sigHash, err := eip712.ComputeTypedDataHash(typedData)
		if err != nil {
			return err
		}

		feePayerSig := extOpt.FeePayerSig
		if len(feePayerSig) != ethcrypto.SignatureLength {
			return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "signature length doesn't match typical [R||S||V] signature 65 bytes")
		}

		// Remove the recovery offset if needed (ie. Metamask eip712 signature)
		if feePayerSig[ethcrypto.RecoveryIDOffset] == 27 || feePayerSig[ethcrypto.RecoveryIDOffset] == 28 {
			feePayerSig[ethcrypto.RecoveryIDOffset] -= 27
		}

		feePayerPubkey, err := secp256k1.RecoverPubkey(sigHash, feePayerSig)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to recover delegated fee payer from sig")
		}

		ecPubKey, err := ethcrypto.UnmarshalPubkey(feePayerPubkey)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to unmarshal recovered fee payer pubkey")
		}

		pk := &ethsecp256k1.PubKey{
			Key: ethcrypto.CompressPubkey(ecPubKey),
		}

		if !pubKey.Equals(pk) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "feePayer pubkey %s is different from transaction pubkey %s", pubKey, pk)
		}

		recoveredFeePayerAcc := sdk.AccAddress(pk.Address().Bytes())

		if !recoveredFeePayerAcc.Equals(feePayer) {
			return sdkerrors.Wrapf(sdkerrors.ErrorInvalidSigner, "failed to verify delegated fee payer %s signature", recoveredFeePayerAcc)
		}

		// VerifySignature of ethsecp256k1 accepts 64 byte signature [R||S]
		// WARNING! Under NO CIRCUMSTANCES try to use pubKey.VerifySignature there
		if !secp256k1.VerifySignature(pubKey.Bytes(), sigHash, feePayerSig[:len(feePayerSig)-1]) {
			return sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "unable to verify signer signature of EIP712 typed data")
		}

		return nil
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrTooManySignatures, "unexpected SignatureData %T", sigData)
	}
}
