<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [ethermint/crypto/v1/ethsecp256k1/keys.proto](#ethermint/crypto/v1/ethsecp256k1/keys.proto)
    - [PrivKey](#ethermint.crypto.v1.ethsecp256k1.PrivKey)
    - [PubKey](#ethermint.crypto.v1.ethsecp256k1.PubKey)
  
- [ethermint/evm/v1/evm.proto](#ethermint/evm/v1/evm.proto)
    - [AccessTuple](#ethermint.evm.v1.AccessTuple)
    - [ChainConfig](#ethermint.evm.v1.ChainConfig)
    - [Log](#ethermint.evm.v1.Log)
    - [Params](#ethermint.evm.v1.Params)
    - [State](#ethermint.evm.v1.State)
    - [TraceConfig](#ethermint.evm.v1.TraceConfig)
    - [TransactionLogs](#ethermint.evm.v1.TransactionLogs)
    - [TxResult](#ethermint.evm.v1.TxResult)
  
- [ethermint/evm/v1/genesis.proto](#ethermint/evm/v1/genesis.proto)
    - [GenesisAccount](#ethermint.evm.v1.GenesisAccount)
    - [GenesisState](#ethermint.evm.v1.GenesisState)
  
- [ethermint/evm/v1/tx.proto](#ethermint/evm/v1/tx.proto)
    - [AccessListTx](#ethermint.evm.v1.AccessListTx)
    - [DynamicFeeTx](#ethermint.evm.v1.DynamicFeeTx)
    - [ExtensionOptionsEthereumTx](#ethermint.evm.v1.ExtensionOptionsEthereumTx)
    - [LegacyTx](#ethermint.evm.v1.LegacyTx)
    - [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx)
    - [MsgEthereumTxResponse](#ethermint.evm.v1.MsgEthereumTxResponse)
  
    - [Msg](#ethermint.evm.v1.Msg)
  
- [ethermint/evm/v1/query.proto](#ethermint/evm/v1/query.proto)
    - [EstimateGasResponse](#ethermint.evm.v1.EstimateGasResponse)
    - [EthCallRequest](#ethermint.evm.v1.EthCallRequest)
    - [QueryAccountRequest](#ethermint.evm.v1.QueryAccountRequest)
    - [QueryAccountResponse](#ethermint.evm.v1.QueryAccountResponse)
    - [QueryBalanceRequest](#ethermint.evm.v1.QueryBalanceRequest)
    - [QueryBalanceResponse](#ethermint.evm.v1.QueryBalanceResponse)
    - [QueryBaseFeeRequest](#ethermint.evm.v1.QueryBaseFeeRequest)
    - [QueryBaseFeeResponse](#ethermint.evm.v1.QueryBaseFeeResponse)
    - [QueryCodeRequest](#ethermint.evm.v1.QueryCodeRequest)
    - [QueryCodeResponse](#ethermint.evm.v1.QueryCodeResponse)
    - [QueryCosmosAccountRequest](#ethermint.evm.v1.QueryCosmosAccountRequest)
    - [QueryCosmosAccountResponse](#ethermint.evm.v1.QueryCosmosAccountResponse)
    - [QueryParamsRequest](#ethermint.evm.v1.QueryParamsRequest)
    - [QueryParamsResponse](#ethermint.evm.v1.QueryParamsResponse)
    - [QueryStorageRequest](#ethermint.evm.v1.QueryStorageRequest)
    - [QueryStorageResponse](#ethermint.evm.v1.QueryStorageResponse)
    - [QueryTraceBlockRequest](#ethermint.evm.v1.QueryTraceBlockRequest)
    - [QueryTraceBlockResponse](#ethermint.evm.v1.QueryTraceBlockResponse)
    - [QueryTraceTxRequest](#ethermint.evm.v1.QueryTraceTxRequest)
    - [QueryTraceTxResponse](#ethermint.evm.v1.QueryTraceTxResponse)
    - [QueryTxLogsRequest](#ethermint.evm.v1.QueryTxLogsRequest)
    - [QueryTxLogsResponse](#ethermint.evm.v1.QueryTxLogsResponse)
    - [QueryValidatorAccountRequest](#ethermint.evm.v1.QueryValidatorAccountRequest)
    - [QueryValidatorAccountResponse](#ethermint.evm.v1.QueryValidatorAccountResponse)
  
    - [Query](#ethermint.evm.v1.Query)
  
- [ethermint/feemarket/v1/feemarket.proto](#ethermint/feemarket/v1/feemarket.proto)
    - [Params](#ethermint.feemarket.v1.Params)
  
- [ethermint/feemarket/v1/genesis.proto](#ethermint/feemarket/v1/genesis.proto)
    - [GenesisState](#ethermint.feemarket.v1.GenesisState)
  
- [ethermint/feemarket/v1/query.proto](#ethermint/feemarket/v1/query.proto)
    - [QueryBaseFeeRequest](#ethermint.feemarket.v1.QueryBaseFeeRequest)
    - [QueryBaseFeeResponse](#ethermint.feemarket.v1.QueryBaseFeeResponse)
    - [QueryBlockGasRequest](#ethermint.feemarket.v1.QueryBlockGasRequest)
    - [QueryBlockGasResponse](#ethermint.feemarket.v1.QueryBlockGasResponse)
    - [QueryParamsRequest](#ethermint.feemarket.v1.QueryParamsRequest)
    - [QueryParamsResponse](#ethermint.feemarket.v1.QueryParamsResponse)
  
    - [Query](#ethermint.feemarket.v1.Query)
  
- [ethermint/types/v1/account.proto](#ethermint/types/v1/account.proto)
    - [EthAccount](#ethermint.types.v1.EthAccount)
  
- [ethermint/types/v1/dynamic_fee.proto](#ethermint/types/v1/dynamic_fee.proto)
    - [ExtensionOptionDynamicFeeTx](#ethermint.types.v1.ExtensionOptionDynamicFeeTx)
  
- [ethermint/types/v1/indexer.proto](#ethermint/types/v1/indexer.proto)
    - [TxResult](#ethermint.types.v1.TxResult)
  
- [ethermint/types/v1/web3.proto](#ethermint/types/v1/web3.proto)
    - [ExtensionOptionsWeb3Tx](#ethermint.types.v1.ExtensionOptionsWeb3Tx)
  
- [vulcanize/auction/v1beta1/types.proto](#vulcanize/auction/v1beta1/types.proto)
    - [Auction](#vulcanize.auction.v1beta1.Auction)
    - [Auctions](#vulcanize.auction.v1beta1.Auctions)
    - [Bid](#vulcanize.auction.v1beta1.Bid)
    - [Params](#vulcanize.auction.v1beta1.Params)
  
- [vulcanize/auction/v1beta1/genesis.proto](#vulcanize/auction/v1beta1/genesis.proto)
    - [GenesisState](#vulcanize.auction.v1beta1.GenesisState)
  
- [vulcanize/auction/v1beta1/query.proto](#vulcanize/auction/v1beta1/query.proto)
    - [AuctionRequest](#vulcanize.auction.v1beta1.AuctionRequest)
    - [AuctionResponse](#vulcanize.auction.v1beta1.AuctionResponse)
    - [AuctionsByBidderRequest](#vulcanize.auction.v1beta1.AuctionsByBidderRequest)
    - [AuctionsByBidderResponse](#vulcanize.auction.v1beta1.AuctionsByBidderResponse)
    - [AuctionsByOwnerRequest](#vulcanize.auction.v1beta1.AuctionsByOwnerRequest)
    - [AuctionsByOwnerResponse](#vulcanize.auction.v1beta1.AuctionsByOwnerResponse)
    - [AuctionsRequest](#vulcanize.auction.v1beta1.AuctionsRequest)
    - [AuctionsResponse](#vulcanize.auction.v1beta1.AuctionsResponse)
    - [BalanceRequest](#vulcanize.auction.v1beta1.BalanceRequest)
    - [BalanceResponse](#vulcanize.auction.v1beta1.BalanceResponse)
    - [BidRequest](#vulcanize.auction.v1beta1.BidRequest)
    - [BidResponse](#vulcanize.auction.v1beta1.BidResponse)
    - [BidsRequest](#vulcanize.auction.v1beta1.BidsRequest)
    - [BidsResponse](#vulcanize.auction.v1beta1.BidsResponse)
    - [QueryParamsRequest](#vulcanize.auction.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#vulcanize.auction.v1beta1.QueryParamsResponse)
  
    - [Query](#vulcanize.auction.v1beta1.Query)
  
- [vulcanize/auction/v1beta1/tx.proto](#vulcanize/auction/v1beta1/tx.proto)
    - [MsgCommitBid](#vulcanize.auction.v1beta1.MsgCommitBid)
    - [MsgCommitBidResponse](#vulcanize.auction.v1beta1.MsgCommitBidResponse)
    - [MsgCreateAuction](#vulcanize.auction.v1beta1.MsgCreateAuction)
    - [MsgCreateAuctionResponse](#vulcanize.auction.v1beta1.MsgCreateAuctionResponse)
    - [MsgRevealBid](#vulcanize.auction.v1beta1.MsgRevealBid)
    - [MsgRevealBidResponse](#vulcanize.auction.v1beta1.MsgRevealBidResponse)
  
    - [Msg](#vulcanize.auction.v1beta1.Msg)
  
- [vulcanize/bond/v1beta1/bond.proto](#vulcanize/bond/v1beta1/bond.proto)
    - [Bond](#vulcanize.bond.v1beta1.Bond)
    - [Params](#vulcanize.bond.v1beta1.Params)
  
- [vulcanize/bond/v1beta1/genesis.proto](#vulcanize/bond/v1beta1/genesis.proto)
    - [GenesisState](#vulcanize.bond.v1beta1.GenesisState)
  
- [vulcanize/bond/v1beta1/query.proto](#vulcanize/bond/v1beta1/query.proto)
    - [QueryGetBondByIDRequest](#vulcanize.bond.v1beta1.QueryGetBondByIDRequest)
    - [QueryGetBondByIDResponse](#vulcanize.bond.v1beta1.QueryGetBondByIDResponse)
    - [QueryGetBondModuleBalanceRequest](#vulcanize.bond.v1beta1.QueryGetBondModuleBalanceRequest)
    - [QueryGetBondModuleBalanceResponse](#vulcanize.bond.v1beta1.QueryGetBondModuleBalanceResponse)
    - [QueryGetBondsByOwnerRequest](#vulcanize.bond.v1beta1.QueryGetBondsByOwnerRequest)
    - [QueryGetBondsByOwnerResponse](#vulcanize.bond.v1beta1.QueryGetBondsByOwnerResponse)
    - [QueryGetBondsRequest](#vulcanize.bond.v1beta1.QueryGetBondsRequest)
    - [QueryGetBondsResponse](#vulcanize.bond.v1beta1.QueryGetBondsResponse)
    - [QueryParamsRequest](#vulcanize.bond.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#vulcanize.bond.v1beta1.QueryParamsResponse)
  
    - [Query](#vulcanize.bond.v1beta1.Query)
  
- [vulcanize/bond/v1beta1/tx.proto](#vulcanize/bond/v1beta1/tx.proto)
    - [MsgCancelBond](#vulcanize.bond.v1beta1.MsgCancelBond)
    - [MsgCancelBondResponse](#vulcanize.bond.v1beta1.MsgCancelBondResponse)
    - [MsgCreateBond](#vulcanize.bond.v1beta1.MsgCreateBond)
    - [MsgCreateBondResponse](#vulcanize.bond.v1beta1.MsgCreateBondResponse)
    - [MsgRefillBond](#vulcanize.bond.v1beta1.MsgRefillBond)
    - [MsgRefillBondResponse](#vulcanize.bond.v1beta1.MsgRefillBondResponse)
    - [MsgWithdrawBond](#vulcanize.bond.v1beta1.MsgWithdrawBond)
    - [MsgWithdrawBondResponse](#vulcanize.bond.v1beta1.MsgWithdrawBondResponse)
  
    - [Msg](#vulcanize.bond.v1beta1.Msg)
  
- [vulcanize/nameservice/v1beta1/attributes.proto](#vulcanize/nameservice/v1beta1/attributes.proto)
    - [ServiceProviderRegistration](#vulcanize.nameservice.v1beta1.ServiceProviderRegistration)
    - [WebsiteRegistrationRecord](#vulcanize.nameservice.v1beta1.WebsiteRegistrationRecord)
    - [X500](#vulcanize.nameservice.v1beta1.X500)
  
- [vulcanize/nameservice/v1beta1/nameservice.proto](#vulcanize/nameservice/v1beta1/nameservice.proto)
    - [AuctionBidInfo](#vulcanize.nameservice.v1beta1.AuctionBidInfo)
    - [AuthorityEntry](#vulcanize.nameservice.v1beta1.AuthorityEntry)
    - [BlockChangeSet](#vulcanize.nameservice.v1beta1.BlockChangeSet)
    - [NameAuthority](#vulcanize.nameservice.v1beta1.NameAuthority)
    - [NameEntry](#vulcanize.nameservice.v1beta1.NameEntry)
    - [NameRecord](#vulcanize.nameservice.v1beta1.NameRecord)
    - [NameRecordEntry](#vulcanize.nameservice.v1beta1.NameRecordEntry)
    - [Params](#vulcanize.nameservice.v1beta1.Params)
    - [Record](#vulcanize.nameservice.v1beta1.Record)
    - [Signature](#vulcanize.nameservice.v1beta1.Signature)
  
- [vulcanize/nameservice/v1beta1/genesis.proto](#vulcanize/nameservice/v1beta1/genesis.proto)
    - [GenesisState](#vulcanize.nameservice.v1beta1.GenesisState)
  
- [vulcanize/nameservice/v1beta1/query.proto](#vulcanize/nameservice/v1beta1/query.proto)
    - [AccountBalance](#vulcanize.nameservice.v1beta1.AccountBalance)
    - [ExpiryQueueRecord](#vulcanize.nameservice.v1beta1.ExpiryQueueRecord)
    - [GetNameServiceModuleBalanceRequest](#vulcanize.nameservice.v1beta1.GetNameServiceModuleBalanceRequest)
    - [GetNameServiceModuleBalanceResponse](#vulcanize.nameservice.v1beta1.GetNameServiceModuleBalanceResponse)
    - [QueryGetAuthorityExpiryQueue](#vulcanize.nameservice.v1beta1.QueryGetAuthorityExpiryQueue)
    - [QueryGetAuthorityExpiryQueueResponse](#vulcanize.nameservice.v1beta1.QueryGetAuthorityExpiryQueueResponse)
    - [QueryGetRecordExpiryQueue](#vulcanize.nameservice.v1beta1.QueryGetRecordExpiryQueue)
    - [QueryGetRecordExpiryQueueResponse](#vulcanize.nameservice.v1beta1.QueryGetRecordExpiryQueueResponse)
    - [QueryListNameRecordsRequest](#vulcanize.nameservice.v1beta1.QueryListNameRecordsRequest)
    - [QueryListNameRecordsResponse](#vulcanize.nameservice.v1beta1.QueryListNameRecordsResponse)
    - [QueryListRecordsRequest](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest)
    - [QueryListRecordsRequest.KeyValueInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.KeyValueInput)
    - [QueryListRecordsRequest.ReferenceInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ReferenceInput)
    - [QueryListRecordsRequest.ValueInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ValueInput)
    - [QueryListRecordsResponse](#vulcanize.nameservice.v1beta1.QueryListRecordsResponse)
    - [QueryLookupCrn](#vulcanize.nameservice.v1beta1.QueryLookupCrn)
    - [QueryLookupCrnResponse](#vulcanize.nameservice.v1beta1.QueryLookupCrnResponse)
    - [QueryParamsRequest](#vulcanize.nameservice.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#vulcanize.nameservice.v1beta1.QueryParamsResponse)
    - [QueryRecordByBondIDRequest](#vulcanize.nameservice.v1beta1.QueryRecordByBondIDRequest)
    - [QueryRecordByBondIDResponse](#vulcanize.nameservice.v1beta1.QueryRecordByBondIDResponse)
    - [QueryRecordByIDRequest](#vulcanize.nameservice.v1beta1.QueryRecordByIDRequest)
    - [QueryRecordByIDResponse](#vulcanize.nameservice.v1beta1.QueryRecordByIDResponse)
    - [QueryResolveCrn](#vulcanize.nameservice.v1beta1.QueryResolveCrn)
    - [QueryResolveCrnResponse](#vulcanize.nameservice.v1beta1.QueryResolveCrnResponse)
    - [QueryWhoisRequest](#vulcanize.nameservice.v1beta1.QueryWhoisRequest)
    - [QueryWhoisResponse](#vulcanize.nameservice.v1beta1.QueryWhoisResponse)
  
    - [Query](#vulcanize.nameservice.v1beta1.Query)
  
- [vulcanize/nameservice/v1beta1/tx.proto](#vulcanize/nameservice/v1beta1/tx.proto)
    - [MsgAssociateBond](#vulcanize.nameservice.v1beta1.MsgAssociateBond)
    - [MsgAssociateBondResponse](#vulcanize.nameservice.v1beta1.MsgAssociateBondResponse)
    - [MsgDeleteNameAuthority](#vulcanize.nameservice.v1beta1.MsgDeleteNameAuthority)
    - [MsgDeleteNameAuthorityResponse](#vulcanize.nameservice.v1beta1.MsgDeleteNameAuthorityResponse)
    - [MsgDissociateBond](#vulcanize.nameservice.v1beta1.MsgDissociateBond)
    - [MsgDissociateBondResponse](#vulcanize.nameservice.v1beta1.MsgDissociateBondResponse)
    - [MsgDissociateRecords](#vulcanize.nameservice.v1beta1.MsgDissociateRecords)
    - [MsgDissociateRecordsResponse](#vulcanize.nameservice.v1beta1.MsgDissociateRecordsResponse)
    - [MsgReAssociateRecords](#vulcanize.nameservice.v1beta1.MsgReAssociateRecords)
    - [MsgReAssociateRecordsResponse](#vulcanize.nameservice.v1beta1.MsgReAssociateRecordsResponse)
    - [MsgRenewRecord](#vulcanize.nameservice.v1beta1.MsgRenewRecord)
    - [MsgRenewRecordResponse](#vulcanize.nameservice.v1beta1.MsgRenewRecordResponse)
    - [MsgReserveAuthority](#vulcanize.nameservice.v1beta1.MsgReserveAuthority)
    - [MsgReserveAuthorityResponse](#vulcanize.nameservice.v1beta1.MsgReserveAuthorityResponse)
    - [MsgSetAuthorityBond](#vulcanize.nameservice.v1beta1.MsgSetAuthorityBond)
    - [MsgSetAuthorityBondResponse](#vulcanize.nameservice.v1beta1.MsgSetAuthorityBondResponse)
    - [MsgSetName](#vulcanize.nameservice.v1beta1.MsgSetName)
    - [MsgSetNameResponse](#vulcanize.nameservice.v1beta1.MsgSetNameResponse)
    - [MsgSetRecord](#vulcanize.nameservice.v1beta1.MsgSetRecord)
    - [MsgSetRecordResponse](#vulcanize.nameservice.v1beta1.MsgSetRecordResponse)
    - [Payload](#vulcanize.nameservice.v1beta1.Payload)
  
    - [Msg](#vulcanize.nameservice.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="ethermint/crypto/v1/ethsecp256k1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/crypto/v1/ethsecp256k1/keys.proto



<a name="ethermint.crypto.v1.ethsecp256k1.PrivKey"></a>

### PrivKey
PrivKey defines a type alias for an ecdsa.PrivateKey that implements
Tendermint's PrivateKey interface.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="ethermint.crypto.v1.ethsecp256k1.PubKey"></a>

### PubKey
PubKey defines a type alias for an ecdsa.PublicKey that implements
Tendermint's PubKey interface. It represents the 33-byte compressed public
key format.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/evm.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/evm.proto



<a name="ethermint.evm.v1.AccessTuple"></a>

### AccessTuple
AccessTuple is the element type of an access list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | hex formatted ethereum address |
| `storage_keys` | [string](#string) | repeated | hex formatted hashes of the storage keys |






<a name="ethermint.evm.v1.ChainConfig"></a>

### ChainConfig
ChainConfig defines the Ethereum ChainConfig parameters using *sdk.Int values
instead of *big.Int.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `homestead_block` | [string](#string) |  | Homestead switch block (nil no fork, 0 = already homestead) |
| `dao_fork_block` | [string](#string) |  | TheDAO hard-fork switch block (nil no fork) |
| `dao_fork_support` | [bool](#bool) |  | Whether the nodes supports or opposes the DAO hard-fork |
| `eip150_block` | [string](#string) |  | EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150) EIP150 HF block (nil no fork) |
| `eip150_hash` | [string](#string) |  | EIP150 HF hash (needed for header only clients as only gas pricing changed) |
| `eip155_block` | [string](#string) |  | EIP155Block HF block |
| `eip158_block` | [string](#string) |  | EIP158 HF block |
| `byzantium_block` | [string](#string) |  | Byzantium switch block (nil no fork, 0 = already on byzantium) |
| `constantinople_block` | [string](#string) |  | Constantinople switch block (nil no fork, 0 = already activated) |
| `petersburg_block` | [string](#string) |  | Petersburg switch block (nil same as Constantinople) |
| `istanbul_block` | [string](#string) |  | Istanbul switch block (nil no fork, 0 = already on istanbul) |
| `muir_glacier_block` | [string](#string) |  | Eip-2384 (bomb delay) switch block (nil no fork, 0 = already activated) |
| `berlin_block` | [string](#string) |  | Berlin switch block (nil = no fork, 0 = already on berlin) |
| `london_block` | [string](#string) |  | London switch block (nil = no fork, 0 = already on london) |
| `arrow_glacier_block` | [string](#string) |  | Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated) |
| `gray_glacier_block` | [string](#string) |  | EIP-5133 (bomb delay) switch block (nil = no fork, 0 = already activated) |
| `merge_netsplit_block` | [string](#string) |  | Virtual fork after The Merge to use as a network splitter |






<a name="ethermint.evm.v1.Log"></a>

### Log
Log represents an protobuf compatible Ethereum Log that defines a contract
log event. These events are generated by the LOG opcode and stored/indexed by
the node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address of the contract that generated the event |
| `topics` | [string](#string) | repeated | list of topics provided by the contract. |
| `data` | [bytes](#bytes) |  | supplied by the contract, usually ABI-encoded |
| `block_number` | [uint64](#uint64) |  | block in which the transaction was included |
| `tx_hash` | [string](#string) |  | hash of the transaction |
| `tx_index` | [uint64](#uint64) |  | index of the transaction in the block |
| `block_hash` | [string](#string) |  | hash of the block in which the transaction was included |
| `index` | [uint64](#uint64) |  | index of the log in the block |
| `removed` | [bool](#bool) |  | The Removed field is true if this log was reverted due to a chain reorganisation. You must pay attention to this field if you receive logs through a filter query. |






<a name="ethermint.evm.v1.Params"></a>

### Params
Params defines the EVM module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evm_denom` | [string](#string) |  | evm denom represents the token denomination used to run the EVM state transitions. |
| `enable_create` | [bool](#bool) |  | enable create toggles state transitions that use the vm.Create function |
| `enable_call` | [bool](#bool) |  | enable call toggles state transitions that use the vm.Call function |
| `extra_eips` | [int64](#int64) | repeated | extra eips defines the additional EIPs for the vm.Config |
| `chain_config` | [ChainConfig](#ethermint.evm.v1.ChainConfig) |  | chain config defines the EVM chain configuration parameters |
| `allow_unprotected_txs` | [bool](#bool) |  | Allow unprotected transactions defines if replay-protected (i.e non EIP155 signed) transactions can be executed on the state machine. |






<a name="ethermint.evm.v1.State"></a>

### State
State represents a single Storage key value pair item.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="ethermint.evm.v1.TraceConfig"></a>

### TraceConfig
TraceConfig holds extra parameters to trace functions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tracer` | [string](#string) |  | custom javascript tracer |
| `timeout` | [string](#string) |  | overrides the default timeout of 5 seconds for JavaScript-based tracing calls |
| `reexec` | [uint64](#uint64) |  | number of blocks the tracer is willing to go back |
| `disable_stack` | [bool](#bool) |  | disable stack capture |
| `disable_storage` | [bool](#bool) |  | disable storage capture |
| `debug` | [bool](#bool) |  | print output during capture end |
| `limit` | [int32](#int32) |  | maximum length of output, but zero means unlimited |
| `overrides` | [ChainConfig](#ethermint.evm.v1.ChainConfig) |  | Chain overrides, can be used to execute a trace using future fork rules |
| `enable_memory` | [bool](#bool) |  | enable memory capture |
| `enable_return_data` | [bool](#bool) |  | enable return data capture |






<a name="ethermint.evm.v1.TransactionLogs"></a>

### TransactionLogs
TransactionLogs define the logs generated from a transaction execution
with a given hash. It it used for import/export data as transactions are not
persisted on blockchain state after an upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  |  |
| `logs` | [Log](#ethermint.evm.v1.Log) | repeated |  |






<a name="ethermint.evm.v1.TxResult"></a>

### TxResult
TxResult stores results of Tx execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | contract_address contains the ethereum address of the created contract (if any). If the state transition is an evm.Call, the contract address will be empty. |
| `bloom` | [bytes](#bytes) |  | bloom represents the bloom filter bytes |
| `tx_logs` | [TransactionLogs](#ethermint.evm.v1.TransactionLogs) |  | tx_logs contains the transaction hash and the proto-compatible ethereum logs. |
| `ret` | [bytes](#bytes) |  | ret defines the bytes from the execution. |
| `reverted` | [bool](#bool) |  | reverted flag is set to true when the call has been reverted |
| `gas_used` | [uint64](#uint64) |  | gas_used notes the amount of gas consumed while execution |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/genesis.proto



<a name="ethermint.evm.v1.GenesisAccount"></a>

### GenesisAccount
GenesisAccount defines an account to be initialized in the genesis state.
Its main difference between with Geth's GenesisAccount is that it uses a
custom storage type and that it doesn't contain the private key field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address defines an ethereum hex formated address of an account |
| `code` | [string](#string) |  | code defines the hex bytes of the account code. |
| `storage` | [State](#ethermint.evm.v1.State) | repeated | storage defines the set of state key values for the account. |






<a name="ethermint.evm.v1.GenesisState"></a>

### GenesisState
GenesisState defines the evm module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [GenesisAccount](#ethermint.evm.v1.GenesisAccount) | repeated | accounts is an array containing the ethereum genesis accounts. |
| `params` | [Params](#ethermint.evm.v1.Params) |  | params defines all the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/tx.proto



<a name="ethermint.evm.v1.AccessListTx"></a>

### AccessListTx
AccessListTx is the data of EIP-2930 access list transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | destination EVM chain ID |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_price` | [string](#string) |  | gas price defines the value for each gas unit |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the unsigned integer value of the transaction amount. |
| `data` | [bytes](#bytes) |  | input defines the data payload bytes of the transaction. |
| `accesses` | [AccessTuple](#ethermint.evm.v1.AccessTuple) | repeated |  |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="ethermint.evm.v1.DynamicFeeTx"></a>

### DynamicFeeTx
DynamicFeeTx is the data of EIP-1559 dinamic fee transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | destination EVM chain ID |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_tip_cap` | [string](#string) |  | gas tip cap defines the max value for the gas tip |
| `gas_fee_cap` | [string](#string) |  | gas fee cap defines the max value for the gas fee |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the the transaction amount. |
| `data` | [bytes](#bytes) |  | input defines the data payload bytes of the transaction. |
| `accesses` | [AccessTuple](#ethermint.evm.v1.AccessTuple) | repeated |  |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="ethermint.evm.v1.ExtensionOptionsEthereumTx"></a>

### ExtensionOptionsEthereumTx







<a name="ethermint.evm.v1.LegacyTx"></a>

### LegacyTx
LegacyTx is the transaction data of regular Ethereum transactions.
NOTE: All non-protected transactions (i.e non EIP155 signed) will fail if the
AllowUnprotectedTxs parameter is disabled.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_price` | [string](#string) |  | gas price defines the value for each gas unit |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the unsigned integer value of the transaction amount. |
| `data` | [bytes](#bytes) |  | input defines the data payload bytes of the transaction. |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="ethermint.evm.v1.MsgEthereumTx"></a>

### MsgEthereumTx
MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [google.protobuf.Any](#google.protobuf.Any) |  | inner transaction data

caches |
| `size` | [double](#double) |  | DEPRECATED: encoded storage size of the transaction |
| `hash` | [string](#string) |  | transaction hash in hex format |
| `from` | [string](#string) |  | ethereum signer address in hex format. This address value is checked against the address derived from the signature (V, R, S) using the secp256k1 elliptic curve |






<a name="ethermint.evm.v1.MsgEthereumTxResponse"></a>

### MsgEthereumTxResponse
MsgEthereumTxResponse defines the Msg/EthereumTx response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | ethereum transaction hash in hex format. This hash differs from the Tendermint sha256 hash of the transaction bytes. See https://github.com/tendermint/tendermint/issues/6539 for reference |
| `logs` | [Log](#ethermint.evm.v1.Log) | repeated | logs contains the transaction hash and the proto-compatible ethereum logs. |
| `ret` | [bytes](#bytes) |  | returned data from evm function (result or data supplied with revert opcode) |
| `vm_error` | [string](#string) |  | vm error is the error returned by vm execution |
| `gas_used` | [uint64](#uint64) |  | gas consumed by the transaction |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.evm.v1.Msg"></a>

### Msg
Msg defines the evm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `EthereumTx` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) | [MsgEthereumTxResponse](#ethermint.evm.v1.MsgEthereumTxResponse) | EthereumTx defines a method submitting Ethereum transactions. | POST|/ethermint/evm/v1/ethereum_tx|

 <!-- end services -->



<a name="ethermint/evm/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/query.proto



<a name="ethermint.evm.v1.EstimateGasResponse"></a>

### EstimateGasResponse
EstimateGasResponse defines EstimateGas response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas` | [uint64](#uint64) |  | the estimated gas |






<a name="ethermint.evm.v1.EthCallRequest"></a>

### EthCallRequest
EthCallRequest defines EthCall request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `args` | [bytes](#bytes) |  | same json format as the json rpc api. |
| `gas_cap` | [uint64](#uint64) |  | the default gas cap to be used |






<a name="ethermint.evm.v1.QueryAccountRequest"></a>

### QueryAccountRequest
QueryAccountRequest is the request type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the account for. |






<a name="ethermint.evm.v1.QueryAccountResponse"></a>

### QueryAccountResponse
QueryAccountResponse is the response type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [string](#string) |  | balance is the balance of the EVM denomination. |
| `code_hash` | [string](#string) |  | code hash is the hex-formatted code bytes from the EOA. |
| `nonce` | [uint64](#uint64) |  | nonce is the account's sequence number. |






<a name="ethermint.evm.v1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the balance for. |






<a name="ethermint.evm.v1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [string](#string) |  | balance is the balance of the EVM denomination. |






<a name="ethermint.evm.v1.QueryBaseFeeRequest"></a>

### QueryBaseFeeRequest
QueryBaseFeeRequest defines the request type for querying the EIP1559 base
fee.






<a name="ethermint.evm.v1.QueryBaseFeeResponse"></a>

### QueryBaseFeeResponse
BaseFeeResponse returns the EIP1559 base fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee` | [string](#string) |  |  |






<a name="ethermint.evm.v1.QueryCodeRequest"></a>

### QueryCodeRequest
QueryCodeRequest is the request type for the Query/Code RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the code for. |






<a name="ethermint.evm.v1.QueryCodeResponse"></a>

### QueryCodeResponse
QueryCodeResponse is the response type for the Query/Code RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code` | [bytes](#bytes) |  | code represents the code bytes from an ethereum address. |






<a name="ethermint.evm.v1.QueryCosmosAccountRequest"></a>

### QueryCosmosAccountRequest
QueryCosmosAccountRequest is the request type for the Query/CosmosAccount RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the account for. |






<a name="ethermint.evm.v1.QueryCosmosAccountResponse"></a>

### QueryCosmosAccountResponse
QueryCosmosAccountResponse is the response type for the Query/CosmosAccount
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cosmos_address` | [string](#string) |  | cosmos_address is the cosmos address of the account. |
| `sequence` | [uint64](#uint64) |  | sequence is the account's sequence number. |
| `account_number` | [uint64](#uint64) |  | account_number is the account numbert |






<a name="ethermint.evm.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/evm parameters.






<a name="ethermint.evm.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/evm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ethermint.evm.v1.Params) |  | params define the evm module parameters. |






<a name="ethermint.evm.v1.QueryStorageRequest"></a>

### QueryStorageRequest
QueryStorageRequest is the request type for the Query/Storage RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the storage state for. |
| `key` | [string](#string) |  | key defines the key of the storage state |






<a name="ethermint.evm.v1.QueryStorageResponse"></a>

### QueryStorageResponse
QueryStorageResponse is the response type for the Query/Storage RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [string](#string) |  | key defines the storage state value hash associated with the given key. |






<a name="ethermint.evm.v1.QueryTraceBlockRequest"></a>

### QueryTraceBlockRequest
QueryTraceBlockRequest defines TraceTx request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) | repeated | txs messages in the block |
| `trace_config` | [TraceConfig](#ethermint.evm.v1.TraceConfig) |  | TraceConfig holds extra parameters to trace functions. |
| `block_number` | [int64](#int64) |  | block number |
| `block_hash` | [string](#string) |  | block hex hash |
| `block_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | block time |






<a name="ethermint.evm.v1.QueryTraceBlockResponse"></a>

### QueryTraceBlockResponse
QueryTraceBlockResponse defines TraceBlock response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |






<a name="ethermint.evm.v1.QueryTraceTxRequest"></a>

### QueryTraceTxRequest
QueryTraceTxRequest defines TraceTx request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) |  | msgEthereumTx for the requested transaction |
| `trace_config` | [TraceConfig](#ethermint.evm.v1.TraceConfig) |  | TraceConfig holds extra parameters to trace functions. |
| `predecessors` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) | repeated | the predecessor transactions included in the same block need to be replayed first to get correct context for tracing. |
| `block_number` | [int64](#int64) |  | block number of requested transaction |
| `block_hash` | [string](#string) |  | block hex hash of requested transaction |
| `block_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | block time of requested transaction |






<a name="ethermint.evm.v1.QueryTraceTxResponse"></a>

### QueryTraceTxResponse
QueryTraceTxResponse defines TraceTx response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | response serialized in bytes |






<a name="ethermint.evm.v1.QueryTxLogsRequest"></a>

### QueryTxLogsRequest
QueryTxLogsRequest is the request type for the Query/TxLogs RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash is the ethereum transaction hex hash to query the logs for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="ethermint.evm.v1.QueryTxLogsResponse"></a>

### QueryTxLogsResponse
QueryTxLogs is the response type for the Query/TxLogs RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `logs` | [Log](#ethermint.evm.v1.Log) | repeated | logs represents the ethereum logs generated from the given transaction. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="ethermint.evm.v1.QueryValidatorAccountRequest"></a>

### QueryValidatorAccountRequest
QueryValidatorAccountRequest is the request type for the
Query/ValidatorAccount RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cons_address` | [string](#string) |  | cons_address is the validator cons address to query the account for. |






<a name="ethermint.evm.v1.QueryValidatorAccountResponse"></a>

### QueryValidatorAccountResponse
QueryValidatorAccountResponse is the response type for the
Query/ValidatorAccount RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_address` | [string](#string) |  | account_address is the cosmos address of the account in bech32 format. |
| `sequence` | [uint64](#uint64) |  | sequence is the account's sequence number. |
| `account_number` | [uint64](#uint64) |  | account_number is the account number |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.evm.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Account` | [QueryAccountRequest](#ethermint.evm.v1.QueryAccountRequest) | [QueryAccountResponse](#ethermint.evm.v1.QueryAccountResponse) | Account queries an Ethereum account. | GET|/ethermint/evm/v1/account/{address}|
| `CosmosAccount` | [QueryCosmosAccountRequest](#ethermint.evm.v1.QueryCosmosAccountRequest) | [QueryCosmosAccountResponse](#ethermint.evm.v1.QueryCosmosAccountResponse) | CosmosAccount queries an Ethereum account's Cosmos Address. | GET|/ethermint/evm/v1/cosmos_account/{address}|
| `ValidatorAccount` | [QueryValidatorAccountRequest](#ethermint.evm.v1.QueryValidatorAccountRequest) | [QueryValidatorAccountResponse](#ethermint.evm.v1.QueryValidatorAccountResponse) | ValidatorAccount queries an Ethereum account's from a validator consensus Address. | GET|/ethermint/evm/v1/validator_account/{cons_address}|
| `Balance` | [QueryBalanceRequest](#ethermint.evm.v1.QueryBalanceRequest) | [QueryBalanceResponse](#ethermint.evm.v1.QueryBalanceResponse) | Balance queries the balance of a the EVM denomination for a single EthAccount. | GET|/ethermint/evm/v1/balances/{address}|
| `Storage` | [QueryStorageRequest](#ethermint.evm.v1.QueryStorageRequest) | [QueryStorageResponse](#ethermint.evm.v1.QueryStorageResponse) | Storage queries the balance of all coins for a single account. | GET|/ethermint/evm/v1/storage/{address}/{key}|
| `Code` | [QueryCodeRequest](#ethermint.evm.v1.QueryCodeRequest) | [QueryCodeResponse](#ethermint.evm.v1.QueryCodeResponse) | Code queries the balance of all coins for a single account. | GET|/ethermint/evm/v1/codes/{address}|
| `Params` | [QueryParamsRequest](#ethermint.evm.v1.QueryParamsRequest) | [QueryParamsResponse](#ethermint.evm.v1.QueryParamsResponse) | Params queries the parameters of x/evm module. | GET|/ethermint/evm/v1/params|
| `EthCall` | [EthCallRequest](#ethermint.evm.v1.EthCallRequest) | [MsgEthereumTxResponse](#ethermint.evm.v1.MsgEthereumTxResponse) | EthCall implements the `eth_call` rpc api | GET|/ethermint/evm/v1/eth_call|
| `EstimateGas` | [EthCallRequest](#ethermint.evm.v1.EthCallRequest) | [EstimateGasResponse](#ethermint.evm.v1.EstimateGasResponse) | EstimateGas implements the `eth_estimateGas` rpc api | GET|/ethermint/evm/v1/estimate_gas|
| `TraceTx` | [QueryTraceTxRequest](#ethermint.evm.v1.QueryTraceTxRequest) | [QueryTraceTxResponse](#ethermint.evm.v1.QueryTraceTxResponse) | TraceTx implements the `debug_traceTransaction` rpc api | GET|/ethermint/evm/v1/trace_tx|
| `TraceBlock` | [QueryTraceBlockRequest](#ethermint.evm.v1.QueryTraceBlockRequest) | [QueryTraceBlockResponse](#ethermint.evm.v1.QueryTraceBlockResponse) | TraceBlock implements the `debug_traceBlockByNumber` and `debug_traceBlockByHash` rpc api | GET|/ethermint/evm/v1/trace_block|
| `BaseFee` | [QueryBaseFeeRequest](#ethermint.evm.v1.QueryBaseFeeRequest) | [QueryBaseFeeResponse](#ethermint.evm.v1.QueryBaseFeeResponse) | BaseFee queries the base fee of the parent block of the current block, it's similar to feemarket module's method, but also checks london hardfork status. | GET|/ethermint/evm/v1/base_fee|

 <!-- end services -->



<a name="ethermint/feemarket/v1/feemarket.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/feemarket.proto



<a name="ethermint.feemarket.v1.Params"></a>

### Params
Params defines the EVM module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `no_base_fee` | [bool](#bool) |  | no base fee forces the EIP-1559 base fee to 0 (needed for 0 price calls) |
| `base_fee_change_denominator` | [uint32](#uint32) |  | base fee change denominator bounds the amount the base fee can change between blocks. |
| `elasticity_multiplier` | [uint32](#uint32) |  | elasticity multiplier bounds the maximum gas limit an EIP-1559 block may have. |
| `enable_height` | [int64](#int64) |  | height at which the base fee calculation is enabled. |
| `base_fee` | [string](#string) |  | base fee for EIP-1559 blocks. |
| `min_gas_price` | [string](#string) |  | min_gas_price defines the minimum gas price value for cosmos and eth transactions |
| `min_gas_multiplier` | [string](#string) |  | min gas denominator bounds the minimum gasUsed to be charged to senders based on GasLimit |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/feemarket/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/genesis.proto



<a name="ethermint.feemarket.v1.GenesisState"></a>

### GenesisState
GenesisState defines the feemarket module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ethermint.feemarket.v1.Params) |  | params defines all the paramaters of the module. |
| `block_gas` | [uint64](#uint64) |  | block gas is the amount of gas wanted on the last block before the upgrade. Zero by default. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/feemarket/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/query.proto



<a name="ethermint.feemarket.v1.QueryBaseFeeRequest"></a>

### QueryBaseFeeRequest
QueryBaseFeeRequest defines the request type for querying the EIP1559 base
fee.






<a name="ethermint.feemarket.v1.QueryBaseFeeResponse"></a>

### QueryBaseFeeResponse
BaseFeeResponse returns the EIP1559 base fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee` | [string](#string) |  |  |






<a name="ethermint.feemarket.v1.QueryBlockGasRequest"></a>

### QueryBlockGasRequest
QueryBlockGasRequest defines the request type for querying the EIP1559 base
fee.






<a name="ethermint.feemarket.v1.QueryBlockGasResponse"></a>

### QueryBlockGasResponse
QueryBlockGasResponse returns block gas used for a given height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas` | [int64](#int64) |  |  |






<a name="ethermint.feemarket.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/evm parameters.






<a name="ethermint.feemarket.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/evm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ethermint.feemarket.v1.Params) |  | params define the evm module parameters. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.feemarket.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#ethermint.feemarket.v1.QueryParamsRequest) | [QueryParamsResponse](#ethermint.feemarket.v1.QueryParamsResponse) | Params queries the parameters of x/feemarket module. | GET|/ethermint/feemarket/v1/params|
| `BaseFee` | [QueryBaseFeeRequest](#ethermint.feemarket.v1.QueryBaseFeeRequest) | [QueryBaseFeeResponse](#ethermint.feemarket.v1.QueryBaseFeeResponse) | BaseFee queries the base fee of the parent block of the current block. | GET|/ethermint/feemarket/v1/base_fee|
| `BlockGas` | [QueryBlockGasRequest](#ethermint.feemarket.v1.QueryBlockGasRequest) | [QueryBlockGasResponse](#ethermint.feemarket.v1.QueryBlockGasResponse) | BlockGas queries the gas used at a given block height | GET|/ethermint/feemarket/v1/block_gas|

 <!-- end services -->



<a name="ethermint/types/v1/account.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/account.proto



<a name="ethermint.types.v1.EthAccount"></a>

### EthAccount
EthAccount implements the authtypes.AccountI interface and embeds an
authtypes.BaseAccount type. It is compatible with the auth AccountKeeper.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [cosmos.auth.v1beta1.BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `code_hash` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/types/v1/dynamic_fee.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/dynamic_fee.proto



<a name="ethermint.types.v1.ExtensionOptionDynamicFeeTx"></a>

### ExtensionOptionDynamicFeeTx
ExtensionOptionDynamicFeeTx is an extension option that specify the maxPrioPrice for cosmos tx


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_priority_price` | [string](#string) |  | the same as `max_priority_fee_per_gas` in eip-1559 spec |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/types/v1/indexer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/indexer.proto



<a name="ethermint.types.v1.TxResult"></a>

### TxResult
TxResult is the value stored in eth tx indexer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | the block height |
| `tx_index` | [uint32](#uint32) |  | cosmos tx index |
| `msg_index` | [uint32](#uint32) |  | the msg index in a batch tx |
| `eth_tx_index` | [int32](#int32) |  | eth tx index, the index in the list of valid eth tx in the block, aka. the transaction list returned by eth_getBlock api. |
| `failed` | [bool](#bool) |  | if the eth tx is failed |
| `gas_used` | [uint64](#uint64) |  | gas used by tx, if exceeds block gas limit, it's set to gas limit which is what's actually deducted by ante handler. |
| `cumulative_gas_used` | [uint64](#uint64) |  | the cumulative gas used within current batch tx |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/types/v1/web3.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/web3.proto



<a name="ethermint.types.v1.ExtensionOptionsWeb3Tx"></a>

### ExtensionOptionsWeb3Tx



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `typed_data_chain_id` | [uint64](#uint64) |  | typed data chain id used only in EIP712 Domain and should match Ethereum network ID in a Web3 provider (e.g. Metamask). |
| `fee_payer` | [string](#string) |  | fee payer is an account address for the fee payer. It will be validated during EIP712 signature checking. |
| `fee_payer_sig` | [bytes](#bytes) |  | fee payer sig is a signature data from the fee paying account, allows to perform fee delegation when using EIP712 Domain. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/auction/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/auction/v1beta1/types.proto



<a name="vulcanize.auction.v1beta1.Auction"></a>

### Auction
Auction represents a sealed-bid on-chain auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `status` | [string](#string) |  |  |
| `owner_address` | [string](#string) |  | Address of the creator of the auction |
| `create_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Timestamp at which the auction was created |
| `commits_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Timestamp at which the commits phase concluded |
| `reveals_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Timestamp at which the reveals phase concluded |
| `commit_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Commit and reveal fees must both be paid when committing a bid Reveal fee is returned only if the bid is revealed |
| `reveal_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `minimum_bid` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Minimum acceptable bid amount for a valid commit |
| `winner_address` | [string](#string) |  | Address of the winner |
| `winning_bid` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Winning bid, i.e., the highest bid |
| `winning_price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Amount the winner pays, i.e. the second highest auction |






<a name="vulcanize.auction.v1beta1.Auctions"></a>

### Auctions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auctions` | [Auction](#vulcanize.auction.v1beta1.Auction) | repeated |  |






<a name="vulcanize.auction.v1beta1.Bid"></a>

### Bid
Bid represents a sealed bid (commit) made during the auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [string](#string) |  |  |
| `bidder_address` | [string](#string) |  |  |
| `status` | [string](#string) |  |  |
| `commit_hash` | [string](#string) |  |  |
| `commit_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `commit_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `reveal_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `reveal_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `bid_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="vulcanize.auction.v1beta1.Params"></a>

### Params
Params defines the auction module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commits_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Duration of the commits phase in seconds |
| `reveals_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Duration of the reveals phase in seconds |
| `commit_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Commit fees |
| `reveal_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Reveal fees |
| `minimum_bid` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Minimum acceptable bid amount |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/auction/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/auction/v1beta1/genesis.proto



<a name="vulcanize.auction.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the genesis state of the auction module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#vulcanize.auction.v1beta1.Params) |  |  |
| `auctions` | [Auction](#vulcanize.auction.v1beta1.Auction) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/auction/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/auction/v1beta1/query.proto



<a name="vulcanize.auction.v1beta1.AuctionRequest"></a>

### AuctionRequest
AuctionRequest is the format for querying a specific auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | Auction ID |






<a name="vulcanize.auction.v1beta1.AuctionResponse"></a>

### AuctionResponse
AuctionResponse returns the details of the queried auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction` | [Auction](#vulcanize.auction.v1beta1.Auction) |  | Auction details |






<a name="vulcanize.auction.v1beta1.AuctionsByBidderRequest"></a>

### AuctionsByBidderRequest
AuctionsByBidderRequest is the format for querying all auctions containing a bidder address


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bidder_address` | [string](#string) |  | Address of the bidder |






<a name="vulcanize.auction.v1beta1.AuctionsByBidderResponse"></a>

### AuctionsByBidderResponse
AuctionsByBidderResponse returns all auctions containing a bidder


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auctions` | [Auctions](#vulcanize.auction.v1beta1.Auctions) |  | List of auctions |






<a name="vulcanize.auction.v1beta1.AuctionsByOwnerRequest"></a>

### AuctionsByOwnerRequest
AuctionsByOwnerRequest is the format for querying all auctions created by an owner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner_address` | [string](#string) |  | Address of the owner |






<a name="vulcanize.auction.v1beta1.AuctionsByOwnerResponse"></a>

### AuctionsByOwnerResponse
AuctionsByOwnerResponse returns all auctions created by an owner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auctions` | [Auctions](#vulcanize.auction.v1beta1.Auctions) |  | List of auctions |






<a name="vulcanize.auction.v1beta1.AuctionsRequest"></a>

### AuctionsRequest
AuctionsRequest is the format for querying all the auctions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination info for the next request |






<a name="vulcanize.auction.v1beta1.AuctionsResponse"></a>

### AuctionsResponse
AuctionsResponse returns the list of all auctions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auctions` | [Auctions](#vulcanize.auction.v1beta1.Auctions) |  | List of auctions |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination info for the next request |






<a name="vulcanize.auction.v1beta1.BalanceRequest"></a>

### BalanceRequest
BalanceRequest is the format to fetch all balances






<a name="vulcanize.auction.v1beta1.BalanceResponse"></a>

### BalanceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Set of all balances within the auction |






<a name="vulcanize.auction.v1beta1.BidRequest"></a>

### BidRequest
BidRequest is the format for querying a specific bid in an auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [string](#string) |  | Auction ID |
| `bidder` | [string](#string) |  | Bidder address |






<a name="vulcanize.auction.v1beta1.BidResponse"></a>

### BidResponse
BidResponse returns the details of the queried bid


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bid` | [Bid](#vulcanize.auction.v1beta1.Bid) |  | Bid details |






<a name="vulcanize.auction.v1beta1.BidsRequest"></a>

### BidsRequest
BidsRequest is the format for querying all bids in an auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [string](#string) |  | Auction ID |






<a name="vulcanize.auction.v1beta1.BidsResponse"></a>

### BidsResponse
BidsResponse returns details of all bids in an auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bids` | [Bid](#vulcanize.auction.v1beta1.Bid) | repeated | List of bids in the auction |






<a name="vulcanize.auction.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the format to query the parameters of the auction module






<a name="vulcanize.auction.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse returns parameters of the auction module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#vulcanize.auction.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="vulcanize.auction.v1beta1.Query"></a>

### Query
Query defines the gRPC querier interface for the auction module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Auctions` | [AuctionsRequest](#vulcanize.auction.v1beta1.AuctionsRequest) | [AuctionsResponse](#vulcanize.auction.v1beta1.AuctionsResponse) | Auctions queries all auctions | GET|/vulcanize/auction/v1beta1/auctions|
| `GetAuction` | [AuctionRequest](#vulcanize.auction.v1beta1.AuctionRequest) | [AuctionResponse](#vulcanize.auction.v1beta1.AuctionResponse) | GetAuction queries an auction | GET|/vulcanize/auction/v1beta1/auctions/{id}|
| `GetBid` | [BidRequest](#vulcanize.auction.v1beta1.BidRequest) | [BidResponse](#vulcanize.auction.v1beta1.BidResponse) | GetBid queries an auction bid | GET|/vulcanize/auction/v1beta1/bids/{auction_id}/{bidder}|
| `GetBids` | [BidsRequest](#vulcanize.auction.v1beta1.BidsRequest) | [BidsResponse](#vulcanize.auction.v1beta1.BidsResponse) | GetBids queries all auction bids | GET|/vulcanize/auction/v1beta1/bids/{auction_id}|
| `AuctionsByBidder` | [AuctionsByBidderRequest](#vulcanize.auction.v1beta1.AuctionsByBidderRequest) | [AuctionsByBidderResponse](#vulcanize.auction.v1beta1.AuctionsByBidderResponse) | AuctionsByBidder queries auctions by bidder | GET|/vulcanize/auction/v1beta1/by-bidder/{bidder_address}|
| `AuctionsByOwner` | [AuctionsByOwnerRequest](#vulcanize.auction.v1beta1.AuctionsByOwnerRequest) | [AuctionsByOwnerResponse](#vulcanize.auction.v1beta1.AuctionsByOwnerResponse) | AuctionsByOwner queries auctions by owner | GET|/vulcanize/auction/v1beta1/by-owner/{owner_address}|
| `QueryParams` | [QueryParamsRequest](#vulcanize.auction.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#vulcanize.auction.v1beta1.QueryParamsResponse) | QueryParams implements the params query command | GET|/vulcanize/auction/v1beta1/params|
| `Balance` | [BalanceRequest](#vulcanize.auction.v1beta1.BalanceRequest) | [BalanceResponse](#vulcanize.auction.v1beta1.BalanceResponse) | Balance queries the auction module account balance | GET|/vulcanize/auction/v1beta1/balance|

 <!-- end services -->



<a name="vulcanize/auction/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/auction/v1beta1/tx.proto



<a name="vulcanize.auction.v1beta1.MsgCommitBid"></a>

### MsgCommitBid
CommitBid defines the message to commit a bid


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [string](#string) |  | Auction ID |
| `commit_hash` | [string](#string) |  | Commit Hash |
| `signer` | [string](#string) |  | Address of the signer |






<a name="vulcanize.auction.v1beta1.MsgCommitBidResponse"></a>

### MsgCommitBidResponse
MsgCommitBidResponse returns the state of the auction after the bid creation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bid` | [Bid](#vulcanize.auction.v1beta1.Bid) |  | Auction details |






<a name="vulcanize.auction.v1beta1.MsgCreateAuction"></a>

### MsgCreateAuction
MsgCreateAuction defines a create auction message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commits_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Duration of the commits phase in seconds |
| `reveals_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Duration of the reveals phase in seconds |
| `commit_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Commit fees |
| `reveal_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Reveal fees |
| `minimum_bid` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Minimum acceptable bid amount |
| `signer` | [string](#string) |  | Address of the signer |






<a name="vulcanize.auction.v1beta1.MsgCreateAuctionResponse"></a>

### MsgCreateAuctionResponse
MsgCreateAuctionResponse returns the details of the created auction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction` | [Auction](#vulcanize.auction.v1beta1.Auction) |  | Auction details |






<a name="vulcanize.auction.v1beta1.MsgRevealBid"></a>

### MsgRevealBid
RevealBid defines the message to reveal a bid


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [string](#string) |  | Auction ID |
| `reveal` | [string](#string) |  | Commit Hash |
| `signer` | [string](#string) |  | Address of the signer |






<a name="vulcanize.auction.v1beta1.MsgRevealBidResponse"></a>

### MsgRevealBidResponse
MsgRevealBidResponse returns the state of the auction after the bid reveal


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction` | [Auction](#vulcanize.auction.v1beta1.Auction) |  | Auction details |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="vulcanize.auction.v1beta1.Msg"></a>

### Msg
Tx defines the gRPC tx interface

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateAuction` | [MsgCreateAuction](#vulcanize.auction.v1beta1.MsgCreateAuction) | [MsgCreateAuctionResponse](#vulcanize.auction.v1beta1.MsgCreateAuctionResponse) | CreateAuction is the command for creating an auction | |
| `CommitBid` | [MsgCommitBid](#vulcanize.auction.v1beta1.MsgCommitBid) | [MsgCommitBidResponse](#vulcanize.auction.v1beta1.MsgCommitBidResponse) | CommitBid is the command for committing a bid | |
| `RevealBid` | [MsgRevealBid](#vulcanize.auction.v1beta1.MsgRevealBid) | [MsgRevealBidResponse](#vulcanize.auction.v1beta1.MsgRevealBidResponse) | RevealBid is the command for revealing a bid | |

 <!-- end services -->



<a name="vulcanize/bond/v1beta1/bond.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/bond/v1beta1/bond.proto



<a name="vulcanize.bond.v1beta1.Bond"></a>

### Bond
Bond represents funds deposited by an account for record rent payments.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id is unique identifier of the bond |
| `owner` | [string](#string) |  | owner of the bond |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balance of the bond |






<a name="vulcanize.bond.v1beta1.Params"></a>

### Params
Params defines the bond module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_bond_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | max_bond_amount is maximum amount to bond |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/bond/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/bond/v1beta1/genesis.proto



<a name="vulcanize.bond.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the bond module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#vulcanize.bond.v1beta1.Params) |  | params defines all the parameters of the module. |
| `bonds` | [Bond](#vulcanize.bond.v1beta1.Bond) | repeated | bonds defines all the bonds |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/bond/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/bond/v1beta1/query.proto



<a name="vulcanize.bond.v1beta1.QueryGetBondByIDRequest"></a>

### QueryGetBondByIDRequest
QueryGetBondByID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="vulcanize.bond.v1beta1.QueryGetBondByIDResponse"></a>

### QueryGetBondByIDResponse
QueryGetBondByIDResponse returns QueryGetBondByID query response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bond` | [Bond](#vulcanize.bond.v1beta1.Bond) |  |  |






<a name="vulcanize.bond.v1beta1.QueryGetBondModuleBalanceRequest"></a>

### QueryGetBondModuleBalanceRequest
QueryGetBondModuleBalanceRequest is request type for bond module balance rpc method






<a name="vulcanize.bond.v1beta1.QueryGetBondModuleBalanceResponse"></a>

### QueryGetBondModuleBalanceResponse
QueryGetBondModuleBalanceResponse is the response type for bond module balance rpc method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="vulcanize.bond.v1beta1.QueryGetBondsByOwnerRequest"></a>

### QueryGetBondsByOwnerRequest
QueryGetBondsByOwnerRequest is request type for Query/GetBondsByOwner RPC Method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.bond.v1beta1.QueryGetBondsByOwnerResponse"></a>

### QueryGetBondsByOwnerResponse
QueryGetBondsByOwnerResponse is response type for Query/GetBondsByOwner RPC Method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bonds` | [Bond](#vulcanize.bond.v1beta1.Bond) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.bond.v1beta1.QueryGetBondsRequest"></a>

### QueryGetBondsRequest
QueryGetBondById queries a bond by bond-id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="vulcanize.bond.v1beta1.QueryGetBondsResponse"></a>

### QueryGetBondsResponse
QueryGetBondsResponse is response type for get the bonds by bond-id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bonds` | [Bond](#vulcanize.bond.v1beta1.Bond) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.bond.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request for query the bond module params






<a name="vulcanize.bond.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse returns response type  of bond module params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#vulcanize.bond.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="vulcanize.bond.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for bond module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#vulcanize.bond.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#vulcanize.bond.v1beta1.QueryParamsResponse) | Params queries bonds module params. | GET|/vulcanize/bond/v1beta1/params|
| `Bonds` | [QueryGetBondsRequest](#vulcanize.bond.v1beta1.QueryGetBondsRequest) | [QueryGetBondsResponse](#vulcanize.bond.v1beta1.QueryGetBondsResponse) | Bonds queries bonds list. | GET|/vulcanize/bond/v1beta1/bonds|
| `GetBondByID` | [QueryGetBondByIDRequest](#vulcanize.bond.v1beta1.QueryGetBondByIDRequest) | [QueryGetBondByIDResponse](#vulcanize.bond.v1beta1.QueryGetBondByIDResponse) | GetBondById | GET|/vulcanize/bond/v1beta1/bonds/{id}|
| `GetBondsByOwner` | [QueryGetBondsByOwnerRequest](#vulcanize.bond.v1beta1.QueryGetBondsByOwnerRequest) | [QueryGetBondsByOwnerResponse](#vulcanize.bond.v1beta1.QueryGetBondsByOwnerResponse) | Get Bonds List by Owner | GET|/vulcanize/bond/v1beta1/by-owner/{owner}|
| `GetBondsModuleBalance` | [QueryGetBondModuleBalanceRequest](#vulcanize.bond.v1beta1.QueryGetBondModuleBalanceRequest) | [QueryGetBondModuleBalanceResponse](#vulcanize.bond.v1beta1.QueryGetBondModuleBalanceResponse) | Get Bonds module balance | GET|/vulcanize/bond/v1beta1/balance|

 <!-- end services -->



<a name="vulcanize/bond/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/bond/v1beta1/tx.proto



<a name="vulcanize.bond.v1beta1.MsgCancelBond"></a>

### MsgCancelBond
MsgCancelBond defines a SDK message for the cancel the bond.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.bond.v1beta1.MsgCancelBondResponse"></a>

### MsgCancelBondResponse
MsgCancelBondResponse defines the Msg/CancelBond response type.






<a name="vulcanize.bond.v1beta1.MsgCreateBond"></a>

### MsgCreateBond
MsgCreateBond defines a SDK message for creating a new bond.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="vulcanize.bond.v1beta1.MsgCreateBondResponse"></a>

### MsgCreateBondResponse
MsgCreateBondResponse defines the Msg/CreateBond response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="vulcanize.bond.v1beta1.MsgRefillBond"></a>

### MsgRefillBond
MsgRefillBond defines a SDK message for refill the amount for bond.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="vulcanize.bond.v1beta1.MsgRefillBondResponse"></a>

### MsgRefillBondResponse
MsgRefillBondResponse defines the Msg/RefillBond response type.






<a name="vulcanize.bond.v1beta1.MsgWithdrawBond"></a>

### MsgWithdrawBond
MsgWithdrawBond defines a SDK message for withdrawing amount from bond.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="vulcanize.bond.v1beta1.MsgWithdrawBondResponse"></a>

### MsgWithdrawBondResponse
MsgWithdrawBondResponse defines the Msg/WithdrawBond response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="vulcanize.bond.v1beta1.Msg"></a>

### Msg
Msg defines the bond Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateBond` | [MsgCreateBond](#vulcanize.bond.v1beta1.MsgCreateBond) | [MsgCreateBondResponse](#vulcanize.bond.v1beta1.MsgCreateBondResponse) | CreateBond defines a method for creating a new bond. | |
| `RefillBond` | [MsgRefillBond](#vulcanize.bond.v1beta1.MsgRefillBond) | [MsgRefillBondResponse](#vulcanize.bond.v1beta1.MsgRefillBondResponse) | RefillBond defines a method for refilling amount for bond. | |
| `WithdrawBond` | [MsgWithdrawBond](#vulcanize.bond.v1beta1.MsgWithdrawBond) | [MsgWithdrawBondResponse](#vulcanize.bond.v1beta1.MsgWithdrawBondResponse) | WithdrawBond defines a method for withdrawing amount from bond. | |
| `CancelBond` | [MsgCancelBond](#vulcanize.bond.v1beta1.MsgCancelBond) | [MsgCancelBondResponse](#vulcanize.bond.v1beta1.MsgCancelBondResponse) | CancelBond defines a method for cancelling a bond. | |

 <!-- end services -->



<a name="vulcanize/nameservice/v1beta1/attributes.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/nameservice/v1beta1/attributes.proto



<a name="vulcanize.nameservice.v1beta1.ServiceProviderRegistration"></a>

### ServiceProviderRegistration



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bond_id` | [string](#string) |  |  |
| `laconic_id` | [string](#string) |  |  |
| `x500` | [X500](#vulcanize.nameservice.v1beta1.X500) |  |  |
| `type` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.WebsiteRegistrationRecord"></a>

### WebsiteRegistrationRecord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `url` | [string](#string) |  |  |
| `repo_registration_record_cid` | [string](#string) |  |  |
| `build_atrifact_cid` | [string](#string) |  |  |
| `TLS_cert_cid` | [string](#string) |  |  |
| `type` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.X500"></a>

### X500



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `common_name` | [string](#string) |  |  |
| `organization_unit` | [string](#string) |  |  |
| `organization_name` | [string](#string) |  |  |
| `locality_name` | [string](#string) |  |  |
| `state_name` | [string](#string) |  |  |
| `country` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/nameservice/v1beta1/nameservice.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/nameservice/v1beta1/nameservice.proto



<a name="vulcanize.nameservice.v1beta1.AuctionBidInfo"></a>

### AuctionBidInfo
AuctionBidInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [string](#string) |  |  |
| `bidder_address` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.AuthorityEntry"></a>

### AuthorityEntry
AuthorityEntry defines the nameservice module AuthorityEntries


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `entry` | [NameAuthority](#vulcanize.nameservice.v1beta1.NameAuthority) |  |  |






<a name="vulcanize.nameservice.v1beta1.BlockChangeSet"></a>

### BlockChangeSet
BlockChangeSet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |
| `records` | [string](#string) | repeated |  |
| `auctions` | [string](#string) | repeated |  |
| `auction_bids` | [AuctionBidInfo](#vulcanize.nameservice.v1beta1.AuctionBidInfo) | repeated |  |
| `authorities` | [string](#string) | repeated |  |
| `names` | [string](#string) | repeated |  |






<a name="vulcanize.nameservice.v1beta1.NameAuthority"></a>

### NameAuthority
NameAuthority


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner_public_key` | [string](#string) |  | Owner public key. |
| `owner_address` | [string](#string) |  | Owner address. |
| `height` | [uint64](#uint64) |  | height at which name/authority was created. |
| `status` | [string](#string) |  |  |
| `auction_id` | [string](#string) |  |  |
| `bond_id` | [string](#string) |  |  |
| `expiry_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="vulcanize.nameservice.v1beta1.NameEntry"></a>

### NameEntry
NameEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `entry` | [NameRecord](#vulcanize.nameservice.v1beta1.NameRecord) |  |  |






<a name="vulcanize.nameservice.v1beta1.NameRecord"></a>

### NameRecord
NameRecord


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `latest` | [NameRecordEntry](#vulcanize.nameservice.v1beta1.NameRecordEntry) |  |  |
| `history` | [NameRecordEntry](#vulcanize.nameservice.v1beta1.NameRecordEntry) | repeated |  |






<a name="vulcanize.nameservice.v1beta1.NameRecordEntry"></a>

### NameRecordEntry
NameRecordEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `height` | [uint64](#uint64) |  |  |






<a name="vulcanize.nameservice.v1beta1.Params"></a>

### Params
Params defines the nameservice module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record_rent` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `record_rent_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `authority_rent` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `authority_rent_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `authority_grace_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `authority_auction_enabled` | [bool](#bool) |  |  |
| `authority_auction_commits_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `authority_auction_reveals_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `authority_auction_commit_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `authority_auction_reveal_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `authority_auction_minimum_bid` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="vulcanize.nameservice.v1beta1.Record"></a>

### Record
Params defines the nameservice module records


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `bond_id` | [string](#string) |  |  |
| `create_time` | [string](#string) |  |  |
| `expiry_time` | [string](#string) |  |  |
| `deleted` | [bool](#bool) |  |  |
| `owners` | [string](#string) | repeated |  |
| `attributes` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `names` | [string](#string) | repeated |  |
| `type` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.Signature"></a>

### Signature
Signature


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sig` | [string](#string) |  |  |
| `pub_key` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/nameservice/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/nameservice/v1beta1/genesis.proto



<a name="vulcanize.nameservice.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the nameservice module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#vulcanize.nameservice.v1beta1.Params) |  | params defines all the params of nameservice module. |
| `records` | [Record](#vulcanize.nameservice.v1beta1.Record) | repeated | records |
| `authorities` | [AuthorityEntry](#vulcanize.nameservice.v1beta1.AuthorityEntry) | repeated | authorities |
| `names` | [NameEntry](#vulcanize.nameservice.v1beta1.NameEntry) | repeated | names |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="vulcanize/nameservice/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/nameservice/v1beta1/query.proto



<a name="vulcanize.nameservice.v1beta1.AccountBalance"></a>

### AccountBalance
AccountBalance is nameservice module account balance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_name` | [string](#string) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="vulcanize.nameservice.v1beta1.ExpiryQueueRecord"></a>

### ExpiryQueueRecord
ExpiryQueueRecord


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `value` | [string](#string) | repeated |  |






<a name="vulcanize.nameservice.v1beta1.GetNameServiceModuleBalanceRequest"></a>

### GetNameServiceModuleBalanceRequest
GetNameServiceModuleBalanceRequest is request type for nameservice module accounts balance






<a name="vulcanize.nameservice.v1beta1.GetNameServiceModuleBalanceResponse"></a>

### GetNameServiceModuleBalanceResponse
GetNameServiceModuleBalanceResponse is response type for nameservice module accounts balance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [AccountBalance](#vulcanize.nameservice.v1beta1.AccountBalance) | repeated |  |






<a name="vulcanize.nameservice.v1beta1.QueryGetAuthorityExpiryQueue"></a>

### QueryGetAuthorityExpiryQueue
QueryGetAuthorityExpiryQueue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="vulcanize.nameservice.v1beta1.QueryGetAuthorityExpiryQueueResponse"></a>

### QueryGetAuthorityExpiryQueueResponse
QueryGetAuthorityExpiryQueueResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorities` | [ExpiryQueueRecord](#vulcanize.nameservice.v1beta1.ExpiryQueueRecord) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.nameservice.v1beta1.QueryGetRecordExpiryQueue"></a>

### QueryGetRecordExpiryQueue
QueryGetRecordExpiryQueue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="vulcanize.nameservice.v1beta1.QueryGetRecordExpiryQueueResponse"></a>

### QueryGetRecordExpiryQueueResponse
QueryGetRecordExpiryQueueResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `records` | [ExpiryQueueRecord](#vulcanize.nameservice.v1beta1.ExpiryQueueRecord) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.nameservice.v1beta1.QueryListNameRecordsRequest"></a>

### QueryListNameRecordsRequest
QueryListNameRecordsRequest is request type for nameservice names records


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="vulcanize.nameservice.v1beta1.QueryListNameRecordsResponse"></a>

### QueryListNameRecordsResponse
QueryListNameRecordsResponse is response type for nameservice names records


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `names` | [NameEntry](#vulcanize.nameservice.v1beta1.NameEntry) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.nameservice.v1beta1.QueryListRecordsRequest"></a>

### QueryListRecordsRequest
QueryListRecordsRequest is request type for nameservice records list


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attributes` | [QueryListRecordsRequest.KeyValueInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.KeyValueInput) | repeated |  |
| `all` | [bool](#bool) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="vulcanize.nameservice.v1beta1.QueryListRecordsRequest.KeyValueInput"></a>

### QueryListRecordsRequest.KeyValueInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [QueryListRecordsRequest.ValueInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ValueInput) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ReferenceInput"></a>

### QueryListRecordsRequest.ReferenceInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ValueInput"></a>

### QueryListRecordsRequest.ValueInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `string` | [string](#string) |  |  |
| `int` | [int64](#int64) |  |  |
| `float` | [double](#double) |  |  |
| `boolean` | [bool](#bool) |  |  |
| `reference` | [QueryListRecordsRequest.ReferenceInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ReferenceInput) |  |  |
| `values` | [QueryListRecordsRequest.ValueInput](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest.ValueInput) | repeated |  |






<a name="vulcanize.nameservice.v1beta1.QueryListRecordsResponse"></a>

### QueryListRecordsResponse
QueryListRecordsResponse is response type for nameservice records list


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `records` | [Record](#vulcanize.nameservice.v1beta1.Record) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.nameservice.v1beta1.QueryLookupCrn"></a>

### QueryLookupCrn
QueryLookupCrn is request type for LookupCrn


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `crn` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryLookupCrnResponse"></a>

### QueryLookupCrnResponse
QueryLookupCrnResponse is response type for QueryLookupCrn


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [NameRecord](#vulcanize.nameservice.v1beta1.NameRecord) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for nameservice params






<a name="vulcanize.nameservice.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for nameservice params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#vulcanize.nameservice.v1beta1.Params) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryRecordByBondIDRequest"></a>

### QueryRecordByBondIDRequest
QueryRecordByBondIdRequest is request type for get the records by bond-id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="vulcanize.nameservice.v1beta1.QueryRecordByBondIDResponse"></a>

### QueryRecordByBondIDResponse
QueryRecordByBondIdResponse is response type for records list by bond-id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `records` | [Record](#vulcanize.nameservice.v1beta1.Record) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="vulcanize.nameservice.v1beta1.QueryRecordByIDRequest"></a>

### QueryRecordByIDRequest
QueryRecordByIDRequest is request type for nameservice records by id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryRecordByIDResponse"></a>

### QueryRecordByIDResponse
QueryRecordByIDResponse is response type for nameservice records by id


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record` | [Record](#vulcanize.nameservice.v1beta1.Record) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryResolveCrn"></a>

### QueryResolveCrn
QueryResolveCrn is request type for ResolveCrn


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `crn` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryResolveCrnResponse"></a>

### QueryResolveCrnResponse
QueryResolveCrnResponse is response type for QueryResolveCrn


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record` | [Record](#vulcanize.nameservice.v1beta1.Record) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryWhoisRequest"></a>

### QueryWhoisRequest
QueryWhoisRequest is request type for Get NameAuthority


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.QueryWhoisResponse"></a>

### QueryWhoisResponse
QueryWhoisResponse is response type for whois request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name_authority` | [NameAuthority](#vulcanize.nameservice.v1beta1.NameAuthority) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="vulcanize.nameservice.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for nameservice module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#vulcanize.nameservice.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#vulcanize.nameservice.v1beta1.QueryParamsResponse) | Params queries the nameservice module params. | GET|/vulcanize/nameservice/v1beta1/params|
| `ListRecords` | [QueryListRecordsRequest](#vulcanize.nameservice.v1beta1.QueryListRecordsRequest) | [QueryListRecordsResponse](#vulcanize.nameservice.v1beta1.QueryListRecordsResponse) | List records | GET|/vulcanize/nameservice/v1beta1/records|
| `GetRecord` | [QueryRecordByIDRequest](#vulcanize.nameservice.v1beta1.QueryRecordByIDRequest) | [QueryRecordByIDResponse](#vulcanize.nameservice.v1beta1.QueryRecordByIDResponse) | Get record by id | GET|/vulcanize/nameservice/v1beta1/records/{id}|
| `GetRecordByBondID` | [QueryRecordByBondIDRequest](#vulcanize.nameservice.v1beta1.QueryRecordByBondIDRequest) | [QueryRecordByBondIDResponse](#vulcanize.nameservice.v1beta1.QueryRecordByBondIDResponse) | Get records by bond id | GET|/vulcanize/nameservice/v1beta1/records-by-bond-id/{id}|
| `GetNameServiceModuleBalance` | [GetNameServiceModuleBalanceRequest](#vulcanize.nameservice.v1beta1.GetNameServiceModuleBalanceRequest) | [GetNameServiceModuleBalanceResponse](#vulcanize.nameservice.v1beta1.GetNameServiceModuleBalanceResponse) | Get nameservice module balance | GET|/vulcanize/nameservice/v1beta1/balance|
| `ListNameRecords` | [QueryListNameRecordsRequest](#vulcanize.nameservice.v1beta1.QueryListNameRecordsRequest) | [QueryListNameRecordsResponse](#vulcanize.nameservice.v1beta1.QueryListNameRecordsResponse) | List name records | GET|/vulcanize/nameservice/v1beta1/names|
| `Whois` | [QueryWhoisRequest](#vulcanize.nameservice.v1beta1.QueryWhoisRequest) | [QueryWhoisResponse](#vulcanize.nameservice.v1beta1.QueryWhoisResponse) | Whois method retrieve the name authority info | GET|/vulcanize/nameservice/v1beta1/whois/{name}|
| `LookupCrn` | [QueryLookupCrn](#vulcanize.nameservice.v1beta1.QueryLookupCrn) | [QueryLookupCrnResponse](#vulcanize.nameservice.v1beta1.QueryLookupCrnResponse) | LookupCrn | GET|/vulcanize/nameservice/v1beta1/lookup|
| `ResolveCrn` | [QueryResolveCrn](#vulcanize.nameservice.v1beta1.QueryResolveCrn) | [QueryResolveCrnResponse](#vulcanize.nameservice.v1beta1.QueryResolveCrnResponse) | ResolveCrn | GET|/vulcanize/nameservice/v1beta1/resolve|
| `GetRecordExpiryQueue` | [QueryGetRecordExpiryQueue](#vulcanize.nameservice.v1beta1.QueryGetRecordExpiryQueue) | [QueryGetRecordExpiryQueueResponse](#vulcanize.nameservice.v1beta1.QueryGetRecordExpiryQueueResponse) | GetRecordExpiryQueue | GET|/vulcanize/nameservice/v1beta1/record-expiry|
| `GetAuthorityExpiryQueue` | [QueryGetAuthorityExpiryQueue](#vulcanize.nameservice.v1beta1.QueryGetAuthorityExpiryQueue) | [QueryGetAuthorityExpiryQueueResponse](#vulcanize.nameservice.v1beta1.QueryGetAuthorityExpiryQueueResponse) | GetAuthorityExpiryQueue | GET|/vulcanize/nameservice/v1beta1/authority-expiry|

 <!-- end services -->



<a name="vulcanize/nameservice/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## vulcanize/nameservice/v1beta1/tx.proto



<a name="vulcanize.nameservice.v1beta1.MsgAssociateBond"></a>

### MsgAssociateBond
MsgAssociateBond


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record_id` | [string](#string) |  |  |
| `bond_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgAssociateBondResponse"></a>

### MsgAssociateBondResponse
MsgAssociateBondResponse






<a name="vulcanize.nameservice.v1beta1.MsgDeleteNameAuthority"></a>

### MsgDeleteNameAuthority
MsgDeleteNameAuthority is SDK message for DeleteNameAuthority


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `crn` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgDeleteNameAuthorityResponse"></a>

### MsgDeleteNameAuthorityResponse
MsgDeleteNameAuthorityResponse






<a name="vulcanize.nameservice.v1beta1.MsgDissociateBond"></a>

### MsgDissociateBond
MsgDissociateBond is SDK message for Msg/DissociateBond


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgDissociateBondResponse"></a>

### MsgDissociateBondResponse
MsgDissociateBondResponse is response type for MsgDissociateBond






<a name="vulcanize.nameservice.v1beta1.MsgDissociateRecords"></a>

### MsgDissociateRecords
MsgDissociateRecords is SDK message for Msg/DissociateRecords


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bond_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgDissociateRecordsResponse"></a>

### MsgDissociateRecordsResponse
MsgDissociateRecordsResponse is response type for MsgDissociateRecords






<a name="vulcanize.nameservice.v1beta1.MsgReAssociateRecords"></a>

### MsgReAssociateRecords
MsgReAssociateRecords is SDK message for Msg/ReAssociateRecords


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_bond_id` | [string](#string) |  |  |
| `old_bond_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgReAssociateRecordsResponse"></a>

### MsgReAssociateRecordsResponse
MsgReAssociateRecordsResponse is response type for MsgReAssociateRecords






<a name="vulcanize.nameservice.v1beta1.MsgRenewRecord"></a>

### MsgRenewRecord
MsgRenewRecord is SDK message for Renew a record


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgRenewRecordResponse"></a>

### MsgRenewRecordResponse
MsgRenewRecordResponse






<a name="vulcanize.nameservice.v1beta1.MsgReserveAuthority"></a>

### MsgReserveAuthority
MsgReserveName


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |
| `owner` | [string](#string) |  | if creating a sub-authority. |






<a name="vulcanize.nameservice.v1beta1.MsgReserveAuthorityResponse"></a>

### MsgReserveAuthorityResponse
MsgReserveNameResponse






<a name="vulcanize.nameservice.v1beta1.MsgSetAuthorityBond"></a>

### MsgSetAuthorityBond
MsgSetAuthorityBond is SDK message for SetAuthorityBond


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `bond_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgSetAuthorityBondResponse"></a>

### MsgSetAuthorityBondResponse
MsgSetAuthorityBondResponse






<a name="vulcanize.nameservice.v1beta1.MsgSetName"></a>

### MsgSetName
MsgSetName


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `crn` | [string](#string) |  |  |
| `cid` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgSetNameResponse"></a>

### MsgSetNameResponse
MsgSetNameResponse






<a name="vulcanize.nameservice.v1beta1.MsgSetRecord"></a>

### MsgSetRecord
MsgSetRecord


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bond_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |
| `payload` | [Payload](#vulcanize.nameservice.v1beta1.Payload) |  |  |






<a name="vulcanize.nameservice.v1beta1.MsgSetRecordResponse"></a>

### MsgSetRecordResponse
MsgSetRecordResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="vulcanize.nameservice.v1beta1.Payload"></a>

### Payload
Payload


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `record` | [Record](#vulcanize.nameservice.v1beta1.Record) |  |  |
| `signatures` | [Signature](#vulcanize.nameservice.v1beta1.Signature) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="vulcanize.nameservice.v1beta1.Msg"></a>

### Msg
Msg

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SetRecord` | [MsgSetRecord](#vulcanize.nameservice.v1beta1.MsgSetRecord) | [MsgSetRecordResponse](#vulcanize.nameservice.v1beta1.MsgSetRecordResponse) | SetRecord will records a new record with given payload and bond id | |
| `RenewRecord` | [MsgRenewRecord](#vulcanize.nameservice.v1beta1.MsgRenewRecord) | [MsgRenewRecordResponse](#vulcanize.nameservice.v1beta1.MsgRenewRecordResponse) | Renew Record will renew the expire record | |
| `AssociateBond` | [MsgAssociateBond](#vulcanize.nameservice.v1beta1.MsgAssociateBond) | [MsgAssociateBondResponse](#vulcanize.nameservice.v1beta1.MsgAssociateBondResponse) | AssociateBond | |
| `DissociateBond` | [MsgDissociateBond](#vulcanize.nameservice.v1beta1.MsgDissociateBond) | [MsgDissociateBondResponse](#vulcanize.nameservice.v1beta1.MsgDissociateBondResponse) | DissociateBond | |
| `DissociateRecords` | [MsgDissociateRecords](#vulcanize.nameservice.v1beta1.MsgDissociateRecords) | [MsgDissociateRecordsResponse](#vulcanize.nameservice.v1beta1.MsgDissociateRecordsResponse) | DissociateRecords | |
| `ReAssociateRecords` | [MsgReAssociateRecords](#vulcanize.nameservice.v1beta1.MsgReAssociateRecords) | [MsgReAssociateRecordsResponse](#vulcanize.nameservice.v1beta1.MsgReAssociateRecordsResponse) | ReAssociateRecords | |
| `SetName` | [MsgSetName](#vulcanize.nameservice.v1beta1.MsgSetName) | [MsgSetNameResponse](#vulcanize.nameservice.v1beta1.MsgSetNameResponse) | SetName will store the name with given crn and name | |
| `ReserveName` | [MsgReserveAuthority](#vulcanize.nameservice.v1beta1.MsgReserveAuthority) | [MsgReserveAuthorityResponse](#vulcanize.nameservice.v1beta1.MsgReserveAuthorityResponse) | Reserve name | |
| `DeleteName` | [MsgDeleteNameAuthority](#vulcanize.nameservice.v1beta1.MsgDeleteNameAuthority) | [MsgDeleteNameAuthorityResponse](#vulcanize.nameservice.v1beta1.MsgDeleteNameAuthorityResponse) | Delete Name method will remove authority name | |
| `SetAuthorityBond` | [MsgSetAuthorityBond](#vulcanize.nameservice.v1beta1.MsgSetAuthorityBond) | [MsgSetAuthorityBondResponse](#vulcanize.nameservice.v1beta1.MsgSetAuthorityBondResponse) | SetAuthorityBond | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

