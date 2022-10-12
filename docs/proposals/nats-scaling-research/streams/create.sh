#!/bin/bash

for i in {1..20}; do
  nats stream create -s nats://127.0.0.1:4222  --user=acc${i} --password=acc --domain acc${i} --config stream.json
  nats stream create -s nats://127.0.0.1:4222  --user=acc${i} --password=acc --domain leaf --config stream.json
done
