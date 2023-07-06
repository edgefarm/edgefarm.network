#!/bin/sh

cp deploy/compositions/consumer.streams.network.edgefarm.io/definition.yaml charts/compositions/templates/definition-consumer.yaml
cp deploy/compositions/consumer.streams.network.edgefarm.io/composition-consumer.yaml charts/compositions/templates/composition-consumer.yaml

cp deploy/compositions/edgenetwork.streams.network.edgefarm.io/composition-edgenetworks.yaml charts/compositions/templates/composition-edgenetworks.yaml
cp deploy/compositions/edgenetwork.streams.network.edgefarm.io/definition.yaml charts/compositions/templates/definition-edgenetwork.yaml

cp deploy/compositions/network.streams.network.edgefarm.io/composition-networks.yaml charts/compositions/templates/composition-networks.yaml
cp deploy/compositions/network.streams.network.edgefarm.io/definition.yaml charts/compositions/templates/definition-network.yaml

cp deploy/compositions/stream.streams.network.edgefarm.io/composition-streams.yaml charts/compositions/templates/composition-streams.yaml
cp deploy/compositions/stream.streams.network.edgefarm.io/definition.yaml charts/compositions/templates/definition-streams.yaml
