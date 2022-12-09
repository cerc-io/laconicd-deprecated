# Build chain

```bash
# it will create binary in build folder with `laconicd`
$ make build
```

# Setup Chain

```bash
./build/laconicd keys add root
./build/laconicd init test-moniker --chain-id ethermint_9000-1
./build/laconicd add-genesis-account $(./build/laconicd keys show root -a) 1000000000000000000aphoton,1000000000000000000stake
./build/laconicd gentx root 1000000000000000000stake --chain-id ethermint_9000-1
./build/laconicd collect-gentxs
./build/laconicd start
```

## Get Params

```bash
$ ./build/laconicd q registry params -o json | jq .
{
  "params": {
    "record_rent": {
      "denom": "stake",
      "amount": "1000000"
    },

    "record_rent_duration": "31536000s",
    "authority_rent": {
      "denom": "stake",
      "amount": "1000000"
    },
    "authority_rent_duration": "31536000s",
    "authority_grace_period": "172800s",
    "authority_auction_enabled": false,
    "authority_auction_commits_duration": "86400s",
    "authority_auction_reveals_duration": "86400s",
    "authority_auction_commit_fee": {
      "denom": "stake",
      "amount": "1000000"
    },
    "authority_auction_reveal_fee": {
      "denom": "stake",
      "amount": "1000000"
    },
    "authority_auction_minimum_bid": {
      "denom": "stake",
      "amount": "5000000"
    }
  }

```

## Create (Set) Record

> First you have to Create bond

```bash
$ ./build/laconicd tx registry set ~/Desktop/examples/records/service_provider_example.yml 95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d --from root --chain-id ethermint_9000-1 --yes -o json
{
  "height": "0",
  "txhash": "BA44ABE1194724694E7CB290F9F3121DB4E63E1A030D95CB84813EEA132CF95F",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}
```

## Get records list

```bash
$ ./build/laconicd q registry list -o json  | jq
[
  {
    "id": "bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae",
    "bondId": "c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440",
    "createTime": "2021-10-04T06:50:06.369025861Z",
    "expiryTime": "2022-10-04T06:50:06.369025861Z",
    "attributes": {
      "attr1": "value1",
      "attr2": "value2",
      "link1": {
        "/": "QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D"
      },
      "link2": {
        "/": "QmP8jTG1m9GSDJLCbeWhVSVgEzCPPwXRdCRuJtQ5Tz9Kc9"
      }
    }
  }
]


```

## Get record by id

```bash
$ ./build/laconicd q registry get bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae -o json | jq .
{
  "record": {
    "id": "bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae",
    "bond_id": "95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d",
    "create_time": "2021-09-27T07:23:25.558111606Z",
    "expiry_time": "2022-09-27T07:23:25.558111606Z",
    "deleted": false,
    "owners": [],
    "attributes": "eyJhdHRyMSI6InZhbHVlMSIsImF0dHIyIjoidmFsdWUyIiwibGluazEiOnsiLyI6IlFtU251V214cHRKWmRMSnBLUmFyeEJNUzJKdTJvQU5WcmdicjJ4V2JpZTliMkQifSwibGluazIiOnsiLyI6IlFtUDhqVEcxbTlHU0RKTENiZVdoVlNWZ0V6Q1BQd1hSZENSdUp0UTVUejlLYzkifX0="
  }
}
```

## Reserve name

```bash
 ./build/laconicd tx registry reserve-name hello --from root --chain-id ethermint_9000-1 --owner $(./build/laconicd key
s show root -a) -y -o json | jq .
{
  "height": "0",
  "txhash": "7EC19157AC89279DEBE840EA3384FC95D1E2A0931C27746CA42AC23AE285B7ED",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}

```

## Query Whois for name authority

```bash
 ./build/laconicd q registry whois hello -o json | jq .
{
  "name_authority": {
    "owner_public_key": "Au3hH1tzL1KgZfXfA71jGYSe5RV9Wg95kwhBWs8V+N+h",
    "owner_address": "ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk",
    "height": "174",
    "status": "active",
    "auction_id": "",
    "bond_id": "",
    "expiry_time": "2021-09-29T07:34:36.304545965Z"
  }
}

```

