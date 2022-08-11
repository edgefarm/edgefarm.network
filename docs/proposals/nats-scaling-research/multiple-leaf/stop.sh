#!/bin/bash
for i in {1..20}; do
  docker kill leaf-acc${i}-nats
done
