#!/bin/bash
set -e
# nodesList=$(kubectl get nodes -l node-role.kubernetes.io/agent= -o jsonpath='{.items[*].metadata.name}')
nodesList=$(kubectl get nodegroups.apps.kubeedge.io -o json | jq -r '.items[].spec.matchLabels[]')
nodesArray=($nodesList)

OUTPUT="edgeApplication.yaml"
cat <<EOF > ${OUTPUT}
apiVersion: apps.kubeedge.io/v1alpha1
kind: EdgeApplication
metadata:
  name: nginx-app
spec:
  workloadTemplate:
    manifests:
      - apiVersion: apps/v1
        kind: Deployment
        metadata:
          name: example-app
          namespace: nodegroup
        spec:
          replicas: 1
          selector:
            matchLabels:
              app: example-app
          template:
            metadata:
              labels:
                app: example-app
            spec:
              containers:
                - env:
                    - name: NODE_NAME
                      valueFrom:
                        fieldRef:
                          fieldPath: spec.nodeName
                  image: ealen/echo-server:0.7.0
                  imagePullPolicy: IfNotPresent
                  name: echo
                - args:
                    - infinity
                  command:
                    - sleep
                  image: curlimages/curl:latest
                  imagePullPolicy: IfNotPresent
                  name: curl
              tolerations:
                - effect: NoExecute
                  key: edgefarm.applications
                  operator: Exists
      - apiVersion: v1
        kind: Service
        metadata:
          name: example-svc
          namespace: nodegroup
        spec:
          ports:
            - port: 8080
              protocol: TCP
              targetPort: 80
          selector:
            app: example-app
  workloadScope:
EOF
echo "    targetNodeGroups:" >> ${OUTPUT}

for node in ${nodesArray[@]}; do
    echo "      - name: ${node}" >> ${OUTPUT}
done

