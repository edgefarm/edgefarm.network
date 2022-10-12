#!/bin/bash
set -e
COMMAND=${1}

TALOS_NODE_CIDR=$(kubectl get nodes -o wide | grep -m1 talos | awk '{print $6}' | cut -d '.' -f1-3 | sed 's/$/\.1/')
TALOS_INTERFACE=$(ip -br -4 a sh | grep ${TALOS_NODE_CIDR} | awk '{print $1}')
export EDGEMESH_NODE_IP=$(kubectl get pods  -n kubeedge -o wide | grep edgemesh-server | awk '{print $6}')

CLOUDCOREPIDFILE=cloudcore-port-forward.pid
EDGEMESHPIDFILE=edgemesh-port-forward.pid

DEFAULT_INTERFACE=$(route | grep default | awk '{print $8}')
HOST_IP=$(ip -4 -br addr show | grep ${DEFAULT_INTERFACE} | awk '{print $3}' | cut -d '/' -f1)

case "$COMMAND" in
    start)
        sudo iptables -t nat -I PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10000 -j DNAT --to-destination 127.0.0.1:10000
        sudo iptables -t nat -I PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10002 -j DNAT --to-destination 127.0.0.1:10002
        sudo iptables -t nat -I PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10003 -j DNAT --to-destination 127.0.0.1:10003
        sudo iptables -t nat -I PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10004 -j DNAT --to-destination 127.0.0.1:10004
        sudo sysctl -w net.ipv4.conf.${DEFAULT_INTERFACE}.route_localnet=1
        kubectl port-forward -n kubeedge svc/cloudcore 10000:10000 10002:10002 10003:10003 10004:10004 &
        PID=$!
        echo $PID > ${CLOUDCOREPIDFILE}
        docker-compose up -d
    ;;
    stop)
        sudo iptables -t nat -D PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10000 -j DNAT --to-destination 127.0.0.1:10000
        sudo iptables -t nat -D PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10002 -j DNAT --to-destination 127.0.0.1:10002
        sudo iptables -t nat -D PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10003 -j DNAT --to-destination 127.0.0.1:10003
        sudo iptables -t nat -D PREROUTING -p tcp -d ${HOST_IP}/24 --dport 10004 -j DNAT --to-destination 127.0.0.1:10004
        sudo kill $(cat ${CLOUDCOREPIDFILE})
        sudo iptables -t nat -v -L PREROUTING -n --line-number
        docker kill proxy-20004
        docker rm proxy-20004
    ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 1
    ;;
esac
