#!/bin/bash
set -e

trap 'killall distrKV' SIGINT

cd $(dirname $0)

killall distrKV || true
sleep 0.1

go install -v

distrKV -loc=shard1.db -addr=127.0.0.1:8080 -config=sharding.toml -shard=shard1 &
distrKV -loc=shard2.db -addr=127.0.0.1:8081 -config=sharding.toml -shard=shard2 &
distrKV -loc=shard3.db -addr=127.0.0.1:8082 -config=sharding.toml -shard=shard3 &

wait