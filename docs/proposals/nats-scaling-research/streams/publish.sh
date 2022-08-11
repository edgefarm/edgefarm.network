#!/bin/bash

for i in {1..20}; do
  nats pub "data.test" -s nats://127.0.0.1:4222 --user=acc${i} --password=acc --count=10000 "this is fucking hell from acc${i}" &
done
