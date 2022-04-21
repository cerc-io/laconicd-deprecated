# Vulcanize chiba-clonk gql

> Browser : http://localhost:9473 for gql

## Start server

```shell
./build/chibaclonkd start --gql-playground --gql-server
```

Basic node status:

```graphql
{
    getStatus {
        version
        node {
            id
            network
            moniker
        }
        sync {
            latest_block_height
            catching_up
        }
        num_peers
        peers {
            is_outbound
            remote_ip
        }
        disk_usage
    }
}
```

Full node status:

```graphql
{
    getStatus {
        version
        node {
            id
            network
            moniker
        }
        sync {
            latest_block_hash
            latest_block_time
            latest_block_height
            catching_up
        }
        validator {
            address
            voting_power
            proposer_priority
        }
        validators {
            address
            voting_power
            proposer_priority
        }
        num_peers
        peers {
            node {
                id
                network
                moniker
            }
            is_outbound
            remote_ip
        }
        disk_usage
    }
}
```

Get records by IDs.

```graphql
{
    getRecordsByIds(ids: ["QmYDtNCKtTu6u6jaHaFAC5PWZXcj7fAmry6NoWwMaixFHz"]) {
        id
        names
        bondId
        createTime
        expiryTime
        owners
        attributes {
            key
            value {
                string
            }
        }
    }
}
```

Query records.

```graphql
{
    queryRecords(attributes: [{ key: "type", value: { string: "crn:bot" } }]) {
        id
        names
        bondId
        createTime
        expiryTime
        owners
        attributes {
            key
            value {
                string
            }
        }
    }
}
```

Get account details:

```graphql
{
    getAccounts(addresses: ["cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094"]) {
        address
        pubKey
        number
        sequence
        balance {
            type
            quantity
        }
    }
}
```

Query bonds:

```graphql
{
    queryBonds(
        attributes: [
            {
                key: "owner"
                value: { string: "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094" }
            }
        ]
    ) {
        id
        owner
        balance {
            type
            quantity
        }
    }
}
```

Get bonds by IDs.

```graphql
{
    getBondsByIds(ids :
    [
        "1c2b677cb2a27c88cc6bf8acf675c94b69051125b40c4fd073153b10f046dd87",
        "c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440"
    ])
    {
        id
        owner
        balance{
            type
            quantity
        }
    }
}
```

Query Bonds by Owner

```graphql
{
    queryBondsByOwner(ownerAddresses: ["ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk"])
    {
        owner
        bonds{
            id
            owner
            balance
            {
                type
                quantity
            }
        }
    }
}
```

Query auctions by ids

```graphql
{
    getAuctionsByIds(ids: ["be98f2073c246194276554eefdb4c95b682a35a0f06fbe619a6da57c10c93e90"]){
        id
        ownerAddress
        createTime
        minimumBid{
            type
            quantity
        }
        commitFee{
            type
            quantity
        }
        commitsEndTime
        revealFee{
            type
            quantity
        }
        revealsEndTime
        winnerBid{
            type
            quantity
        }
        winnerPrice{
            type
            quantity
        }
        winnerAddress
        bids{
            bidderAddress
            commitHash
            commitTime
            commitFee{
                type
                quantity
            }
            revealFee{
                type
                quantity
            }
            revealTime
            bidAmount{
                type
                quantity
            }
            status
        }
    }
}
```

LookUp Authorities

```graphql
{
    lookupAuthorities(names: []){
        ownerAddress
        ownerAddress
        height
        bondId
        status
        expiryTime
        auction {
            id
            ownerAddress
            createTime
            minimumBid{
                type
                quantity
            }
            commitFee{
                type
                quantity
            }
            commitsEndTime
            revealFee{
                type
                quantity
            }
            revealsEndTime
            winnerBid{
                type
                quantity
            }
            winnerPrice{
                type
                quantity
            }
            winnerAddress
            bids{
                bidderAddress
                commitHash
                commitTime
                commitFee{
                    type
                    quantity
                }
                revealFee{
                    type
                    quantity
                }
                revealTime
                bidAmount{
                    type
                    quantity
                }
                status
            }
        }
    }
}
```

LookUp Names

```graphql
{
    lookupNames(names: ["crn://hello/test"]){
        latest{
            id
            height
        }
        history{
            id
            height
        }
    }
}
```

Resolve Names

```graphql
{
    resolveNames(names: ["asd"]){
        id
        names
        bondId
        createTime
        expiryTime
        owners
        attributes {
            key
            value {
                string
            }
        }
    }
}
```