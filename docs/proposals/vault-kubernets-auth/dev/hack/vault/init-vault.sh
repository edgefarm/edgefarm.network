#!/bin/bash

source './common.sh'

login

vault secrets enable kv

echo -e "${RED}Writing policies${NC}"
vault policy write auth-example - << EOF
path "secret/data/creds" {
    capabilities = [ "read" ]
}
EOF

echo -e "${RED}Enabling kubernets auth${NC}"
vault auth enable kubernetes

echo -e "${RED}Setting kubernets auth config${NC}"
TEMPDIR=$(mktemp -d)
IP=$(kubectl get svc -n default kubernetes -o jsonpath='{.spec.clusterIP}')
KUBERNETES_HOST_AND_PORT=$(echo ${IP}:443)
openssl s_client -showcerts -connect 127.0.0.1:34550 </dev/null 2>/dev/null |openssl x509 -outform PEM > ${TEMPDIR}/ca.crt
vault write auth/kubernetes/config kubernetes_host=https://${KUBERNETES_HOST_AND_PORT} kubernetes_ca_crt=${TEMPDIR}/ca.crt
rm -r ${TEMPDIR}

echo -e "${RED}Create the binding role for the auth-example service account${NC}"
vault write auth/kubernetes/role/auth-example \
bound_service_account_names=auth-example \
bound_service_account_namespaces=auth-example \
token_policies=auth-example \
alias_name_source=serviceaccount_name \
token_no_default_policy=true

