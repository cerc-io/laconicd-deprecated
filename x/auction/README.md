# Auction Module CLI Commands

## Build the Chain

The following command builds the Ethermint daemon and places the binary in the `build` directory
```
make build

```

## Setup the Chain

The following steps need to be followed only before running the chain for the first time.

1. Add the root key:
	```
	./build/ethermintd keys add root
	```
	Keep a note of the keyring passphrase if you set it.
2. Init the chain:
	```
	./build/ethermintd init test-moniker --chain-id ethermint_9000-1
	```
3. Add genesis account:
	```
	./build/ethermintd add-genesis-account $(./build/ethermintd keys show root -a) 1000000000000000000aphoton,1000000000000000000stake
	```
4. Make a genesis tx:
	```
	./build/ethermintd gentx root 1000000000000000000stake --chain-id ethermint_9000-1 
	```
5. Collect gentxs:
	```
	./build/ethermintd collect-gentxs
	```

The chain can now be started using:
```
./build/ethermintd start
```

## Querying the Params

The following command will dislay the default params for the `auction` module:
```
# ./build/ethermintd q auction params -o json | jq

{
  "params": {
    "commits_duration": "0s",
    "reveals_duration": "0s",
    "commit_fee": {
      "denom": "",
      "amount": "0"
    },
    "reveal_fee": {
      "denom": "",
      "amount": "0"
    },
    "minimum_bid": {
      "denom": "",
      "amount": "0"
    }
  }
}
```

## Auction TX CLI Commands

### Create Auction

```
# ./build/ethermintd tx auction create 100s 100s 10aphoton 10aphoton 1000aphoton --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)

Enter keyring passphrase:

{"body":{"messages":[{"@type":"/vulcanize.auction.v1beta1.MsgCreateAuction","commits_duration":"100s","reveals_duration":"100s","commit_fee":{"denom":"aphoton","amount":"10"},"reveal_fee":{"denom":"aphoton","amount":"10"},"minimum_bid":{"denom":"aphoton","amount":"1000"},"signer":"ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: ECAD6DF1ECA763FBD26EB7C2C0B77425FFE2FBEA2BEC57CE0FBC173AE0F45298
```

### Commit Bid

```
# ./build/ethermintd tx auction commit-bid e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d 2000aphoton --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)

Enter keyring passphrase:

{"body":{"messages":[{"@type":"/vulcanize.auction.v1beta1.MsgCommitBid","auction_id":"e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d","commit_hash":"bafyreibt4twofrc3xi2es27cfrroy346iy6lr3gkw33i5dltkqqarlyltm","signer":"ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: 71D8CF34026E32A3A34C2C2D4ADF25ABC8D7943A4619761BE27F196603D91B9D
```

### Reveal Bid

```
# ./build/ethermintd tx auction reveal-bid e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d root-bafyreibt4twofrc3xi2es27cfrroy346iy6lr3gkw33i5dltkqqarlyltm.json --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)

Enter keyring passphrase:

{"body":{"messages":[{"@type":"/vulcanize.auction.v1beta1.MsgRevealBid","auction_id":"e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d","reveal":"7b2261756374696f6e4964223a2265376431346337653761366437353337636264623866626536326632326231353533633265663463653337303561646137633238663830666166326662653064222c22626964416d6f756e74223a22323030306170686f746f6e222c2262696464657241646472657373223a226574686d316c37637374777466326c76657632376b613637633233796b376d6d6a38616437746574707163222c22636861696e4964223a2265746865726d696e745f393030302d31222c226e6f697365223a22636c69666620737566666572206472616d6120676f7370656c2077656173656c207061706572206c696272617279206469736f726465722063757276652073706f74206375727461696e207a6562726120696e76657374206465766f74652072656e64657220636c6970207377616c6c6f77206d6f6e6b6579206f62736572766520726573706f6e7365206c696e6b206372616e6520766961626c6520736576656e227d","signer":"ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: 4D1C0B3DDA4050F9BB32240FBD5234229E5C32543C1A0A78033B9531EB0CF8BA
```

