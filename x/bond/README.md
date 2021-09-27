# Build chain
```bash 
# it will create binary in build folder with `ethermintd` 
$ make build 
```
# Setup Chain
```bash
./build/ethermintd keys add root 
./build/ethermintd init test-moniker --chain-id ethermint_9000-1
./build/ethermintd add-genesis-account $(./build/ethermintd keys show root -a) 1000000000000000000aphoton,1000000000000000000stake
./build/ethermintd gentx root 1000000000000000000stake --chain-id ethermint_9000-1 
./build/ethermintd collect-gentxs
./build/ethermintd start
```

# Params
```
$ ./build/ethermintd q bond params -o json | jq .

{
  "params": {
    "max_bond_amount": {
      "denom": "stake",
      "amount": "100000000000"
    }
  }
}
```

# Create Bond
```
 $ ./build/ethermintd tx bond create 100aphoton --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)

{"body":{"messages":[{"@type":"/vulcanize.bond.v1beta1.MsgCreateBond","signer":"ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk","coins":[{"denom":"aphoton","amount":"100"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

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
txhash: C6D362E11D4C9FB06D620F3CCAF363A69A074297A00E9CAECBDA5BE1CC302EB8
```
# Refill Bond
```
 $ ./build/ethermintd tx bond refill c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440 1000aphoton --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)
{"body":{"messages":[{"@type":"/vulcanize.bond.v1beta1.MsgRefillBond","id":"c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440","signer":"ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk","coins":[{"denom":"aphoton","amount":"1000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

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
txhash: 025B04E2C923EFE2299CD171B40829CB1FE4D1A69DA7563C64AAC3D5C27BDFC9

```
# Withdraw from bond
```
 $ ./build/ethermintd tx bond withdraw c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440 1000aphoton --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)
{"body":{"messages":[{"@type":"/vulcanize.bond.v1beta1.MsgWithdrawBond","id":"c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440","signer":"ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk","coins":[{"denom":"aphoton","amount":"1000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

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
txhash: 7C5E2FE674577CD6BFFF9F92FCCBC61EDFE8F1A00CE29AC84D58FB8F013D2F03
```
# List Bond
```
 $ ./build/ethermintd q bond list -o json | jq .
{
  "bonds": [
    {
      "id": "c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440",
      "owner": "ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk",
      "balance": [
        {
          "denom": "aphoton",
          "amount": "100"
        }
      ]
    }
  ],
  "pagination": null
}

```
# Get Bond by Id
```
$ ./build/ethermintd q bond get c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440  -o json | jq .
{
  "bond": {
    "id": "c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440",
    "owner": "ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk",
    "balance": [
      {
        "denom": "aphoton",
        "amount": "100"
      }
    ]
  }
}

```

# Get Bond Module Balance
```
$ ./build/ethermintd q bond balance -o json | jq .                                                              
{
  "balance": [
    {
      "denom": "aphoton",
      "amount": "100"
    }
  ]
}
```

# Get Bonds By Owner
```
$ ./build/ethermintd q bond query-by-owner ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk -o json | jq .
{
  "bonds": [
    {
      "id": "c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440",
      "owner": "ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk",
      "balance": [
        {
          "denom": "aphoton",
          "amount": "100"
        }
      ]
    }
  ],
  "pagination": null
}
```

# Cancel the bond
```
 $ ./build/ethermintd tx bond cancel c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440  --from root --chain-id $(./build/ethermintd status | jq .NodeInfo.network -r)           
{"body":{"messages":[{"@type":"/vulcanize.bond.v1beta1.MsgCancelBond","id":"c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440","signer":"ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

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
txhash: 06440D0B35A862E3A6783E147F0E1CDD3374870DAE07E471D465E2830DAF7044

```

Note : After the bond create and withdraw bond and cancel bond , account amount needs change as per module
```
$ ./build/ethermintd q bank balances $(./build/ethermintd keys show root -a) -o json | jq .
{
  "balances": [
    {
      "denom": "aphoton",
      "amount": "1000000000000000000"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "0"
  }
}
``` 