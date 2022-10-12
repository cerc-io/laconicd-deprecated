#!/bin/bash

echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
laconicd validate-genesis --home /laconic

echo "starting ethermint node $ID in background ..."
laconicd start \
--home /laconic \
--keyring-backend test

echo "started ethermint node"
tail -f /dev/null