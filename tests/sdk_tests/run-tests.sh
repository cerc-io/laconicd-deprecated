#!/usr/bin/env bash
if [ -n "$CERC_SCRIPT_DEBUG" ]; then
    set -x
fi
# Get the key from laconicd
laconicd_key=$( docker compose exec laconicd echo y | docker compose exec laconicd laconicd keys export mykey --unarmored-hex --unsafe )
# Set parameters for the test suite
cosmos_chain_id=laconic_9000-1
laconicd_rest_endpoint=http://laconicd:1317
laconicd_gql_endpoint=http://laconicd:9473/api
# Run tests
docker network inspect sdk_tests_default
docker compose exec sdk-test-runner sh -c "COSMOS_CHAIN_ID=${cosmos_chain_id} LACONICD_REST_ENDPOINT=${laconicd_rest_endpoint} LACONICD_GQL_ENDPOINT=${laconicd_gql_endpoint} PRIVATE_KEY=${laconicd_key} yarn test"
