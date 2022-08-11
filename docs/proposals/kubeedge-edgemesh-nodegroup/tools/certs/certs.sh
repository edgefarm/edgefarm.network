#!/bin/bash
set -e

COMMAND=${1}
SUBCOMMAND=${2}
OUTPUT=${3:-./}

subject="/C=DE/CN=edgefarm.io"
RED='\033[0;31m'
NC='\033[0m' # No Color
DEFAULT_INTERFACE=$(route | grep default | awk '{print $8}')
HOST_IP=$(ip -4 -br addr show | grep ${DEFAULT_INTERFACE} | awk '{print $3}' | cut -d '/' -f1)

createKeys() {
    name=${1}
    createPrivateKey ${name}
    createPublicKey ${name}
}

createPublicKey() {
    name=${1}
    openssl ec -in ${name}.key -pubout -out ${name}-pubkey.pem
}

createPrivateKey() {
    name=${1}
    openssl ecparam -name prime256v1 -genkey -noout -out ${OUTPUT}/${name}.key
}

createRootCert() {
    openssl ecparam -name prime256v1 -genkey -noout -out ${OUTPUT}/rootCA.key
    openssl req -x509 -new -nodes -key ${OUTPUT}/rootCA.key -days 3650 \
    -subj ${subject} -out ${OUTPUT}/rootCA.pem
}

createClientCert() {
    name=${1}
    openssl req -new -key ${OUTPUT}/${name}.key -subj ${subject} -out ${OUTPUT}/${name}.csr
    
    echo "subjectAltName=DNS:${HOST_IP},DNS:*.nip.io,DNS:*.kubeedge.svc.cluster.local" > /tmp/server-extfile.cnf
    openssl x509 -req -in ${OUTPUT}/${name}.csr -out ${OUTPUT}/${name}.pem -CAcreateserial -days 360  \
    -CA ${OUTPUT}/rootCA.pem -CAkey ${OUTPUT}/rootCA.key \
    -extfile /tmp/server-extfile.cnf
    rm rootCA.srl
    rm ${name}.csr
}

case "$COMMAND" in
    cert)
        case "$SUBCOMMAND" in
            create-all)
                echo -e "${RED}Creating CA certificate${NC}"
                createKeys "rootCA"
                createRootCert
                
                echo  -e "${RED}Creating Server certificate${NC}"
                createKeys "server"
                createClientCert "server"
                
                echo  -e "${RED}Creating Node certificate${NC}"
                createKeys "node"
                createClientCert "node"
            ;;
            ca)
                echo -e "${RED}Creating CA certificate${NC}"
                createKeys "rootCA"
                createRootCert
            ;;
            server)
                echo  -e "${RED}Creating Server certificate${NC}"
                createKeys "server"
                createClientCert "server"
            ;;
            node)
                echo  -e "${RED}Creating Node certificate${NC}"
                createKeys "node"
                createClientCert "node"
            ;;
            clean)
                echo  -e "${RED}Cleaning up${NC}"
                rm -f *.pem *.key *.csr *.srl
            ;;
            *)
                echo "Usage: $0 cert {create-all|ca|server|node|clean} <subcommand> [output-dir]"
                exit 1
            ;;
        esac
    ;;
    secret)
        case "$SUBCOMMAND" in
            create)
                name=${3}
                namespace=${4}
                pem=${5}
                key=${6}
                kubectl create secret tls ${name} -n ${namespace} --cert=${pem} --key=${key}
            ;;
            delete)
                name=${3}
                namespace=${4}
                kubectl delete secret ${name} -n ${namespace}
            ;;
            *)
                echo "Usage: $0 secret {create|delete} <secret-name> <namespace> [pem-file] [key-file]"
                exit 1
            ;;
        esac
    ;;
    *)
        echo "Usage: $0 {cert|secret}"
        exit 1
    ;;
    
esac