## Query the registry module balance

```bash
$ ./build/laconicd q registry  balance -o json | jq .
{
  "balances": [
    {
      "account_name": "record_rent",
      "balance": [
        {
          "denom": "stake",
          "amount": "1000000"
        }
      ]
    }
  ]
}

```

## add bond to the authority

```bash
$ ./build/laconicd tx registry authority-bond [Authority Name] [Bond ID ]  --from root --chain-id ethermint_9000-1  -y -o json | jq .
$ ./build/laconicd tx registry authority-bond hello 95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d  --from root --chain-id ethermint_9000-1  -y -o json | jq .
```

## Query the records by associate bond id

```bash
$ ./build/laconicd q registry query-by-bond 95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d -o json | jq .
{
  "records": [
    {
      "id": "bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae",
      "bond_id": "95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d",
      "create_time": "2021-09-27T08:25:32.893155609Z",
      "expiry_time": "2022-09-27T08:25:32.893155609Z",
      "deleted": false,
      "owners": [],
      "attributes": "eyJhdHRyMSI6InZhbHVlMSIsImF0dHIyIjoidmFsdWUyIiwibGluazEiOnsiLyI6IlFtU251V214cHRKWmRMSnBLUmFyeEJNUzJKdTJvQU5WcmdicjJ4V2JpZTliMkQifSwibGluazIiOnsiLyI6IlFtUDhqVEcxbTlHU0RKTENiZVdoVlNWZ0V6Q1BQd1hSZENSdUp0UTVUejlLYzkifX0="
    }
  ],
  "pagination": null
}

```

## dissociate bond from record

```bash
$ ./build/laconicd tx registry dissociate-bond bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae  --from root --chain-id ethermint_9000-1
{"body":{"messages":[{"@type":"/vulcanize.registry.v1beta1.MsgDissociateBond","record_id":"bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae","signer":"ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

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
txhash: 7AFEF524CB0D92D6576FC08601A787786E802449888FD8DDAA7635698CC85060

```

## Associate bond with record

```bash
./build/laconicd tx registry associate-bond bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440 --from root --chain-id ethermint_9000-1  -y -o json | jq .
{
  "height": "0",
  "txhash": "F75C2BF2FE73668AE1332E1237F924AC549E31E822A56394DE5AC17200B199F9",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}

```

## dissociate-records => remove all record from bond

```bash
$./build/laconicd tx registry dissociate-records c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440 --from root --chain-id ethermint_9000-1  -y -o json | jq .
{
  "height": "0",
  "txhash": "0316F503E5DEA47CB108AE6C7C7FFAF3F71CC56BC22F63CB97322E1BE48B33B9",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}
```

## Renew a record

> When a record is expires , needs to renew record

```bash
$ ./build/laconicd tx registry renew-record bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae --from root --chain-id ethermint_9000-1

```

## Set the authority name

```bash
$ ./build/laconicd tx registry set-name crn://hello/test test_hello_cid  --from root --chain-id ethermint_9000-1 -y -o json | jq .
{
  "height": "0",
  "txhash": "66A63C73B076EEE9A2F7605354448EDEB161F0115D4D03AF68C01BA28DB97486",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}
```

## Delete the name

```bash
$./build/laconicd tx registry delete-name crn://hello/test --from root --chain-id ethermint_9000-1 -y
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
txhash: A3FF4C46BAC7BD6E54BBB743A49830AE8C6F6FE59282384789CBA323C1FE540C

```

## List of Authorities Expire Queue

```bash
$ ./build/laconicd q registry authority-expiry  -o json | jq .
{
  "authorities": [],
  "pagination": null
}

```

## List of Records Expire Queue

```bash
$  ./build/laconicd q registry record-expiry -o json | jq .
{
  "records": [],
  "pagination": null
}

```
