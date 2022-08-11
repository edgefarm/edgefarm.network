# Node Group example

The EdgeApplication deploys a Pod with a `http-echo` server container and a `curl` container that can make requests to the server.

Expected Result:
Requests stay on the node if a corresponding `echo` server is running on the same node.

## Deploy examples using virtual nodes

```sh
# View the nodes in the cluster
$ kubectl get nodes                             
NAME                                 STATUS   ROLES           AGE     VERSION
edge0                               Ready    agent,edge      34m    v1.22.6-kubeedge-v1.11.1
talos-default-controlplane-1        Ready    control-plane   103m   v1.22.9
talos-default-controlplane-2        Ready    control-plane   103m   v1.22.9
talos-default-worker-1              Ready    <none>          103m   v1.22.9
talos-default-worker-2              Ready    <none>          103m   v1.22.9
virtual-edge-node-68bd6dfc9-5wbcn   Ready    agent,edge      36m    v1.22.6-kubeedge-v1.11.1
virtual-edge-node-68bd6dfc9-cbzsx   Ready    agent,edge      36m    v1.22.6-kubeedge-v1.11.1
virtual-edge-node-68bd6dfc9-zsqst   Ready    agent,edge      36m    v1.22.6-kubeedge-v1.11.1

# Create namespace where everything happens
$ kubectl create namespace nodegroup

# Create yaml files for the NodeGroups for all edge nodes (virtual and physical)
$ ./create_node_groups.sh

# Apply the nodeGroups
$ for yaml in $(ls node-group*.yaml); do kubectl apply -f $yaml; done

# Create yaml files for the EdgeApplications containing all nodegroups
$ ./create_application.sh

# Apply the edgeApplication
$ kubectl apply -f edgeApplication.yaml

# Listing resources generated previously...
$ kubectl get edgeapplication
NAME        AGE
nginx-app   4s

$ kubectl get nodegroups
NAME                                 AGE
virtual-edge-node-86bf8d4fb9-bgwdv   17s
virtual-edge-node-86bf8d4fb9-gdj5z   17s
virtual-edge-node-86bf8d4fb9-kcmct   16s

$ kubectl get deployments.apps                                
NAME                                             READY   UP-TO-DATE   AVAILABLE   AGE
example-app-virtual-edge-node-86bf8d4fb9-bgwdv   1/1     1            1           18s
example-app-virtual-edge-node-86bf8d4fb9-gdj5z   1/1     1            1           18s
example-app-virtual-edge-node-86bf8d4fb9-kcmct   1/1     1            1           18s

# ... and the resulting resources from the deployment and service
$ kubectl get services
NAME          TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
example-svc   ClusterIP   10.108.27.1   <none>        8080/TCP   35s

$ kubectl get pods -o wide
example-app-edge0-69869b9cf5-r824j                               2/2     Running   0          33s   172.17.0.4   edge0                               <none>           <none>
example-app-virtual-edge-node-68bd6dfc9-5wbcn-589c949c5c-96hjb   2/2     Running   0          33s   172.17.0.4   virtual-edge-node-68bd6dfc9-5wbcn   <none>           <none>
example-app-virtual-edge-node-68bd6dfc9-cbzsx-749d647d98-c2vvg   2/2     Running   0          33s   172.17.0.4   virtual-edge-node-68bd6dfc9-cbzsx   <none>           <none>
example-app-virtual-edge-node-68bd6dfc9-zsqst-858c5bcd6c-9xcl8   2/2     Running   0          33s   172.17.0.4   virtual-edge-node-68bd6dfc9-zsqst   <none>           <none>

```

## Testing the setup 

This shows that requests stay within the node.

