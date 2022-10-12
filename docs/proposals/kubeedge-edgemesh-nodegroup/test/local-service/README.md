# local-service test

This test shows the basic edgemesh functionality. 
Expeted result:

Accessing services on edge nodes within the same network will just work as expected. Accessing cloud services will work find too.

## Deploying

```sh
# Deploy the local-service example
kubectl apply -f namespace.yaml
kubectl apply -f daemonset.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

$ kubectl get pods -o wide
NAME                    READY   STATUS    RESTARTS   AGE     IP           NODE                                NOMINATED NODE   READINESS GATES
alpine-2rnf8            1/1     Running   0          2m30s   172.17.0.2   virtual-edge-node-68bd6dfc9-cbzsx   <none>           <none>
alpine-2sz56            1/1     Running   0          2m30s   172.17.0.3   edge0                               <none>           <none>
alpine-2w9bw            1/1     Running   0          2m30s   172.17.0.2   virtual-edge-node-68bd6dfc9-5wbcn   <none>           <none>
alpine-hx7kv            1/1     Running   0          2m30s   172.17.0.3   virtual-edge-node-68bd6dfc9-zsqst   <none>           <none>
echo-7sph2              1/1     Running   0          22s     172.17.0.4   edge0                               <none>           <none>
echo-j2mbc              1/1     Running   0          22s     172.17.0.4   virtual-edge-node-68bd6dfc9-zsqst   <none>           <none>
echo-pb6ht              1/1     Running   0          22s     172.17.0.4   virtual-edge-node-68bd6dfc9-5wbcn   <none>           <none>
echo-zx5w2              1/1     Running   0          22s     172.17.0.4   virtual-edge-node-68bd6dfc9-cbzsx   <none>           <none>
echo-5888cb78c6-zff9x   1/1     Running   0          47s     10.244.2.4   talos-default-worker-1              <none>           <none>
```

## Manual testing

The alpine pod is used to access the service. The service connets to all echo pods that are running on different edge nodes.
Accessing the pods will result in a round robin behavior of the service, accessing all `echo` pods one after the other.
The result of the HTTP request returns the node name of the `echo` pod that answered the request.

```sh
# Repeatedly access the service on the edge node. This run uses the alpine pod that runs on the physical edge node.
$ kubectl -n local-service exec -it alpine-2sz56 -- curl http://echo-svc.local-service.svc.cluster.local:8080 | jq -r '.environment.NODE_NAME'
edge0

$ kubectl -n local-service exec -it alpine-2sz56 -- curl http://echo-svc.local-service.svc.cluster.local:8080 | jq -r '.environment.NODE_NAME'
virtual-edge-node-68bd6dfc9-zsqst

$ kubectl -n local-service exec -it alpine-2sz56 -- curl http://echo-svc.local-service.svc.cluster.local:8080 | jq -r '.environment.NODE_NAME'
virtual-edge-node-68bd6dfc9-cbzsx

$ kubectl -n local-service exec -it alpine-2sz56 -- curl http://echo-svc.local-service.svc.cluster.local:8080 | jq -r '.environment.NODE_NAME'
virtual-edge-node-68bd6dfc9-5wbcn

$ kubectl -n local-service exec -it alpine-2sz56 -- curl http://echo-svc.local-service.svc.cluster.local:8080 | jq -r '.environment.NODE_NAME'
talos-default-worker-1

$ kubectl -n local-service exec -it alpine-2sz56 -- curl http://echo-svc.local-service.svc.cluster.local:8080 | jq -r '.environment.NODE_NAME'
edge0
```

## Result

The result is as expected. Pods that are running on edge nodes can access the service no matter where it runs (cloud node, virtual edge node, physical edge node).
