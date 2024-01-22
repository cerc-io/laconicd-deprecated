#!/bin/bash

KEY="mykey"
CHAINID="laconic_9000-1"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="${LOGLEVEL:-info}"
# trace evm
TRACE="--trace"
# TRACE=""

if [ "$1" == "clean" ] || [ ! -d "$HOME/.laconicd/data/blockstore.db" ]; then
  # validate dependencies are installed
  command -v jq > /dev/null 2>&1 || {
    echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
    exit 1
  }

  # remove existing daemon and client
  rm -rf $HOME/.laconicd/*
  rm -rf $HOME/.laconic/*

  if [ -n "`which make`" ]; then
    make install
  fi

  laconicd config keyring-backend $KEYRING
  laconicd config chain-id $CHAINID

  # if $KEY exists it should be deleted
  laconicd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

  # Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
  laconicd init $MONIKER --chain-id $CHAINID

  update_genesis() {
  jq "$1" $HOME/.laconicd/config/genesis.json > $HOME/.laconicd/config/tmp_genesis.json &&
    mv $HOME/.laconicd/config/tmp_genesis.json $HOME/.laconicd/config/genesis.json
  }

  # Change parameter token denominations to aphoton
  update_genesis '.app_state["staking"]["params"]["bond_denom"]="aphoton"'
  update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="aphoton"'
  update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aphoton"'
  update_genesis '.app_state["mint"]["params"]["mint_denom"]="aphoton"'
  # Custom modules
  update_genesis '.app_state["registry"]["params"]["record_rent"]["denom"]="aphoton"'
  update_genesis '.app_state["registry"]["params"]["authority_rent"]["denom"]="aphoton"'
  update_genesis '.app_state["registry"]["params"]["authority_auction_commit_fee"]["denom"]="aphoton"'
  update_genesis '.app_state["registry"]["params"]["authority_auction_reveal_fee"]["denom"]="aphoton"'
  update_genesis '.app_state["registry"]["params"]["authority_auction_minimum_bid"]["denom"]="aphoton"'

  if [[ "$TEST_REGISTRY_EXPIRY" == "true" ]]; then
    echo "Setting timers for expiry tests."

    update_genesis '.app_state["registry"]["params"]["record_rent_duration"]="60s"'
    update_genesis '.app_state["registry"]["params"]["authority_grace_period"]="60s"'
    update_genesis '.app_state["registry"]["params"]["authority_rent_duration"]="60s"'
  fi

  if [[ "$TEST_AUCTION_ENABLED" == "true" ]]; then
    echo "Enabling auction and setting timers."

    update_genesis '.app_state["registry"]["params"]["authority_auction_enabled"]=true'
    update_genesis '.app_state["registry"]["params"]["authority_rent_duration"]="60s"'
    update_genesis '.app_state["registry"]["params"]["authority_grace_period"]="300s"'
    update_genesis '.app_state["registry"]["params"]["authority_auction_commits_duration"]="60s"'
    update_genesis '.app_state["registry"]["params"]["authority_auction_reveals_duration"]="60s"'
  fi

  # increase block time (?)
  update_genesis '.consensus_params["block"]["time_iota_ms"]="1000"'

  # Set gas limit in genesis
  update_genesis '.consensus_params["block"]["max_gas"]="10000000"'

  # disable produce empty block
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.laconicd/config/config.toml
    else
      sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.laconicd/config/config.toml
  fi

  if [[ "$1" == "pending" ]]; then
    alias sed-i="sed -i"
    if [[ "$OSTYPE" == "darwin"* ]]; then
      alias sed-i="sed -i ''"
    fi
    sed-i \
      -e 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' \
      -e 's/timeout_propose = "3s"/timeout_propose = "30s"/g' \
      -e 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' \
      -e 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' \
      -e 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' \
      -e 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' \
      -e 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' \
      -e 's/timeout_commit = "5s"/timeout_commit = "150s"/g' \
      -e 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' \
      $HOME/.laconicd/config/config.toml
  fi

  # Allocate genesis accounts (cosmos formatted addresses)
  laconicd add-genesis-account $KEY 100000000000000000000000000aphoton --keyring-backend $KEYRING

  # Sign genesis transaction
  laconicd gentx $KEY 1000000000000000000000aphoton --keyring-backend $KEYRING --chain-id $CHAINID

  # Collect genesis tx
  laconicd collect-gentxs

  # Run this to ensure everything worked and that the genesis file is setup correctly
  laconicd validate-genesis

  if [[ "$1" == "pending" ]]; then
    echo "pending mode is on, please wait for the first block committed."
  fi
else
  echo "Using existing database at $HOME/.laconicd.  To replace, run '`basename $0` clean'"
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
laconicd start \
  --pruning=nothing \
  --evm.tracer=json $TRACE \
  --log_level $LOGLEVEL \
  --minimum-gas-prices=0.0001aphoton \
  --json-rpc.api eth,txpool,personal,net,debug,web3,miner \
  --api.enable \
  --gql-server --gql-playground
