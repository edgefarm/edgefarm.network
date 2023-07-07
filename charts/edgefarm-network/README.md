# edgefarm.network

This helm chart installs edgefarm.network components. It installs several components:
  - provider-kubernetes
  - provider-nats
  - provider-natssecrets
  - network-compositions
  - network-resource-info
  - edgenetwork-operator
  - network-dependencies-webhook
  - nats

## Prerequisites

    Kubernetes 1.22+
    Helm 3.2.0+
    Crossplane 1.11.3+
    Vault with vault-plugin-secrets-nats 1.3.2+

# Needed Vault configuration

To be able to make this work a few things things have to be done in the vault configuration

## 1. configure kubernetes auth

See https://developer.hashicorp.com/vault/docs/auth/kubernetes and 
https://developer.hashicorp.com/vault/docs/auth/kubernetes#use-local-service-account-token-as-the-reviewer-jwt for more information.

If your Vault service runs in the same cluster, you simply can put the `kubernetes_host` to `https://10.96.0.1`.

```
$ KUBE_CA_CERT=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.certificate-authority-data}' | base64 --decode)
$ vault write auth/kubernetes/config kubernetes_host=https://10.96.0.1 kubernetes_ca_cert="$KUBE_CA_CERT" disable_local_ca_jwt="true"
```

Otherwise, you need to modify `kubernetes_host` and pass the correct CA cert.

## 2. create policy 

Replace each occurance of `<operatorName>` with the name of your operator:

```console
vault policy write nats-auth-config - <<EOF
path "nats-secrets/jwt/operator/<operatorName>" {
   capabilities = ["read"]
}
path "nats-secrets/nkey/operator/<operatorName>/account/sys" {
   capabilities = ["read"]
}
path "nats-secrets/jwt/operator/<operatorName>/account/sys" { 
   capabilities = ["read"]
}
EOF
```

## 3. create role

By creating the kubernetes role, the service account for `nats-auth-config` is allowed to access paths specified in the policy `nats-auth-config`.

``` console
vault write auth/kubernetes/role/nats-auth-config  bound_service_account_names=nats-auth-config  bound_service_account_namespaces="*" policies=nats-auth-config ttl=24h
```

# Chart configuration 

You can deploy backend and core components independently by enabling/disabling them:

| Component                 | Description                                                      | Default value |
| ------------------------- | ---------------------------------------------------------------- | ------------- |
| operatorName              | Specifies the name of the nats operator                          | true          |
| natsAuthConfig.enabled    | Specifies if the backend cluster parts should be deployed        | true          |
| networkBaseConfig.enabled | Specifies if the network base config configmap should be created | true          |
| operator.enabled          | Specifies if the operator should be created                      | true          |
| sysAccount.enabled        | Specifies if the sys-account should be created                   | true          |


## natsAuthConfig configuration

| Component                                              | Description                                                                   | Default value                  |
| ------------------------------------------------------ | ----------------------------------------------------------------------------- | ------------------------------ |
| natsAuthConfig.resolver.address                        | Specifies the address of the nats server                                      | "nats://nats.default.svc:4222" |
| natsAuthConfig.resolver.config.type                    | Specifies the type of the nats resolver                                       | full                           |
| natsAuthConfig.resolver.config.dir                     | Specifies the directory to cache JWTs                                         | "/data/jwt"                    |
| natsAuthConfig.resolver.config.allow_delete            | Specifies if account information can be deleted                               | true                           |
| natsAuthConfig.resolver.config.interval                | Specifies the resolver interval                                               | "2m"                           |
| natsAuthConfig.resolver.config.timeout                 | Specifies the resolver timeout                                                | "1.9s"                         |
| natsAuthConfig.nats.authConfigmapDestination.name      | Specifies the name of the configmap where the auth config will be stored      | nats-auth-config               |
| natsAuthConfig.nats.authConfigmapDestination.namespace | Specifies the namespace of the configmap where the auth config will be stored | nats                           |

