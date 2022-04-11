package rest

// clientrest "github.com/cosmos/cosmos-sdk/client/rest"
// "github.com/cosmos/cosmos-sdk/types/rest"
// authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

// // RegisterTxRoutes - Central function to define routes that get registered by the main application
// func RegisterTxRoutes(clientCtx client.Context, rtr *mux.Router) {
// 	r := testutil.GetRequestWithHeaders(rtr)
// 	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(clientCtx)).Methods("GET")
// 	r.HandleFunc("/txs", testutil.QueryTxsRequestHandlerFn(clientCtx)).Methods("GET")
// 	r.HandleFunc("/txs/decode", testutil.DecodeTxRequestHandlerFn(clientCtx)).Methods("POST")
// }

// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
// func QueryTxRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		hashHexStr := vars["hash"]

// 		clientCtx, ok := testutil.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		ethHashPrefix := "0x"
// 		// if !strings.HasPrefix(hashHexStr, ethHashPrefix) {
// 		// 	authrest.QueryTxRequestHandlerFn(clientCtx)
// 		// 	return
// 		// }

// 		// eth Tx
// 		ethHashPrefixLength := len(ethHashPrefix)
// 		output, err := getEthTransactionByHash(clientCtx, hashHexStr[ethHashPrefixLength:])
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		rest.PostProcessResponseBare(w, clientCtx, output)
// 	}
// }

// // GetTransactionByHash returns the transaction identified by hash.
// func getEthTransactionByHash(clientCtx client.Context, hashHex string) ([]byte, error) {
// 	hash, err := hex.DecodeString(hashHex)
// 	if err != nil {
// 		return nil, err
// 	}
// 	node, err := clientCtx.GetNode()
// 	if err != nil {
// 		return nil, err
// 	}

// 	tx, err := node.Tx(context.Background(), hash, false)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Can either cache or just leave this out if not necessary
// 	block, err := node.Block(context.Background(), &tx.Height)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := feemarkettypes.NewQueryClient(clientCtx)
// 	res, err := client.BaseFee(context.Background(), &feemarkettypes.QueryBaseFeeRequest{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var baseFee *big.Int
// 	if res.BaseFee != nil {
// 		baseFee = res.BaseFee.BigInt()
// 	}

// 	blockHash := common.BytesToHash(block.Block.Header.Hash())

// 	ethTxs, err := rpctypes.RawTxToEthTx(clientCtx, tx.Tx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	height := uint64(tx.Height)

// 	for _, ethTx := range ethTxs {
// 		if common.HexToHash(ethTx.Hash) == common.BytesToHash(hash) {
// 			rpcTx, err := rpctypes.NewRPCTransaction(ethTx.AsTransaction(), blockHash, height, uint64(tx.Index), baseFee)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return json.Marshal(rpcTx)
// 		}
// 	}

// 	return nil, errors.New("eth tx not found")
// }
