#!/bin/sh

echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
chibaclonkd validate-genesis --home /chibaclonk

echo "starting chibaclonk node $ID in background ..."
chibaclonkd start \
--home /chibaclonk \
--keyring-backend test

echo "started ethermint node"
tail -f /dev/null