vault-kubernes-auth
===================

This part of the proposals docs is to show how the vault-kubernetes-auth plugin works. 
To spin up the environment simply run the following commands:

```sh
# Spin up the environment
$ devspace run init
$ devspace run build-example
$ devspace deploy
$ devspace run init-vault

# Check that the auth-example-pods are running
$ kubectl get pods -n auth-example
NAME                 READY   STATUS    RESTARTS   AGE   IP          NODE                        NOMINATED NODE   READINESS GATES
auth-example-94ljq   1/1     Running   0          99s   10.42.0.4   k3d-auth-example-server-0   <none>           <none>
auth-example-98w9l   1/1     Running   0          99s   10.42.1.5   k3d-auth-example-agent-0    <none>           <none>

# Check that the auth-example-pods are running not finding any secrets
$ kubectl logs auth-example-94ljq
error: unable to read secret: no secret found at secret/data/creds
sleeping for 5 seconds, then trying again
error: unable to read secret: no secret found at secret/data/creds
sleeping for 5 seconds, then trying again
error: unable to read secret: no secret found at secret/data/creds
sleeping for 5 seconds, then trying again

# Write the secret to vault and check that the pods are recieving it
$ devspace run write-secret

# Check that the auth-example-pods are running not finding any secrets
$ kubectl logs auth-example-94ljq
error: unable to read secret: no secret found at secret/data/creds
sleeping for 5 seconds, then trying again
secret: YouNeverGuessThis
sleeping for 5 seconds, then trying again
secret: YouNeverGuessThis
sleeping for 5 seconds, then trying again
```

