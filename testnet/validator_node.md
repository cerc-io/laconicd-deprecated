### Create Validator Post Genesis

1. Run Full Node
2. Create Account and Get test tokens
3. Create Validator

### 1.Run Full Node

- Check "[Run Full Node](full_node.md)" section to Run a Full Node

### 2. Create Account & Get test tokens

```
chibaclonkd keys add <key-name> --keyring-backend test
```

NOTE: Save `mnemonic` and related account details (public key). You will need to use the need mnemonic/private key to
recover accounts at a later point in time.

##### Get Test tokens from faucet

- Faucet website link will be updated
- 1 CHK =  1 * 10e^18 achk

### 3.Create Validator

- ##### Check full node sync status

  `chibaclonkd status 2>&1 | jq -r ".SyncInfo"`

  `catching_up: false` means node is completely synced
- ##### Create validator

`Note:`  Only execute below transaction after complete sync of your full node

Please replace `key_name` with your key name, amount with staking amount, validator description and `moniker` also

```
chibaclonkd tx staking create-validator \
  --amount="AMOUNT" \
  --pubkey=$(chibaclonkd tendermint show-validator) \
  --moniker="my-moniker" \
  --website="https://myweb.site" \
  --details="description of your validator" \
  --chain-id="chibaclonk_81337-2" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-adjustment="1.2" \
  --gas-prices="0.025achk" \
  --from=<key_name>
```