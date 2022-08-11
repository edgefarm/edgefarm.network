#!/bin/bash
RED='\033[0;31m'
RED_BOLD='\033[1;31m'
NC='\033[0m' # No Color

NS=vault-system
FULL_URI=http://127.0.0.1:48200

wait_for_token() {
    while true; do
        POD=$(kubectl get pods -n ${NS} --no-headers -o custom-columns=":metadata.name")
        ROOT_TOKEN=$(kubectl logs -n ${NS} ${POD} | grep "Root Token")
        if [ ! -z "${ROOT_TOKEN}" ]; then
            RET=$(echo ${ROOT_TOKEN} | awk -F " " '{print $3}')
            echo $RET
            break
        fi
        sleep 1
    done
}

login() {
    export VAULT_TOKEN=$(wait_for_token)
    export VAULT_ADDR=${FULL_URI}
    vault login ${VAULT_TOKEN}
}
