#!/bin/bash
set -e

nodesList=$(kubectl get nodes -l node-role.kubernetes.io/agent= -o jsonpath='{.items[*].metadata.name}')
nodesArray=($nodesList)

for node in ${nodesArray[@]}; do
    echo ${node}
    kubectl label node ${node} node.edgefarm.io=${node} --overwrite
cat <<EOF > node-group-${node}.yaml
apiVersion: apps.kubeedge.io/v1alpha1
kind: NodeGroup
metadata:
  name: ${node}
  labels:
    group: ${node}
spec:
  matchLabels:
    node.edgefarm.io: ${node}
EOF
done
