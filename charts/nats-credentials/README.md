# nats-credentials

This helm chart is part of edgefarm.network and manages the operator, sys-account and nats auth config.
An operator, sys-account and sys-account user will be created and the auth-config for nats will be created and stored in a configmap.

Note, that for this to work you need to have Crossplane (>=v1.11.3), Hashicorp Vault with plugin natssecrets and provider-natssecrets installed and configured.

# Prerequisites

    Kubernetes 1.22+
    Helm 3.2.0+
    Crossplane 1.11.3+
    Hashicorp Vault with plugin natssecrets 1.3.4+ (https://github.com/edgefarm/vault-plugin-secrets-nats)
    provider-natsecrets v0.2.2+ (https://github.com/edgefarm/provider-natssecrets)

# Configuration 

You can deploy backend and core components independently by enabling/disabling them:

| Component       | Description                                               | Default value |
| --------------- | --------------------------------------------------------- | ------------- |
| operatorName    | Specifies the name of the nats operator                   | true          |
| core.enabled    | Specifies if the core cluster parts should be deployed    | true          |
| backend.enabled | Specifies if the backend cluster parts should be deployed | true          |

## backend configuration

| Component                                       | Description                                                                   | Default value                  |
| ----------------------------------------------- | ----------------------------------------------------------------------------- | ------------------------------ |
| backend.createOperator                          | Specifies if the operator should be created                                   | true                           |
| backend.createSysaccount                        | Specifies if the sys-account should be created                                | true                           |
| backend.resolver.address                        | Specifies the address of the nats server                                      | "nats://nats.default.svc:4222" |
| backend.resolver.config.type                    | Specifies the type of the nats resolver                                       | full                           |
| backend.resolver.config.dir                     | Specifies the directory to cache JWTs                                         | "/data/jwt"                    |
| backend.resolver.config.allow_delete            | Specifies if account information can be deleted                               | true                           |
| backend.resolver.config.interval                | Specifies the resolver interval                                               | "2m"                           |
| backend.resolver.config.timeout                 | Specifies the resolver timeout                                                | "1.9s"                         |
| backend.nats.authConfigmapDestination.name      | Specifies the name of the configmap where the auth config will be stored      | nats-auth-config               |
| backend.nats.authConfigmapDestination.namespace | Specifies the namespace of the configmap where the auth config will be stored | nats                           |