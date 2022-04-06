### Create Validator Post Genesis

1. Run Full Node
2. Create Account and Get test tokens
3. Create Validator

### 1.Run Full Node

- Check "[Run Full Node](full-node.md)" section to Run a Full Node

### 2. Create Account & Get test tokens

```
chibaclonkd keys add <key-name> --keyring-backend test
```

NOTE: Save `mnemonic` and related account details (public key). You will need to use the need mnemonic/private key to
recover accounts at a later point in time.

##### Get Test tokens from faucet

- Open this link : http://167.172.173.94:1314/ and paste your account
- 1 gnt = 10^18 agnt
- Each Transaction you will get 500gnt
- Total Tokens 5000gnt for account

### 3.Create Validator

- ##### Check full node sync status

  `chibaclonkd status 2>&1 | jq -r ".SyncInfo"`

  `catching_up: false` means node is completely synced
- ##### Create validator

`Note:`  Only execute below transaction after complete sync of your full node

Please replace `key_name` with your key name and `moniker` also

```
chibaclonkd tx staking create-validator \
  --amount=4500000000000000000000agnt \
  --pubkey=$(chibaclonkd tendermint show-validator) \
  --moniker="my-moniker" \
  --website="https://myweb.site" \
  --details="description of your validator" \
  --chain-id="chibaclonk_81337-1" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-adjustment="1.2" \
  --gas-prices="0.025agnt" \
  --from=<key_name>
```