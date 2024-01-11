#!/bin/bash

##
## This script generates a new account keypair, transfers funds to it, and creates
## a bond with those funds.  The amount of can be optionally specified, else a
## default value estimated to be sufficient for about 5000 records.
##
## The `laconic` CLI and a valid config file for it are required.  The default location
## for the config file is `~/.laconic/config.yml`, but this can be controlled with
## the environment variable LACONIC_CONFIG.  A `userKey` for a source account with
## sufficient funds available must be present in this file for the funds transfer
## to succeed.
##
## Example:
##
## ❯ scripts/create-and-fund-account.sh 1000000000
## {
##   "name": "68334d7175fd4f86befa4902657e5270",
##   "type": "local",
##   "address": "ethm15r5x94km0swq55aszwd7hnr9wksq7wmr38xes7",
##   "pubkey": "AuKqlSldJJXj4gYMFt2HeX9DJ3aUosYA7n6zBz9Tg7/i",
##   "mnemonic": "umbrella bean special unaware accident giant distance mix ghost feel possible cost road grant endless man maple derive rebuild learn mask water attract resist",
##   "bond": "3d3a73f09115d289d330781455e6eac217305dc4a20e19bde808011fe3775a93",
##   "balance": 1000000000,
##   "privkey": "480880fde7aff1461da584b436cb3a84692413c84623fda78e127bb4e704ce76"
## }
##

AVG_RECORD_PHOTON=1000000
NUM_RECORDS=5000
BOND_OVERHEAD=200000
KEYNAME=`uuidgen | tr -d '-'`
KEYRING_DIR=`mktemp -d`
KEYRING=test
LACONIC_CONFIG=${LACONIC_CONFIG:-$HOME/.laconic/config.yml}
BOND_AMOUNT=${1:-$((AVG_RECORD_PHOTON * NUM_RECORDS))}

ACCOUNT_JSON=$(laconicd keys add $KEYNAME --keyring-backend $KEYRING --algo eth_secp256k1 --keyring-dir $KEYRING_DIR --output json)
PRIVATE_KEY=$(yes | laconicd keys export $KEYNAME --keyring-backend $KEYRING --keyring-dir $KEYRING_DIR --unarmored-hex --unsafe)
PUB_KEY=$(echo $ACCOUNT_JSON | jq -r ".pubkey | fromjson | .key")

laconicd keys delete $KEYNAME --keyring-backend $KEYRING --keyring-dir $KEYRING_DIR -y 2> /dev/null
rm -rf $KEYRING_DIR

laconic -c $LACONIC_CONFIG cns tokens send --address $(echo $ACCOUNT_JSON | jq -r '.address') --type aphoton --quantity $((BOND_AMOUNT + BOND_OVERHEAD)) > /dev/null
BOND_ID=$(laconic -c $LACONIC_CONFIG cns bond create --user-key $PRIVATE_KEY --type aphoton --quantity $BOND_AMOUNT | jq -r '.bondId')

echo $ACCOUNT_JSON | jq ".bond = \"$BOND_ID\"" | jq ".balance = $BOND_AMOUNT" | jq ".privkey = \"$PRIVATE_KEY\"" | jq ".pubkey = \"$PUB_KEY\""
