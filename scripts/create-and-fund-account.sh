#!/bin/bash

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
