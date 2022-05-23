package middleware_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/middleware"
	ante "github.com/tharsis/ethermint/app/middleware"
	"github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/tests"
	"github.com/tharsis/ethermint/x/evm/statedb"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func (suite MiddlewareTestSuite) TestEthSigVerificationDecorator() {
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewEthSigVerificationMiddleware(suite.app.EvmKeeper))
	addr, privKey := tests.NewAddrKey()

	signedTx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	signedTx.From = addr.Hex()
	err := signedTx.Sign(suite.ethSigner, tests.NewSigner(privKey))
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		tx        sdk.Tx
		reCheckTx bool
		expPass   bool
	}{
		{"ReCheckTx", &invalidTx{}, true, false},
		{"invalid transaction type", &invalidTx{}, false, false},
		{
			"invalid sender",
			evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 1, &addr, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			false,
			false,
		},
		{"successful signature verification", signedTx, false, true},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite MiddlewareTestSuite) TestNewEthAccountVerificationDecorator() {
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewEthAccountVerificationMiddleware(
		suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.EvmKeeper,
	))

	addr := tests.GenerateAddress()

	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		tx       sdk.Tx
		malleate func()
		checkTx  bool
		expPass  bool
	}{
		{"not CheckTx", nil, func() {}, false, true},
		{"invalid transaction type", &invalidTx{}, func() {}, true, false},
		{
			"sender not set to msg",
			evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			func() {},
			true,
			false,
		},
		{
			"sender not EOA",
			tx,
			func() {
				// set not as an EOA
				vmdb.SetCode(addr, []byte("1"))
			},
			true,
			false,
		},
		{
			"not enough balance to cover tx cost",
			tx,
			func() {
				// reset back to EOA
				vmdb.SetCode(addr, nil)
			},
			true,
			false,
		},
		{
			"success new account",
			tx,
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			true,
			true,
		},
		{
			"success existing account",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			true,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vmdb = suite.StateDB()
			tc.malleate()
			suite.Require().NoError(vmdb.Commit())

			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx.WithIsCheckTx(tc.checkTx)), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite MiddlewareTestSuite) TestEthNonceVerificationDecorator() {
	suite.SetupTest()
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewEthIncrementSenderSequenceMiddleware(suite.app.AccountKeeper))
	addr := tests.GenerateAddress()

	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()

	testCases := []struct {
		name      string
		tx        sdk.Tx
		malleate  func()
		reCheckTx bool
		expPass   bool
	}{
		{"ReCheckTx", &invalidTx{}, func() {}, true, false},
		{"invalid transaction type", &invalidTx{}, func() {}, false, false},
		{"sender account not found", tx, func() {}, false, false},
		{
			"sender nonce missmatch",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			false,
			false,
		},
		{
			"success",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.Require().NoError(acc.SetSequence(1))
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			false,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()
			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite MiddlewareTestSuite) TestEthGasConsumeDecorator() {
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewEthGasConsumeMiddleware(suite.app.EvmKeeper, config.DefaultMaxTxGasWanted))

	addr := tests.GenerateAddress()

	txGasLimit := uint64(1000)
	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), txGasLimit, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()

	tx2GasLimit := uint64(1000000)
	tx2 := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), tx2GasLimit, big.NewInt(1), nil, nil, nil, &ethtypes.AccessList{{Address: addr, StorageKeys: nil}})
	tx2.From = addr.Hex()

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		tx       sdk.Tx
		gasLimit uint64
		malleate func()
		expPass  bool
		expPanic bool
	}{
		{"invalid transaction type", &invalidTx{}, 0, func() {}, false, false},
		{
			"sender not found",
			evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			0,
			func() {},
			false, false,
		},
		{
			"gas limit too low",
			tx,
			0,
			func() {},
			false, false,
		},
		{
			"not enough balance for fees",
			tx2,
			0,
			func() {},
			false, false,
		},
		{
			"not enough tx gas",
			tx2,
			0,
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			false, true,
		},
		{
			"not enough block gas",
			tx2,
			0,
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))

				suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1))
			},
			false, true,
		},
		{
			"success",
			tx2,
			config.DefaultMaxTxGasWanted, // it's capped
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))

				suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(10000000000000000000))
			},
			true, false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vmdb = suite.StateDB()
			tc.malleate()
			suite.Require().NoError(vmdb.Commit())

			if tc.expPanic {
				suite.Require().Panics(func() {
					_, _, _ = txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx.WithIsCheckTx(true).WithGasMeter(sdk.NewGasMeter(1))), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})
				})
				return
			}

			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx.WithIsCheckTx(true).WithGasMeter(sdk.NewInfiniteGasMeter())), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite MiddlewareTestSuite) TestCanTransferDecorator() {
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewCanTransferMiddleware(suite.app.EvmKeeper))

	addr, privKey := tests.NewAddrKey()

	suite.app.FeeMarketKeeper.SetBaseFee(suite.ctx, big.NewInt(100))

	tx := evmtypes.NewTxContract(
		suite.app.EvmKeeper.ChainID(),
		1,
		big.NewInt(10),
		1000,
		big.NewInt(150),
		big.NewInt(200),
		nil,
		nil,
		&ethtypes.AccessList{},
	)
	tx2 := evmtypes.NewTxContract(
		suite.app.EvmKeeper.ChainID(),
		1,
		big.NewInt(10),
		1000,
		big.NewInt(150),
		big.NewInt(200),
		nil,
		nil,
		&ethtypes.AccessList{},
	)

	tx.From = addr.Hex()

	err := tx.Sign(suite.ethSigner, tests.NewSigner(privKey))
	suite.Require().NoError(err)

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		tx       sdk.Tx
		malleate func()
		expPass  bool
	}{
		{"invalid transaction type", &invalidTx{}, func() {}, false},
		{"AsMessage failed", tx2, func() {}, false},
		{
			"evm CanTransfer failed",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			false,
		},
		{
			"success",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vmdb = suite.StateDB()
			tc.malleate()
			suite.Require().NoError(vmdb.Commit())

			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx.WithIsCheckTx(true)), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite MiddlewareTestSuite) TestEthIncrementSenderSequenceDecorator() {
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewEthIncrementSenderSequenceMiddleware(suite.app.AccountKeeper))
	addr, privKey := tests.NewAddrKey()

	contract := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 0, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	contract.From = addr.Hex()
	err := contract.Sign(suite.ethSigner, tests.NewSigner(privKey))
	suite.Require().NoError(err)

	to := tests.GenerateAddress()
	tx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 0, &to, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()
	err = tx.Sign(suite.ethSigner, tests.NewSigner(privKey))
	suite.Require().NoError(err)

	tx2 := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 1, &to, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx2.From = addr.Hex()
	err = tx2.Sign(suite.ethSigner, tests.NewSigner(privKey))
	suite.Require().NoError(err)

	testCases := []struct {
		name     string
		tx       sdk.Tx
		malleate func()
		expPass  bool
		expPanic bool
	}{
		{
			"invalid transaction type",
			&invalidTx{},
			func() {},
			false, false,
		},
		{
			"no signers",
			evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 1, &to, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			func() {},
			false, false,
		},
		{
			"account not set to store",
			tx,
			func() {},
			false, false,
		},
		{
			"success - create contract",
			contract,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			true, false,
		},
		{
			"success - call",
			tx2,
			func() {},
			true, false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()

			if tc.expPanic {
				suite.Require().Panics(func() {
					_, _, _ = txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})
				})
				return
			}

			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

			if tc.expPass {
				suite.Require().NoError(err)
				msg := tc.tx.(*evmtypes.MsgEthereumTx)

				txData, err := evmtypes.UnpackTxData(msg.Data)
				suite.Require().NoError(err)

				nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, addr)
				suite.Require().Equal(txData.GetNonce()+1, nonce)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite MiddlewareTestSuite) TestEthSetupContextDecorator() {
	txHandler := middleware.ComposeMiddlewares(noopTxHandler, ante.NewEthSetUpContextMiddleware(suite.app.EvmKeeper))
	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)

	testCases := []struct {
		name    string
		tx      sdk.Tx
		expPass bool
	}{
		{"invalid transaction type - does not implement GasTx", &invalidTx{}, false},
		{
			"success - transaction implement GasTx",
			tx,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, _, err := txHandler.CheckTx(sdk.WrapSDKContext(suite.ctx), txtypes.Request{Tx: tc.tx}, txtypes.RequestCheckTx{})
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
