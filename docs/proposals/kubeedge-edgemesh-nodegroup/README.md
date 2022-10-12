# kubeedge edgemesh nodegroup evaluation

The goal of this evaluation is to verify the nodegroup and edgeapplication features of kubeedge.
This is needed, because edge nodes gain the ability to resolve locally to kubernets service addresses and stay locally on the device if a pod is scheduled that matches the service.

The old fashion of edgefarm.applications defined a DaemonSet that matched using a combination of tolerations and affinity rules to schedule the pods on the specific nodes.
Now, the idea is to put one nodegroup per edge node and define an EdgeApplication containing a Deployment and Service. The big benefit of this is that the old component `node-dns` can die, because from now on it would be possible to use service addresses for DNS resolution directly.

## Preparations

First, let's create cluster and the certificates needed. For this step roor permissions are needed.

```sh
# Note, that it takes the network interface that has the default route to determine the IP address the physical edge nodes will connect to.
# So please make sure that the any edge device you want to bring to this cluster is part of this network too.
# If you are using virtual-edge-nodes, just go on and ignore everything you've read in this comment before.
devspace run init
```

## Deployment

```sh
$ devspace deploy
```

Now you have a fully operational cluster with installed kubeedge and edgemesh.

## Provisioning Nodes

### Physical edge nodes

copy the following files to the edge nodes:

Example devices has IP `192.168.1.100`

```sh
# Provision edgemesh certs
IP=192.168.1.100
scp tools/certs/node* root@${IP}:/etc/kubeedge/edgemesh/certs
scp tools/certs/rootCA.pem root@${IP}:/etc/kubeedge/edgemesh/ca

# Create certs for edge node
mkdir -p nodes/ca
mkdir -p nodes/certs
kubectl get secrets -n kubeedge kubeedge-ca -o 'go-template={{index .data "tls.crt"}}' | base64 -d > nodes/ca/rootCA.pem
kubectl get secrets -n kubeedge kubeedge-ca -o 'go-template={{index .data "tls.key"}}' | base64 -d > nodes/ca/rootCA-key.pem
export CAROOT=nodes/ca OUTPUT=nodes/certs
mkcert -client -cert-file ${OUTPUT}/node.pem -key-file ${OUTPUT}/node.key "*.nip.io"

# Provision edge node certs
scp nodes/ca/rootCA.pem nodes/certs/* root@${IP}:/etc/kubeedge/certs
```

You need to start the port forwarding to allow access to the physical edge nodes to the services that are running on your local cluster.

```sh
devspace run port-forward-start
```

Check for your edge nodes by listing all cluster nodes.
    
```sh
$ kubectl get nodes -o wide                                                                        
NAME                                STATUS   ROLES           AGE    VERSION                    INTERNAL-IP     EXTERNAL-IP   OS-IMAGE                                                 KERNEL-VERSION          CONTAINER-RUNTIME
edge0                               Ready    agent,edge      13s    v1.22.6-kubeedge-v1.11.1   192.168.1.100   <none>        Poky (Yocto Project Reference Distro) 3.1.14 (dunfell)   5.4.72-v8               docker://19.3.8
talos-default-controlplane-1        Ready    control-plane   68m    v1.22.9                    10.5.0.2        <none>        Talos (v1.2.0-beta.0)                                    5.8.18-050818-generic   containerd://1.6.7
talos-default-controlplane-2        Ready    control-plane   68m    v1.22.9                    10.5.0.3        <none>        Talos (v1.2.0-beta.0)                                    5.8.18-050818-generic   containerd://1.6.7
talos-default-worker-1              Ready    <none>          68m    v1.22.9                    10.5.0.4        <none>        Talos (v1.2.0-beta.0)                                    5.8.18-050818-generic   containerd://1.6.7
talos-default-worker-2              Ready    <none>          68m    v1.22.9                    10.5.0.5        <none>        Talos (v1.2.0-beta.0)                                    5.8.18-050818-generic   containerd://1.6.7
```

Now start edgecore and see the device showing up as node.
*Note: if you already provisioned a device, you need to remove everything related to edgecore and try again.*

### Virtual edge nodes

```sh
# Clone the repo containing the virtual edge nodes
git clone https://github.com/edgefarm/virtual-edge-node.git
cd virtual-edge-node

# Create the virtual edge nodes certs
kubectl get secrets -n kubeedge kubeedge-ca -o 'go-template={{index .data "tls.crt"}}' | base64 -d > example/.env/ca/rootCA.pem
kubectl get secrets -n kubeedge kubeedge-ca -o 'go-template={{index .data "tls.key"}}' | base64 -d > example/.env/ca/rootCA-key.pem
export CAROOT=example/.env/ca OUTPUT=example/.env/node/
mkcert -client -cert-file ${OUTPUT}/node.pem -key-file ${OUTPUT}/node.key "*.nip.io" "cloudcore.kubeedge.svc.cluster.local"

# Deploy the nodes
kustomize build example | kubectl apply -f -
```

Check for your edge nodes by listing all cluster nodes.

```sh
$ kubectl get nodes -o wide                                                                        
NAME                                STATUS   ROLES           AGE    VERSION                    INTERNAL-IP     EXTERNAL-IP   OS-IMAGE                KERNEL-VERSION          CONTAINER-RUNTIME
talos-default-controlplane-1        Ready    control-plane   68m    v1.22.9                    10.5.0.2        <none>        Talos (v1.2.0-beta.0)   5.8.18-050818-generic   containerd://1.6.7
talos-default-controlplane-2        Ready    control-plane   68m    v1.22.9                    10.5.0.3        <none>        Talos (v1.2.0-beta.0)   5.8.18-050818-generic   containerd://1.6.7
talos-default-worker-1              Ready    <none>          68m    v1.22.9                    10.5.0.4        <none>        Talos (v1.2.0-beta.0)   5.8.18-050818-generic   containerd://1.6.7
talos-default-worker-2              Ready    <none>          68m    v1.22.9                    10.5.0.5        <none>        Talos (v1.2.0-beta.0)   5.8.18-050818-generic   containerd://1.6.7
virtual-edge-node-68bd6dfc9-5wbcn   Ready    agent,edge      118s   v1.22.6-kubeedge-v1.11.1   10.244.3.2      <none>        Ubuntu 20.04 LTS        5.8.18-050818-generic   docker://19.3.11
virtual-edge-node-68bd6dfc9-cbzsx   Ready    agent,edge      117s   v1.22.6-kubeedge-v1.11.1   10.244.2.3      <none>        Ubuntu 20.04 LTS        5.8.18-050818-generic   docker://19.3.11
virtual-edge-node-68bd6dfc9-zsqst   Ready    agent,edge      118s   v1.22.6-kubeedge-v1.11.1   10.244.3.3      <none>        Ubuntu 20.04 LTS        5.8.18-050818-generic   docker://19.3.11
```

## Testing

See the testing docs for more information.
[test/local-service](test/local-service/README.md)
[test/nodegroup](test/nodegroup/README.md)

## Cleanup

```sh
# If you've started the port forwarding, you can stop it now.
$ devspace run port-forward-stop

# Remove the cluster
$ devspace run purge
```