## Auction Query CLI Commands

### List Auctions

```
# ./build/ethermintd q auction list

auctions:
  auctions:
  - commit_fee:
      amount: "10"
      denom: aphoton
    commits_end_time: "2021-09-30T07:57:07.933412800Z"
    create_time: "2021-09-30T07:55:27.933412800Z"
    id: e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d
    minimum_bid:
      amount: "1000"
      denom: aphoton
    owner_address: ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc
    reveal_fee:
      amount: "10"
      denom: aphoton
    reveals_end_time: "2021-09-30T07:58:47.933412800Z"
    status: commit
    winner_address: ""
    winning_bid:
      amount: "0"
      denom: ""
    winning_price:
      amount: "0"
      denom: ""
pagination: null
```

### Get Bid

```
# ./build/ethermintd q auction get-bid e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0e ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc

bid:
  auction_id: e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d
  bid_amount:
    amount: "0"
    denom: ""
  bidder_address: ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc
  commit_fee:
    amount: "10"
    denom: aphoton
  commit_hash: bafyreibt4twofrc3xi2es27cfrroy346iy6lr3gkw33i5dltkqqarlyltm
  commit_time: "2021-09-30T08:49:48.358878200Z"
  reveal_fee:
    amount: "10"
    denom: aphoton
  reveal_time: "0001-01-01T00:00:00Z"
  status: commit
```

### Get All Bids for an Auction

```
./build/ethermintd q auction get-bids e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d

bids:
- auction_id: e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d
  bid_amount:
    amount: "0"
    denom: ""
  bidder_address: ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc
  commit_fee:
    amount: "10"
    denom: aphoton
  commit_hash: bafyreibt4twofrc3xi2es27cfrroy346iy6lr3gkw33i5dltkqqarlyltm
  commit_time: "2021-09-30T08:49:48.358878200Z"
  reveal_fee:
    amount: "10"
    denom: aphoton
  reveal_time: "0001-01-01T00:00:00Z"
  status: commit
```

### Get Auction by AuctionID

```
# ./build/ethermintd q auction get e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d

auction:
  commit_fee:
    amount: "10"
    denom: aphoton
  commits_end_time: "2021-09-30T07:57:07.933412800Z"
  create_time: "2021-09-30T07:55:27.933412800Z"
  id: e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d
  minimum_bid:
    amount: "1000"
    denom: aphoton
  owner_address: ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc
  reveal_fee:
    amount: "10"
    denom: aphoton
  reveals_end_time: "2021-09-30T07:58:47.933412800Z"
  status: commit
  winner_address: ""
  winning_bid:
    amount: "0"
    denom: ""
  winning_price:
    amount: "0"
    denom: ""

```

### Get Auction by Bidder

```
# ./build/ethermintd q auction query-by-owner ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc

auctions:
  auctions:
  - commit_fee:
      amount: "10"
      denom: aphoton
    commits_end_time: "2021-09-30T07:57:07.933412800Z"
    create_time: "2021-09-30T07:55:27.933412800Z"
    id: e7d14c7e7a6d7537cbdb8fbe62f22b1553c2ef4ce3705ada7c28f80faf2fbe0d
    minimum_bid:
      amount: "1000"
      denom: aphoton
    owner_address: ethm1l7cstwtf2lvev27ka67c23yk7mmj8ad7tetpqc
    reveal_fee:
      amount: "10"
      denom: aphoton
    reveals_end_time: "2021-09-30T07:58:47.933412800Z"
    status: commit
    winner_address: ""
    winning_bid:
      amount: "0"
      denom: ""
    winning_price:
      amount: "0"
      denom: ""
```

### Query Account Balance

```
# ./build/ethermintd q auction balance                                                   

balance:
- amount: "20"
  denom: aphoton
```
