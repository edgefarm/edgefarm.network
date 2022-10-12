#!/bin/bash

# create docker network 'nats' if it does not exist yet
docker network inspect nats >/dev/null 2>&1 || \
    docker network create nats

for i in {1..20}; do
  PORT=$(printf "43%02d\n" $i)
  docker run --rm -d -it --network nats --name leaf-acc${i}-nats -v $(pwd)/store_leaf/leaf-acc${i}:/store_leaf -v $(pwd)/../common-config/:/common-config/ -v $(pwd)/config/leaf-acc${i}.config:/config/leaf.conf -p ${PORT}:4111 nats:latest -c /config/leaf.conf
done
