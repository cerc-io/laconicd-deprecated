#!/bin/sh

echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
laconicd validate-genesis --home /laconic

echo "starting laconic node $ID in background ..."
laconicd start \
--home /laconic \
--keyring-backend test \ 
--mode validator

echo "started ethermint node"
tail -f /dev/null