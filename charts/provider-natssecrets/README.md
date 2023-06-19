# provider-natssecrets

This helm chart is part of edgefarm.network and contains provider-natssecrets, a Crossplane provider for [vault-plugin-secrets-nats](https://github.com/edgefarm/vault-plugin-secrets-nats) to manage NATS secrets resources like operator, account, user, service account and service account user nkey, JWT and creds file resources.
Note, that for this to work you need to have Crossplane (>=v1.11.3) installed and configured.

## Prerequisites

    Kubernetes 1.22+
    Helm 3.2.0+
    Crossplane 1.11.3+
    Vault with vault-plugin-secrets-nats 1.3.2+