```sh
# (NODE A) First check logs of echo server if it received a request yet
$ kubectl logs example-app-virtual-edge-node-86bf8d4fb9-bgwdv-697b47b85c-f5hmv -c echo
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.

# (NODE A) Do the actual request
$ kubectl exec -it example-app-edge0-69869b9cf5-r824j -c curl --  curl http://example-svc.nodegroup.svc.cluster.local:8080/\?echo_env_body\=NODE_NAME  
"edge0"

# (NODE A) Check again if the request got logged in the echo server
$ kubectl logs example-app-edge0-69869b9cf5-r824j -c echo
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.
{"name":"echo-server","hostname":"example-app-edge0-69869b9cf5-r824j","pid":1,"level":30,"host":{"hostname":"example-svc.nodegroup.svc.cluster.local","ip":"::ffff:172.17.0.1","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/?echo_env_body=NODE_NAME","protocol":"http"},"request":{"params":{},"query":{"echo_env_body":"NODE_NAME"},"cookies":{},"body":{},"headers":{"host":"example-svc.nodegroup.svc.cluster.local:8080","user-agent":"curl/7.84.0-DEV","accept":"*/*"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"example-app-edge0-69869b9cf5-r824j","NODE_NAME":"edge0","NODE_VERSION":"16.16.0","YARN_VERSION":"1.22.19","HOME":"/root"},"msg":"Thu, 25 Aug 2022 10:14:23 GMT | [GET] - http://example-svc.nodegroup.svc.cluster.local:8080/?echo_env_body=NODE_NAME","time":"2022-08-25T10:14:23.617Z","v":0}


# (NODE A) Send another two requests to the same server
$ kubectl exec -it example-app-edge0-69869b9cf5-r824j -c curl --  curl http://example-svc.nodegroup.svc.cluster.local:8080/\?echo_env_body\=NODE_NAME  
"edge0"
$ kubectl exec -it example-app-edge0-69869b9cf5-r824j -c curl --  curl http://example-svc.nodegroup.svc.cluster.local:8080/\?echo_env_body\=NODE_NAME  
"edge0"

$ kubectl logs example-app-virtual-edge-node-86bf8d4fb9-bgwdv-697b47b85c-vndbv -c echo 
$ kubectl logs example-app-edge0-69869b9cf5-r824j -c echo
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.
{"name":"echo-server","hostname":"example-app-edge0-69869b9cf5-r824j","pid":1,"level":30,"host":{"hostname":"example-svc.nodegroup.svc.cluster.local","ip":"::ffff:172.17.0.1","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/?echo_env_body=NODE_NAME","protocol":"http"},"request":{"params":{},"query":{"echo_env_body":"NODE_NAME"},"cookies":{},"body":{},"headers":{"host":"example-svc.nodegroup.svc.cluster.local:8080","user-agent":"curl/7.84.0-DEV","accept":"*/*"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"example-app-edge0-69869b9cf5-r824j","NODE_NAME":"edge0","NODE_VERSION":"16.16.0","YARN_VERSION":"1.22.19","HOME":"/root"},"msg":"Thu, 25 Aug 2022 10:14:23 GMT | [GET] - http://example-svc.nodegroup.svc.cluster.local:8080/?echo_env_body=NODE_NAME","time":"2022-08-25T10:14:23.617Z","v":0}
{"name":"echo-server","hostname":"example-app-edge0-69869b9cf5-r824j","pid":1,"level":30,"host":{"hostname":"example-svc.nodegroup.svc.cluster.local","ip":"::ffff:172.17.0.1","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/?echo_env_body=NODE_NAME","protocol":"http"},"request":{"params":{},"query":{"echo_env_body":"NODE_NAME"},"cookies":{},"body":{},"headers":{"host":"example-svc.nodegroup.svc.cluster.local:8080","user-agent":"curl/7.84.0-DEV","accept":"*/*"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"example-app-edge0-69869b9cf5-r824j","NODE_NAME":"edge0","NODE_VERSION":"16.16.0","YARN_VERSION":"1.22.19","HOME":"/root"},"msg":"Thu, 25 Aug 2022 10:14:23 GMT | [GET] - http://example-svc.nodegroup.svc.cluster.local:8080/?echo_env_body=NODE_NAME","time":"2022-08-25T10:14:23.617Z","v":0}

# (NODE B) Send requests on another node B and check back to node A after that
$ kubectl  logs example-app-virtual-edge-node-68bd6dfc9-5wbcn-589c949c5c-96hjb -c echo                                                    
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.

# (NODE B) Send request
$ kubectl  exec -it example-app-virtual-edge-node-68bd6dfc9-5wbcn-589c949c5c-96hjb -c curl --  curl http://example-svc.nodegroup.svc.cluster.local:8080/\?echo_env_body\=NODE_NAME
"virtual-edge-node-68bd6dfc9-5wbcn"

# (NODE B) Check logs of http server
$ kubectl logs example-app-virtual-edge-node-68bd6dfc9-5wbcn-589c949c5c-96hjb -c echo                                                                                            
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.
{"name":"echo-server","hostname":"example-app-virtual-edge-node-68bd6dfc9-5wbcn-589c949c5c-96hjb","pid":1,"level":30,"host":{"hostname":"example-svc.nodegroup.svc.cluster.local","ip":"::ffff:172.17.0.1","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/?echo_env_body=NODE_NAME","protocol":"http"},"request":{"params":{},"query":{"echo_env_body":"NODE_NAME"},"cookies":{},"body":{},"headers":{"host":"example-svc.nodegroup.svc.cluster.local:8080","user-agent":"curl/7.84.0-DEV","accept":"*/*"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"example-app-virtual-edge-node-68bd6dfc9-5wbcn-589c949c5c-96hjb","NODE_NAME":"virtual-edge-node-68bd6dfc9-5wbcn","NODE_VERSION":"16.16.0","YARN_VERSION":"1.22.19","HOME":"/root"},"msg":"Thu, 25 Aug 2022 10:28:46 GMT | [GET] - http://example-svc.nodegroup.svc.cluster.local:8080/?echo_env_body=NODE_NAME","time":"2022-08-25T10:28:46.530Z","v":0}

# (NODE A) Check back to http server logs on NODE A
$ kubectl logs example-app-edge0-69869b9cf5-r824j -c echo
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.
{"name":"echo-server","hostname":"example-app-edge0-69869b9cf5-r824j","pid":1,"level":30,"host":{"hostname":"example-svc.nodegroup.svc.cluster.local","ip":"::ffff:172.17.0.1","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/?echo_env_body=NODE_NAME","protocol":"http"},"request":{"params":{},"query":{"echo_env_body":"NODE_NAME"},"cookies":{},"body":{},"headers":{"host":"example-svc.nodegroup.svc.cluster.local:8080","user-agent":"curl/7.84.0-DEV","accept":"*/*"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"example-app-edge0-69869b9cf5-r824j","NODE_NAME":"edge0","NODE_VERSION":"16.16.0","YARN_VERSION":"1.22.19","HOME":"/root"},"msg":"Thu, 25 Aug 2022 10:14:23 GMT | [GET] - http://example-svc.nodegroup.svc.cluster.local:8080/?echo_env_body=NODE_NAME","time":"2022-08-25T10:14:23.617Z","v":0}
{"name":"echo-server","hostname":"example-app-edge0-69869b9cf5-r824j","pid":1,"level":30,"host":{"hostname":"example-svc.nodegroup.svc.cluster.local","ip":"::ffff:172.17.0.1","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/?echo_env_body=NODE_NAME","protocol":"http"},"request":{"params":{},"query":{"echo_env_body":"NODE_NAME"},"cookies":{},"body":{},"headers":{"host":"example-svc.nodegroup.svc.cluster.local:8080","user-agent":"curl/7.84.0-DEV","accept":"*/*"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"example-app-edge0-69869b9cf5-r824j","NODE_NAME":"edge0","NODE_VERSION":"16.16.0","YARN_VERSION":"1.22.19","HOME":"/root"},"msg":"Thu, 25 Aug 2022 10:14:23 GMT | [GET] - http://example-svc.nodegroup.svc.cluster.local:8080/?echo_env_body=NODE_NAME","time":"2022-08-25T10:14:23.617Z","v":0}

# (NODE C) The third node was not used for the requests, so no logs are available
$ kubectl logs example-app-virtual-edge-node-68bd6dfc9-zsqst-858c5bcd6c-9xcl8 -c echo
(node:1) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
Listening on port 80.
```

## Results

The results match the expected. Requests from a Pod to a Service that were both deployed via the edgeApplication stay on the node. The Pod can resolve to the `echo` server through the Service address.