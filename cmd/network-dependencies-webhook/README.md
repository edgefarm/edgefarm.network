# network-dependencies-webhook

This validation webhook takes care of the order of deleteion of the network components that have dependencies between them.

## How it works

There are a number of resources that take care of the network components in the cluster. These resources are:

1. account (accounts.issue.natssecrets.crossplane.io/v1alpha1)
2. system-user (users.issue.natssecrets.crossplane.io/v1alpha1)
3. sys-account-user (users.issue.natssecrets.crossplane.io/v1alpha1)
4. custom users (users.issue.natssecrets.crossplane.io/v1alpha1)
5. providerConfig for nats (providerconfigs.nats.crossplane.io/v1alpha1)
6. edgeNetwork creating leaf-nats pods (xedgenetworks.streams.network.edgefarm.io/v1alpha1)
7. streams (xstreams.streams.network.edgefarm.io/v1alpha1)
8. consumers (xconsumers.streams.network.edgefarm.io/v1alpha1)

The webhook intercepts delete requests for each resource and takes allows or blocks the deletion based on the dependencies between the resources. The dependencies are as follows:

- always allow delete custom users
- always allow delete consumers
- no consumers for streams? allow delete streams
- no streams? allow delete pods
- no streams? allow delete providerConfig
- no providerConfig? allow delete system-user
- no system-user and no custom user? allow delete account
- no account? allow delete sys-account-user

