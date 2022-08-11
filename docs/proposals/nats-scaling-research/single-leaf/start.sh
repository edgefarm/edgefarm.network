#!/bin/bash

# create docker network 'nats' if it does not exist yet
docker network inspect nats >/dev/null 2>&1 || \
    docker network create nats

docker run --rm -d -it --network nats --name single-leaf-nats -v $(pwd)/store_leaf:/store_leaf -v $(pwd)/../common-config/:/common-config/ -v $(pwd)/config/:/config/ -p 4111:4111 nats:latest -c /config/leaf.conf
