#!/bin/bash

# create docker network 'nats' if it does not exist yet
docker network inspect nats >/dev/null 2>&1 || \
    docker network create nats

docker run --rm -d --network nats -it --name main-nats -v $(pwd)/../common-config/:/common-config/ -v $(pwd)/config/:/config/ -v $(pwd)/creds:/creds -p 4222:4222 -p 7422:7422 nats:latest -c /config/server.conf
